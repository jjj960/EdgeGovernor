package modulesControl

import (
	"EdgeGovernor/modules/comm/gRPC"
	"EdgeGovernor/modules/workload"
	"EdgeGovernor/pkg/constants"
	"log"
	"time"
)

var CMC *CommonModulesController

type CommonModulesController struct {
	commServer *gRPC.GRPCServer
	workload   *workload.WorkloadCollector
}

func NewCMC() *CommonModulesController {
	commServer := gRPC.NewGRPCServer(50051)
	workload1 := workload.NewWorkloadCollector(time.Duration(constants.CollectTime))
	return &CommonModulesController{
		workload:   workload1,
		commServer: commServer,
	}
}

func (cmc *CommonModulesController) Start() {
	log.Println("CMC enabled.")
	go cmc.workload.Start()
	go cmc.commServer.Start()
}

func (cmc *CommonModulesController) Close() {
	log.Println("CMC turned off.")
	cmc.workload.Close()
	cmc.commServer.Close()
}
