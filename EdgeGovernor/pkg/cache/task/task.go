package task

import (
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
)

func GetNodeTaskNum(nodeName string) int {
	mutex.Lock()
	defer mutex.Unlock()

	return nodeSimpleTaskData[nodeName].TaskCount
}

func GetNodeAllTaskNum(nodeName string) []string {
	nodeInfo, _ := nodeSimpleTaskData[nodeName]

	// 收集所有任务名称
	taskNames := make([]string, 0, len(nodeInfo.Tasks))
	for taskName := range nodeInfo.Tasks {
		taskNames = append(taskNames, taskName)
	}

	return taskNames
}

func AddNodeTask(nodeName, taskName string, deployed bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	node, ok := nodeSimpleTaskData[nodeName]
	if !ok {
		node = NodeInfo{
			Tasks: make(map[string]TaskDeployment),
		}
	}

	node.Tasks[taskName] = TaskDeployment{TaskName: taskName, Deployed: deployed}
	node.TaskCount++

	nodeSimpleTaskData[nodeName] = node // 更新节点信息

	err := BackUpMap()
	if err != nil {
		return fmt.Errorf("NodeMap backup fail: %s\n", err)
	}

	return nil
}

// UpdateTaskDeployment 更新指定节点的任务部署状态
func UpdateTaskDeployment(nodeName, taskName string, deployed bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	if node, ok := nodeSimpleTaskData[nodeName]; ok {
		if task, ok := node.Tasks[taskName]; ok {
			task.Deployed = deployed
			node.Tasks[taskName] = task
			nodeSimpleTaskData[nodeName] = node // 更新节点信息
		} else {
			return fmt.Errorf("Task %s not found in Node %s\n", taskName, nodeName)
		}

	} else {
		return fmt.Errorf("Node %s not found\n", nodeName)
	}

	err := BackUpMap()
	if err != nil {
		return fmt.Errorf("NodeMap backup fail: %s\n", err)
	}

	return nil
}

// DeleteNodeTask 删除指定节点的任务信息
func DeleteNodeTask(nodeName, taskName string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if node, ok := nodeSimpleTaskData[nodeName]; ok {
		if _, ok := node.Tasks[taskName]; ok {
			delete(node.Tasks, taskName)        // 从节点任务列表中删除指定任务
			node.TaskCount--                    // 更新任务数量
			nodeSimpleTaskData[nodeName] = node // 更新节点信息
		} else {
			return fmt.Errorf("Task %s not found in Node %s\n", taskName, nodeName)
		}
	} else {
		return fmt.Errorf("Node %s not found\n", nodeName)
	}
	err := BackUpMap()
	if err != nil {
		return fmt.Errorf("NodeMap backup fail: %s\n", err)
	}

	return nil
}

func CheckTaskName(taskName string) bool {
	for key := range nodeSimpleTaskData {
		if taskName == key {
			return true
		}
	}
	return false
}

func BackUpMap() error { //备份关键map
	data, err := utils.Jsoniter.Marshal(nodeSimpleTaskData)

	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %s", err)
	}

	encryptMap := sec.Safer.Encrypt(data) //加密map

	_, err = utils.ETCDCli.Put(context.Background(), "/menet/backup/nodeSimpleTaskMap", string(encryptMap))
	if err != nil {
		return fmt.Errorf("failed to put map in etcd: %w", err)
	}

	return nil
}

func GetNodeSimpleTaskMap() error {
	data, err := utils.ETCDCli.Get(context.Background(), "/menet/backup/nodeSimpleTaskMap")
	if err != nil {
		return fmt.Errorf("failed to get map in etcd: %s", err)
	}

	for _, kv := range data.Kvs {
		decryptMap := sec.Safer.Decrypt(kv.Value)

		var tempMap map[string]NodeInfo
		err := utils.Jsoniter.Unmarshal([]byte(decryptMap), &tempMap)
		if err != nil {
			return fmt.Errorf("Error unmarshalling Map JSON: %s", err)
		}

		for key, value := range tempMap {
			nodeSimpleTaskData[key] = value
		}

	}

	return nil
}
