package utils

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/etcd"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/docker/docker/client"
	_ "github.com/marcboeker/go-duckdb"
	"go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

var DockerCli *client.Client
var ETCDCli *clientv3.Client
var DuckDBCli *sql.DB
var K8sClientset *kubernetes.Clientset

func GetDuckDBCli() {
	db, err := sql.Open("duckdb", fmt.Sprintf("/data/menet/%s/db/%sdb", constants.Hostname, constants.Hostname))
	if err != nil {
		log.Fatal(err)
	}
	DuckDBCli = db
}

func GetDockerCli() {
	//cli客户端对象
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	DockerCli = cli
}

func GetK8sCli() {
	config, err := clientcmd.BuildConfigFromFlags("", "/etc/kubernetes/kubelet.conf")
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connect this cluster's k8s successfully.")
	}
	K8sClientset = clientset
}

func GetETCDCli() {
	dialTimeout := 5 * time.Second

	cert, err := tls.LoadX509KeyPair(constants.EtcdCert, constants.EtcdCertKey)
	if err != nil {
		log.Fatal(err)
	}

	caData, err := ioutil.ReadFile(constants.EtcdCa)
	if err != nil {
		log.Fatal(err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	cfg := clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   etcd.GetAvailableETCDEndPoints(),
		TLS:         _tlsConfig,
	}

	cli, err := clientv3.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	ETCDCli = cli
}
