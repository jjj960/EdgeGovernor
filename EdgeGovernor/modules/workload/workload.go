package workload

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/logging"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"log"
	"sync"
	"time"
)

type WorkloadCollector struct {
	stopChan chan struct{}
	mu       sync.Mutex
	interval time.Duration
	started  bool
}

func NewWorkloadCollector(interval time.Duration) *WorkloadCollector {
	return &WorkloadCollector{
		stopChan: make(chan struct{}),
		interval: interval * time.Second,
		started:  false,
	}
}

func (wc *WorkloadCollector) Start() {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	if wc.started {
		return
	}

	wc.started = true
	go func() {
		ticker := time.NewTicker(wc.interval)
		for {
			select {
			case <-ticker.C:
				result := logging.GetHostWorkload()
				logging.GetContainerWorkload()
				var data models.Hostload
				err := utils.Jsoniter.Unmarshal(result, &data)
				if err != nil {
					log.Println("Parsing JSON failed:", err)
				}
				selfRole, _ := utils.NodeTables.GetNodeRole(constants.Hostname)
				if selfRole == "Follower" {
					duckdb.InsertHostload(data)
					ip, _ := utils.NodeTables.GetNodeIP(constants.Leader)
					utils.SingleSend(ip, constants.Leader, "machine workload", string(result))
				} else {
					duckdb.InsertHostload(data)

				}

			case <-wc.stopChan:
				return
			}
		}
	}()
}

func (wc *WorkloadCollector) Close() {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	if !wc.started {
		return
	}

	close(wc.stopChan)
	wc.started = false
}
