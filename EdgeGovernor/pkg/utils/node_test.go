package utils

import (
	"EdgeGovernor/modules/comm/gRPC"
	"EdgeGovernor/modules/web/app/dao"
	"fmt"
	"testing"
)

func TestNodesTable_Iter(t *testing.T) {
	StorageNodeMsg()
	fmt.Println(NodeTables.GetAllNodeName())
	node, _ := dao.GetTaskNum()
	server := gRPC.NewGRPCServer(50051)
	go server.Start()
	cpu, mem, disk := dao.GetBarData(node)
	fmt.Println(cpu, mem, disk)
}
