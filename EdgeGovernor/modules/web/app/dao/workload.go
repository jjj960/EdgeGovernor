package dao

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/logging"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"log"
)

func GetBarData(nodes []string) ([]string, []string, []string) {
	var cpuSize []string
	var memSize []string
	var diskSize []string
	var data models.Hostload

	for i := 0; i < len(nodes); i++ {
		if nodes[i] == constants.Hostname {
			local := logging.GetHostWorkload()
			err := utils.Jsoniter.Unmarshal(local, &data)
			if err != nil {
				log.Println("Parsing JSON failed:", err)
			}
			cpuSize = append(cpuSize, utils.Float64toString(data.CPUUsagePercent))
			memSize = append(memSize, utils.Float64toString(data.MemoryUsedPercent))
			diskSize = append(diskSize, utils.Float64toString(data.DiskUsedPercent))
		} else {
			ip, _ := utils.NodeTables.GetNodeIP(nodes[i])
			result, err := utils.SingleSend(ip, nodes[i], "workload report", "")
			if err != nil {
				cpuSize = append(cpuSize, "0")
				memSize = append(memSize, "0")
				diskSize = append(diskSize, "0")
				continue
			}
			err = utils.Jsoniter.Unmarshal([]byte(result), &data)
			if err != nil {
				log.Println("Parsing JSON failed:", err)
			}
			cpuSize = append(cpuSize, utils.Float64toString(data.CPUUsagePercent))
			memSize = append(memSize, utils.Float64toString(data.MemoryUsedPercent))
			diskSize = append(diskSize, utils.Float64toString(data.DiskUsedPercent))
		}
	}

	return cpuSize, memSize, diskSize

}
