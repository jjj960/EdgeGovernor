package modulesControl

import (
	statusMonitor2 "EdgeGovernor/modules/statusMonitor"
	"EdgeGovernor/modules/taskDeploy/simple"
	"EdgeGovernor/modules/web/app/routers"

	//"EdgeGovernor/modules/taskDeploy/simple"
	//"EdgeGovernor/modules/web/app/routers"
	"EdgeGovernor/pkg/constants"
	"log"
	"time"
)

var MMC *MasterModulesController

type MasterModulesController struct {
	statusMonitor *statusMonitor2.HeartbeatChecker
	taskMonitor   *simple.TaskModule
	webServer     *routers.WebServer
}

func NewMMC() *MasterModulesController {
	heartBeatChecker := statusMonitor2.NewHeartbeatChecker(time.Duration(constants.HeartBeatTimes))
	taskMonitor := simple.NewTaskModule()
	webServer := routers.NewWebServer()

	return &MasterModulesController{
		taskMonitor:   taskMonitor,
		statusMonitor: heartBeatChecker,
		webServer:     webServer,
	}
}

func (mmc *MasterModulesController) Start() {
	log.Println("MMC enabled.")
	go mmc.taskMonitor.Start()
	go mmc.webServer.Start("5000")
	go mmc.statusMonitor.Start()
}

func (mmc *MasterModulesController) Close() {
	log.Println("MMC turned off.")
	mmc.taskMonitor.Close()
	mmc.webServer.Close()
	mmc.statusMonitor.Close()
}
