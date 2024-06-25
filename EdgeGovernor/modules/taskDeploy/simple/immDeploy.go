package simple

import (
	"EdgeGovernor/pkg/alogorithm"
	algorithm2 "EdgeGovernor/pkg/cache/algorithm"
	task2 "EdgeGovernor/pkg/cache/task"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/etcd/resource"
	"EdgeGovernor/pkg/k8s"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func taskPublish(task *models.Task) error {
	fmt.Println("Attempting to deploy task:", task.Name)
	maxHost := task.DeployNode //最优调度节点
	var nodeMsg [][]string
	var taskMsg []string
	nodeMsg = collectFollowerMachineStatus()
	//requestMsg = append(requestMsg, nodeMsg...)
	taskMsg = append(taskMsg, task.Name, task.DeployNode,
		utils.Int64toString(task.RequestCPU),
		utils.Int64toString(task.RequestMem),
		utils.Int64toString(task.RequestNet),
		utils.Int64toString(task.RequestDisk),
		strconv.Itoa(task.Scalable))

	if task.DeployNode == "" { //部署节点为空,说明采用全自动部署
		//请求调度算法
		//var requestMsg [][]string
		maxHost, _ = GetSchedulerScore(nodeMsg, task.SchedulingAlgorithm)
		fmt.Println("调度的最优节点为: ", maxHost)
	}
	check := alogorithm.CheckNodeResource(maxHost, nodeMsg, taskMsg)
	if !check {
		id := utils.SnowFlake.Generate()
		utils.AlarmMsgChannel <- models.Msg{
			ID:           id.Int64(),
			GenerateTime: id.Time(),
			Tpye:         "Task deployment failed",
			Detail:       []byte(fmt.Sprintf("Cluster has insufficient resources, Task %s deployment failed", task.Name)),
			Status:       false,
		}

		return fmt.Errorf("Task %s deployment failed: The cluster cannot meet the task requirements", task.Name)
	} else {
		totalRequestCPU, _, _, _ := GetTotalRequestResource()

		cpu, mem := alogorithm.ResourceAllocation(maxHost, nodeMsg, taskMsg, totalRequestCPU) //资源分配
		var re string
		var errs error
		if constants.ClusterStatus == "coordination" {
			task.Type = "Pod"
			re, errs = k8s.CreateResourcePod(maxHost, task.Name, task.Image, utils.Int64toString(cpu), utils.Int64toString(mem),
				task.PersistData, utils.Int64toString(task.RequestDisk), nil)
		} else {
			task.Type = "Container"
			ip, _ := utils.NodeTables.GetNodeIP(maxHost)

			jsonData, err := utils.Jsoniter.Marshal(task)
			if err != nil {
				fmt.Println("Error:", err)
			}
			_, err = utils.SingleSend(ip, maxHost, "task deployment", string(jsonData))
			if err != nil {
				return fmt.Errorf("Task %s deployment failed: %s", task.Name, err)
			}

		}
		if re == "success" {
			fmt.Println("结果为:", re)
			fmt.Println(task.Name, maxHost, cpu)
			task2.UpdateTaskDeployment(maxHost, task.Name, true)
			resource.UpdateTaskMsg(task.Name, maxHost, cpu, task.Type, "Deployed")

		} else {
			id := utils.SnowFlake.Generate()
			utils.AlarmMsgChannel <- models.Msg{
				ID:           id.Int64(),
				GenerateTime: id.Time(),
				Tpye:         "Task deployment failed",
				Detail:       []byte(fmt.Sprintf("Task %s deployment failed, %s", task.Name, errs.Error())),
				Status:       false,
			}
			return fmt.Errorf("Task %s deployment failed: %s", task.Name, errs)
		}
	}

	return nil
}

func collectFollowerMachineStatus() [][]string {
	nodeMsg := utils.NodeTables.GetLiveNode()
	for i, node := range nodeMsg {
		if node[0] == constants.Hostname && constants.ClusterStatus == "coordination" {
			continue
		}
		var data models.Hostload
		result, _ := utils.SingleSend(node[1], node[0], "workload report", "")

		err := utils.Jsoniter.Unmarshal([]byte(result), &data)
		if err != nil {
			log.Println("解析JSON失败：", err)
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
	//fmt.Println(nodeMsg)
	return nodeMsg
}

func GetSchedulerScore(nodeMsg [][]string, algorithm string) (string, error) {

	result, err := utils.SchedulingRequest(algorithm2.GetAlgorithmURL(algorithm, "Schedule"), nodeMsg)
	if err != nil {
		return "", fmt.Errorf("Failed to obtain scheduling data:", err)
	}

	return strings.ReplaceAll(result, "\"", ""), nil
}
