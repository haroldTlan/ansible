// traserver
package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"snmpserver/snmp"
	"strconv"
	"strings"
	"time"
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

type Base struct {
	Name   interface{} `json:"name"`
	Status interface{} `json:"status"`
}

type NodeStatus struct {
	Ip   string        `json:"ip"`
	Type []interface{} `json:"type"`
}

type Node struct {
	Ip     string      `json:"ip"`
	Status interface{} `json:"status"`
}

type StatInfo struct {
	Exports  string `json:"exports"`
	Storages string `json:"storages"`
}

type View struct {
	Disk       []Disk
	NumOfDisks int64
	Raid       []Raid
	NumOfRaids int64

	Vol       []Volume
	NumOfVols int64

	Fs      []Filesystems
	NumOfFs int64

	Initiator       []Initiator
	NumOfInitiators int64

	Jours      []ResJournals `json:"journals"`
	NumOfJours int64
}

type StoreView struct {
	RestDisks   []ResDisks       `json:"disks"`
	RestRaids   []ResRaids       `json:"raids"`
	RestVolumes []ResVols        `json:"volumes"`
	RestFs      []ResFilesystems `json:"filesystems"`
	RestInits   []ResInitiators  `json:"initiators"`
	RestJours   []ResJournals    `json:"journals"`
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

func RefreshAllViews() (View, error) {
	var views View
	disks, disks_num, err := SelectDisks()
	if err != nil {
		return views, err
	}
	raids, raids_num, err := SelectRaids()
	if err != nil {
		return views, err
	}
	vols, vols_num, err := SelectVolumes()
	if err != nil {
		return views, err
	}
	fs, fs_num, err := SelectFilesystems()
	if err != nil {
		return views, err
	}
	inits, inits_num, err := SelectInitiators()
	if err != nil {
		return views, err
	}
	jours, jours_num, err := SelectJournals()
	if err != nil {
		return views, err
	}
	views = View{Disk: disks, NumOfDisks: disks_num, Raid: raids, NumOfRaids: raids_num, Vol: vols, NumOfVols: vols_num, Fs: fs, NumOfFs: fs_num, Initiator: inits, NumOfInitiators: inits_num, Jours: jours, NumOfJours: jours_num}
	//views := View{NumOfDisks: disks_num, NumOfRaids: raids_num, NumOfVols: vols_num, NumOfFs: fs_num, NumOfInitiators: inits_num}

	return views, nil
}

func Refreshing() {
	go func() {
		time.Sleep(4 * time.Second)
		ones := make([]Machine, 0)
		if _, err := o.QueryTable("machine").All(&ones); err != nil {
			fmt.Println(ones, err)
		}
		if len(ones) > 0 {
			for _, one := range ones {
				RefreshStores(one.Uuid)
			}
		}
	}()

}

func RefreshStores(uuid string) error {
	var one Machine
	err := o.Using("default")
	if err != nil {
		return err
	}

	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return err
	}
	if err := RefreshReJournals(uuid); err != nil {
		return err
	}

	if one.Devtype == "export" {
		return nil
	}

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
	jours, err := resJournals(uuid)
	if err != nil {
		return store, err
	}

	store.RestDisks = disks
	store.RestRaids = raids
	store.RestVolumes = vols
	store.RestFs = fs
	store.RestInits = inits
	store.RestJours = jours

	return store, nil
}

func refreshSetRozofs(settingtype string, ip string, export string) (int, int, string, error) {

	var one []Device

	if settingtype == "export" || settingtype == "client" {
		return 0, 0, "", nil
	}

	cluNow, _ := o.QueryTable("device").Filter("devtype", settingtype).Filter("export", export).Filter("status", 0).Exclude("cid", 0).All(&one)
	fmt.Println(cluNow)
	fmt.Println(9)
	if cluNow > 0 {
		cid := one[len(one)-1].Cid
		sid := one[len(one)-1].Sid + 1
		slot := strconv.Itoa(cid) + "_" + strconv.Itoa(sid)

		return cid, sid, slot, nil

	} else {
		cluBefore, _ := o.QueryTable("device").Filter("devtype", settingtype).Filter("export", export).Filter("status", 1).Exclude("cid", 0).All(&one)
		if cluBefore > 0 {
			expanded, _ := o.QueryTable("device").Filter("devtype", settingtype).Filter("export", export).Filter("expand", 1).All(&one)
			if expanded > 0 {
				cid := one[len(one)-1].Cid + 1
				slot := strconv.Itoa(cid) + "_" + "1"
				return cid, 1, slot, nil

			} else {
				cid := one[len(one)-1].Cid
				sid := one[len(one)-1].Sid + 1
				slot := strconv.Itoa(cid) + "_" + strconv.Itoa(sid)
				return cid, sid, slot, nil

			}

		} else {
			return 1, 1, "1_1", nil
		}

	}

	return 1, 1, "0", nil

}

func refreshSetStorages(ip string, export string, cid int) (int, string, error) {
	var one []Device
	var sid int
	var slot string

	cluBefore, _ := o.QueryTable("device").Filter("export", export).Exclude("cid", 0).Filter("expand", 0).All(&one)
	if cluBefore > 0 {
		sid = one[len(one)-1].Sid + 1
		slot = strconv.Itoa(cid) + "_" + strconv.Itoa(sid)
	} else {
		sid = 1
		slot = strconv.Itoa(cid) + "_" + strconv.Itoa(sid)

	}

	return sid, slot, nil
}

func RefreshStatRemove(uuid string) { //auto delete info.yml  monitoring
	var one Machine
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		fmt.Println(err)
	}
	ip := one.Ip

	//path:="/home/monitor/info/vars/info.yml"
	path := "/root/code/yml/vars/info.yml"
	str := read(path)

	var stat StatInfo
	yaml.Unmarshal([]byte(str), &stat)

	fmt.Printf("%+v", one)
	if one.Devtype == "export" {
		stat.Exports = strings.Replace(stat.Exports, ip, "", -1)
	} else {
		stat.Storages = strings.Replace(stat.Storages, ip, "", -1)
	}
	down, _ := yaml.Marshal(&stat)
	write(path, fmt.Sprintf("---\n%s\n", string(down)))
}

func RefreshStatAdd(ip string) { //auto add info.yml
	var one Machine
	if _, err := o.QueryTable("machine").Filter("ip", ip).All(&one); err != nil {
		fmt.Println(err)
	}

	//path:="/home/monitor/info/vars/info.yml"
	path := "/root/code/yml/vars/info.yml"
	str := read(path)
	fmt.Printf("%+v", one)
	var stat StatInfo
	yaml.Unmarshal([]byte(str), &stat)
	if strings.Contains(stat.Storages, ip) || strings.Contains(stat.Exports, ip) {

	} else {
		if one.Devtype == "export" {
			stat.Exports = stat.Exports + "," + ip
		} else {
			stat.Storages = stat.Storages + "," + ip
		}
	}

	down, _ := yaml.Marshal(&stat)
	write(path, fmt.Sprintf("---\n%s\n", string(down)))
}

func read(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func write(path string, str string) {
	yaml := []byte(str)

	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	err = ioutil.WriteFile(path, yaml, 0666)
	if err != nil {
		panic(err)
	}

}
