package utils

import (
	"EdgeGovernor/pkg/constants"
	"github.com/cornelk/hashmap"
	"log"
)

var NodeTables *NodesTable
var NodeList = []string{"cloud", "edge1", "edge2", "edge3"}

type Node struct {
	ID        int64
	Hostname  string
	IP        string
	Status    string
	Role      string
	CPU       float64
	Memory    int64
	DiskSpace int64
	Net       int64
}

type NodesTable struct { //节点池
	m *hashmap.HashMap
}

func NewNodeTable() *NodesTable {
	return &NodesTable{
		m: &hashmap.HashMap{},
	}
}

func (n *NodesTable) GetNodeEntry(nodeName string) (*Node, bool) {
	v, ok := n.m.Get(nodeName)
	if !ok {
		return nil, false
	}
	return v.(*Node), true
}

func (n *NodesTable) GetNodeIP(nodeName string) (string, bool) {
	v, ok := n.m.Get(nodeName)
	if !ok {
		return "", false
	}
	return v.(*Node).IP, true
}

func (n *NodesTable) GetNodeID(nodeName string) (int64, bool) {
	v, ok := n.m.Get(nodeName)
	if !ok {
		return 0, false
	}
	return v.(*Node).ID, true
}

func (n *NodesTable) GetNodeStatus(nodeName string) (string, bool) {
	v, ok := n.m.Get(nodeName)
	if !ok {
		return "", false
	}
	return v.(*Node).Status, true
}

func (n *NodesTable) GetNodeResource(nodeName string) (float64, int64, int64, int64, bool) {
	v, ok := n.m.Get(nodeName)
	if !ok {
		return 0, 0, 0, 0, false
	}
	return v.(*Node).CPU, v.(*Node).Memory, v.(*Node).DiskSpace, v.(*Node).Net, true
}

func (n *NodesTable) GetNodeRole(nodeName string) (string, bool) {
	v, ok := n.m.Get(nodeName)
	if !ok {
		return "", false
	}
	return v.(*Node).Role, true
}

func (n *NodesTable) Set(nodeName string, entry *Node) {
	if entry == nil {
		log.Panic("you can't enter nil entry")
	}
	n.m.Set(nodeName, entry)
}

func (n *NodesTable) Del(name string) {
	n.m.Del(name)
}

func (n *NodesTable) Len() int {
	return n.m.Len()
}

func (n *NodesTable) Clear() {
	n.m = &hashmap.HashMap{}
}

func (n *NodesTable) GetAllNodeName() []string {
	var nodeNames []string
	iter := n.m.Iter()
	// 遍历哈希映射中的键值对
	for kv := range iter {
		a, _ := InterfaceToString(kv.Key)
		nodeNames = append(nodeNames, a)
	}

	return nodeNames
}

func (n *NodesTable) UpdateNodeStatus(nodeName string, status string) bool {
	entry, ok := n.GetNodeEntry(nodeName)
	if !ok {
		return false
	}
	if entry.Status != status {
		if status == "Online" {
			constants.LiveNodeCount++
		} else {
			constants.LiveNodeCount--
		}
		entry.Status = status
		NodeTables.Set(nodeName, entry)
	}

	return true
}

func (n *NodesTable) UpdateNodeRole(nodeName string, role string) bool {
	entry, ok := n.GetNodeEntry(nodeName)
	if !ok {
		return false
	}
	if entry.Role != role {
		entry.Role = role
		NodeTables.Set(nodeName, entry)
	}
	return true
}

func (n *NodesTable) GetLiveNode() [][]string {
	var nodes [][]string
	for node := range NodeTables.Iter() {
		if node.Value.Status == "Offline" || node.Key == constants.Leader {
			continue
		}
		subList := []string{node.Key, node.Value.IP}
		nodes = append(nodes, subList)
	}
	return nodes
}

func (n *NodesTable) GetFollowerNode() [][]string {
	var nodes [][]string
	for node := range NodeTables.Iter() {
		if node.Key == constants.Leader {
			continue
		}
		subList := []string{node.Key, node.Value.IP, node.Value.Status}
		nodes = append(nodes, subList)
	}
	return nodes
}

func (n *NodesTable) Iter() <-chan struct {
	Key   string
	Value *Node
} {
	ch := make(chan struct {
		Key   string
		Value *Node
	})

	go func() {
		defer close(ch)
		for kv := range n.m.Iter() {
			ch <- struct {
				Key   string
				Value *Node
			}{Key: kv.Key.(string), Value: kv.Value.(*Node)}
		}
	}()

	return ch
}

func StorageNodeMsg() {
	NodeTables = NewNodeTable()
	NodeTables.Set("cloud", &Node{
		ID:        1,
		IP:        "192.168.47.152",
		Status:    "Online",
		Role:      "Leader",
		CPU:       8000,
		Memory:    8192,
		DiskSpace: 40,
		Net:       1000,
	})

	NodeTables.Set("edge1", &Node{
		ID:        2,
		IP:        "192.168.47.153",
		Status:    "Online",
		Role:      "Follower",
		CPU:       4000,
		Memory:    4096,
		DiskSpace: 40,
		Net:       1000,
	})

	NodeTables.Set("edge2", &Node{
		ID:        3,
		IP:        "192.168.47.154",
		Status:    "Online",
		Role:      "Follower",
		CPU:       1000,
		Memory:    2048,
		DiskSpace: 40,
		Net:       1000,
	})

	NodeTables.Set("edge3", &Node{
		ID:        4,
		IP:        "192.168.47.155",
		Status:    "Online",
		Role:      "Follower",
		CPU:       1000,
		Memory:    2048,
		DiskSpace: 40,
		Net:       1000,
	})
}
