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

type LogInfo struct {
	Level     string
	LogType   string
	ChLogType string
	Result    string
	ChResult  string
}

func getLogConfig(logtype string, result bool) LogInfo {
	var config LogInfo

	switch result {
	case true:
		config.Result = "successfully"
		config.ChResult = "成功"

	case false:
		config.Result = "unsuccessfully"
		config.ChResult = "失败"

	default:
		config.Result = "bullshit"
		config.ChResult = "不可能，要是出现我吃翔三斤"

	}

	switch logtype {
	case "set":
		config.Level = "info"
		config.LogType = logtype
		config.ChLogType = "配置"
		return config
	case "check":
		config.Level = "info"
		config.LogType = logtype
		config.ChLogType = "检查"
		return config
	case "unset":
		config.Level = "warning"
		config.LogType = logtype
		config.ChLogType = "解除"
		return config
	case "error":
		config.Level = "error"
		config.LogType = logtype
		config.ChLogType = "错误"
		return config
	case "over":
		config.Level = "warning"
		config.LogType = logtype
		config.ChLogType = "警告"
		return config
	default:
		return config

	}

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
		//UpdateDisk(disk.Uuid, disk.Location, disk.MachineId, disk.Status, disk.Role, disk.Raid, disk.Size)
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

func RefreshViews(uuid string) error {
	fmt.Println("ok")
	if err := RefreshReDisks(uuid); err != nil {
		return err
	}
	if err := RefreshReRaids(uuid); err != nil {
		return err
	}
	if err := RefreshReVolumes(uuid); err != nil {
		return err
	}
	if err := RefreshReFilesystems(uuid); err != nil {
		return err
	}
	if err := RefreshReInitiators(uuid); err != nil {
		return err
	}

	if err := RefreshReRaidVolumes(uuid); err != nil {
		return err
	}
	if err := RefreshReInitVolumes(uuid); err != nil {
		return err
	}
	if err := RefreshReNetInits(uuid); err != nil {
		return err
	}

	return nil
}

func restApi(uuid string) (StoreView, error) {
	var store StoreView

	disks, err := resDisks(uuid)
	if err != nil {
		return store, err
	}

	raids, err := resRaids(uuid)
	if err != nil {
		return store, err
	}

	vols, err := resVols(uuid)
	if err != nil {
		return store, err
	}

	fs, err := resFs(uuid)
	if err != nil {
		return store, err
	}

	inits, err := resInits(uuid)
	if err != nil {
		return store, err
	}

	store.RestDisks = disks
	store.RestRaids = raids
	store.RestVolumes = vols
	store.RestFs = fs
	store.RestInits = inits

	return store, nil
}

func RefreshAuto() {
	go func() {
		machines, _ := SelectAllMachines()

		if len(machines) > 0 {
			for _, val := range machines {
				if err := RefreshViews(val.Uuid); err != nil {
					fmt.Println(err)
				}
			}
		} else {
		}
	}()

}
