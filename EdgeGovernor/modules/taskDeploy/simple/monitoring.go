package simple

import (
	task2 "EdgeGovernor/pkg/cache/task"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/database/etcd/resource"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"errors"
	"fmt"
	"log"
	"time"
)

func AddTask(task *models.Task, isTime bool) error {
	task.Name, _ = getUniqueTaskName(task.Name)
	resource.InsertTask(*task, "Undeployed")
	task2.AddNodeTask(task.DeployNode, task.Name, false)
	id, _ := utils.GetID()
	logEntry := models.OperationLog{
		ID:            id,
		NodeName:      constants.Hostname,
		NodeIP:        constants.IP,
		OperationType: "add task",
		Description:   fmt.Sprintf("Task %s added", task.Name),
		Result:        true,
		CreatedAt:     time.Now(),
	}
	duckdb.InsertOperationLog(logEntry)
	log.Printf("Add task: %s,execution time:%s\n", task.Name, task.PublishTime)

	if isTime { //是定时发布
		err := PushTask(task)
		if err != nil {
			return errors.New("Task added to queue failed")
		}
	} else { //即时发布
		err := taskPublish(task)
		if err != nil {
			return errors.New("Task publishing failed")
		}
	}
	return nil
}

func getUniqueTaskName(taskName string) (string, error) {
	newTaskName := taskName
	isExists := task2.CheckTaskName(newTaskName)

	for isExists {
		randomStr, _ := utils.GenerateTaskNameRandomString(5)
		newTaskName = taskName + "-" + randomStr
		isExists = task2.CheckTaskName(newTaskName)
	}

	return newTaskName, nil
}
