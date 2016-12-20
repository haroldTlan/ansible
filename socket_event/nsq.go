package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crackcomm/nsqueue/consumer"
	"os"
	"time"
)

var (
	maxInFlight = 30
	nsqdAddr    = "127.0.0.1:4150"
)

func NsqConsumerInit() {
	go func() {
		consumer.Register("CloudEvent", "consume82", maxInFlight, handle)
		consumer.Connect(nsqdAddr)
		consumer.Start(true)
	}()
}

func handle(msg *consumer.Message) {
	var dat map[string]interface{}
	if err := json.Unmarshal(msg.Body, &dat); err != nil {
		message := fmt.Sprintf("nsq json " + string(msg.Body))
		AddLogtoChan(message, err)
		WriteConf("unknown.conf", "\nerror(json):"+string(msg.Body))
		return
	}

	result := newEvent(dat)
	if value, ok := result.(error); ok {
		message := fmt.Sprintf("nsq newEvent " + string(msg.Body))
		AddLogtoChan(message, value)
		WriteConf("unknown.conf", "\nerror(event):"+string(msg.Body))
		return
	}

	if err := refreshOverViews(dat["ip"].(string), dat["event"].(string)); err != nil {
		AddLogtoChan("nsq refreshOver ", err)
	}

	eventTopic.Publish(result)
	fmt.Printf("%+v\n", result)
	msg.Success()
}

func newEvent(values map[string]interface{}) interface{} {
	machineId, err := analyze(values["ip"].(string))
	if err != nil {
		return err
	}

	switch values["event"].(string) {
	case "ping.offline", "ping.online", "databox.created", "databox.removed":
		return HeartBeat{Event: values["event"].(string),
			Ip:        values["ip"].(string),
			MachineId: machineId}

	case "fs.removed", "fs.created":
		return FsSystem{Event: values["event"].(string),
			Volume:    values["volume"].(string),
			Type:      values["type"].(string),
			MachineId: machineId,
			Ip:        values["ip"].(string)}

	case "disk.unplugged":
		return DiskUnplugged{Event: values["event"].(string),
			Uuid:      values["uuid"].(string),
			Location:  values["location"].(string),
			DevName:   values["dev_name"].(string),
			MachineId: machineId,
			Ip:        values["ip"].(string)}

	case "disk.plugged", "raid.created", "volume.created", "volume.removed", "raid.degraded", "raid.failed", "volume.failed", "volume.normal", "raid.normal":
		return DiskPlugged{Event: values["event"].(string),
			Uuid:      values["uuid"].(string),
			MachineId: machineId,
			Ip:        values["ip"].(string)}

	case "raid.removed":
		disks := values["raid_disks"].([]interface{})
		var ones []string
		for _, val := range disks {
			disk := val.(string)
			ones = append(ones, disk)
		}
		return RaidRemove{Event: values["event"].(string),
			Uuid:      values["uuid"].(string),
			RaidDisks: ones,
			MachineId: machineId,
			Ip:        values["ip"].(string)}
	}
	return nil
}

func refreshOverViews(ip, event string) error {
	InsertJournals(event, ip)
	_, one, err := SelectMachine(ip)
	if err != nil {
		return err
	}

	if event == "ping.offline" {
		DelOutlineMachine(one.Uuid)

	} else if event == "ping.online" {
		time.Sleep(15 * time.Second)
		if err := RefreshStores(one.Uuid); err != nil {
			return err
		}

	} else {
		time.Sleep(15 * time.Second)
		if err := RefreshStores(one.Uuid); err != nil {
			return err
		}
		time.Sleep(4 * time.Second)
	}
	return nil
}

func analyze(machine string) (string, error) {
	if num, one, err := SelectMachine(machine); err == nil && num > 0 {
		return one.Uuid, nil
	}

	return "", errors.New("Machine is not being monitored")
}

func WriteConf(path string, str string) {
	fi, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		AddLogtoChan("nsq read ", err)
	}
	defer fi.Close()

	final := []byte(str + "\n")
	if fi.Write(final); err != nil {
		AddLogtoChan("nsq write ", err)
	}
}
