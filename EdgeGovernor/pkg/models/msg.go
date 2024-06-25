package models

type Msg struct {
	ID           int64  `json:"id"`
	GenerateTime int64  `json:"generate_time"`
	Tpye         string `json:"tpye"`
	Detail       []byte `json:"detail"`
	Status       bool   `json:"status"`
}

var (
	Msgs = []Msg{}
)

type ClusterMsg struct {
	Candidate     string `json:"candidate"`
	Leader        string `json:"leader"`
	NodeCount     int    `json:"nodeCount"`
	LiveNodeCount int    `json:"liveNodeCount"`
}

type SingleSendMsg struct {
	Addr       string `json:"addr"`
	TargetNode string `json:"target_node"`
	Types      string `json:"types"`
	Details    string `json:"details"`
}

type JobResultReport struct {
	Host         string `json:"host"`
	JobID        string `json:"job_id"`
	WorkflowName string `json:"workflow_name"`
	Status       string `json:"status"`
}
