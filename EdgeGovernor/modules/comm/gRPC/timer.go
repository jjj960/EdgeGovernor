package gRPC

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"time"
)

type CustomTimer struct {
	ticker *time.Ticker
	second int
}

var HeartBeatTimer *CustomTimer

func NewCustomTimer(duration time.Duration) *CustomTimer {
	return &CustomTimer{
		ticker: time.NewTicker(duration),
		second: 0,
	}
}

func (ct *CustomTimer) Start() {
	for {
		select {
		case <-ct.ticker.C:
			ct.second++
			if ct.second > constants.HeartBeatTimeout {
				if constants.Candidate == constants.Hostname {
					fmt.Println("The current Leader node is malfunctioning, starting to execute Leader election")
					constants.ClusterStatus = "selfGovernment"
					candidateIP, _ := utils.NodeTables.GetNodeIP(constants.Candidate)
					utils.SingleSend(candidateIP, constants.Candidate, "leader election", "")
					ct.second = 0
				} else {
					fmt.Println("The current Leader node is malfunctioning, waiting for a new Leader to be generated")
				}
			}
		}
	}
}

func (ct *CustomTimer) Reset() {
	ct.second = 0
}

func (ct *CustomTimer) Stop() {
	ct.second = 0
	ct.ticker.Stop()
}
