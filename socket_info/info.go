package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Config struct {
	Interval int
	Ansible  int
}

type Stats struct {
	Exports  []Device `json:"exports"`
	Storages []Device `json:"storages"`
}

type Device struct {
	Dev []StoreView `json:"info"`
	Ip  string      `json:"ip"`
}

type StoreView struct {
	Dfs       []Df    `json:"df"`
	Cpu       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Temp      float64 `json:"temp"`
	Write     float64 `json:"write_mb"`
	Read      float64 `json:"read_mb"`
	TimeStamp float64 `json:"timestamp"`
	CacheT    float64 `json:"cache_total"`
	CacheU    float64 `json:"cache_used"`
	W_Vol     float64 `json:"write_vol"`
	R_Vol     float64 `json:"read_vol"`
}

type Df struct {
	Name      string  `json:"name"`
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Used_per  float64 `json:"used_per"`
}

type Log struct {
	Message    string `json:"message"`
	Created_at int64  `json:"created_at"`
	Level      string `json:"level"`
	Source     string `json:"scource"`
}

func InfoStat() {
	go func() {
		for {
			var conf Config
			infoConf := readConf("/etc/ansible/info/info.conf")
			if err := yaml.Unmarshal([]byte(infoConf), &conf); err != nil {
				AddLogtoChan("infostat yaml info.conf", err)
			}
			time.Sleep(time.Duration(conf.Interval) * time.Second)

			var allInfo Stats
			allInfo.Exports = make([]Device, 0)
			allInfo.Storages = make([]Device, 0)

			str := readConf("/etc/ansible/info/static")
			var dat []map[string]interface{}
			if err := json.Unmarshal([]byte(str), &dat); err != nil {
				AddLogtoChan("json static", err)
			}

			for i := 0; i < len(dat); i++ {
				if dat[i]["status"] == "success" {
					success := transform(dat[i])
					if dat[i]["type"] == "masterInfo" {
						allInfo.Exports = append(allInfo.Exports, success)
					} else {
						allInfo.Storages = append(allInfo.Storages, success)
					}

				}

			}

			statTopic.Publish(allInfo)

		}
	}()

}

func dfAssert(items []interface{}) []Df {
	var dfs []Df
	for i := 0; i < len(items); i++ {
		var d Df

		df := items[i].(map[string]interface{})
		str1 := strconv.FormatFloat(df["total"].(float64)/1024.0/1024.0, 'f', 2, 64)
		d.Total, _ = strconv.ParseFloat(str1, 64)
		str2 := strconv.FormatFloat(df["used"].(float64)/1024.0/1024.0, 'f', 2, 64)
		d.Available, _ = strconv.ParseFloat(str2, 64)
		d.Used_per = df["used_per"].(float64)
		d.Name = df["name"].(string)
		dfs = append(dfs, d)
	}
	return dfs
}

func transform(items map[string]interface{}) Device {
	var masterAll Device

	masterAll.Ip = items["ip"].(string)
	infos := make([]StoreView, 0)
	results := items["result"].([]interface{})

	if len(results) > 0 {
		var info StoreView

		for _, j := range results {
			item_dfs := make([]Df, 0)
			vals := j.(map[string]interface{})
			item_float64 := map[string]float64{}

			for k, v := range vals {
				if _, ok := v.([]interface{}); ok {
					item_dfs = dfAssert(v.([]interface{}))
				} else if _, ok := v.(float64); ok {
					item_float64[k] = v.(float64)
				}

			}

			info.Dfs = item_dfs
			info.Cpu = item_float64["cpu"]
			info.Mem = item_float64["mem"]
			info.Temp = item_float64["temp"]
			info.Write = item_float64["nwrite_mb"]
			info.Read = item_float64["nread_mb"]
			info.TimeStamp = item_float64["timestamp"]
			infos = append(infos, info)
		}
	}
	masterAll.Dev = infos[len(infos)-1 : len(infos)]
	return masterAll
}

func ansible() {
	go func() {
		for {
			var conf Config
			infoConf := readConf("/etc/ansible/info/info.conf")
			if err := yaml.Unmarshal([]byte(infoConf), &conf); err != nil {
				AddLogtoChan("ansible yaml info.conf", err)
			}
			time.Sleep(time.Duration(conf.Ansible) * time.Second)
			if _, err := exec.Command("python", "/etc/ansible/info/device.py").Output(); err != nil {
				AddLogtoChan("ansible", err)
			}
		}
	}()
}

func AddLogtoChan(apiName string, err error) {
	var message string
	var log Log
	if err == nil {
		message = fmt.Sprintf("[STATIS]statistics success")
		log = Log{Level: "INFO", Message: message}
	} else {
		pc, fn, line, _ := runtime.Caller(1)
		message = fmt.Sprintf("[STATIS][%s %s:%d] %s, %s", runtime.FuncForPC(pc).Name(), fn, line, apiName, err)
		log = Log{Level: "ERROR", Message: message}
	}

	ChanLogInfo <- log
	return
}

func readConf(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		AddLogtoChan("Open", err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		AddLogtoChan("Read", err)
	}
	return string(fd)
}
