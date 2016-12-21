package main

import (
	"time"
)

type HeartBeat struct {
	Event     string `json:"event"`
	Ip        string `json:"ip"`
	MachineId string `json:"machineId"`
	//Status    string `json:"stauts"`
}

type DiskUnplugged struct {
	Event     string `json:"event"`
	Uuid      string `json:"uuid"`
	Location  string `json:"location"`
	DevName   string `json:"dev_name"`
	Ip        string `json:"ip"`
	MachineId string `json:"machineId"`
}

type DiskPlugged struct {
	Event     string `json:"event"`
	Uuid      string `json:"uuid"`
	Ip        string `json:"ip"`
	MachineId string `json:"machineId"`
}

type RaidRemove struct {
	Event     string   `json:"event"`
	Uuid      string   `json:"uuid"`
	RaidDisks []string `json:"raid_disks"`
	Ip        string   `json:"ip"`
	MachineId string   `json:"machineId"`
}

type FsSystem struct {
	Event     string `json:"event"`
	Type      string `json:"type"`
	Volume    string `json:"volume"`
	Ip        string `json:"ip"`
	MachineId string `json:"machineId"`
}

type Machine struct {
	Uid     int       `orm:"pk" json:"uid"`
	Uuid    string    `json:"uuid"`
	Ip      string    `json:"ip"`
	Devtype string    `json:"ip"`
	Slotnr  int       `json:"slotnr"`
	Created time.Time `orm:"index" json:"created"`
	Status  bool      `json:"status"`
}

type Disks struct { //the table is remote disk
	Uuid      string `orm:"pk"  json:"id"`
	Health    string `json:"health"`
	Role      string `json:"role"`
	Location  string `json:"location"`
	Raid      string `json:"raid"`
	CapSector int64  `json:"cap"`
	Vendor    string `json:"vendor"`
	Model     string `json:"model"`
	Sn        string `json:"sn"`
}

type Disk struct { //the table is local disk
	Disks
	MachineId string `orm:"column(machineId)"  json:"machineId"`
}

type Initiators struct { //the table'name is initiators
	Wwn    string `orm:"pk"                 json:"wwn"`
	Target string `orm:"column(target_wwn)" json:"target"`
}

type Initiator struct { //the table is local initiator
	Initiators
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type Raids struct { //the table'name is disk
	Uuid    string `orm:"pk"  json:"id"`
	Health  string `json:"health"`
	Level   string `json:"level"`
	Name    string `json:"name"`
	Cap     int    `json:"cap"`
	Used    int    `orm:"column(used_cap)" json:"used"`
	Deleted int    `json:"deleted"`
}

type Raid struct { //the table is local raid
	Raids
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type Volumes struct { //the table'name is disk
	Uuid    string `orm:"pk" json:"id"`
	Health  string `json:"health"`
	Name    string `json:"name"`
	Cap     int64  `json:"cap"`
	Used    int    `json:"used"`
	Type    string `orm:"column(owner_type)" json:"type"`
	Deleted int    `json:"deleted"`
}

type Volume struct { //the table is local vol
	Volumes
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type Xfs struct { //the table'name is xfs
	Uuid       string `orm:"pk" json:"id""`
	Volume     string `json:"volume"`
	Name       string `json:"name"`
	Chunk      string `orm:"column(chunk_kb)" json:"chunk"`
	Type       string `json:"type"`
	MountPoint string `orm:"column(mountpoint)" json:"mountpoint"`
}

type Filesystems struct { //the table is local fs
	Xfs
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type RaidVolumes struct {
	Id     int    `json:"id"`
	Raid   string `json:"raid"`
	Volume string `json:"volume"`
	Type   string `json:"type"`
}

type RaidVolume struct {
	RaidVolumes
	MachineId string `orm:"column(machineId)"`
}

type InitiatorVolumes struct {
	Id        int    `json:"id"`
	Initiator string `json:"initiator"`
	Volume    string `json:"volume"`
}

type InitiatorVolume struct {
	InitiatorVolumes
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type NetworkInitiators struct {
	Id        int    `json:"id"`
	Initiator string `json:"initiator"`
	Eth       string `json:"eth"`
	Port      int    `json:"port"`
}

type NetworkInitiator struct {
	NetworkInitiators
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type Journals struct {
	Id         int
	Created_at time.Time `orm:"index" json:"created"`
	Updated_at time.Time `orm:"index" json:"updated"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
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

type Mail struct {
	Uid     int    `orm:"pk" json:"uid"`
	Address string `json:"address"`
	Level   int    `json:"level"`
	Ttl     int    `json:"ttl"`
}

type Log struct {
	Message    string `json:"message"`
	Created_at int64  `json:"created_at"`
	Level      string `json:"level"`
	Source     string `json:"scource"`
}
