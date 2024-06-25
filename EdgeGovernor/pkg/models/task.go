package models

import "time"

type Task struct {
	Name                string    `json:"name"`
	PublishTime         time.Time `json:"publish_time"`
	DeployNode          string    `json:"deploy_node"`
	SchedulingAlgorithm string    `json:"scheduling_algorithm"`
	PreDeploy           bool      `json:"pre_deploy"`
	Image               string    `json:"image"`
	OriRequestCPU       int64     `json:"ori_request_cpu"`
	RequestCPU          int64     `json:"request_cpu"`
	RequestMem          int64     `json:"request_mem"`
	OriRequestNet       int64     `json:"ori_request_net"`
	RequestNet          int64     `json:"request_net"`
	RequestDisk         int64     `json:"request_disk"`
	Priority            int       `json:"priority"`
	Migrate             int       `json:"migrate"`
	PersistData         int       `json:"persist_data"`
	DataDir             []string  `json:"data_dir"`
	Scalable            int       `json:"scalable"`
	ServiceDiscovery    string    `json:"service_discovery"`
	Type                string    `json:"type"`
	Status              string    `json:"status"`
}

type MigrationTask struct {
	Name       string `json:"name"`
	DeployNode string `json:"deploy_node"`
	TargetNode string `json:"target_node"`
}
