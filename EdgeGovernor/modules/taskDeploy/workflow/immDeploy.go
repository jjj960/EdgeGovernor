package workflow

import (
	algorithm2 "EdgeGovernor/pkg/cache/algorithm"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/k8s"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"log"
	"strings"
)

func workflowPublish(wf models.Workflow) error {
	fmt.Println("Attempting to deploy workflow:", wf.WorkflowName)

	rootJob, err := GetDAGRootJob(wf.WorkflowName)
	if err != nil {
		return err
	}

	for key, value := range rootJob {
		fmt.Println("Attempting to deploy job:", key)
		JobPublish(value, wf.WorkflowName, wf.SchedulingAlgorithm)
	}

	return nil
}

func JobPublish(job models.Job, workflowName string, schedulingAlgorithm string) error {
	maxHost := job.DeployNode //最优调度节点
	var nodeMsg [][]string
	nodeMsg = collectFollowerMachineStatus()

	if job.DeployNode == "" { //部署节点为空,说明采用全自动部署
		maxHost, _ = GetSchedulerScore(nodeMsg, schedulingAlgorithm)
		fmt.Println("The optimal node for scheduling is:", maxHost)
		job.DeployNode = maxHost
	}
	dataDir := []string{fmt.Sprintf("/data/menet/workflow/%s/%s", workflowName, job.Name)}

	if constants.ClusterStatus == "coordination" {
		_, errs := k8s.CreatePod(maxHost, job.Name, job.Image, workflowName, dataDir)
		if errs != nil {
			return errs
		}
		job.WorkflowName = workflowName
		jsonData, err := utils.Jsoniter.Marshal(job)
		if err != nil {
			fmt.Println("Error:", err)
		}
		ip, _ := utils.NodeTables.GetNodeIP(maxHost)

		utils.SingleSend(ip, maxHost, "workflow job send", string(jsonData))

	} else {
		ip, _ := utils.NodeTables.GetNodeIP(maxHost)
		job.WorkflowName = workflowName
		jsonData, err := utils.Jsoniter.Marshal(job)
		if err != nil {
			fmt.Println("Error:", err)
		}
		utils.SingleSend(ip, maxHost, "workflow job deployment", string(jsonData))

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
		//fmt.Println(result)
		//result := logging.Getmachineworkload()
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
