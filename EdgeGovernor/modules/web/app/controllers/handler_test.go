package controllers

import (
	"EdgeGovernor/modules/web/app/dao"
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"testing"
)

func TestGetNodeName(t *testing.T) {
	utils.StorageNodeMsg()
	nodeName := utils.NodeTables.GetAllNodeName()

	response := map[string]interface{}{
		"node": nodeName,
	}

	fmt.Println(response)
}

func TestGetNodeMsg(t *testing.T) {
	utils.StorageNodeMsg()
	sec.Safer = sec.NewMSecLayer("2TXQfxBE#TrULn6FZx") //Initialize encryptor
	start := (1 - 1) * 4
	end := 1 * 4
	//if end > constants.NodeCount {
	//	end = constants.NodeCount
	//}

	nodeMessages := dao.GetNodesMsg(start, end)

	response := map[string]interface{}{
		"nodeMessage": nodeMessages,
		"total":       constants.NodeCount,
	}
	fmt.Println(response)
}

func TestSearchNode(t *testing.T) {
	utils.StorageNodeMsg()
	sec.Safer = sec.NewMSecLayer("2TXQfxBE#TrULn6FZx") //Initialize encryptor
	nodeMessages := dao.SearchNodeMsg("edge1")

	response := map[string]interface{}{
		"nodeMessage": nodeMessages,
	}
	fmt.Println(response)
}
