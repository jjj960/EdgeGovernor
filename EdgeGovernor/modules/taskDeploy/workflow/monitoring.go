package workflow

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func AddWorkflow(wf *models.PostWorkflow, isTime bool) error {
	wf.WorkflowName, _ = getUniqueTaskName(wf.WorkflowName)
	d, adag, err := ReadWorkflow(wf.DataType, wf.Data)
	if err != nil {
		return err
	}

	var DiskOcc int64
	// 打印解析结果
	for _, job := range adag.Jobs {
		for _, file := range job.Uses {
			data, _ := strconv.Atoi(file.Size)
			DiskOcc = DiskOcc + int64(data)
		}
	}

	newWf := models.Workflow{WorkflowName: wf.WorkflowName,
		DeployTime: wf.DeployTime,
		DAG:        *d,
		JobNum:     d.GetSize(),
		DiskOcc:    DiskOcc,
		Status:     false,
	}

	log.Printf("Add workflow: %s,execution time:%s\n", newWf.WorkflowName, newWf.DeployTime)

	id, _ := utils.GetID()
	logEntry := models.OperationLog{
		ID:            id,
		NodeName:      constants.Hostname,
		NodeIP:        constants.IP,
		OperationType: "add workflow",
		Description:   fmt.Sprintf("Workflow %s added", wf.WorkflowName),
		Result:        true,
		CreatedAt:     time.Now(),
	}
	duckdb.InsertOperationLog(logEntry)

	if isTime { //是定时发布
		err := PushWorkflow(newWf)
		if err != nil {
			return errors.New("Workflow added to queue failed")
		}
	} else { //即时发布
		err := workflowPublish(newWf)
		if err != nil {
			return errors.New("Workflow publishing failed")
		}
	}

	return nil
}

func getUniqueTaskName(workflowName string) (string, error) {
	isExists := CheckWorkflowName(workflowName)

	for isExists {
		randomStr, _ := utils.GenerateTaskNameRandomString(5)
		workflowName = workflowName + "-" + randomStr
		isExists = CheckWorkflowName(workflowName)
	}

	return workflowName, nil
}
