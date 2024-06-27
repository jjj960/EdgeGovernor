package dao

import (
	"EdgeGovernor/pkg/cache/task"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/database/etcd/resource"
	docker "EdgeGovernor/pkg/docker/resource"
	"EdgeGovernor/pkg/k8s"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"strconv"
	"time"
)

func CheckTaskNode(taskName string, taskNode string) bool {
	task, err := resource.GetTaskMsg(taskName)
	if err != nil {
		return false
	}
	if task.Migrate == 1 && task.DeployNode == taskNode {
		return true
	}
	return false
}

func GetNodeTasks(Hostname string) []map[string]string {
	tasks := task.GetNodeAllTaskNum(Hostname)
	var result []map[string]string
	for i := 0; i < len(tasks); i++ {
		task, _ := resource.GetTaskMsg(tasks[i])

		if Hostname == "" {
			Hostname = "-"
		}
		node1 := make(map[string]string)
		node1["nodeName"] = Hostname
		node1["taskName"] = task.Name
		node1["taskMirror"] = task.Image
		node1["cpuSize"] = strconv.FormatInt(task.RequestCPU, 10)
		node1["memorySize"] = strconv.FormatInt(task.RequestMem, 10)
		node1["diskSize"] = strconv.FormatInt(task.RequestDisk, 10)
		node1["status"] = task.Status
		node1["publishTime"] = task.PublishTime.Format("2006-01-02 15:04:05")
		result = append(result, node1)
	}

	return result
}

func DeleteTask(taskName string) {

	tasks, _ := resource.GetTaskMsg(taskName)
	if tasks.Status == "Deployed" {
		if tasks.Type == "Container" {
			if tasks.PersistData == 1 {
				docker.DeleteContainer(taskName, true)
			} else {
				docker.DeleteContainer(taskName, false)
			}
		} else {
			k8s.DeletePod(taskName, tasks.PersistData)
		}

	}
	resource.DeleteTask(taskName)
	task.DeleteNodeTask(tasks.DeployNode, tasks.Name)
	id, _ := utils.GetID()
	logEntry := models.OperationLog{
		ID:            id,
		NodeName:      constants.Hostname,
		NodeIP:        constants.IP,
		OperationType: "delete task",
		Description:   fmt.Sprintf("Task %s deleted", taskName),
		Result:        true,
		CreatedAt:     time.Now(),
	}
	duckdb.InsertOperationLog(logEntry)
}

func UpdateTaskMsg(taskName string, targetNode string) {
	task, _ := resource.GetTaskMsg(taskName)
	resource.UpdateTaskMsg(taskName, targetNode, task.RequestCPU, task.Type, task.Status)
}

func GetTaskNum() ([]string, []string) {
	nodeName := utils.NodeTables.GetAllNodeName()
	var nodeValue []string
	for _, node := range nodeName {
		count := task.GetNodeTaskNum(node)
		nodeValue = append(nodeValue, strconv.Itoa(count))
	}
	return nodeName, nodeValue
}
