package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	//"strconv"
	//"strings"
	"time"
)

type Stats struct {
	Exports  []interface{} `json:"exports"`
	Storages []interface{} `json:"storages"`
}

type Master struct {
	Info []Info `json:"info"`
	Ip   string `json:"ip"`
}

/*type MasterView struct {
	Cpu       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Temp      float64 `json:"temp"`
	Write     float64 `json:"write_mb"`
	Read      float64 `json:"temp_mb"`
	TimeStamp float64 `json:"timestamp"`
	Tmp       `json:"cache"`
	Pro       []Process `json:"process"`
}*/

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
	Cpu       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Temp      float64 `json:"temp"`
	Write     float64 `json:"write_mb"`
	Read      float64 `json:"read_mb"`
	TimeStamp float64 `json:"timestamp"`
}

//func main() {
func InfoStat() {
	go func() {
		for {
			time.Sleep(4 * time.Second)
			var allInfo Stats
			//		allInfo := make(Stats.Exports, 0)
			var str string
			_, err := exec.Command("python", "temp.py").Output()
			if err == nil {
				str = read3("static")
			}
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

					//				allInfo = append(allInfo, success)
				}

			}
			fmt.Printf("%+v", allInfo)
			statTopic.Publish(allInfo)

		}
	}()

}

func transform(results map[string]interface{}) interface{} {

	if results["type"] == "masterInfo" {
		var masterAll Master

		masterAll.Ip = results["ip"].(string)
		result := results["result"].([]interface{}) //statistics
		masters := make([]Info, 1)

		if len(result) > 0 {
			for _, j := range result {

				vals := j.(map[string]interface{})
				item := map[string]float64{}
				var master Info

				for k, v := range vals {
					item[k] = v.(float64)
				}

				master.Cpu = item["cpu"]
				master.Mem = item["mem"]
				master.Temp = item["temp"]
				master.Write = item["nwrite_mb"]
				master.Read = item["nread_mb"]
				master.TimeStamp = item["timestamp"]
				masters = append(masters, master)
			}

		}
		masterAll.Info = masters[len(masters)-1 : len(masters)]

		return masterAll

	} else if results["type"] == "storeInfo" {
		var storeAll Store

		storeAll.Ip = results["ip"].(string)
		result := results["result"].([]interface{}) //statistics
		stores := make([]StoreView, 1)

		if len(result) > 0 {
			for _, j := range result {

				vals := j.(map[string]interface{})
				item := map[string]float64{}
				var store StoreView

				for k, v := range vals {
					item[k] = v.(float64)
				}

				store.Cpu = item["cpu"]
				store.Mem = item["mem"]
				store.Temp = item["temp"]
				store.Write = item["nwrite_mb"]
				store.Read = item["nread_mb"]
				store.TimeStamp = item["timestamp"]
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
