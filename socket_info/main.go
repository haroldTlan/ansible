package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	//"strings"
	"time"
)

type Config struct {
	Interval int
	Ansible  int
}

type Stats struct {
	Exports  []interface{} `json:"exports"`
	Storages []interface{} `json:"storages"`
}

type Master struct {
	Info []Info `json:"info"`
	Ip   string `json:"ip"`
}

type Process struct {
	Name string  `json:"protype"`
	Cpu  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
}

type Tmp struct {
	Total   float64 `json:"total"`
	Used    float64 `json:"used"`
	UsedPer float64 `json:"persentage"`
}

type Store struct {
	Dev []StoreView `json:"info"`
	Ip  string      `json:"ip"`
}

type Info struct {
	Cpu       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Temp      float64 `json:"temp"`
	Write     float64 `json:"write_mb"`
	Read      float64 `json:"read_mb"`
	TimeStamp float64 `json:"timestamp"`
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
	Name     string  `json:"name"`
	Total    float64 `json:"total"`
	Used     float64 `json:"used"`
	Used_per float64 `json:"used_per"`
}

//func main() {
func InfoStat() {
	go func() {
		for {
			var conf Config
			infoConf := read3("info.conf")
			yaml.Unmarshal([]byte(infoConf), &conf)
			time.Sleep(time.Duration(conf.Interval) * time.Second)

			var allInfo Stats
			str := read3("static")
			var dat []map[string]interface{}
			json.Unmarshal([]byte(str), &dat)
			for i := 0; i < len(dat); i++ {
				if dat[i]["status"] == "success" {
					success := transform(dat[i])
					if dat[i]["type"] == "masterInfo" {
						allInfo.Exports = append(allInfo.Exports, success)
					} else {
						allInfo.Storages = append(allInfo.Storages, success)
					}

					//			allInfo = append(allInfo, success)
				}

			}
			//		fmt.Printf("%+v\n", allInfo)

			statTopic.Publish(allInfo)

		}
	}()

}

func ansible() {
	go func() {
		for {
			var conf Config
			infoConf := read3("info.conf")
			yaml.Unmarshal([]byte(infoConf), &conf)
			time.Sleep(time.Duration(conf.Ansible) * time.Second)
			_, err := exec.Command("python", "device.py").Output()
			if err != nil {
				fmt.Println("\n\n\n\n", err)
			}
		}
	}()
}

func transform(results map[string]interface{}) interface{} {

	if results["type"] == "masterInfo" {
		var masterAll Store

		masterAll.Ip = results["ip"].(string)
		result := results["result"].([]interface{}) //statistics
		masters := make([]StoreView, 1)

		if len(result) > 0 {
			for _, j := range result {
				vals := j.(map[string]interface{})
				item_float64 := map[string]float64{}

				var master StoreView
				var item_array []Df
				for k, v := range vals {
					if _, ok := v.([]interface{}); ok {
						item := v.([]interface{})
						for i := 0; i < len(item); i++ {
							var d Df
							df := item[i].(map[string]interface{})
							str1 := strconv.FormatFloat(df["total"].(float64)/1024.0/1024.0, 'f', 2, 64)
							d.Total, _ = strconv.ParseFloat(str1, 64)
							str2 := strconv.FormatFloat(df["used"].(float64)/1024.0/1024.0, 'f', 2, 64)
							d.Used, _ = strconv.ParseFloat(str2, 64)
							d.Used_per = df["used_per"].(float64)
							d.Name = df["name"].(string)
							item_array = append(item_array, d)
						}
					} else if _, ok := v.(float64); ok {
						item_float64[k] = v.(float64)
					}
				}
				master.Dfs = item_array
				master.Cpu = item_float64["cpu"]
				master.Mem = item_float64["mem"]
				master.Temp = item_float64["temp"]
				master.Write = item_float64["nwrite_mb"]
				master.Read = item_float64["nread_mb"]
				master.TimeStamp = item_float64["timestamp"]
				masters = append(masters, master)
			}

		}
		masterAll.Dev = masters[len(masters)-1 : len(masters)]

		return masterAll

	} else if results["type"] == "storeInfo" {
		var storeAll Store

		storeAll.Ip = results["ip"].(string)
		result := results["result"].([]interface{}) //statistics
		stores := make([]StoreView, 1)

		if len(result) > 0 {
			for _, j := range result {

				vals := j.(map[string]interface{})
				item_float64 := map[string]float64{}

				var store StoreView
				var item_array []Df
				for k, v := range vals {
					if _, ok := v.([]interface{}); ok {
						item := v.([]interface{})
						for i := 0; i < len(item); i++ {
							var d Df
							df := item[i].(map[string]interface{})
							str1 := strconv.FormatFloat(df["total"].(float64)/1024.0/1024.0, 'f', 2, 64)
							d.Total, _ = strconv.ParseFloat(str1, 64)
							str2 := strconv.FormatFloat(df["used"].(float64)/1024.0/1024.0, 'f', 2, 64)
							d.Used, _ = strconv.ParseFloat(str2, 64)
							d.Used_per = df["used_per"].(float64)
							d.Name = df["name"].(string)
							item_array = append(item_array, d)
						}
					} else if _, ok := v.(float64); ok {
						item_float64[k] = v.(float64)
					}
				}
				store.Dfs = item_array
				store.Cpu = item_float64["cpu"]
				store.Mem = item_float64["mem"]
				store.Temp = item_float64["temp"]
				store.Write = item_float64["nwrite_mb"]
				store.Read = item_float64["nread_mb"]
				store.TimeStamp = item_float64["timestamp"]
				store.CacheT = item_float64["rdcache_t"]
				store.CacheU = item_float64["rdcache_u"]
				store.W_Vol = item_float64["rvol_mb"]
				store.R_Vol = item_float64["wvol_mb"]
				stores = append(stores, store)
			}
		}
		//	storeAll.Dev = stores
		storeAll.Dev = stores[len(stores)-1 : len(stores)]
		return storeAll
	}

	return nil

}

func getProp(d interface{}) (interface{}, bool) {
	fmt.Println(reflect.TypeOf(d))
	fmt.Println(reflect.ValueOf(d).Len())
	a := reflect.ValueOf(d)
	b := a.Index(1)

	c := reflect.ValueOf(b)
	fmt.Println(c)

	return nil, false
}
func read3(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}
