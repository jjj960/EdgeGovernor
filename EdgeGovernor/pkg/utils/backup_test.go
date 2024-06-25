package utils

import (
	"EdgeGovernor/pkg/sec"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	GetETCDCli()
	sec.Safer = sec.NewMSecLayer("2TXQfxBE#TrULn6FZx")
	NodeTables = NewNodeTable()
	NodeTables.Set("cloud", &Node{
		IP:        "192.168.47.152",
		Status:    "Online",
		Role:      "Leader",
		CPU:       8000,
		Memory:    8192,
		DiskSpace: 40,
		Net:       1000,
	})

	NodeTables.Set("edge1", &Node{
		IP:        "192.168.47.153",
		Status:    "Online",
		Role:      "Follower",
		CPU:       4000,
		Memory:    4096,
		DiskSpace: 40,
		Net:       1000,
	})

	BackUpNodeMap()
	GetNodeTablesMap()
	fmt.Println(NodeTables.GetNodeIP("edge1"))
}
