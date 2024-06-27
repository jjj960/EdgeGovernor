package logging

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"io/ioutil"
	"log"
	"os"

	"strconv"
	"strings"
	"syscall"
	"time"
)

func GetHostWorkload() []byte { //获取主机负载
	cpupercent, cpucap, cpure := getCPU()
	mempercent, memcap, memre := getMemory()
	percent, free, total, _ := getDisk()
	byterecv, bytesent, bandwidth := getNetIO("ens33") //网速下行和上行,单位为Kb/s
	hostname, _ := os.Hostname()
	data := models.Hostload{
		Timestamp:         time.Now().UnixNano() / int64(time.Millisecond),
		Hostname:          hostname,
		CPUUsagePercent:   cpupercent[0],
		CPUCapacity:       cpucap,
		CPUResidue:        int64(cpure),
		MemoryUsedPercent: mempercent,
		MemoryCapacity:    memcap,
		MemoryResidue:     memre,
		DiskUsedPercent:   percent,
		DiskCapacity:      total,
		DiskResidue:       free,
		BytesRecv:         byterecv,
		BytesSent:         bytesent,
		BandWidth:         bandwidth,
	}
	selfRole, _ := utils.NodeTables.GetNodeRole(constants.Hostname)
	if selfRole == "Leader" && (cpupercent[0] >= 80.0 || mempercent >= 80.0) {
		id, time := utils.GetID()
		utils.AlarmMsgChannel <- models.Msg{
			ID:           id,
			GenerateTime: time,
			Tpye:         "Host Resource warning",
			Detail:       []byte(fmt.Sprintf("Node %s has insufficient resources, CPU load is %f, memory load is %f", hostname, cpupercent, mempercent)),
			Status:       false,
		}
	}
	// 将结构体转换为JSON格式
	jsonData, err := utils.Jsoniter.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	//fmt.Println(string(jsonData))

	return jsonData
}

func getNetIO(dev string) (float64, float64, float64) {
	down, up, err := totalFlowByDevice(dev)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(time.Second * 1)
	down2, up2, err := totalFlowByDevice(dev)
	if err != nil {
		log.Println(err)
	}

	downStr := float64((down2 - down) / 1024) //以kb为单位
	upStr := float64((up2 - up) / 1024)
	bandwidth := constants.BandWidth
	return downStr, upStr, bandwidth
}

func totalFlowByDevice(dev string) (models.ReceiveBytes, models.TransmitBytes, error) {
	devInfo, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}
	var receive = -1
	var transmit = -1
	var receiveBytes uint64
	var transmitBytes uint64
	lines := strings.Split(string(devInfo), "\n")
	for _, line := range lines {
		if strings.Contains(line, dev) {
			i := 0
			fields := strings.Split(line, ":")
			for _, field := range fields {
				if strings.Contains(field, dev) {
					i = 1
				} else {
					values := strings.Fields(field)
					for _, value := range values {
						//logger.Debug(value)
						if receive == i {
							bytes, _ := strconv.ParseInt(value, 10, 64)
							receiveBytes = uint64(bytes)
						} else if transmit == i {
							bytes, _ := strconv.ParseInt(value, 10, 64)
							transmitBytes = uint64(bytes)
						}
						i++
					}
				}
			}
		} else if strings.Contains(line, "face") {
			index := 0
			tag := false
			fields := strings.Split(line, "|")
			for _, field := range fields {
				if strings.Contains(field, "face") {
					index = 1
				} else if strings.Contains(field, "bytes") {
					values := strings.Fields(field)
					for _, value := range values {
						//logger.Debug(value)
						if strings.Contains(value, "bytes") {
							if !tag {
								tag = true
								receive = index
							} else {
								transmit = index
							}
						}
						index++
					}
				}
			}
		}
	}
	//log.Printf("receive_bytes :", receiveBytes)
	//log.Printf("transmit_bytes :", transmitBytes)
	return models.ReceiveBytes(receiveBytes), models.TransmitBytes(transmitBytes), nil
}

func getDisk() (float64, int64, int64, int64) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, 0, 0, 0
	}

	total, _ := utils.BytesToMB(int64(stat.Blocks) * stat.Bsize) // 总空间
	free, _ := utils.BytesToMB(int64(stat.Bfree) * stat.Bsize)   // 可用空间
	used := total - free                                         // 已使用空间
	percent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(used)/float64(total)*100), 64)
	return percent, free, total, used

}

func getCPU() ([]float64, int64, float64) {
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	corenum, _ := cpu.Counts(true)
	residue := float64(corenum) * 10 * (100 - percent[0])
	return percent, int64(corenum * 1000), residue
}

func getMemory() (float64, int64, int64) {
	memInfo, _ := mem.VirtualMemory()
	memTotal, _ := utils.BytesToMB(int64(memInfo.Total))
	memAva, _ := utils.BytesToMB(int64(memInfo.Available))
	return memInfo.UsedPercent, memTotal, memAva
}
