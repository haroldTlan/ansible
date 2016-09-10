// traserver
package main

import (
	"fmt"
	"regexp"
	"snmpserver/snmp"
	"strings"
)

type DiskInfo struct {
	Uuid      string
	Location  string
	MachineId string
	Status    string
	Role      string
	Raid      string
	Size      string
}

func NewDiskInfo(machineId string) *DiskInfo {
	return &DiskInfo{MachineId: machineId}
}

func RefreshDisks(ip string, machineId string) error {
	oid := "1.3.6.1.4.1.8888.1.1.0"
	out, err := snmp.Get(ip, oid)
	fmt.Println("out!!!!!!!!!:%v", out)
	if err != nil {
		return err
	}

	disks := extractDisks(out, machineId)

	for _, disk := range disks {
		fmt.Println(disk.Uuid, disk.Location, disk.MachineId, disk.Status, disk.Role, disk.Raid, disk.Size)
		UpdateDisk(disk.Uuid, disk.Location, disk.MachineId, disk.Status, disk.Role, disk.Raid, disk.Size)
	}

	return err
}

func extractDisks(out string, machineId string) []*DiskInfo {
	disks := make([]*DiskInfo, 0)

	disks_tmp := strings.Split(out, "[Disk_location]")

	for _, disk_tmp := range disks_tmp[1:] {
		disk := extractSingleDisk(disk_tmp, machineId)
		disks = append(disks, disk)
	}

	return disks
}

func extractSingleDisk(out string, machineId string) *DiskInfo {
	disk := NewDiskInfo(machineId)
	regLocation := regexp.MustCompile(`:(\d.\d.\d+)`)
	regUuid := regexp.MustCompile(`\[Disk_uuid\]:\s*(\S+)`)
	regStatus := regexp.MustCompile(`\[Disk_status\]:\s*(\S+)`)
	regRole := regexp.MustCompile(`\[Disk_role\]:\s*(\S+)`)
	regRaid := regexp.MustCompile(`\[Disk_raid\]:\s*(\S+)`)
	regSize := regexp.MustCompile(`\[Disk_size\]:\s*(\S+)`)
	disk.Location = regLocation.FindStringSubmatch(out)[1]
	disk.Uuid = regUuid.FindStringSubmatch(out)[1]
	disk.Status = regStatus.FindStringSubmatch(out)[1]
	disk.Role = regRole.FindStringSubmatch(out)[1]
	disk.Raid = regRaid.FindStringSubmatch(out)[1]
	disk.Size = regSize.FindStringSubmatch(out)[1]
	return disk
}
