package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"os/exec"
	"strings"
)

func refreshRozoCheck(checktype string, ip string) (interface{}, error) {
	var result []interface{}

	aim := strings.Split(checktype, "?")[0]
	option := strings.Split(checktype, "?")[1]
	fmt.Println(aim, option)

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("rozo %s %s -E %s", aim, option, ip))

	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	cmd.Stdout = w
	if err := cmd.Run(); strings.Contains(string(w.Bytes()), "FAILED") {
		return string(w.Bytes()), err
	}

	pure := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(w.Bytes()), &pure)

	if aim == "node" {
		if option == "list" {
			var status Node
			for k, v := range pure {
				status.Ip = k.(string)
				status.Status = assert(v)
				result = append(result, status)
			}
		} else if option == "status" {
			var status Node
			for k, v := range pure {
				status.Ip = k.(string)
				status.Status = assert(v)
				result = append(result, status)

			}
		} else if option == "config" {
		}
	}
	fmt.Printf("%+v", result)
	return result, nil
}

func assert(items interface{}) interface{} {
	var base Base
	if value, ok := items.(map[interface{}]interface{}); ok {
		var bases []Base
		//  fmt.Println(value)
		for v, k := range value {
			base.Name = v
			base.Status = assert(k)
			bases = append(bases, base)
		}
		return bases
	}

	if value, ok := items.([]interface{}); ok {
		var bases []interface{}
		for i := 0; i < len(value); i++ {
			bases = append(bases, assert(value[i]))
		}
		return bases
	}

	if value, ok := items.(string); ok {
		return value
	}
	return nil
}
