package modulesControl

import (
	"EdgeGovernor/pkg/utils"
	"fmt"
)

func ListenModulesChannel() {
	for {
		select {
		case val := <-utils.ModuleControlChannel:
			if val { //true为Leader到Follower
				fmt.Println("Module from Leader to Follower")
				MMC.Close()
				FMC = NewFMC()
				FMC.Start()
			} else { //false为Follower到Leader
				fmt.Println("Module from Follower to Leader")
				FMC.Close()
				MMC = NewMMC()
				MMC.Start()
			}
		}
	}
}
