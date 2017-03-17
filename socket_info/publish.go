package main

import (
	"encoding/json"
	"fmt"
	"github.com/crackcomm/nsqueue/producer"
)

var (
	amount   = 30
	nsqdAddr = "127.0.0.1:4150"
	topics   = "CloudEvent"
)

type Ping struct {
	Event  string  `json:"event"`
	Type   string  `json:"type"`
	Ip     string  `json:"ip"`
	Value  float64 `json:"value"`
	Status string  `json:"status"`
}

func NsqInit() {
	go func() {
		producer.Connect(nsqdAddr)
	}()
}

func NsqRequest(event, sendtype, ip, status string, value float64) {
	buffer := eventType(event, sendtype, ip, status, value)
	fmt.Println(event, ip, status, topics, value)
	fmt.Printf("\n\n\n%+v\n\n\n", string(buffer))
	producer.PublishAsync(topics, buffer, nil)
}

func eventType(event, sendtype, ip, status string, value float64) []byte {
	var ping Ping

	ping.Event = event
	ping.Type = sendtype
	ping.Ip = ip
	ping.Value = value
	ping.Status = status
	stringSlice, _ := json.Marshal(ping)
	buffer := []byte(stringSlice)
	return buffer
}
