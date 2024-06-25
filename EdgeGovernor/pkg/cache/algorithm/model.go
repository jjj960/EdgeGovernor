package algorithm

import "sync"

type AlgorithmStatus struct {
	Name   string
	URL    string
	Type   string
	Status string
}

var AlgorithmStatusMap = make(map[string]AlgorithmStatus)
var mutex = &sync.RWMutex{} // 用于保护并发访问
