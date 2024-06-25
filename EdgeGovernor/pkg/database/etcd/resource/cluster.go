package resource

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/etcd"
	"EdgeGovernor/pkg/utils"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

/*
etcd数据库集群相关操作
*/

func GetETCDMemberMsg() error { //获取etcd集群成员的详细信息
	ctx := context.Background()

	response, err := utils.ETCDCli.MemberList(ctx)
	if err != nil {
		return fmt.Errorf("Failed to obtain member information: %s", err)
	}

	for _, member := range response.Members {
		fmt.Printf("ID: %d\n", member.ID)
		fmt.Printf("Name: %s\n", member.Name)
		fmt.Printf("PeerURLs: %s\n", member.PeerURLs)
		fmt.Printf("ClientURLs: %s\n", member.ClientURLs)
		fmt.Println("---")
	}

	return nil
}

type nodeHealth struct {
	Node   string `json:"node"`
	Health bool   `json:"health"`
}

func GetETCDClusterHealth() (bool, []string) { //检查etcd集群是否健康
	cert, err := tls.LoadX509KeyPair(constants.EtcdCert, constants.EtcdCertKey)
	if err != nil {
		log.Fatal(err)
	}

	// 加载根证书
	caCert, err := ioutil.ReadFile(constants.EtcdCa)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 配置TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	//

	// 执行节点健康检查
	var wg sync.WaitGroup
	nodeHealthCh := make(chan nodeHealth, 3)

	for _, ep := range etcd.GetAvailableETCDEndPoints() {
		wg.Add(1)
		go func(ep string) {
			defer wg.Done()
			epCli, err := clientv3.New(clientv3.Config{
				Endpoints:   []string{ep},
				DialTimeout: 2 * time.Second,
				TLS:         tlsConfig,
			})
			if err != nil {
				nodeHealthCh <- nodeHealth{Node: ep, Health: false}
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err = epCli.Get(ctx, "health")
			cancel()

			if err == nil {
				nodeHealthCh <- nodeHealth{Node: ep, Health: true}
			} else {
				nodeHealthCh <- nodeHealth{Node: ep, Health: false}
			}
		}(ep)
	}

	wg.Wait()
	close(nodeHealthCh)

	unhealthyNodes := []string{}
	for nh := range nodeHealthCh {
		if !nh.Health {
			unhealthyNodes = append(unhealthyNodes, nh.Node)
		}
	}

	if len(unhealthyNodes) > 0 {
		return false, unhealthyNodes //集群不健康,并返回不健康的节点信息
	} else {
		return true, nil //集群健康
	}
}

func GetETCDLeader() (uint64, error) { //获取ETCD集群的Leader ID
	ctx := context.Background()

	response, err := utils.ETCDCli.Status(ctx, etcd.GetAvailableETCDEndPoints()[0])
	if err != nil {
		return 0, fmt.Errorf("Failed to obtain ETCD Leader: %s", err)
	}

	return response.Leader, nil
}

func GetETCDAllKeys() error { //获取数据库中所有key
	ctx := context.Background()
	prefix := "/"

	// 执行查询
	response, err := utils.ETCDCli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("Failed to obtain all keys: %s", err)
	}

	// 遍历查询结果
	for _, kv := range response.Kvs {
		fmt.Printf("Key: %s\n", kv.Key)
	}

	return nil
}

