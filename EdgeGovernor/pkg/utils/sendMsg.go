package utils

import (
	"EdgeGovernor/modules/comm/gRPC/proto"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/sec"
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func SchedulingRequest(url string, nodeMsg [][]string) (string, error) {
	//url := "http://192.168.47.128:50052/scheduler" // 目标URL

	jsonData, err := Jsoniter.Marshal(nodeMsg) // 将二维数组转换为JSON格式
	if err != nil {
		return "", fmt.Errorf("Error converting to JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) // 发送POST请求
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}

	result := string(body)

	return result, nil
}

func SingleSend(addr string, targetNode string, types string, details string) (string, error) { //信息单点发送
	conn, err := grpc.Dial(addr+":50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		NodeTables.UpdateNodeStatus(targetNode, "Offline")
		return "", err
	}
	defer conn.Close()

	c := proto.NewNodeCommClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.Messaging(ctx, &proto.Request{
		Type:       sec.Safer.Encrypt([]byte(types)),
		SourceNode: sec.Safer.Encrypt([]byte(constants.Hostname)),
		TargetNode: sec.Safer.Encrypt([]byte(targetNode)),
		Detail:     sec.Safer.Encrypt([]byte(details)),
	})
	if err != nil {
		log.Printf("Failed to receive messages from %s: %v", targetNode, err)
		return "", err
	}

	nodeOldStatus, _ := NodeTables.GetNodeStatus(targetNode)
	if nodeOldStatus == "Offline" && targetNode != "cloud" {
		ClusterUpdate(targetNode)
	}
	NodeTables.UpdateNodeStatus(targetNode, "Online")
	return r.GetMessage(), nil
}

func Broadcast(types string, details string) { //向集群所有存活的节点发送消息
	for node := range NodeTables.Iter() {
		if node.Key == constants.Hostname || node.Value.Status == "Offline" {
			continue
		}
		_, _ = SingleSend(node.Value.IP, node.Key, types, details)
	}
}

func HeartBeat() { //心跳检测,向所有除自己以外的其他节点发送消息
	for node := range NodeTables.Iter() {
		if node.Key == constants.Hostname {
			continue
		}
		_, err := SingleSend(node.Value.IP, node.Key, "survival testing", "")
		if err == nil && node.Key == "cloud" { //cloud节点恢复通信
			clusterData := &models.ClusterMsg{
				LiveNodeCount: constants.LiveNodeCount,
				NodeCount:     constants.NodeCount,
				Leader:        constants.Leader,
				Candidate:     constants.Candidate,
			}
			jsonData, err := Jsoniter.Marshal(clusterData)
			if err != nil {
				fmt.Println("Error:", err)
			}
			SingleSend(node.Value.IP, node.Key, "leader restore", string(jsonData))
			ModuleControlChannel <- true
		}
	}
}

func ClusterUpdate(targetNode string) {
	ip, _ := NodeTables.GetNodeIP(targetNode)
	clusterData := &models.ClusterMsg{
		LiveNodeCount: constants.LiveNodeCount,
		NodeCount:     constants.NodeCount,
		Leader:        constants.Leader,
		Candidate:     constants.Candidate,
	}
	jsonData, err := Jsoniter.Marshal(clusterData)
	if err != nil {
		fmt.Println("Error:", err)
	}
	SingleSend(ip, targetNode, "message synchronization", string(jsonData))
}
