package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crackcomm/nsqueue/consumer"
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
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Body, &data); err != nil {
		AddLogtoChan(err)
		return
	}

	result := eventJugde(data)

	eventTopic.Publish(result)
	fmt.Printf("%+v\n", result)
	msg.Success()
}

func eventJugde(values map[string]interface{}) interface{} {
	switch values["event"].(string) {
	case "safety.created":
		return Safety{Event: values["event"].(string),
			Ip: values["ip"].(string)}

	case "info.warning", "info.normal":
		result := Warning{Event: values["event"].(string),
			Type:   values["type"].(string),
			Ip:     values["ip"].(string),
			Value:  values["value"].(float64),
			Status: values["status"].(string)}
		if result.Status == "true" {
			if err := RefreshInfoMail(result); err != nil {
				AddLogtoChan(err)
			}
		}
		return result

	default:
		machineId, err := analyze(values["ip"].(string))
		if err != nil {
			fmt.Println(err)
			return err
		}
		result := newEvent(values, machineId)
		if value, ok := result.(error); ok {
			AddLogtoChan(value)
			return nil
		}
		if err := RefreshOverViews(values["ip"].(string), values["event"].(string)); err != nil {
			AddLogtoChan(err)
		}
		return result
	}
}

func newEvent(values map[string]interface{}, machineId string) interface{} {
	switch values["event"].(string) {
	case "ping.offline", "ping.online":
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

	case "machine.created":
		InitSingleRemote(values["ip"].(string))
		return HeartBeat{Event: values["event"].(string),
			Ip:        values["ip"].(string),
			MachineId: machineId}

	}
	return nil
}

func analyze(machine string) (string, error) {
	if num, one, err := SelectMachine(machine); err == nil && num > 0 {
		return one.Uuid, nil
	}

	return "", errors.New("Machine is not being monitored")
}
