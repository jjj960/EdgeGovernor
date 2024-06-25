package workflow

import (
	"EdgeGovernor/pkg/models"
	"errors"
	"github.com/heimdalr/dag"
)

func GenerateDAG(adag models.Adag) (*dag.DAG, error) { // 将解析的作业信息存入DAG图中
	d := dag.NewDAG()
	for _, job := range adag.Jobs {
		err := d.AddVertexByID(job.ID, job)
		if err != nil {
			return nil, errors.New("DAG failed to add vertex")
		}
	}

	for _, child := range adag.Children {
		for _, parent := range child.Parents {
			err := d.AddEdge(parent.Ref, child.Ref)
			if err != nil {
				return nil, errors.New("DAG failed to add edge")
			}
		}
	}

	return d, nil
}

func GetDAGRootJob(workflowName string) (map[string]models.Job, error) { //获取工作流根任务
	workflow, err := GetWorkflowByName(workflowName)
	if err != nil {
		return nil, err
	}

	data := workflow.DAG.GetRoots()
	if len(data) == 0 {
		return nil, errors.New("There are no jobs in this workflow")
	}

	var rootJobMap map[string]models.Job

	for key, value := range workflow.DAG.GetRoots() {
		vertex, ok := value.(models.Job)
		if !ok {
			return nil, errors.New("Error reading Job data")
		}
		rootJobMap[key] = vertex
	}

	return rootJobMap, nil
}

func UpdateJobStatus(workflowName string, jobID string, status bool) error {
	workflow, err := GetWorkflowByName(workflowName)
	if err != nil {
		return err
	}

	job, err := workflow.DAG.GetVertex(jobID)
	if err != nil {
		return err
	}

	vertex, ok := job.(models.Job)
	if !ok {
		return errors.New("Error reading Job data")
	}

	vertex.Status = status
	return nil
}

func GetJobChildren(workflowName string, jobID string) (map[string]models.Job, error) {
	workflow, err := GetWorkflowByName(workflowName)
	if err != nil {
		return nil, err
	}

	childrens, err := workflow.DAG.GetChildren(jobID)
	if err != nil {
		return nil, err
	}

	var childrenMap map[string]models.Job

	for key, value := range childrens {
		vertex, ok := value.(models.Job)
		if !ok {
			return nil, errors.New("Error reading Job data")
		}
		childrenMap[key] = vertex
	}

	return childrenMap, nil
}

func GetJobParent(workflowName string, jobID string) (map[string]models.Job, error) {
	workflow, err := GetWorkflowByName(workflowName)
	if err != nil {
		return nil, err
	}

	childrens, err := workflow.DAG.GetParents(jobID)
	if err != nil {
		return nil, err
	}

	var parentMap map[string]models.Job

	for key, value := range childrens {
		vertex, ok := value.(models.Job)
		if !ok {
			return nil, errors.New("Error reading Job data")
		}
		parentMap[key] = vertex
	}

	return parentMap, nil
}

func CheckJobStatus(workflowName string, jobID string) (bool, error) {
	workflow, err := GetWorkflowByName(workflowName)
	if err != nil {
		return false, err
	}

	job, err := workflow.DAG.GetVertex(jobID)
	if err != nil {
		return false, err
	}

	vertex, ok := job.(models.Job)
	if !ok {
		return false, errors.New("Error reading Job data")
	}

	return vertex.Status, nil
}

func GetJobMsg(workflowName string, jobID string) (models.Job, error) {
	workflow, err := GetWorkflowByName(workflowName)
	if err != nil {
		return models.Job{}, err
	}

	job, err := workflow.DAG.GetVertex(jobID)
	if err != nil {
		return models.Job{}, err
	}

	vertex, ok := job.(models.Job)
	if !ok {
		return models.Job{}, errors.New("Error reading Job data")
	}

	return vertex, nil
}
