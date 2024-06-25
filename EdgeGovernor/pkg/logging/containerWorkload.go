package logging

import (
	"EdgeGovernor/pkg/constants"
	"EdgeGovernor/pkg/database/duckdb"
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"strings"
	"time"
)

func GetContainerWorkload() {
	//获取ctx
	ctx := context.Background()
	//获取容器id
	containers, err := utils.DockerCli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		result := strings.Join(container.Names, "")[1:]
		if strings.Contains(result, "k8s") && strings.Contains(result, "menet") && !strings.Contains(result, "POD") {
			if err != nil {
				panic(err)
			}
			getContainerStats(result)
			//jsonData, err := utils.Jsoniter.Marshal(a)
			//if err != nil {
			//	fmt.Println("Error:", err)
			//}
		}
	}
	////fmt.Println(workload)
	//return workload
}

func getContainerStats(containerId string) (*models.ContainerLoad, error) {
	ctx := context.Background()

	containerStats, err := utils.DockerCli.ContainerStats(ctx, containerId, false)
	if err != nil {
		return nil, err
	}
	defer containerStats.Body.Close()

	var v *types.StatsJSON
	dec := utils.Jsoniter.NewDecoder(containerStats.Body)
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}

	return getCollectorMetrics(v), nil
}

// 处理types.StatsJSON数据
func getCollectorMetrics(stats *types.StatsJSON) *models.ContainerLoad {
	var (
		memPercent        = 0.0
		cpuPercent        = 0.0
		blkRead, blkWrite uint64
		mem               = 0.0
		memLimit          = 0.0
		pids              uint64
		netRx             = 0.0
		netTx             = 0.0
	)
	//cpu
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
	cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	//memory
	mem = float64(stats.MemoryStats.Usage)
	memLimit = float64(stats.MemoryStats.Limit)
	if stats.MemoryStats.Limit != 0 {
		memPercent = float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit) * 100.0
	}
	//network
	for _, v := range stats.Networks {
		netRx += float64(v.RxBytes)
		netTx += float64(v.TxBytes)
	}
	//block
	var blkio = stats.BlkioStats
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = blkRead + bioEntry.Value
		case "write":
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	//pidsCurrent
	pids = stats.PidsStats.Current
	containerSubName, _ := utils.GetSubString(strings.ReplaceAll(stats.Name[1:], "-", ""))
	selfRole, _ := utils.NodeTables.GetNodeRole(constants.Hostname)
	if selfRole == "master" && (cpuPercent >= 80.0 || memPercent >= 80.0) {
		id := utils.SnowFlake.Generate()
		utils.AlarmMsgChannel <- models.Msg{
			ID:           id.Int64(),
			GenerateTime: id.Time(),
			Tpye:         "Container Resource warning",
			Detail:       []byte(fmt.Sprintf("Container %s has insufficient resources, CPU load is %f, memory load is %f", containerSubName, cpuPercent, memPercent)),
			Status:       false,
		}
	}
	load := &models.ContainerLoad{
		Timestamp:        time.Now().UnixNano() / int64(time.Millisecond),
		Name:             containerSubName,
		ID:               stats.ID,
		CPUPercentage:    cpuPercent,
		Memory:           int64(mem),
		MemoryLimit:      int64(memLimit),
		MemoryPercentage: memPercent,
		NetworkRx:        netRx,
		NetworkTx:        netTx,
		BlockRead:        float64(blkRead),
		BlockWrite:       float64(blkWrite),
		PidsCurrent:      int64(pids),
	}

	duckdb.CreateContainerloadTable(containerSubName)
	duckdb.InsertContainerload(containerSubName, *load)
	return load
}
