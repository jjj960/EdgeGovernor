package resource

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/logging"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	"github.com/buger/jsonparser"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"

	"log"
)

func InsertNode(node utils.Node) error {
	key := "/menet/node/" + node.Hostname
	nodeData, err := utils.Jsoniter.Marshal(node)
	if err != nil {
		return fmt.Errorf("JSON encoding failed: %s", err)
	}
	value := string(nodeData)

	_, err = utils.ETCDCli.Put(context.Background(), key, value)
	if err != nil {
		return fmt.Errorf("Node data insertion failed: %s", err)
	}

	return nil
}

func GetNodeMsg(nodeName string) (utils.Node, error) {
	key := "/menet/node/" + nodeName

	resp, err := utils.ETCDCli.Get(context.Background(), key)
	if err != nil {
		return utils.Node{}, fmt.Errorf("failed to retrieve data from etcd: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return utils.Node{}, fmt.Errorf("failed to find node from etcd: %v", err)
	}

	value := string(resp.Kvs[0].Value)

	var nodeMsg utils.Node
	err = utils.Jsoniter.Unmarshal([]byte(value), &nodeMsg)
	if err != nil {
		return utils.Node{}, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return nodeMsg, nil
}

func GetAllNodeName() ([]string, error) {
	prefix := "/menet/node"

	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain value: %v", err)
	}

	var NodeList []string

	for _, kv := range resp.Kvs {
		NodeList = append(NodeList, string(kv.Key))
	}

	return NodeList, nil
}

func GetNodeNum() (int, error) {
	prefix := "/menet/node"

	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return 0, fmt.Errorf("Failed to obtain value: %v", err)
	}

	return len(resp.Kvs), nil
}

func GetLiveFollowerNode() ([][]string, error) {
	prefix := "/menet/node"

	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain value: %v", err)
	}

	var NodeList [][]string

	for _, kv := range resp.Kvs {
		var nodeMsg utils.Node
		value := string(kv.Value)
		err = utils.Jsoniter.Unmarshal([]byte(value), &nodeMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
		}
		if nodeMsg.Role != "Follower" && nodeMsg.Status != "Online" {
			continue
		}

		node := []string{nodeMsg.Hostname, nodeMsg.IP}
		NodeList = append(NodeList, node)
	}

	return NodeList, nil
}

func GetFollowerNode() ([][]string, error) {
	prefix := "/menet/node"

	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain value: %v", err)
	}

	var NodeList [][]string

	for _, kv := range resp.Kvs {
		var nodeMsg utils.Node
		value := string(kv.Value)
		err = utils.Jsoniter.Unmarshal([]byte(value), &nodeMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
		}
		if nodeMsg.Role != "Follower" {
			continue
		}

		node := []string{nodeMsg.Hostname, nodeMsg.IP}
		NodeList = append(NodeList, node)
	}

	return NodeList, nil
}

