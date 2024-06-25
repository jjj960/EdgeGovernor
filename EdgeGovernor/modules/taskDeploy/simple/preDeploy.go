package simple

import (
	algorithm2 "EdgeGovernor/pkg/cache/algorithm"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"strings"
)

//func ImagePrePull(task models.Task) error {
//	if task.DeployNode != "" {
//		err := imagePull(task.DeployNode, task.Image)
//		if err != nil {
//			return fmt.Errorf("Mirror pull failure for task %s", task.Name)
//		}
//	} else {
//		nodeMsg, err := getWorkloadForcast()
//		if err != nil {
//			return errors.New("Failed to obtain cluster load information")
//		}
//		fmt.Println("Attempting to pre deploy task:", task.Name)
//		//请求调度算法
//		//var requestMsg [][]string
//		var taskMsg []string
//		//requestMsg = append(requestMsg, nodeMsg...)
//		taskMsg = append(taskMsg, task.Name, task.DeployNode,
//			utils.Int64toString(task.RequestCPU),
//			utils.Int64toString(task.RequestMem),
//			utils.Int64toString(task.RequestNet),
//			utils.Int64toString(task.RequestDisk),
//			strconv.Itoa(task.Scalable))
//
//		maxHost, _ := GetSchedulerScore(nodeMsg, task.SchedulingAlgorithm)
//		fmt.Println("预调度的最优节点为: ", maxHost)
//		check := alogorithm.CheckNodeResource(maxHost, nodeMsg, taskMsg)
//		if !check {
//			return fmt.Errorf("The task %s may not meet the deployment conditions, and the requested memory or disk resources are too large.", task.Name)
//		} else {
//			totalRequestCPU, _, _, _ := GetTotalRequestResource()
//
//			cpu, mem := alogorithm.ResourceAllocation(maxHost, nodeMsg, taskMsg, totalRequestCPU) //资源分配
//			var nodeip string
//			fmt.Println("节点信息为：", nodeMsg)
//			for i, record := range nodeMsg {
//				if record[0] == maxHost {
//					fmt.Println("节点信息为：", nodeMsg[i][1])
//					nodeip = nodeMsg[i][1]
//					orReCPU, _ := utils.StringtoInt64(nodeMsg[i][4])   //原始剩余CPU数量
//					orReMem, _ := utils.StringtoInt64(nodeMsg[i][7])   //原始剩余内存数量
//					orNet, _ := utils.StringtoInt64(nodeMsg[i][11])    //原始网络带宽
//					orReDisk, _ := utils.StringtoInt64(nodeMsg[i][10]) //原始剩余磁盘数量
//
//					capCPU, _ := utils.StringtoInt64(nodeMsg[i][3])
//					capMem, _ := utils.StringtoInt64(nodeMsg[i][6])
//					capDisk, _ := utils.StringtoInt64(nodeMsg[i][9])
//
//					newReCPU := orReCPU - cpu
//					newReMem := orReMem - mem
//					newNet := orNet + task.RequestNet
//					newReDisk := orReDisk - task.RequestDisk
//
//					newPercentCPU := 1 - float64(newReCPU)/float64(capCPU)
//					newPercentMem := 1 - float64(newReMem)/float64(capMem)
//					newPercentDisk := 1 - float64(newReDisk)/float64(capDisk)
//					// 修改相应的字段值
//					nodeMsg[i][2] = utils.Float64toString(newPercentCPU)  //更新CPU占比
//					nodeMsg[i][4] = utils.Int64toString(newReCPU)         //更新剩余CPU数量
//					nodeMsg[i][5] = utils.Float64toString(newPercentMem)  //更新内存占比
//					nodeMsg[i][7] = utils.Int64toString(newReMem)         //更新剩余内存数量
//					nodeMsg[i][8] = utils.Float64toString(newPercentDisk) //更新磁盘占比
//					nodeMsg[i][10] = utils.Int64toString(newReDisk)       //更新剩余磁盘数量
//					nodeMsg[i][11] = utils.Int64toString(newNet)          //更新网络带宽
//
//					break // 找到对应记录后退出循环
//				}
//			}
//			log.Println("ip地址为：", nodeip)
//			log.Println("开始进行镜像预缓存，缓存镜像：", task.Image, nodeip, maxHost)
//			_, err = utils.SingleSend(nodeip, maxHost, "image cache", task.Image) //通知目标节点提前缓存镜像
//			if err != nil {
//				return fmt.Errorf("The task %s may not meet the deployment conditions, and the requested memory or disk resources are too large.", task.Name)
//			}
//			log.Printf("Mirror pull failure for task %s", task.Name)
//		}
//	}
//	return nil
//}
//
//func Predeploy(taskNames []string) { //预调度功能
//	var tasks []models.Task
//	var taskDetails []string
//	for _, taskName := range taskNames {
//		taskdetail := mariadb.GetTaskMsg(taskName) //获取该任务的详细信息
//		taskDetails = append(taskDetails, taskdetail)
//	}
//	for _, jsonStr := range taskDetails {
//		var task models.Task
//
//		err := utils.Jsoniter.Unmarshal([]byte(jsonStr), &task)
//		if err != nil {
//			fmt.Println("Failed to unmarshal JSON:", err)
//			continue
//		}
//		tasks = append(tasks, task)
//	}
//
//	// 按照 优先级 值进行排序
//	sort.Slice(tasks, func(i, j int) bool {
//		return tasks[i].Priority > tasks[j].Priority
//	})
//
//	nodeMsg, err := getWorkloadForcast()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	// 根据优先级,从大到小输出taskName
//	for i, task := range tasks {
//		fmt.Println("Attempting to pre deploy task:", task.Name)
//		//请求调度算法
//		//var requestMsg [][]string
//		var taskMsg []string
//		//requestMsg = append(requestMsg, nodeMsg...)
//		taskMsg = append(taskMsg, task.Name, task.DeployNode,
//			utils.Int64toString(task.RequestCPU),
//			utils.Int64toString(task.RequestMem),
//			utils.Int64toString(task.RequestNet),
//			utils.Int64toString(task.RequestDisk),
//			strconv.Itoa(task.Scalable))
//
//		maxHost, _ := GetSchedulerScore(nodeMsg, task.SchedulingAlgorithm)
//		fmt.Println("预调度的最优节点为: ", maxHost)
//		check := alogorithm.CheckNodeResource(maxHost, nodeMsg, taskMsg)
//		if !check {
//			log.Printf("The task %s may not meet the deployment conditions, and the requested memory or disk resources are too large.", task.Name)
//		} else {
//			var totalRequestCPU int64
//			for c := i + 1; c < len(tasks); c++ {
//				totalRequestCPU += tasks[c].RequestCPU
//			}
//			cpu, mem := alogorithm.ResourceAllocation(maxHost, nodeMsg, taskMsg, totalRequestCPU) //资源分配
//			var nodeip string
//			fmt.Println("节点信息为：", nodeMsg)
//			for i, record := range nodeMsg {
//				if record[0] == maxHost {
//					fmt.Println("节点信息111为：", nodeMsg[i][1])
//					nodeip = nodeMsg[i][1]
//					orReCPU, _ := utils.StringtoInt64(nodeMsg[i][4])   //原始剩余CPU数量
//					orReMem, _ := utils.StringtoInt64(nodeMsg[i][7])   //原始剩余内存数量
//					orNet, _ := utils.StringtoInt64(nodeMsg[i][11])    //原始网络带宽
//					orReDisk, _ := utils.StringtoInt64(nodeMsg[i][10]) //原始剩余磁盘数量
//
//					capCPU, _ := utils.StringtoInt64(nodeMsg[i][3])
//					capMem, _ := utils.StringtoInt64(nodeMsg[i][6])
//					capDisk, _ := utils.StringtoInt64(nodeMsg[i][9])
//
//					newReCPU := orReCPU - cpu
//					newReMem := orReMem - mem
//					newNet := orNet + task.RequestNet
//					newReDisk := orReDisk - task.RequestDisk
//
//					newPercentCPU := 1 - float64(newReCPU)/float64(capCPU)
//					newPercentMem := 1 - float64(newReMem)/float64(capMem)
//					newPercentDisk := 1 - float64(newReDisk)/float64(capDisk)
//					// 修改相应的字段值
//					nodeMsg[i][2] = utils.Float64toString(newPercentCPU)  //更新CPU占比
//					nodeMsg[i][4] = utils.Int64toString(newReCPU)         //更新剩余CPU数量
//					nodeMsg[i][5] = utils.Float64toString(newPercentMem)  //更新内存占比
//					nodeMsg[i][7] = utils.Int64toString(newReMem)         //更新剩余内存数量
//					nodeMsg[i][8] = utils.Float64toString(newPercentDisk) //更新磁盘占比
//					nodeMsg[i][10] = utils.Int64toString(newReDisk)       //更新剩余磁盘数量
//					nodeMsg[i][11] = utils.Int64toString(newNet)          //更新网络带宽
//
//					break // 找到对应记录后退出循环
//				}
//			}
//			log.Println("ip地址为：", nodeip)
//			log.Println("开始进行镜像预缓存，缓存镜像：", task.Image, nodeip, maxHost)
//			utils.SingleSend(nodeip, maxHost, "image cache", task.Image) //通知目标节点提前缓存镜像
//			log.Printf("Successfully cached the image %s in advance", task.Image)
//		}
//	}
//}

