package dao

import (
	"EdgeGovernor/pkg/cache/task"
	"EdgeGovernor/pkg/logging"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"log"
	"strconv"
)

func SearchNodeMsg(Hostname string) map[string]string { //获取单个节点的详细信息
	nodeMsg, _ := utils.NodeTables.GetNodeEntry(Hostname)
	IP := nodeMsg.IP
	Role := nodeMsg.Role
	Status := nodeMsg.Status

	node1 := make(map[string]string)
	// 遍历结果集
	var CPU, Mem, Disk string
	if Status == "Offline" {
		CPU, Mem, Disk = "-", "-", "-"
	} else {
		if Role == "Leader" {
			result := logging.GetHostWorkload()
			var data models.Hostload
			err := utils.Jsoniter.Unmarshal(result, &data)
			if err != nil {
				log.Println("Parsing JSON failed:", err)
			}
			CPU = utils.Int64toString(data.CPUCapacity-data.CPUResidue) + " / " + utils.Int64toString(data.CPUCapacity)
			Mem = utils.Int64toString(data.MemoryCapacity-data.MemoryResidue) + " / " + utils.Int64toString(data.MemoryCapacity)
			Disk = utils.Int64toString(data.DiskCapacity-data.DiskResidue) + " / " + utils.Int64toString(data.DiskCapacity)
		} else {
			CPU, Mem, Disk = collectFollowerMachineStatus(Hostname, IP)
		}

		node1["nodeName"] = Hostname
		node1["cpuSize"] = CPU
		node1["memorySize"] = Mem
		node1["diskSize"] = Disk
		node1["equipment"] = Role
		node1["status"] = Status
		node1["taskNum"] = strconv.Itoa(task.GetNodeTaskNum(Hostname))
	}

	return node1
}

func GetNodesMsg(start int, end int) []map[string]string { //获取一些 节点的详细信息
	var result []map[string]string
	for i := start; i <= end-1; i++ {
		var Hostname, IP, Role, Status string
		nodeMsg, _ := utils.NodeTables.GetNodeEntry(utils.NodeList[i])
		fmt.Println(nodeMsg)
		Hostname, IP, Role, Status = nodeMsg.Hostname, nodeMsg.IP, nodeMsg.Role, nodeMsg.Status
		var CPU, Mem, Disk string
		if Status == "Offline" {
			CPU, Mem, Disk = "-", "-", "-"
		} else {
			if Role == "Leader" {
				result := logging.GetHostWorkload()
				var data models.Hostload
				err := utils.Jsoniter.Unmarshal(result, &data)
				if err != nil {
					log.Println("Parsing JSON failed:", err)
				}
				CPU = utils.Int64toString(data.CPUCapacity-data.CPUResidue) + " / " + utils.Int64toString(data.CPUCapacity)
				Mem = utils.Int64toString(data.MemoryCapacity-data.MemoryResidue) + " / " + utils.Int64toString(data.MemoryCapacity)
				Disk = utils.Int64toString(data.DiskCapacity-data.DiskResidue) + " / " + utils.Int64toString(data.DiskCapacity)
			} else {
				CPU, Mem, Disk = collectFollowerMachineStatus(Hostname, IP)
			}
		}
		node1 := make(map[string]string)
		node1["nodeName"] = Hostname
		node1["cpuSize"] = CPU
		node1["memorySize"] = Mem
		node1["diskSize"] = Disk
		node1["equipment"] = Role
		node1["status"] = Status
		node1["taskNum"] = strconv.Itoa(task.GetNodeTaskNum(Hostname))
		result = append(result, node1)
	}

	return result
}

func collectFollowerMachineStatus(Hostname string, IP string) (string, string, string) {
	var data models.Hostload
	result, _ := utils.SingleSend(IP, Hostname, "workload report", "")
	//fmt.Println(result)
	//result := logging.Getmachineworkload()
	err := utils.Jsoniter.Unmarshal([]byte(result), &data)
	if err != nil {
		log.Println("解析JSON失败：", err)
	}
	CPU := utils.Int64toString(data.CPUCapacity-data.CPUResidue) + " / " + utils.Int64toString(data.CPUCapacity)
	Mem := utils.Int64toString(data.MemoryCapacity-data.MemoryResidue) + " / " + utils.Int64toString(data.MemoryCapacity)
	Disk := utils.Int64toString(data.DiskCapacity-data.DiskResidue) + " / " + utils.Int64toString(data.DiskCapacity)

	//fmt.Println(nodeMsg)
	return CPU, Mem, Disk
}
