package task

import "sync"

// TaskDeployment 表示任务的部署情况
type TaskDeployment struct {
	TaskName string `json:"taskName"`
	Deployed bool   `json:"deployed"`
	// 其他任务相关的信息
}

// NodeInfo 表示节点的任务信息
type NodeInfo struct {
	TaskCount int                       `json:"taskCount"`
	Tasks     map[string]TaskDeployment `json:"tasks"`
	// 其他节点相关的信息
}

var nodeSimpleTaskData = make(map[string]NodeInfo)
var mutex = &sync.RWMutex{} // 用于保护并发访问
