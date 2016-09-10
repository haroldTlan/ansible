// trapserver
package main

import (
	"fmt"
	"net"
	"snmpserver/snmp"
	"snmpserver/topic"
)

var trapTopic = topic.New()

type DiskEvent struct {
	Name                 string
	Uuid                 string
	Location             string
	MachineId            string
	RawReadErrorRate     string
	SpinUpTime           string
	StartStopCount       string
	ReallocatedSectorCt  string
	SeekErrorRate        string
	PowerOnHours         string
	SpinRetryCount       string
	PowerCycleCount      string
	PowerOffRetractCount string
	LoadCycleCount       string
	CurrentPendingSector string
	OfflineUncorrectable string
	UDMACRCErrorCount    string
	Status               string
	Role                 string
	Raid                 string
	Size                 string
}

const (
	DISKEVENTTYPE        = ".1.3.6.1.4.1.8888"
	EVENT                = ".1.3.6.1.4.1.8888.1.1"
	UUID                 = ".1.3.6.1.4.1.8888.1.3"
	LOCATION             = ".1.3.6.1.4.1.8888.1.2"
	MACHINEID            = ".1.3.6.1.4.1.8888.1.4"
	RAWREADERRORRATE     = ".1.3.6.1.4.1.8888.1.5.1"
	SpinUpTime           = ".1.3.6.1.4.1.8888.1.5.2"
	StartStopCount       = ".1.3.6.1.4.1.8888.1.5.3"
	ReallocatedSectorCt  = ".1.3.6.1.4.1.8888.1.5.4"
	SeekErrorRate        = ".1.3.6.1.4.1.8888.1.5.5"
	PowerOnHours         = ".1.3.6.1.4.1.8888.1.5.6"
	SpinRetryCount       = ".1.3.6.1.4.1.8888.1.5.7"
	PowerCycleCount      = ".1.3.6.1.4.1.8888.1.5.8"
	PowerOffRetractCount = ".1.3.6.1.4.1.8888.1.5.9"
	LoadCycleCount       = ".1.3.6.1.4.1.8888.1.5.10"
	CurrentPendingSector = ".1.3.6.1.4.1.8888.1.5.11"
	OfflineUncorrectable = ".1.3.6.1.4.1.8888.1.5.12"
	UDMACRCErrorCount    = ".1.3.6.1.4.1.8888.1.5.13"
)

func newDiskEvent(values map[string]interface{}) DiskEvent {
	return DiskEvent{Name: values[EVENT].(string),
		Uuid:                 values[UUID].(string),
		Location:             values[LOCATION].(string),
		MachineId:            values[MACHINEID].(string),
		RawReadErrorRate:     values[RAWREADERRORRATE].(string),
		SpinUpTime:           values[SpinUpTime].(string),
		StartStopCount:       values[StartStopCount].(string),
		ReallocatedSectorCt:  values[ReallocatedSectorCt].(string),
		SeekErrorRate:        values[SeekErrorRate].(string),
		PowerOnHours:         values[PowerOnHours].(string),
		SpinRetryCount:       values[SpinRetryCount].(string),
		PowerCycleCount:      values[PowerCycleCount].(string),
		PowerOffRetractCount: values[PowerOffRetractCount].(string),
		LoadCycleCount:       values[LoadCycleCount].(string),
		CurrentPendingSector: values[CurrentPendingSector].(string),
		OfflineUncorrectable: values[OfflineUncorrectable].(string),
		UDMACRCErrorCount:    values[UDMACRCErrorCount].(string)}
}

func TrapServer() {
	go func() {
		fmt.Println("hello world!")
		socket, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 162})
		if err != nil {
			panic(err)
		}
		defer socket.Close()

		for {
			buf := make([]byte, 2048)
			read, from, _ := socket.ReadFromUDP(buf)
			fmt.Println("Get msg from ", from.IP)
			HandleUdp(buf[:read])
		}
	}()
}

func HandleUdp(data []byte) {
	trap, err := snmp.ParseUdp(data)
	if err != nil {
		fmt.Println("Err", err.Error())
	}
	fmt.Println(trap.Version, trap.Community, trap.EnterpriseId, trap.Address)
	for k, v := range trap.Values {
		fmt.Printf("%s = %s\n", k, v)
	}

	var event DiskEvent
	if trap.EnterpriseId == DISKEVENTTYPE {
		event = newDiskEvent(trap.Values)
		fmt.Println(event)
	}

	trapTopic.Publish(event)
	//fmt.Printf("From trapserver.go:%s\n",event)
}
