package workflow

import (
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"sort"
)

type workflowQueue []models.Workflow

var (
	Wq workflowQueue
)

func (Wq workflowQueue) Len() int { return len(Wq) }

func (Wq workflowQueue) Less(i, j int) bool {
	return Wq[i].DeployTime.Before(Wq[j].DeployTime)
}

func (Wq workflowQueue) Swap(i, j int) { Wq[i], Wq[j] = Wq[j], Wq[i] }

func (Wq *workflowQueue) Push(x interface{}) {
	item := x.(models.Workflow)
	*Wq = append(*Wq, item)
}

func (Wq *workflowQueue) Pop() interface{} {
	old := *Wq
	n := len(old)
	item := old[n-1]
	*Wq = old[0 : n-1]
	return item
}

func PeekFirstWorkflow() (models.Workflow, error) {
	if Wq.Len() == 0 {
		return models.Workflow{}, errors.New("The workflow queue is empty")
	}

	return (Wq)[0], nil
}

func GetWorkflowByName(workflowName string) (models.Workflow, error) {
	for _, workflow := range Wq {
		if workflow.WorkflowName == workflowName {
			return workflow, nil
		}
	}

	return models.Workflow{}, errors.New("Workflow not found in queue")
}

func PushWorkflow(workflow models.Workflow) error {
	_, err := GetWorkflowByName(workflow.WorkflowName)
	if err != nil {
		heap.Push(&Wq, workflow)
		sort.Stable(Wq) // 重新排序消息队列
	} else {
		return errors.New("Duplicate workflow name")
	}
	BackupworkflowQueue()

	return nil
}

func PopWorkflow() (models.Workflow, error) {
	if Wq.Len() == 0 {
		return models.Workflow{}, errors.New("The workflow queue is empty")
	}
	workflow := heap.Pop(&Wq).(models.Workflow)
	BackupworkflowQueue()
	return workflow, nil
}

func DeleteWorkflow(workflowName string) error {
	for i, workflow := range Wq {
		if workflow.WorkflowName == workflowName {
			Wq = append((Wq)[:i], (Wq)[i+1:]...) // 从队列中删除指定任务

			return nil
		}
	}
	BackupworkflowQueue()
	return errors.New("Workflow not found in queue")
}

func ChangeWorkflowStatus(workflowName string, status bool) error {
	for i, workflow := range Wq {
		if workflow.WorkflowName == workflowName {
			(Wq)[i].Status = status
			return nil
		}
	}
	BackupworkflowQueue()
	return errors.New("Workflow not found in queue")
}

func ClearQueue() { //清除队列
	Wq = make(workflowQueue, 0)
}

func GetUnexecutedWorkflows() []models.Workflow {
	var result []models.Workflow
	for _, workflow := range Wq {
		if workflow.Status == false {
			result = append(result, workflow)
		}
	}

	return result
}

func GetExecutedMsgs() []models.Workflow {
	var result []models.Workflow
	for _, workflow := range Wq {
		if workflow.Status == true {
			result = append(result, workflow)
		}
	}

	return result
}

func GetAllWorkflow() []models.Workflow {

	var result []models.Workflow
	for _, workflow := range Wq {
		result = append(result, workflow)
	}

	return result
}

func CheckWorkflowName(workflowName string) bool {
	for _, workflow := range Wq {
		if workflow.WorkflowName == workflowName {
			return true
		}
	}
	return false
}

func BackupworkflowQueue() error {
	jsonData, err := utils.Jsoniter.Marshal(Wq)
	if err != nil {
		return fmt.Errorf("Error marshalling tasks to JSON: %s", err)
	}

	encryptQueue := sec.Safer.Encrypt(jsonData) //加密map

	_, err = utils.ETCDCli.Put(context.Background(), "/menet/backup/workflowQueue", string(encryptQueue))
	if err != nil {
		return fmt.Errorf("failed to put queue in etcd: %w", err)
	}

	return nil
}

func GetBackupworkflowQueue() error {
	data, err := utils.ETCDCli.Get(context.Background(), "/menet/backup/workflowQueue")
	if err != nil {
		return fmt.Errorf("failed to get map in etcd: %s", err)
	}

	for _, kv := range data.Kvs {
		decryptMap := sec.Safer.Decrypt(kv.Value)

		var newWq workflowQueue
		err = utils.Jsoniter.Unmarshal([]byte(decryptMap), &newWq)
		if err != nil {
			return fmt.Errorf("Error unmarshalling workflowQueue data: %s", err)
		}
		Wq = newWq
	}

	return nil
}
