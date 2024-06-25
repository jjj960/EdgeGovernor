package statusMonitor

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"time"
)

var term = 0

type HeartbeatChecker struct {
	interval time.Duration
	stop     chan struct{}
}

// NewHeartbeatChecker 创建一个新的 HeartbeatChecker 实例
func NewHeartbeatChecker(interval time.Duration) *HeartbeatChecker {
	return &HeartbeatChecker{
		interval: interval * time.Second,
		stop:     make(chan struct{}),
	}
}

// Start 方法用于开始心跳检测
func (hb *HeartbeatChecker) Start() {
	ticker := time.NewTicker(hb.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if term >= constants.Term && constants.ClusterStatus == "selfGovernment" {
				fmt.Println("The term of the Leader has expired, and the process of changing the Leader has begun")
				if constants.Candidate == constants.Hostname { //当前候选人仍然为本机,则本机继续当选Leader,并重置任期
					fmt.Println("Reset term")
					term = 0
				} else {
					fmt.Println("Leader has been transferred to:", constants.Candidate)
					IP, _ := utils.NodeTables.GetNodeIP(constants.Candidate)
					utils.SingleSend(IP, constants.Candidate, "leader elected", "")

					utils.ModuleControlChannel <- true
				}
			}
			fmt.Println("Send heartbeat packets")
			utils.HeartBeat()
			term++
		case <-hb.stop:
			fmt.Println("Stop sending heartbeat packets")
			return
		}
	}
}

// Close 方法用于停止心跳检测
func (hb *HeartbeatChecker) Close() {
	close(hb.stop)
}
