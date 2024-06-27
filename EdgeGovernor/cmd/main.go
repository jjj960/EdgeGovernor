package main

import (
	"EdgeGovernor/modules/comm/gRPC"
	"EdgeGovernor/modules/resourceCheck"
	"EdgeGovernor/modulesControl"
	"EdgeGovernor/pkg/cache/algorithm"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"github.com/bwmarrin/snowflake"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	modulesControl.CMC = modulesControl.NewCMC()
	modulesControl.FMC = modulesControl.NewFMC()
	modulesControl.MMC = modulesControl.NewMMC()

	modulesControl.CMC.Start()

	if !utils.IsFirstRun() && constants.Hostname == "cloud" {
		modulesControl.FMC.Start()
	} else if constants.Hostname == constants.Leader {
		modulesControl.MMC.Start()
	} else {
		modulesControl.FMC.Start()
	}

	go resourceCheck.JobCheckMonitor()
	go modulesControl.ListenModulesChannel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}

func init() {

	utils.StorageNodeMsg()
	constants.Hostname, _ = os.Hostname()
	constants.IP, _ = utils.GetIP()
	constants.BandWidth, _ = utils.GetNetbandwidth("ens33")
	constants.LiveNodeCount = utils.NodeTables.Len()
	constants.NodeCount = utils.NodeTables.Len()

	id, _ := utils.NodeTables.GetNodeID(constants.Hostname)
	utils.SnowFlake, _ = snowflake.NewNode(id)

	constants.Leader = "cloud"
	constants.Candidate = "cloud"
	constants.ClusterStatus = "coordination"
	constants.CollectTime = 5

	utils.GetDuckDBCli() //Initialize Database
	utils.GetETCDCli()
	utils.GetDockerCli()
	if constants.Hostname == "cloud" {
		utils.GetK8sCli()
	}

	duckdb.CreateInitTable()

	//k8s.CreateNamespace()
	//
	sec.Safer = sec.NewMSecLayer("2TXQfxBE#TrULn6FZx") //Initialize encryptor
	utils.ModuleControlChannel = make(chan bool)

	gRPC.HeartBeatTimer = gRPC.NewCustomTimer(1 * time.Second)
	algorithm.AddAlgorithmStatus("TOPSIS", "http://192.168.47.128:50052/scheduler", "Schedule", "Enable")
	algorithm.AddAlgorithmStatus("NodeLoadAssessmentAlgorithm", "http://192.168.47.128:50052/scheduler", "NodeLoadAssessment", "Enable")
	algorithm.AddAlgorithmStatus("BRA", "http://192.168.47.128:50052/scheduler", "Schedule", "Enable")
	algorithm.AddAlgorithmStatus("LRP", "http://192.168.47.128:50052/scheduler", "Schedule", "Enable")
}
