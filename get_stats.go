package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func getCPUPercent() string {
	pcts, err := cpu.Percent(0, false)
	cpuPercent := 0.0

	if err != nil {
		log.Println("Ошибка в getCPUPercent: ", err)
		return fmt.Sprintf(`"cpu":{"loadPercent":%f}`, cpuPercent)
	}
	if len(pcts) == 0 {
		cpuPercent = 0
	} else {
		cpuPercent = pcts[0]
	}

	return fmt.Sprintf(`"cpu":{"loadPercent":%f}`, cpuPercent)
}

func getLoadAvg() string {
	loadAvg, err := load.Avg()
	if err != nil {
		log.Println("Ошибка в getLoadAvg: ", err)
		return `"loadAvg":{"1min":0, "5min":0, "15min": 0}`
	}

	return fmt.Sprintf(`"loadAvg":{"1min":%f, "5min":%f, "15min": %f}`, loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
}

func getMem() string {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Ошибка в getMem: ", err)
		return `"memory":{"loadPercent":0}`
	}
	return fmt.Sprintf(`"memory":{"loadPercent":%f}`, memStat.UsedPercent)
}

func getDisk() string {
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Println("Ошибка в getDisk: ", err)
		return `"disk":{"loadPercent":0}`
	}
	return fmt.Sprintf(`"disk":{"loadPercent":%f}`, diskStat.UsedPercent)
}

func getTemps() string {
	tempStats, err := host.SensorsTemperatures()
	if err != nil {
		log.Println("Ошибка в getTemps: ", err)
		return `"temp":{"temp":0}`
	}
	return fmt.Sprintf(`"temp":{"temp":%f}`, tempStats[0].Temperature)
}

func getStatsJson() string {
	var wg sync.WaitGroup
	results := make(chan string, 5)

	wg.Add(5)

	go func() {
		defer wg.Done()
		results <- getCPUPercent()
	}()

	go func() {
		defer wg.Done()
		results <- getLoadAvg()
	}()

	go func() {
		defer wg.Done()
		results <- getMem()
	}()

	go func() {
		defer wg.Done()
		results <- getDisk()
	}()

	go func() {
		defer wg.Done()
		results <- getTemps()
	}()

	wg.Wait()
	close(results)

	jsonString := "{"
	for res := range results {
		jsonString += fmt.Sprintf("%s,", res)
	}
	jsonString = jsonString[:len(jsonString)-1] + "}"

	return jsonString
}
