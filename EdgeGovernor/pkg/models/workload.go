package models

type Hostload struct {
	Timestamp         int64   `json:"timestamp"`
	Hostname          string  `json:"hostname"`
	CPUUsagePercent   float64 `json:"cpu_usage_percent"`
	CPUCapacity       int64   `json:"cpu_capacity"`
	CPUResidue        int64   `json:"cpu_residue"`
	MemoryUsedPercent float64 `json:"memory_used_percent"`
	MemoryCapacity    int64   `json:"memory_capacity"`
	MemoryResidue     int64   `json:"memory_residue"`
	DiskUsedPercent   float64 `json:"disk_used_percent"`
	DiskCapacity      int64   `json:"disk_capacity"`
	DiskResidue       int64   `json:"disk_residue"`
	BytesRecv         float64 `json:"bytes_recv"`
	BytesSent         float64 `json:"bytes_sent"`
	BandWidth         float64 `json:"band-width"`
}

type ReceiveBytes uint64

type TransmitBytes uint64

type ContainerLoad struct {
	Timestamp        int64   `json:"timestamp"`
	Name             string  `json:"name"`
	ID               string  `json:"id"`
	CPUPercentage    float64 `json:"cpu_percentage"`
	Memory           int64   `json:"memory"`
	MemoryLimit      int64   `json:"memory_limit"`
	MemoryPercentage float64 `json:"memory_percentage"`
	NetworkRx        float64 `json:"network_rx"`
	NetworkTx        float64 `json:"network_tx"`
	BlockRead        float64 `json:"block_read"`
	BlockWrite       float64 `json:"block_write"`
	PidsCurrent      int64   `json:"pids_current"`
}
