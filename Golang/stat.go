package main

import (
	"hwraid/topic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var statTopic = topic.New()

type Stat struct {
	CPU       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Temp      float64 `json:"temp"`
	ReadMB    float64 `json:"read_mb"`
	WriteMB   float64 `json:"write_mb"`
	FReadMB   float64 `json:"fread_mb"`
	FWriteMB  float64 `json:"fwrite_mb"`
	NReadMB   float64 `json:"nread_mb"`
	NWriteMB  float64 `json:"nwrite_mb"`
	Timestamp int64   `json:"timestamp"`
}

func prettyFloat(f float64) float64 {
	return float64(int64(f*100)) / 100
}

func ServeStat() {
	go func() {
		var stat Stat
		lastNStats, _ := net.NetIOCounters(false)
		for {
			n1 := time.Now()
			percents, err := cpu.CPUPercent(10*time.Second, false)
			if err != nil {
				stat.CPU = 0
			}
			stat.CPU = prettyFloat(float64(percents[0]))

			vm, err := mem.VirtualMemory()
			if err != nil {
				stat.Mem = 0
			}
			stat.Mem = prettyFloat(vm.UsedPercent)

			stat.Temp = 40

			nstats, err := net.NetIOCounters(false)
			if err != nil {
				nstats = lastNStats
			}
			n2 := time.Now()
			delta := n2.Sub(n1)
			stat.NWriteMB = prettyFloat(float64(nstats[0].BytesSent-lastNStats[0].BytesSent) / delta.Seconds() / (1024 * 1024))
			stat.NReadMB = prettyFloat(float64(nstats[0].BytesRecv-lastNStats[0].BytesRecv) / delta.Seconds() / (1024 * 1024))
			stat.Timestamp = n2.Unix()

			lastNStats = nstats
			statTopic.Publish(stat)
		}
	}()
}
