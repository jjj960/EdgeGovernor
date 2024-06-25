package models

import (
	"github.com/heimdalr/dag"
	"time"
)

type PostWorkflow struct {
	WorkflowName string      `json:"workflow_name"`
	DeployTime   time.Time   `json:"deploy_time"`
	DataType     string      `json:"data_type"`
	Data         interface{} `json:"data"`
}

type Workflow struct {
	WorkflowName        string    `json:"workflow_name"`
	DeployTime          time.Time `json:"deploy_time"`
	DAG                 dag.DAG   `json:"dag"`
	JobNum              int       `json:"job_num"`
	DiskOcc             int64     `json:"disk_occ"`
	SchedulingAlgorithm string    `json:"scheduling_algorithm"`
	CompletedJob        []string  `json:"completed_job"`
	Status              bool      `json:"status"`
}

type JobOperation struct {
	Type string `json:"type"`
	Job  Job    `json:"job"`
}
