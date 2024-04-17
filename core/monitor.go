package core

import (
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// Monitor ...
func Monitor() (int64, int64, int64, int64, int64) {
	memInfo, _ := mem.VirtualMemory()
	cpuInfo, _ := cpu.Percent(0, true)
	diskInfo, _ := disk.Usage("/")
	netInfo, _ := net.IOCounters(true)

	totalCPU := 0.0
	for _, v := range cpuInfo {
		totalCPU += v
	}
	totalCPU = totalCPU / float64(len(cpuInfo))

	sendByte := 0.0
	recvByte := 0.0
	for _, v := range netInfo {
		sendByte += float64(v.BytesSent)
		recvByte += float64(v.BytesRecv)
	}

	return int64(memInfo.UsedPercent), int64(totalCPU), int64(diskInfo.UsedPercent), int64(sendByte), int64(recvByte)
}

type MonitorAttr struct {
	Mem       int64 `json:"memory"`
	Cpu       int64 `json:"cpu"`
	Disk      int64 `json:"disk"`
	NetSend   int64 `json:"net_send"`
	NetRecv   int64 `json:"net_recv"`
	Goroutine int64 `json:"goroutine"`
}

// Monitor ...
func MonitorEx() *MonitorAttr {
	memInfo, _ := mem.VirtualMemory()
	cpuInfo, _ := cpu.Percent(0, true)
	diskInfo, _ := disk.Usage("/")
	netInfo, _ := net.IOCounters(true)
	goroutine := runtime.NumGoroutine()
	//runtime.ReadMemStats()

	totalCPU := 0.0
	for _, v := range cpuInfo {
		totalCPU += v
	}
	totalCPU = totalCPU / float64(len(cpuInfo))

	sendByte := 0.0
	recvByte := 0.0
	for _, v := range netInfo {
		sendByte += float64(v.BytesSent)
		recvByte += float64(v.BytesRecv)
	}

	return &MonitorAttr{
		Mem:       int64(memInfo.UsedPercent),
		Cpu:       int64(totalCPU),
		Disk:      int64(diskInfo.UsedPercent),
		NetSend:   int64(sendByte),
		NetRecv:   int64(recvByte),
		Goroutine: int64(goroutine),
	}
}
