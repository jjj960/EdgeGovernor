package resourceCheck

import (
	"EdgeGovernor/modules/taskDeploy/workflow"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/docker/resource"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"time"
)

var CheckJobMap map[string]models.Job

func JobCheckMonitor() {
	go func() {
		for {
			select {
			case operation := <-utils.JobOperationChannel:
				switch operation.Type {
				case "check":
					jobCheck(operation.Job)
				case "publish":
					jobPublish(operation.Job)
				}

			case <-time.After(time.Second):

			}
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if len(CheckJobMap) == 0 {
				continue
			} else {
				for kv, value := range CheckJobMap {
					result := resource.CheckContainerStatus(kv)
					switch result {
					case "exited":
						ip, _ := utils.NodeTables.GetNodeIP(constants.Leader)
						msg := models.JobResultReport{
							JobID:        value.ID,
							WorkflowName: value.WorkflowName,
							Host:         constants.Hostname,
							Status:       "exited",
						}

						jsonData, err := utils.Jsoniter.Marshal(msg)
						if err != nil {
							fmt.Println("Error:", err)
						}
						utils.SingleSend(ip, constants.Leader, "job completion report", string(jsonData))

						delete(CheckJobMap, kv)
					}
				}
			}
		}
	}()
}

func jobCheck(job models.Job) {
	result := resource.CheckContainerStatus(job.Name)
	switch result {
	case "running":
		CheckJobMap[job.Name] = job
	case "exited":
		ip, _ := utils.NodeTables.GetNodeIP(constants.Leader)
		msg := models.JobResultReport{
			JobID:        job.ID,
			WorkflowName: job.WorkflowName,
			Host:         constants.Hostname,
			Status:       "exited",
		}

		jsonData, err := utils.Jsoniter.Marshal(msg)
		if err != nil {
			fmt.Println("Error:", err)
		}
		utils.SingleSend(ip, constants.Leader, "job completion report", string(jsonData))

	case "error":
		id := utils.SnowFlake.Generate()
		utils.AlarmMsgChannel <- models.Msg{
			ID:           id.Int64(),
			GenerateTime: id.Time(),
			Tpye:         "Job error",
			Detail:       []byte(fmt.Sprintf("Task %s in workflow %s runs incorrectly", job.Name, job.WorkflowName)),
			Status:       false,
		}
	}
}

func jobPublish(job models.Job) {
	workflow.UpdateJobStatus(job.WorkflowName, job.ID, true)
	parentMap, _ := workflow.GetJobParent(job.WorkflowName, job.ID)
	childrenMap, _ := workflow.GetJobChildren(job.WorkflowName, job.ID)

	if len(parentMap) > 0 {
		for _, value := range parentMap {
			msg, _ := workflow.GetJobMsg(value.WorkflowName, value.ID)
			if !msg.Status {
				return
			}
		}
	}

	wf, _ := workflow.GetWorkflowByName(job.WorkflowName)
	for _, value := range childrenMap {
		err := workflow.JobPublish(value, value.WorkflowName, wf.SchedulingAlgorithm)
		if err != nil {
			id := utils.SnowFlake.Generate()
			utils.AlarmMsgChannel <- models.Msg{
				ID:           id.Int64(),
				GenerateTime: id.Time(),
				Tpye:         "Job Publish fail",
				Detail:       []byte(fmt.Sprintf("Task %s in workflow %s failed to publish: %s", value.Name, value.WorkflowName, err.Error())),
				Status:       false,
			}
		}
	}
}
