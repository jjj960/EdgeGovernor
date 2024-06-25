package constants

var (
	ClusterStatus string //表示集群运行状态,分为selfGovernment和coordination两种
)

var ( //用于检测集群节点状态的相关变量
	Timeout          = 3 //选举超时时间（单位：秒）
	HeartBeatTimeout = 4 //心跳检测超时时间
	HeartBeatTimes   = 3 //心跳检测频率（单位：秒）
	Term             = 10
)

var (
	CollectTime int //负载信息收集间隔
)

var ( //集群角色管理
	Candidate string //集群当前候选人
	Leader    string //集群当前领导人
)

var ( //集群节点数量管理
	NodeCount     int
	LiveNodeCount int
)

var MirrorCount int //镜像仓库数量

var WebKey = "12315.hys"
