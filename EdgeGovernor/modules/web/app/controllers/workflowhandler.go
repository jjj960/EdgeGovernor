package controllers

import (
	"EdgeGovernor/modules/taskDeploy/workflow"
	"EdgeGovernor/modules/web/app/dao"
	"EdgeGovernor/pkg/models"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

func AddWorkflow(c echo.Context) error {
	taskName := c.FormValue("taskName")
	taskpostTime1 := c.FormValue("taskpostTime1")
	taskpostTime2 := c.FormValue("taskpostTime2")
	pmValue := c.FormValue("pmValue")
	textArea := c.FormValue("textarea")
	fileList := c.FormValue("fileList")
	wf := new(models.PostWorkflow)

	if pmValue == "text" {
		wf.WorkflowName = taskName
		wf.DataType = "text"
		wf.DeployTime, _ = time.Parse("2006-01-02 15:04:05", (taskpostTime1 + " " + taskpostTime2))
		wf.Data = textArea
	} else {
		wf.WorkflowName = taskName
		wf.DataType = "file"
		wf.DeployTime, _ = time.Parse("2006-01-02 15:04:05", (taskpostTime1 + " " + taskpostTime2))
		wf.Data = fileList
	}

	err := workflowProcess(wf)
	if err != nil {
		return err
	}
	// 返回成功响应
	return c.JSON(http.StatusOK, map[string]string{"message": "Workflow added successfully"})
}

func workflowProcess(wf *models.PostWorkflow) error {
	if wf.DeployTime.IsZero() {
		wf.DeployTime = time.Now()
		err := workflow.AddWorkflow(wf, false)
		if err != nil {
			return err
		}
	} else {
		isPassed, err := isExecutionTimePassed(wf.DeployTime)
		if err != nil {
			log.Printf("Error parsing execution time: %v", err)
			return errors.New("Error parsing execution time")
		}

		if isPassed {
			err = workflow.AddWorkflow(wf, true)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Refuse adding task: %s, execution time: %s has passed", wf.WorkflowName, wf.DeployTime)
			return errors.New("Refuse adding task, execution time has passed")
		}
	}

	return nil
}

func SearchWorkflow(c echo.Context) error {
	workflowName := c.FormValue("workflowName")
	wf, _ := dao.GetWorkflowMsg(workflowName)

	return c.JSON(http.StatusOK, wf)
}

func DeleteWorkflow(c echo.Context) error {
	workflowName := c.FormValue("workflowName")
	dao.DeleteWorkflow(workflowName)
	return c.String(http.StatusOK, "success")
}

func GetWorkflowNum(c echo.Context) error {
	nodeNum := dao.GetWorkflowNum()

	response := map[string]interface{}{
		"nodeNum": nodeNum,
	}
	fmt.Println(response)
	return c.JSON(http.StatusOK, response)
}
