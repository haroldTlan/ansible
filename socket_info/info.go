package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type SysInfo struct {
	Masters Master  `json:"master"`
	Stores  []Store `json:"store"`
}

type Master struct {
	Process []ProcessInfo `json:"process"`
	NetFlow Flow          `json:"netflow"`
	Cache   Tmp           `json:"cache"`
}

type Store struct {
	Flow        //string  `json:"flow"`
	Ip   string `json:"ip"`
}

type Flow struct {
	FlowType string  `json:"flowtype"`
	Sent     float64 `json:"sent"`
	Receive  float64 `json:"receive"`
}

type Tmp struct {
	Total   float64 `json:"total"`
	Used    float64 `json:"used"`
	UsedPer float64 `json:"persentage"`
}

type ProcessInfo struct {
	ProType string  `json:"protype"`
	Cpu     float64 `json:"cpu"`
	Mem     float64 `json:"mem"`
}

func InfoStat() {
	//func main() {
	go func() {
		for {
			time.Sleep(2 * time.Second)

			cmd, err := exec.Command("python", "temp.py").Output()
			if err != nil {
				fmt.Println(err)
			}

			allInfo := infoTrans(string(cmd))
			statTopic.Publish(allInfo)
			fmt.Printf("%+v", allInfo)
		}
	}()
	/*
		cmd, err := exec.Command("python", "temp.py").Output()
		if err != nil {
			fmt.Println(err)
		}

		a := infoTrans(string(cmd))
		fmt.Printf("%+v", a)
		time.Sleep(2 * time.Second)*/

}

func infoTrans(items string) SysInfo {
	var cols []map[string]string

	items = strings.Split(items, "[{")[1]
	items = strings.Split(items, "}]")[0]
	clearQuo := strings.Replace(items, "'", "", -1)
	infoList := strings.Split(clearQuo, "}, {")

	for _, val := range infoList {

		slices := strings.Split(val, ", ")
		item := map[string]string{}

		for _, key := range slices {
			slice := strings.Split(key, ": ")
			/*cols = append(cols, map[string]string{
				slice[0]: slice[1],
			})*/
			item[slice[0]] = slice[1]

		}
		cols = append(cols, item)
		/*a := getInfo(item)
		fmt.Printf("%+v", a)*/

	}
	all := getInfo(cols)

	return all
}

func getInfo(items []map[string]string) SysInfo {
	var store Store
	var master Master
	var sys SysInfo
	for _, val := range items {

		if val["status"] == "success" {
			if val["type"] == "masterInfo" {

				stat_flow, _ := flowMode(val["result"], 1)
				stat_process, _ := processMode(val["result"], 0)
				stat_cache, _ := tmpMode(val["result"], 2)

				master.NetFlow = stat_flow
				master.Process = stat_process
				master.Cache = stat_cache

			} else if val["type"] == "storeInfo" {
				stat_flow, _ := flowMode(val["result"], 1)
				store.Flow = stat_flow
				store.Ip = val["ip"]
				//store.Flow = append(store.Flow, stat_flow)
			}

		}
		sys.Stores = append(sys.Stores, store)
	}
	sys.Masters = master

	return sys

}

func flowMode(out string, infoType int) (Flow, error) {
	/*var bareinfo string
	if infoType == "masterInfo" {
		bareinfo = strings.Split(out, "?")[3]
	} else if infoType == "storeInfo" {
		fmt.Println(out)
		bareinfo = strings.Split(string(out), "?")[1]
		fmt.Println(bareinfo)
	}*/
	var stat Flow
	//var flow []interface{}
	bareinfo := strings.Split(string(out), "?")[infoType]

	args := strings.Split(bareinfo, "\\n\\t")

	sent, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println(args[1], err)
		sent = 0
	}

	args[2] = strings.Replace(args[2], "\\n", "", -1)
	rec, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		fmt.Println(args[2], err)
		rec = 0
	}

	stat.Sent = sent
	stat.Receive = rec

	return stat, nil
}

func tmpMode(out string, infoType int) (Tmp, error) {
	var stat Tmp

	bareinfo := strings.Split(string(out), "?")[infoType]

	args := strings.Split(bareinfo, "\\n\\t")

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

	args[3] = strings.Replace(args[3], "\\n", "", -1)
	per, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		fmt.Println(args[3], err)
		per = 0
	}

	stat.Total = total
	stat.Used = used
	stat.UsedPer = per

	return stat, nil
}

func processMode(out string, infoType int) ([]ProcessInfo, error) {
	var stat ProcessInfo
	var process []ProcessInfo

	bareinfo := strings.Split(string(out), "?")[infoType]
	infos := strings.Split((bareinfo), "[")

	for _, args := range infos[1:] {
		arg := strings.Split(args, "\\n\\t")
		stat.ProType = arg[0]

		cpu, err := strconv.ParseFloat(arg[1], 64)
		if err != nil {
			fmt.Println(arg[1], err)
			cpu = 0
		}
		arg[2] = strings.Replace(arg[2], "\\n", "", -1)
		mem, err := strconv.ParseFloat(arg[2], 64)
		if err != nil {
			fmt.Println(arg[2], err)
			mem = 0
		}

		stat.Cpu = cpu
		stat.Mem = mem

		process = append(process, stat)
	}

	return process, nil
}
