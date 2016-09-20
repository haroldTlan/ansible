package main

import (
	"fmt"
	//"hwraid/topic"
	"os/exec"
	"strings"

	"strconv"
	"time"
)

//var statTopic = topic.New()

type SysInfo struct {
	//Process ProcessInfo
	NetFlow Flow
	//Cache   Tmp
}

type Flow struct {
	FlowType string  `json:"flowtype"`
	Sent     float64 `json:"sent"`
	Receive  float64 `json:"receive"`
}
type ProcessInfo struct {
	ProType string  `json:"protype"`
	Cpu     float64 `json:"cpu"`
	Mem     float64 `json:"men"`
}
type Tmp struct {
	Total   float64 `json:"total"`
	Used    float64 `json:"used"`
	UsedPer float64 `json:"persentage"`
}

func InfoStat() {
	var sys SysInfo

	go func() {
		for {
			var pub []interface{}
			time.Sleep(2 * time.Second)
			out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("python /root/code/temp.py")).Output()
			if err != nil {
				fmt.Println(err)
			}

			//stat_process, _ := processMode(string(out))
			stat_flow, _ := flowMode(string(out))
			//stat_cache, _ := tmpMode(string(out))

			sys.NetFlow = stat_flow
			//sys.Process = stat_process
			//sys.Cache = stat_cache

			pub = append(pub, sys)

			statTopic.Publish(pub)
		}
	}()
}

//func processMode(out string) (interface{}, error) {
func processMode(out string) (ProcessInfo, error) {
	var stat ProcessInfo
	//var process []interface{}

	bareinfo := strings.Split(string(out), "?")[2]
	infos := strings.Split((bareinfo), "[")

	for _, args := range infos[1:] {
		arg := strings.Split(args, "\n\t")
		stat.ProType = arg[0]

		cpu, err := strconv.ParseFloat(arg[1], 64)
		if err != nil {
			fmt.Println(arg[1], err)
			cpu = 0
		}
		arg[2] = strings.Replace(arg[2], "\n", "", -1)
		mem, err := strconv.ParseFloat(arg[2], 64)
		if err != nil {
			fmt.Println(arg[2], err)
			mem = 0
		}

		stat.Cpu = cpu
		stat.Mem = mem

		//process = append(process, stat)
	}
	return stat, nil
}

func flowMode(out string) (Flow, error) {
	var stat Flow
	//var flow []interface{}

	bareinfo := strings.Split(string(out), "?")[3]

	args := strings.Split(bareinfo, "\n\t")

	sent, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println(args[1], err)
		sent = 0
	}

	args[2] = strings.Replace(args[2], "\n", "", -1)
	rec, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		fmt.Println(args[2], err)
		rec = 0
	}

	stat.Sent = sent
	stat.Receive = rec

	//flow = append(flow, stat)

	return stat, nil
}

func tmpMode(out string) (Tmp, error) {
	var stat Tmp
	//var tmp []interface{}

	bareinfo := strings.Split(string(out), "?")[4]

	args := strings.Split(bareinfo, "\n\t")

	total, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println(args[1], err)
		total = 0
	}

	used, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		fmt.Println(args[2], err)
		used = 0
	}

	args[3] = strings.Replace(args[3], "\n", "", -1)
	per, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		fmt.Println(args[3], err)
		per = 0
	}

	stat.Total = total
	stat.Used = used
	stat.UsedPer = per

	//tmp = append(tmp, stat)

	return stat, nil
}
