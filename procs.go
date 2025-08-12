package main

import (
	"encoding/json"
	"log"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	PID         int32   `json:"pid"`
	Name        string  `json:"name"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryMB    float64 `json:"memory_mb"`
	MemoryUsage float32 `json:"memory_percent"`
}

func getProcs() string {
	procs, err := process.Processes()
	if err != nil {
		log.Fatalf("Ошибка получения списка процессов: %v", err)
	}

	var processList []ProcessInfo

	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			name = "ERR"
		}

		cpu, err := p.CPUPercent()
		if err != nil {
			cpu = 0
		}

		memPercent, err := p.MemoryPercent()
		if err != nil {
			memPercent = 0
		}

		memInfo, err := p.MemoryInfo()
		var memoryMb float64
		if err != nil {
			memoryMb = 0
		} else {
			memoryMb = float64(memInfo.RSS) / 1024.0 / 1024.0
		}

		processList = append(processList, ProcessInfo{
			PID:         p.Pid,
			Name:        name,
			CPUPercent:  cpu,
			MemoryMB:    memoryMb,
			MemoryUsage: memPercent,
		})
	}

	jsonData, err := json.MarshalIndent(processList, "", "  ")
	if err != nil {
		log.Fatalf("Ошибка сериализации JSON: %v", err)
	}

	return string(jsonData)
}
