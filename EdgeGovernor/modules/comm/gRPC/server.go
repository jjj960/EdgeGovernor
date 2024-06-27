package gRPC

import (
	"EdgeGovernor/modules/comm/gRPC/proto"
	"EdgeGovernor/modules/taskDeploy/workflow"
	"EdgeGovernor/pkg/cache/algorithm"
	"EdgeGovernor/pkg/cache/task"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/docker/resource"
	"EdgeGovernor/pkg/logging"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var NodeWorkloads [][]string
var electing = false

type server struct {
	proto.UnimplementedNodeCommServer
}

type GRPCServer struct {
	port   int
	server *grpc.Server
}

func (s *server) Messaging(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	//fmt.Println(, )
	typeMsg, sourceNode, _, detail := msgProcess(in)
	//fmt.Println(typeMsg)
	var replyMsg string

	switch typeMsg {
	case "get leader":
		if detail == constants.WebKey {
			replyMsg = constants.Leader
		} else {
			replyMsg = "Password error!"
		}

	case "survival testing":
		HeartBeatTimer.Reset()
		replyMsg = "survival"

	case "task deployment":
		var task models.Task
		err := utils.Jsoniter.Unmarshal([]byte(detail), &task)
		if err != nil {
			log.Println("解析JSON失败：", err)
		}

		re, errs := resource.CreateResourceContainer(task.Name, task.Image, task.RequestCPU, task.RequestMem, task.PersistData, task.DataDir)
		if errs != nil {
			replyMsg = re
		}

	case "workflow job send":
		var data models.Job
		err := utils.Jsoniter.Unmarshal([]byte(detail), &data)
		if err != nil {
			log.Println("解析JSON失败：", err)
		}
		runtimeFloat, _ := strconv.ParseFloat(data.Runtime, 64)

		go func() {
			seconds := int(runtimeFloat)                                  // 转换为整数秒数
			milliseconds := int((runtimeFloat - float64(seconds)) * 1000) // 获取小数部分转换为毫秒
			time.Sleep(time.Duration(seconds)*time.Second + time.Duration(milliseconds)*time.Millisecond)

			utils.JobOperationChannel <- models.JobOperation{
				Type: "check",
				Job:  data,
			}
		}()

	case "job completion report":
		var data models.Job
		err := utils.Jsoniter.Unmarshal([]byte(detail), &data)
		if err != nil {
			log.Println("解析JSON失败：", err)
		}

		utils.JobOperationChannel <- models.JobOperation{
			Type: "publish",
			Job:  data,
		}

	case "workflow job deployment":
		var data models.Job
		err := utils.Jsoniter.Unmarshal([]byte(detail), &data)
		if err != nil {
			log.Println("解析JSON失败：", err)
		}
		dataDir := []string{fmt.Sprintf("/data/menet/workflow/%s/%s", data.WorkflowName, data.Name)}

		re, errs := resource.CreateContainer(data.Name, data.Image, data.WorkflowName, dataDir)
		if errs != nil {
			replyMsg = errs.Error()
		} else {
			replyMsg = re
			runtimeFloat, _ := strconv.ParseFloat(data.Runtime, 64)

			go func() {
				seconds := int(runtimeFloat)                                  // 转换为整数秒数
				milliseconds := int((runtimeFloat - float64(seconds)) * 1000) // 获取小数部分转换为毫秒
				time.Sleep(time.Duration(seconds)*time.Second + time.Duration(milliseconds)*time.Millisecond)

				utils.JobOperationChannel <- models.JobOperation{
					Type: "check",
					Job:  data,
				}
			}()
		}

	case "candidate update": //更新集群的最新候选人
		constants.Candidate = detail
		fmt.Println("Update Candidate to:", detail)
		id, _ := utils.GetID()
		logEntry := models.OperationLog{
			ID:            id,
			NodeName:      constants.Hostname,
			NodeIP:        constants.IP,
			OperationType: "candidate update",
			Description:   fmt.Sprintf("candidate update to: %s", detail),
			Result:        true,
			CreatedAt:     time.Now(),
		}
		duckdb.InsertOperationLog(logEntry)

	case "leader election": //该消息只有Candidate才能收到，收到消息后负责收集各个节点的负载数据，选举出新领导者

		if !electing { //正在选举中,消息忽略
			fmt.Println("Start conducting leader elections")
			constants.ClusterStatus = "selfGovernment"
			electing = true
			nodeWorkloads := collectFollowerMachineStatus()
			result, err := NodeLoadAssessmentRequest(algorithm.GetAlgorithmURL("NodeLoadAssessmentAlgorithm", "NodeLoadAssessment"), nodeWorkloads)
			if err != nil {
				fmt.Println("Failed to obtain scheduling data:", err)
			}
			newLeader := strings.ReplaceAll(result, "\"", "")
			fmt.Println("After evaluation, the new Leader is:", newLeader)
			id, _ := utils.GetID()
			logEntry := models.OperationLog{
				ID:            id,
				NodeName:      constants.Hostname,
				NodeIP:        constants.IP,
				OperationType: "leader election",
				Description:   fmt.Sprintf("Start the leader election initiated by node %s", constants.Hostname),
				Result:        true,
				CreatedAt:     time.Now(),
			}
			duckdb.InsertOperationLog(logEntry)
			if constants.Hostname == newLeader { //如果选举出的新领导者为自己
				utils.Broadcast("leader change", newLeader)
				utils.ModuleControlChannel <- false
			} else { //如果选举出的领导者不是自己，则发送消息通知该节点
				IP, _ := utils.NodeTables.GetNodeIP(newLeader)
				utils.SingleSend(IP, newLeader, "leader elected", "")
			}
			electing = false
		}

		replyMsg = "accept"

	case "leader elected": //当选举出新领导者后，通知该节点，该节点会收到消息
		fmt.Println("This node has received a message and has been selected as the new leader")
		id, _ := utils.GetID()
		logEntry := models.OperationLog{
			ID:            id,
			NodeName:      constants.Hostname,
			NodeIP:        constants.IP,
			OperationType: "leader elected",
			Description:   fmt.Sprintf("Node %s has become a new leader in the cluster", constants.Hostname),
			Result:        true,
			CreatedAt:     time.Now(),
		}
		duckdb.InsertOperationLog(logEntry)

		workflow.GetBackupworkflowQueue()
		utils.GetNodeTablesMap()
		algorithm.GetAlgorithmStatusMap()
		task.GetNodeSimpleTaskMap()
		if constants.Hostname != "cloud" { //本机当选领导者后，本机并不是初始领导者cloud
			constants.ClusterStatus = "selfGovernment"
		} else {
			constants.ClusterStatus = "coordination"
		}
		utils.NodeTables.UpdateNodeRole(constants.Leader, "Follower") //将原先的Leader更新为Follower
		constants.Leader = constants.Hostname
		utils.NodeTables.UpdateNodeRole(constants.Hostname, "Leader") //将本机更新为Leader

		utils.ModuleControlChannel <- false //模块切换

		utils.Broadcast("leader change", constants.Hostname) //向集群广播本机为Leader

		replyMsg = "accept"

	case "leader restore":

		var data models.ClusterMsg
		err := utils.Jsoniter.Unmarshal([]byte(detail), &data)
		if err != nil {
			log.Println("Parsing JSON failed:", err)
		}
		constants.LiveNodeCount = data.LiveNodeCount
		constants.NodeCount = data.NodeCount
		constants.Leader = constants.Hostname
		constants.Candidate = data.Candidate

		utils.ModuleControlChannel <- false

		utils.Broadcast("leader change", constants.Hostname)

	case "leader change": //集群中所有节点都会收到该Leader转换消息
		fmt.Println("Received new Leader notification")
		HeartBeatTimer.Reset()
		if detail != "cloud" { //当前Leader不为初始Leader
			constants.ClusterStatus = "selfGovernment"
		} else {
			constants.ClusterStatus = "coordination"
		}

		utils.NodeTables.UpdateNodeRole(constants.Leader, "Follower")
		constants.Leader = detail
		utils.NodeTables.UpdateNodeRole(constants.Leader, "Leader")

		replyMsg = "accept"

	case "message synchronization": //集群状态同步
		fmt.Println("Synchronizing cluster status from node:", sourceNode)
		var data models.ClusterMsg
		err := utils.Jsoniter.Unmarshal([]byte(detail), &data)
		if err != nil {
			log.Println("Parsing JSON failed:", err)
		}

		constants.LiveNodeCount = data.LiveNodeCount
		constants.NodeCount = data.NodeCount
		constants.Leader = data.Leader
		constants.Candidate = data.Candidate

	case "machine workload":

		var data models.Hostload
		err := utils.Jsoniter.Unmarshal([]byte(detail), &data)
		if err != nil {
			log.Println("Parsing JSON failed:", err)
		}
		var nodes []string
		if data.CPUUsagePercent >= 80.0 || data.MemoryUsedPercent >= 80.0 {
			id := utils.SnowFlake.Generate()
			utils.AlarmMsgChannel <- models.Msg{
				ID:           id.Int64(),
				GenerateTime: id.Time(),
				Tpye:         "Host Resource warning",
				Detail:       []byte(fmt.Sprintf("Node %s has insufficient resources, CPU load is %f, memory load is %f", data.Hostname, data.CPUUsagePercent, data.MemoryUsedPercent)),
				Status:       false,
			}
		}
		sourceIP, _ := utils.NodeTables.GetNodeIP(sourceNode)
		nodes = append(nodes, sourceNode, sourceIP, utils.Float64toString(data.CPUUsagePercent),
			utils.Int64toString(data.CPUCapacity),
			utils.Int64toString(data.CPUResidue),
			utils.Float64toString(data.MemoryUsedPercent),
			utils.Int64toString(data.MemoryCapacity),
			utils.Int64toString(data.MemoryResidue),
			utils.Float64toString(data.DiskUsedPercent),
			utils.Int64toString(data.DiskCapacity),
			utils.Int64toString(data.DiskResidue),
			utils.Float64toString(data.BytesRecv),
			utils.Float64toString(data.BytesSent),
			utils.Float64toString(data.BandWidth))
		NodeWorkloads = append(NodeWorkloads, nodes)
		if len(NodeWorkloads) == constants.LiveNodeCount {
			result, err := utils.SchedulingRequest("http://192.168.47.128:50052/scheduler", NodeWorkloads)
			if err != nil {
				fmt.Println("Failed to obtain scheduling data:", err)
			}
			candidate := strings.ReplaceAll(result, "\"", "")
			NodeWorkloads = NodeWorkloads[:0]
			constants.Candidate = candidate
			fmt.Println("Update Candidate to:", candidate)
			utils.Broadcast("candidate update", candidate)
		}

		replyMsg = "Confirm receipt"

	case "workload report": //接收来自其他节点的负载状态请求,获取负载状态后返回
		result := logging.GetHostWorkload()
		replyMsg = string(result)

	default:
		fmt.Println(detail)
	}

	return &proto.Response{Message: replyMsg}, nil
}

func msgProcess(in *proto.Request) (string, string, string, string) {
	ty := in.GetType()
	sn := in.GetSourceNode()
	tn := in.GetTargetNode()
	details := in.GetDetail()

	typeMsg := string(sec.Safer.Decrypt(ty))
	sourceNode := string(sec.Safer.Decrypt(sn))
	targetNode := string(sec.Safer.Decrypt(tn))
	detail := string(sec.Safer.Decrypt(details))

	if typeMsg == "get leader" {
		return typeMsg, sourceNode, targetNode, detail
	}

	return typeMsg, sourceNode, targetNode, detail
}

func NewGRPCServer(port int) *GRPCServer {
	return &GRPCServer{
		port:   port,
		server: grpc.NewServer(),
	}
}

func (gs *GRPCServer) Start() {
	log.SetFlags(log.Ltime | log.Llongfile)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gs.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	reflection.Register(gs.server)
	proto.RegisterNodeCommServer(gs.server, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := gs.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (gs *GRPCServer) Close() {
	gs.server.GracefulStop()
}

func collectFollowerMachineStatus() [][]string {
	nodeMsg := utils.NodeTables.GetLiveNode()
	for i, node := range nodeMsg {
		if node[0] == constants.Hostname {
			var data models.Hostload
			result := logging.GetHostWorkload()
			err := utils.Jsoniter.Unmarshal(result, &data)
			if err != nil {
				log.Println("Parsing JSON failed:", err)
				return nil
			}
			node = append(node, utils.Float64toString(data.CPUUsagePercent),
				utils.Int64toString(data.CPUCapacity),
				utils.Int64toString(data.CPUResidue),
				utils.Float64toString(data.MemoryUsedPercent),
				utils.Int64toString(data.MemoryCapacity),
				utils.Int64toString(data.MemoryResidue),
				utils.Float64toString(data.DiskUsedPercent),
				utils.Int64toString(data.DiskCapacity),
				utils.Int64toString(data.DiskResidue),
				utils.Float64toString(data.BytesRecv),
				utils.Float64toString(data.BytesSent),
				utils.Float64toString(data.BandWidth))

			nodeMsg[i] = node
		} else {
			fmt.Println("Collecting information from other nodes:", node[0])
			var data models.Hostload
			result, _ := utils.SingleSend(node[1], node[0], "workload report", "")
			fmt.Println("The collection result is:", result)
			err := utils.Jsoniter.Unmarshal([]byte(result), &data)
			if err != nil {
				log.Println("Parsing JSON failed:", err)
				return nil
			}
			node = append(node, utils.Float64toString(data.CPUUsagePercent),
				utils.Int64toString(data.CPUCapacity),
				utils.Int64toString(data.CPUResidue),
				utils.Float64toString(data.MemoryUsedPercent),
				utils.Int64toString(data.MemoryCapacity),
				utils.Int64toString(data.MemoryResidue),
				utils.Float64toString(data.DiskUsedPercent),
				utils.Int64toString(data.DiskCapacity),
				utils.Int64toString(data.DiskResidue),
				utils.Float64toString(data.BytesRecv),
				utils.Float64toString(data.BytesSent),
				utils.Float64toString(data.BandWidth))

			nodeMsg[i] = node
		}
	}
	//fmt.Println(nodeMsg)
	return nodeMsg
}

func NodeLoadAssessmentRequest(url string, nodeMsg [][]string) (string, error) {
	//url := "http://192.168.47.128:50052/scheduler" // 目标URL

	jsonData, err := utils.Jsoniter.Marshal(nodeMsg) // 将二维数组转换为JSON格式
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
