package main

type HeartBeat struct {
	Event     string `json:"event"`
	Ip        string `json:"ip"`
	MachineId string `json:"machineId"`
	Status    string `json:"stauts"`
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

type Machine struct {
	Uid     int    `orm:"pk"  json:"uid"`
	Uuid    string `json:"uuid"`
	Ip      string `json:"ip"`
	Devtype string `json:"ip"`
	Slotnr  int    `json:"slotnr"`
	Status  bool   `json:"status"`
}
