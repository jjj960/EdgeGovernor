package dao

import (
	"EdgeGovernor/modules/comm/gRPC"
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"testing"
)

func TestAddWorkflow(t *testing.T) {
	utils.StorageNodeMsg()
	fmt.Println(utils.NodeTables.GetAllNodeName())
	node, _ := GetTaskNum()
	server := gRPC.NewGRPCServer(50051)
	go server.Start()
	sec.Safer = sec.NewMSecLayer("2TXQfxBE#TrULn6FZx") //Initialize encryptor

	cpu, mem, disk := GetBarData(node)
	fmt.Println(cpu, mem, disk)
}
