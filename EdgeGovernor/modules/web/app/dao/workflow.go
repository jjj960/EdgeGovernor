package dao

import (
	"EdgeGovernor/modules/taskDeploy/workflow"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"time"
)

func DeleteWorkflow(workflowName string) {
	workflow.DeleteWorkflow(workflowName)
	id, _ := utils.GetID()
	logEntry := models.OperationLog{
		ID:            id,
		NodeName:      constants.Hostname,
		NodeIP:        constants.IP,
		OperationType: "delete workflow",
		Description:   fmt.Sprintf("Workflow %s deleted", workflowName),
		Result:        true,
		CreatedAt:     time.Now(),
	}
	duckdb.InsertOperationLog(logEntry)
}

func GetWorkflowNum() int {
	wfs := workflow.GetAllWorkflow()

	return len(wfs)
}

func GetWorkflowMsg(workflowName string) (models.Workflow, error) {
	return workflow.GetWorkflowByName(workflowName)
}
