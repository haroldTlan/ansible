package main

import (
	"encoding/json"
	"fmt"
	"github.com/crackcomm/nsqueue/consumer"
	"hwraid/topic"
	//"io/ioutil"
	"errors"
	"os"
)

var (
	eventTopic  = topic.New()
	nsqdAddr    = "192.168.2.83:4150"
	maxInFlight = 30
)

func nsqConsumerInit() {
	go func() {
		consumer.Register("CloudEvent", "consume82", maxInFlight, handle)
		consumer.Connect(nsqdAddr)
		consumer.Start(true)
	}()
}

func handle(msg *consumer.Message) {
	var dat map[string]interface{}
	if err := json.Unmarshal(msg.Body, &dat); err != nil {
		WriteConf("unknown.conf", string(msg.Body))
		return
	}

	result := newEvent(dat)
	if result == nil {
		WriteConf("unknown.conf", string(msg.Body))
		return
	}

	eventTopic.Publish(result)
	fmt.Printf("%+v\n", result)
	msg.Success()
}

func newEvent(values map[string]interface{}) interface{} {
	if err := analyze(values["ip"].(string)); err != nil {
		return nil
	}

	switch values["event"].(string) {

	case "ping.offline", "ping.online":
		InsertJournals(values["event"].(string), values["ip"].(string))
		return HeartBeat{Event: values["event"].(string),
			Ip:     values["ip"].(string),
			Status: values["status"].(string)}

	case "disk.unplugged":
		InsertJournals(values["event"].(string), values["ip"].(string))
		return DiskUnplugged{Event: values["event"].(string),
			Uuid:     values["uuid"].(string),
			Location: values["location"].(string),
			DevName:  values["dev_name"].(string),
			Ip:       values["ip"].(string)}

	case "disk.plugged", "raid.created", "volume.created", "volume.removed", "raid.degraded", "raid.failed", "volume.failed":
		InsertJournals(values["event"].(string), values["ip"].(string))
		return DiskPlugged{Event: values["event"].(string),
			Uuid: values["uuid"].(string),
			Ip:   values["ip"].(string)}

	case "raid.removed":
		InsertJournals(values["event"].(string), values["ip"].(string))
		disks := values["raid_disks"].([]interface{})
		var ones []string
		for _, val := range disks {
			disk := val.(string)
			ones = append(ones, disk)
		}
		return RaidRemove{Event: values["event"].(string),
			Uuid:      values["uuid"].(string),
			RaidDisks: ones,
			Ip:        values["ip"].(string)}
	}
	return nil
}

func analyze(machine string) error {
	var one Machine
	num, err := o.QueryTable("machine").Filter("ip", machine).All(&one)
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("Machine is not being monitored")
	}
	return nil
}

func WriteConf(path string, str string) {
	fi, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	final := []byte(str + "\n")
	if fi.Write(final); err != nil {
		panic(err)
	}

}