func MemberAdd(name string, url string) error {
	cert, err := tls.LoadX509KeyPair(constants.EtcdCert, constants.EtcdCertKey)
	if err != nil {
		return fmt.Errorf("Certificate loading failed in adding ETCD cluster members: %s", err)
	}

	// 加载根证书
	caCert, err := ioutil.ReadFile(constants.EtcdCa)
	if err != nil {
		return fmt.Errorf("Certificate loading failed in adding ETCD cluster members: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 配置TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// 新成员的监听地址和通信地址
	newMemberName := name
	newMemberListenAddr := "https://" + url + ":2380"
	newMemberPeerAddr := "https://" + url + ":2380"

	// 创建etcd客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcd.GetAvailableETCDEndPoints(),
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
	})
	if err != nil {
		return fmt.Errorf("Failed to create client in adding ETCD cluster members: %s", err)
	}
	defer cli.Close()

	// 获取初始成员列表
	resp, err := cli.MemberList(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to obtain initial member list in adding ETCD cluster members: %s", err)
	}

	// 检查新成员是否已经加入到集群
	for _, m := range resp.Members {
		if m.Name == newMemberName {
			return fmt.Errorf("Member %s already exists in the cluster", newMemberName)
		}
	}

	// 获取初始成员的URL列表
	existingPeerURLs := make([]string, len(resp.Members))
	for i, m := range resp.Members {
		existingPeerURLs[i] = m.PeerURLs[0]
	}

	// 将新成员加入集群
	addResp, err := cli.MemberAdd(context.Background(), []string{newMemberPeerAddr})
	if err != nil {
		return fmt.Errorf("Adding ETCD cluster members failed: %s", err)
	}

	// 获取新成员的ID
	newMemberID := addResp.Member.ID

	// 添加新成员的监听地址和通信地址
	_, err = cli.MemberUpdate(context.Background(), newMemberID, []string{newMemberListenAddr})
	if err != nil {
		log.Fatal(err)
	}

	// 更新初始成员的URL列表，包括新成员的URL
	allPeerURLs := append(existingPeerURLs, newMemberPeerAddr)

	// 将更新后的URL列表应用到所有成员
	for _, m := range resp.Members {
		_, err = cli.MemberUpdate(context.Background(), m.ID, allPeerURLs)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("New member %s added to the cluster", newMemberName)

	//TODO: 通知目标节点启动etcd数据库
	err = etcdClusterSystemctlStart(name, url)
	if err != nil {
		return err
	}

	return nil
}

func etcdClusterSystemctlStart(name string, url string) error { //启动linux的service
	ip := url
	peerPort := 2380
	listenPort := 2379

	data := `[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
WorkingDirectory=/data/etcd/
ExecStart=/data/etcd/bin/etcd \
 --name=%s \
 --cert-file=/data/etcd/ssl/server.pem \
 --key-file=/data/etcd/ssl/server-key.pem \
 --peer-cert-file=/data/etcd/ssl/peer.pem \
 --peer-key-file=/data/etcd/ssl/peer-key.pem \
 --trusted-ca-file=/data/etcd/ssl/ca.pem \
 --peer-trusted-ca-file=/data/etcd/ssl/ca.pem \
 --initial-advertise-peer-urls=%s \
 --listen-peer-urls=%s \
 --listen-client-urls=%s \
 --advertise-client-urls=%s \
 --initial-cluster-token=etcd-cluster-0 \
 --initial-cluster=etcd1=%s,etcd2=%s,etcd3=%s \
 --initial-cluster-state=%s \
 --data-dir=/data/etcd \
 --snapshot-count=50000 \
 --auto-compaction-retention=1 \
 --max-request-bytes=10485760 \
 --quota-backend-bytes=8589934592
Restart=always
RestartSec=15
LimitNOFILE=65536
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target`

	// 替换URL中的占位符
	execStart := fmt.Sprintf(data,
		name,
		fmt.Sprintf("https://%s:%d", ip, peerPort),
		fmt.Sprintf("https://%s:%d", ip, peerPort),
		fmt.Sprintf("https://%s:%d", ip, listenPort),
		fmt.Sprintf("https://%s:%d", ip, listenPort),
		fmt.Sprintf(constants.Etcd1URL),
		fmt.Sprintf(constants.Etcd2URL),
		fmt.Sprintf(constants.Etcd3URL),
		fmt.Sprintf("existing"),
	)

	err := etcd.ServiceCreate("etcd", execStart)
	if err != nil {
		return fmt.Errorf("Etcd service creation failed: %s", err)
	}

	return nil
}

//func EtcdSingleSystemctlStart(name string) error {
//	data := `[Unit]
//Description=Etcd Server
//After=network.target
//After=network-online.target
//Wants=network-online.target
//
//[Service]
//Type=notify
//WorkingDirectory=/data/etcd/
//ExecStart=/data/etcd/bin/etcd \
// --name=%s \
// --cert-file=/data/etcd/ssl/server.pem \
// --key-file=/data/etcd/ssl/server-key.pem \
// --peer-cert-file=/data/etcd/ssl/peer.pem \
// --peer-key-file=/data/etcd/ssl/peer-key.pem \
// --trusted-ca-file=/data/etcd/ssl/ca.pem \
// --peer-trusted-ca-file=/data/etcd/ssl/ca.pem \
// --initial-advertise-peer-urls=%s \
// --listen-peer-urls=%s \
// --listen-client-urls=%s \
// --advertise-client-urls=%s \
// --data-dir=/data/etcd \
// --snapshot-count=50000 \
// --auto-compaction-retention=1 \
// --max-request-bytes=10485760 \
// --quota-backend-bytes=8589934592
//Restart=always
//RestartSec=15
//LimitNOFILE=65536
//OOMScoreAdjust=-999
//
//[Install]
//WantedBy=multi-user.target`
//
//	// 替换URL中的占位符
//	execStart := fmt.Sprintf(data,
//		name,
//		fmt.Sprintf("https://%s:%d", ip, 2380),
//		fmt.Sprintf("https://%s:%d", ip, 2380),
//		fmt.Sprintf("https://%s:%d", ip, 2379),
//		fmt.Sprintf("https://%s:%d", ip, 2379),
//	)
//
//	err := etcd.ServiceCreate("etcd", execStart)
//	if err != nil {
//		return fmt.Errorf("Etcd service creation failed: %s", err)
//	}
//
//	return nil
//}

func MemberRemove(nodeName string) error {
	// 获取成员列表
	resp, err := utils.ETCDCli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 检查要移除的成员是否存在
	var memberToRemove *clientv3.Member
	for _, m := range resp.Members {
		if m.Name == nodeName {
			memberToRemove = (*clientv3.Member)(m)
			break
		}
	}

	if memberToRemove == nil {
		return fmt.Errorf("Member does not exist in the etcd cluster")
	}

	// 从集群中移除成员
	_, err = utils.ETCDCli.MemberRemove(context.Background(), memberToRemove.ID)
	if err != nil {
		return fmt.Errorf("Removing ETCD cluster members failed: %s", err)
	}

	return nil
}

func MemberNameToMemberID(memberName string) (uint64, error) {
	// 查询成员信息
	resp, err := utils.ETCDCli.MemberList(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve member list: %w", err)
	}

	// 根据成员名称查找ID
	var memberID uint64
	for _, member := range resp.Members {
		if member.Name == memberName {
			memberID = member.ID
			break
		}
	}

	if memberID == 0 {
		return 0, fmt.Errorf("member not found: %s", memberName)
	}

	return memberID, nil
}

func MemberIDToMemberName(memberID uint64) (string, error) {
	resp, err := utils.ETCDCli.MemberList(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to retrieve member list: %w", err)
	}

	var memberName string
	for _, member := range resp.Members {
		if member.ID == memberID {
			memberName = member.Name
			break
		}
	}

	if memberName == "" {
		return "", fmt.Errorf("member not found with ID: %d", memberID)
	}

	return memberName, nil
}

func EtcdSnapshot() error {
	currentTime := time.Now().Format("20060102_150405") // 根据需要的时间格式定义时间字符串

	filename := fmt.Sprintf("etcd_%s.db", currentTime) // 根据时间字符串构建文件名
	cmd := exec.Command("etcdctl",
		"snapshot", "save", "/data/etcd/backup/"+filename)
	cmd.Env = append(os.Environ(), "ETCDCTL_API=3")
	cmd.Args = append(cmd.Args,
		"--cacert=/data/etcd/ssl/ca.pem",
		"--cert=/data/etcd/ssl/server.pem",
		"--key=/data/etcd/ssl/server-key.pem",
		"--endpoints=https://192.168.47.147:2379")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ETCD backup failed: %s", err)
	}

	return nil
}

func EtcdRestore(filename string) error {
	cmd := exec.Command("etcdctl",
		"snapshot", "save", "/data/etcd/backup/"+filename)
	cmd.Env = append(os.Environ(), "ETCDCTL_API=3")
	cmd.Args = append(cmd.Args,
		"--cacert=/data/etcd/ssl/ca.pem",
		"--cert=/data/etcd/ssl/server.pem",
		"--key=/data/etcd/ssl/server-key.pem",
		"--endpoints=https://192.168.47.147:2379")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ETCD backup failed: %s", err)
	}

	return nil
}
