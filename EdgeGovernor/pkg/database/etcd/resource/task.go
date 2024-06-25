package resource

import (
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InsertTask(task models.Task, taskStatus string) error {

	key := "/menet/task/" + task.Name
	task.Status = taskStatus
	taskData, err := utils.Jsoniter.Marshal(task)
	if err != nil {
		return fmt.Errorf("JSON encoding failed: %s", err)
	}

	value := string(taskData)

	_, err = utils.ETCDCli.Put(context.Background(), key, value)
	if err != nil {
		return fmt.Errorf("Node data insertion failed: %s", err)
	}

	return nil
}

func DeleteTask(taskName string) {
	key := "/menet/task/" + taskName

	utils.ETCDCli.Delete(context.Background(), key)
}

func GetTaskMsg(taskName string) (models.Task, error) {
	key := "/menet/task/" + taskName

	resp, err := utils.ETCDCli.Get(context.Background(), key)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to retrieve data from etcd: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return models.Task{}, fmt.Errorf("failed to find node from etcd: %v", err)
	}

	value := string(resp.Kvs[0].Value)

	var taskMsg models.Task
	err = utils.Jsoniter.Unmarshal([]byte(value), &taskMsg)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return taskMsg, nil
}

func UpdateTaskMsg(taskName string, deployNode string, cpu int64, types, status string) error {
	taskMsg, err := GetTaskMsg(taskName)
	if err != nil {
		return err
	}

	taskMsg.DeployNode = deployNode
	taskMsg.RequestCPU = cpu
	taskMsg.Type = types

	err = InsertTask(taskMsg, status)
	if err != nil {
		return err
	}

	return nil
}

func CheckTaskName(taskName string) (bool, error) {
	prefix := "/menet/task"

	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return false, fmt.Errorf("Failed to obtain value: %v", err)
	}

	check := false

	for _, kv := range resp.Kvs {
		if string(kv.Key) == taskName {
			check = true
		}
	}

	return check, nil
}

func CheckTaskNode(taskName string, taskNode string) bool {
	taskMsg, err := GetTaskMsg(taskName)

	return (taskMsg.Migrate == 1 && err == nil && taskMsg.DeployNode == taskNode)
}
