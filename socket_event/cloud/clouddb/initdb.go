package clouddb

import (
	"time"
)

type Machine struct {
	Uid     int       `orm:"pk"` // json:"uid"`
	Uuid    string    // `json:"uuid"`
	Ip      string    // `json:"ip"`
	Devtype string    // `json:"ip"`
	Slotnr  int       // `json:"slotnr"`
	Created time.Time `orm:"index"` //json:"created"`
	Status  bool      //  `json:"status"`
}

type Disks struct { //the table is remote disk
	Uuid      string `orm:"pk"` // json:"id"`
	Health    string //`json:"health"`
	Role      string //`json:"role"`
	Location  string //`json:"location"`
	Raid      string //`json:"raid"`
	CapSector int64  //`json:"cap"`
	Vendor    string //`json:"vendor"`
	Model     string //`json:"model"`
	Sn        string //`json:"sn"`
}

type Disk struct { //the table is local disk
	Disks
	MachineId string `orm:"column(machineId)"` // json:"machineId"`
}

type Initiators struct { //the table'name is initiators
	Wwn    string `orm:"pk"`                 //json:"wwn"`
	Target string `orm:"column(target_wwn)"` //json:"target"`
}

type Initiator struct { //the table is local initiator
	Initiators
	MachineId string `orm:"column(machineId)"` //json:"machineId"`
}

type Raids struct { //the table'name is disk
	Uuid    string `orm:"pk"` // json:"id"`
	Health  string //`json:"health"`
	Level   string //`json:"level"`
	Name    string //`json:"name"`
	Cap     int    //`json:"cap"`
	Used    int    `orm:"column(used_cap)"` //json:"used"`
	Deleted int    //`json:"deleted"`
}

type Raid struct { //the table is local raid
	Raids
	MachineId string `orm:"column(machineId)"` //json:"machineId"`
}

type Volumes struct { //the table'name is disk
	Uuid    string `orm:"pk"` //json:"id"`
	Health  string //`json:"health"`
	Name    string //`json:"name"`
	Cap     int64  //`json:"cap"`
	Used    int    //`json:"used"`
	Type    string `orm:"column(owner_type)"` //json:"type"`
	Deleted int    //`json:"deleted"`
}

type Volume struct { //the table is local vol
	Volumes
	MachineId string `orm:"column(machineId)"` //json:"machineId"`
}

type Xfs struct { //the table'name is xfs
	Uuid       string `orm:"pk"` //json:"id""`
	Volume     string //`json:"volume"`
	Name       string //`json:"name"`
	Chunk      string `orm:"column(chunk_kb)"` //json:"chunk"`
	Type       string //`json:"type"`
	MountPoint string `orm:"column(mountpoint)"` //json:"mountpoint"`
}

type Filesystems struct { //the table is local fs
	Xfs
	MachineId string `orm:"column(machineId)"` //json:"machineId"`
}

type RaidVolumes struct {
	Id     int    //`json:"id"`
	Raid   string //`json:"raid"`
	Volume string //`json:"volume"`
	Type   string //`json:"type"`
}

type RaidVolume struct {
	RaidVolumes
	MachineId string `orm:"column(machineId)"`
}

type InitiatorVolumes struct {
	Id        int    //`json:"id"`
	Initiator string //`json:"initiator"`
	Volume    string //`json:"volume"`
}

type InitiatorVolume struct {
	InitiatorVolumes
	MachineId string `orm:"column(machineId)"` //json:"machineId"`
}

type NetworkInitiators struct {
	Id        int    //`json:"id"`
	Initiator string //`json:"initiator"`
	Eth       string //`json:"eth"`
	Port      int    //`json:"port"`
}

type NetworkInitiator struct {
	NetworkInitiators
	MachineId string `orm:"column(machineId)"` //json:"machineId"`
}

type Journals struct {
	Id         int
	Created_at time.Time `orm:"index" json:"created"`
	Updated_at time.Time `orm:"index" json:"updated"`
	//Machine    string    `json:"machine"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Journal struct {
	Journals
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type Emergency struct {
	Uid            int       `orm:"pk" json:"uid"`
	Created_at     time.Time `orm:"index" json:"created"`
	Updated_at     time.Time `orm:"index" json:"updated"`
	Ip             string    `json:"ip"`
	Event          string    `json:"event"`
	Level          string    `json:"level"`
	Message        string    `json:"message"`
	ChineseMessage string    `json:"chinese_message"`
	Status         bool      `json:"status"`
}

type Setting struct {
	Uid         int    `orm:"pk"`                  //json:"uid"`
	Settingtype string `orm:"column(settingtype)"` //json:"settingtype"`
	Ip          string //`json:"ip"`
	Status      bool   //`json:"status"`
}

type RozofsSetting struct {
	Uid         int    `orm:"pk"                  json:"uid"`
	Settingtype string `orm:"column(settingtype)" json:"settingtype"`
	Ip          string `json:"ip"`
	Expand      string `json:"expand"`
	Status      bool   `json:"status"`
}

type LocalJournals struct {
	Uid            int       `orm:"pk" json:"uid"`
	Created_at     time.Time `orm:"index" json:"created"`
	Updated_at     time.Time `orm:"index" json:"updated"`
	Level          string    `json:"level"`
	ChineseMessage string    `json:"chinese_message"`
	Status         bool      `json:"status"`
}

type Export struct {
	Uid     int       `orm:"pk" json:"uid"`
	Uuid    string    `json:"uuid"`
	Ip      string    `json:"ip"`
	Version string    `json:"version"`
	Size    string    `json:"size"`
	Status  bool      `json:"status"`
	Devtype string    `json:"devtype"`
	Created time.Time `orm:"index" json:"created"`
}

type Storage struct {
	Uid     int       `orm:"pk" json:"uid"`
	Uuid    string    `json:"uuid"`
	Ip      string    `json:"ip"`
	Version string    `json:"version"`
	Size    string    `json:"size"`
	Export  string    `json:"export"`
	Cid     int       `json:"cid"`
	Sid     int       `json:"sid"`
	Slot    string    `json:"slot"`
	Status  bool      `json:"status"`
	Devtype string    `json:"devtype"`
	Created time.Time `orm:"index" json:"created"`
}

type Client struct {
	Uid     int       `orm:"pk" json:"uid"`
	Uuid    string    `json:"uuid"`
	Ip      string    `json:"ip"`
	Version string    `json:"version"`
	Size    string    `json:"size"`
	Status  bool      `json:"status"`
	Devtype string    `json:"devtype"`
	Created time.Time `orm:"index" json:"created"`
}

type Threshhold struct {
	Uid     int    `orm:"pk" json:"uid"`
	Type    string `json:"type"`
	Normal  int    `json:"normal"`
	Warning int    `json:"warning"`
}

type Mail struct {
	Uid     int    `orm:"pk" json:"uid"`
	Address string `json:"address"`
}
