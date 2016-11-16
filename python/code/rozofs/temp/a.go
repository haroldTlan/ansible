package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	//	"encoding/json"
	"os/exec"
	"strings"
)

type NodeStatus struct {
	Ip   string        `yaml:"ip"`
	Type []interface{} `yaml:"type"`
}

type NodeList struct {
	Ip     string
	Status []Base
}

type Base struct {
	Name   interface{} `yaml:"name"`
	Status interface{} `yaml:"status"`
}

func main() {
	//var status NodeStatus
	var status NodeList
	var result []interface{}
	ip := "192.168.2.186"
	nodeType := "status"
	//nodeType := "list"
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("rozo node %s -E %s", nodeType, ip))

	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		fmt.Printf("Run returns: %s\n", err)
	}

	fmt.Println(strings.Contains(string(w.Bytes()), "FAILED"))

	//out, _ := exec.Command("/bin/sh", "-c", "ls").Output()
	fmt.Println(string(w.Bytes()))
	dat := make(map[interface{}]interface{})
	//var dat map[string]interface{}
	yaml.Unmarshal([]byte(w.Bytes()), &dat)

	fmt.Println(dat)
	fmt.Printf("%+v", temp(dat))
	fmt.Println("\n\n\n\n")
	fmt.Println(result, status)

}

func temp(items interface{}) interface{} {
	var base Base
	if value, ok := items.(map[interface{}]interface{}); ok {
		var bases []Base
		fmt.Println(value)
		for v, k := range value {
			base.Name = v
			base.Status = temp(k)
			bases = append(bases, base)
		}
		return bases
	}

	if value, ok := items.([]interface{}); ok {
		var bases []interface{}
		for i := 0; i < len(value); i++ {
			bases = append(bases, temp(value[i]))
		}
		return bases
	}

	if value, ok := items.(string); ok {
		return value
	}

	return nil
}

//fmt.Println(dat["192.168.2.190"].([]map[string]interface{}))

/*for k, v := range a {
status.Ip = k.(string)
var base Base
types := v.([]interface{})
for i := 0; i < len(types); i++ {
	a := types[i].(map[interface{}]interface{})
	for key, val := range a {
		fmt.Println(key, "a", val)
		base.Name = key.(string)
		base.Status = methodSwitch(val)*/
/*			switch val.(type) {
			case string:
				base.Status = val.(string)
			case []map[string]string:
				fmt.Println("!!!!!!!!!!!!!!!")
				a := val.([]interface{})[0]
				fmt.Println(a)
			default:
				fmt.Println("!!!!!!!!!!!!!!!")
				w := val.([]interface{})[0]
				fmt.Println(w.(map[interface{}]interface{}))
				//				for k, v := range w {
				//					fmt.Println(k.(string), v.(string))
				//				}
				//				fmt.Println(a)

			}*/

/*		}
		status.Status = append(status.Status, base)
	}
	//fmt.Println(types[0])
	//status.Type = types
	result = append(result, status)
}
fmt.Printf("%+v", result)

fmt.Println("\n\n")*/

func methodSwitch(val interface{}) interface{} {
	switch val.(type) {
	case string:
		return val.(string)
	case map[string]string:
		fmt.Println("asdasd")
		var base Base
		for key, val := range val.(map[string]string) {
			base.Name = key
			base.Status = val
		}
		return base
	case map[interface{}]interface{}:
		fmt.Println("asdasd\n\n\n")
		//fmt.Println(val.(map[string]string))
		var base Base
		/*for k, v := range val {
			base.Name = methodSwitch(key)
			base.Status = methodSwitch(val)
		}*/
		return base

	case []interface{}:
		var arrs []interface{}
		arr := val.([]interface{})
		for i := 0; i < len(arr); i++ {
			arrs = append(arrs, methodSwitch(arr[i]))
			/*		a := arr[i].(map[interface{}]interface{})
					fmt.Println(a)
					for k, v := range a {
						fmt.Println(k.(string))
						fmt.Println(v.(string))
					}*/

		}

		return arrs
	case []map[string]string:
		var base Base
		for key, val := range val.([]map[string]string)[0] {
			base.Name = key
			base.Status = val
		}
		fmt.Println("map")
		return base
	default:
		fmt.Println("zxc")

		return val
	}
}