//func getWorkloadForcast() ([][]string, error) {
//	str, err := utils.ForecastRequest() //请求负载调度微服务,获取未来5分钟的负载情况
//	if err != nil {
//		return nil, fmt.Errorf("Failed to obtain workload prediction data:", err)
//	}
//	cleanedStr := strings.ReplaceAll(str, "[", "")
//	cleanedStr = strings.ReplaceAll(cleanedStr, "]", "")
//	arr := strings.Split(cleanedStr, " ")
//	var nodeData [][]string
//	size := 14
//
//	for i := 0; i < len(arr); i += size {
//		end := i + size
//
//		if end > len(arr) {
//			end = len(arr)
//		}
//
//		nodeData = append(nodeData, arr[i:end])
//	}
//
//	return nodeData, nil
//}

func GetSchedulerScore(nodeMsg [][]string, algorithm string) (string, error) {

	result, err := utils.SchedulingRequest(algorithm2.GetAlgorithmURL(algorithm, "Schedule"), nodeMsg)
	if err != nil {
		return "", fmt.Errorf("Failed to obtain scheduling data:", err)
	}

	return strings.ReplaceAll(result, "\"", ""), nil
}

//func imagePull(nodeName string, image string) error {
//	nodeip := mariadb.GetNodeIP(nodeName)
//	result, _ := utils.SingleSend(nodeip, nodeName, "image cache", image) //通知目标节点提前缓存镜像
//	if result != "success pull" {
//		return errors.New(result)
//	} else {
//		log.Printf("Successfully cached the image %s in advance", image)
//	}
//	return nil
//}
