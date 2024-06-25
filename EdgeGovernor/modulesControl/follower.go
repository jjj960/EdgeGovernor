package modulesControl

import (
	"EdgeGovernor/modules/comm/gRPC"
	"EdgeGovernor/pkg/constants"
	"log"
	"time"
)

var FMC *FollowerModulesController

type FollowerModulesController struct {
	heartBeatTimer *gRPC.CustomTimer
}

func NewFMC() *FollowerModulesController {

	gRPC.HeartBeatTimer = gRPC.NewCustomTimer(time.Duration(constants.HeartBeatTimeout) * time.Second)

	return &FollowerModulesController{
		heartBeatTimer: gRPC.HeartBeatTimer,
	}
}

func (fmc *FollowerModulesController) Start() {
	log.Println("FMC enabled.")
	go fmc.heartBeatTimer.Start()
}

func (fmc *FollowerModulesController) Close() {
	log.Println("FMC turned off.")
	fmc.heartBeatTimer.Stop()
}
