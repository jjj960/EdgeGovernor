package dao

import (
	"EdgeGovernor/modules/taskDeploy/workflow"
	"EdgeGovernor/pkg/models"
)

func DeleteWorkflow(workflowName string) {
	workflow.DeleteWorkflow(workflowName)
}

func GetWorkflowNum() int {
	wfs := workflow.GetAllWorkflow()

	return len(wfs)
}

func GetWorkflowMsg(workflowName string) (models.Workflow, error) {
	return workflow.GetWorkflowByName(workflowName)
}