func UpdateNodeStatus(nodeName, status string) error {
	nodeMsg, err := GetNodeMsg(nodeName)
	if err != nil {
		return err
	}

	if nodeMsg.Status != status {
		if status == "Online" {
			constants.LiveNodeCount++
		} else {
			constants.LiveNodeCount--
		}
		nodeMsg.Status = status
		err = InsertNode(nodeMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateNodeRole(nodeName, role string) error {
	nodeMsg, err := GetNodeMsg(nodeName)
	if err != nil {
		return err
	}

	if nodeMsg.Role != role {
		nodeMsg.Role = role
		err = InsertNode(nodeMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

func SearchNodeMsg(Hostname string) map[string]string {
	nodeMsg, _ := GetNodeMsg(Hostname)
	node1 := make(map[string]string)
	taskNums, _ := getNodeTaskNum()
	var CPU, Mem, Disk string
	if nodeMsg.Status == "deactive" {
		CPU, Mem, Disk = "-", "-", "-"
	} else {
		if nodeMsg.Role == "leader" {
			result := logging.GetHostWorkload()
			var data models.Hostload
			err := utils.Jsoniter.Unmarshal(result, &data)
			if err != nil {
				log.Println("解析JSON失败：", err)
			}
			CPU = utils.Int64toString(data.CPUCapacity-data.CPUResidue) + " / " + utils.Int64toString(data.CPUCapacity)
			Mem = utils.Int64toString(data.MemoryCapacity-data.MemoryResidue) + " / " + utils.Int64toString(data.MemoryCapacity)
			Disk = utils.Int64toString(data.DiskCapacity-data.DiskResidue) + " / " + utils.Int64toString(data.DiskCapacity)
		} else {
			CPU, Mem, Disk = collectFollowerMachineStatus(Hostname, nodeMsg.IP)
		}
	}

	node1["nodeName"] = Hostname
	node1["cpuSize"] = CPU
	node1["memorySize"] = Mem
	node1["diskSize"] = Disk
	node1["equipment"] = nodeMsg.Role
	node1["status"] = nodeMsg.Status
	node1["taskNum"] = strconv.Itoa(taskNums[Hostname])

	return node1
}

func GetAllNodeMsg() ([]map[string]string, error) {
	prefix := "/menet/node"

	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain value: %v", err)
	}

	var nodeMsg utils.Node
	taskNums, _ := getNodeTaskNum()
	var result []map[string]string
	var CPU, Mem, Disk string

	for _, kv := range resp.Kvs {
		err = utils.Jsoniter.Unmarshal(kv.Value, &nodeMsg)
		if nodeMsg.Status == "Offline" {
			CPU, Mem, Disk = "-", "-", "-"
		} else {
			if nodeMsg.Role == "Leader" {
				result := logging.GetHostWorkload()
				var data models.Hostload
				err := utils.Jsoniter.Unmarshal(result, &data)
				if err != nil {
					log.Println("解析JSON失败：", err)
				}
				CPU = utils.Int64toString(data.CPUCapacity-data.CPUResidue) + " / " + utils.Int64toString(data.CPUCapacity)
				Mem = utils.Int64toString(data.MemoryCapacity-data.MemoryResidue) + " / " + utils.Int64toString(data.MemoryCapacity)
				Disk = utils.Int64toString(data.DiskCapacity-data.DiskResidue) + " / " + utils.Int64toString(data.DiskCapacity)
			} else {
				CPU, Mem, Disk = collectFollowerMachineStatus(nodeMsg.Hostname, nodeMsg.IP)
			}
		}
		node1 := make(map[string]string)
		node1["nodeName"] = nodeMsg.Hostname
		node1["cpuSize"] = CPU
		node1["memorySize"] = Mem
		node1["diskSize"] = Disk
		node1["equipment"] = nodeMsg.Role
		node1["status"] = nodeMsg.Status
		node1["taskNum"] = strconv.Itoa(taskNums[nodeMsg.Hostname])

		result = append(result, node1)
	}

	return result, nil
}

func GetNodesMsg(start int, end int) ([]map[string]string, error) {
	prefix := "/menet/node"
	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain value: %v", err)
	}
	var nodeMsg utils.Node
	taskNums, _ := getNodeTaskNum()
	var result []map[string]string
	var CPU, Mem, Disk string
	counter := 0
	for _, kv := range resp.Kvs {
		counter++
		if counter >= start && counter <= end {
			err = utils.Jsoniter.Unmarshal(kv.Value, &nodeMsg)
			if nodeMsg.Status == "Offline" {
				CPU, Mem, Disk = "-", "-", "-"
			} else {
				if nodeMsg.Role == "Leader" {
					result := logging.GetHostWorkload()
					var data models.Hostload
					err := utils.Jsoniter.Unmarshal(result, &data)
					if err != nil {
						log.Println("解析JSON失败：", err)
					}
					CPU = utils.Int64toString(data.CPUCapacity-data.CPUResidue) + " / " + utils.Int64toString(data.CPUCapacity)
					Mem = utils.Int64toString(data.MemoryCapacity-data.MemoryResidue) + " / " + utils.Int64toString(data.MemoryCapacity)
					Disk = utils.Int64toString(data.DiskCapacity-data.DiskResidue) + " / " + utils.Int64toString(data.DiskCapacity)
				} else {
					CPU, Mem, Disk = collectFollowerMachineStatus(nodeMsg.Hostname, nodeMsg.IP)
				}
			}
			node1 := make(map[string]string)
			node1["nodeName"] = nodeMsg.Hostname
			node1["cpuSize"] = CPU
			node1["memorySize"] = Mem
			node1["diskSize"] = Disk
			node1["equipment"] = nodeMsg.Role
			node1["status"] = nodeMsg.Status
			node1["taskNum"] = strconv.Itoa(taskNums[nodeMsg.Hostname])

			result = append(result, node1)
		}
		if counter > end {
			break
		}
	}

	return result, nil
}

func getNodeTaskNum() (map[string]int, error) {
	prefix := "/menet/task"
	taskNums := make(map[string]int)
	resp, err := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain value: %v", err)
	}

	for _, kv := range resp.Kvs {
		deployNode, _, _, _ := jsonparser.Get(kv.Value, "DeployNode")
		status, _, _, _ := jsonparser.Get(kv.Value, "Status")
		if string(status) == "Deployed" {
			taskNums[string(deployNode)]++
		}
	}
	return taskNums, nil
}

func collectFollowerMachineStatus(Hostname string, IP string) (string, string, string) {
	var data models.Hostload
	result, _ := utils.SingleSend(IP, Hostname, "simpleworkload report", "")
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
