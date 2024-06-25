package alogorithm

import "EdgeGovernor/pkg/utils"

func CheckNodeResource(node string, nodemsg [][]string, msg []string) bool {
	requestCPU, _ := utils.StringtoInt64(msg[2])
	requestMem, _ := utils.StringtoInt64(msg[3])
	requestDisk, _ := utils.StringtoInt64(msg[5])
	for _, item := range nodemsg {
		if node == item[0] {
			reCPU, _ := utils.StringtoInt64(item[4])
			reMem, _ := utils.StringtoInt64(item[7])
			reDisk, _ := utils.StringtoInt64(item[10])
			if requestMem > reMem || requestDisk > reDisk { //如果任务请求内存和磁盘大小大于节点剩余磁盘大小,直接拒绝
				return false
			} else {
				if requestCPU > reCPU && msg[6] == "0" {
					return false
				}
			}
		}
	}
	return true
}

func ResourceAllocation(host string, nodemsg [][]string, taskmsg []string, totalRequestCPU int64) (int64, int64) {
	requestCPU, _ := utils.StringtoInt64(taskmsg[2])    //任务请求的CPU数量
	var nodeReCPU int64                                 //节点剩余的CPU数量
	var totalReCPU int64                                //集群剩余的总CPU数量
	allocationMem, _ := utils.StringtoInt64(taskmsg[3]) //为该任务分配的内存
	var allocationCPU int64                             //为该任务分配的cpu
	for _, item := range nodemsg {
		reCPU, _ := utils.StringtoInt64(item[4])
		totalReCPU += reCPU
	}
	for _, item := range nodemsg {
		if host == item[0] {
			nodeReCPU, _ = utils.StringtoInt64(item[4])
		}
	}

	if totalRequestCPU != 0 {
		//(1)如果集群节点的剩余资源大于当前节点所需的资源任务，可以生成任务容器
		if totalRequestCPU <= totalReCPU {
			if requestCPU < nodeReCPU {
				//集群资源满足
				allocationCPU = requestCPU
			} else {
				if requestCPU >= nodeReCPU {
					allocationCPU = nodeReCPU * 8 / 10
				} else {
					if requestCPU < nodeReCPU {
						allocationCPU = requestCPU
					} else {
						allocationCPU = nodeReCPU * 8 / 10
					}
				}
			}
		} else {
			//(2) 要创建的任务的CPU资源不足,则按比例压缩任务的CPU资源。
			if totalRequestCPU > totalReCPU {
				//压缩CPU资源
				allocationCPUTemp := requestCPU * totalReCPU / totalRequestCPU
				if allocationCPUTemp < nodeReCPU {
					//集群资源满足
					allocationCPU = allocationCPUTemp
				} else {
					if allocationCPUTemp >= nodeReCPU {
						allocationCPU = nodeReCPU * 8 / 10
					} else {
						if allocationCPUTemp < nodeReCPU {
							allocationCPU = allocationCPUTemp
						} else {
							allocationCPU = nodeReCPU * 8 / 10
						}
					}
				}
			}
		}
	}
	return allocationCPU, allocationMem
}
