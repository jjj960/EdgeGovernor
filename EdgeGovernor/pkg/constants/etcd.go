package constants

var ( //etcd集群信息
	EtcdCert    = "/data/etcd/ssl/client.pem" //证书路径
	EtcdCertKey = "/data/etcd/ssl/client-key.pem"
	EtcdCa      = "/data/etcd/ssl/ca.pem"

	Etcd1Enable = true //表示该数据库是否可用
	Etcd1Name   = "etcd1"
	Etcd1ID     uint64
	Etcd1URL    = "https://192.168.47.152:2379" //数据库地址

	Etcd2Enable = true
	Etcd2Name   = "etcd2"
	Etcd2ID     uint64
	Etcd2URL    = "https://192.168.47.153:2379"

	Etcd3Enable = true
	Etcd3Name   = "etcd3"
	Etcd3ID     uint64
	Etcd3URL    = "https://192.168.47.154:2379"

	Etcd4Enable = true
	Etcd4Name   = "etcd4"
	Etcd4ID     uint64
	Etcd4URL    = "https://192.168.47.155:2379"

	OriginalEndPoints = []string{"https://192.168.47.152:2379", "https://192.168.47.153:2379", "https://192.168.47.154:2379", "https://192.168.47.155:2379"}
)
