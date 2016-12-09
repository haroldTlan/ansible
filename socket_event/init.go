package main

type HeartBeat struct {
	Event  string `json:"event"`
	Ip     string `json:"ip"`
	Status string `json:"stauts"`
}

type DiskUnplugged struct {
	Event    string `json:"event"`
	Uuid     string `json:"uuid"`
	Location string `json:"location"`
	DevName  string `json:"dev_name"`
	Ip       string `json:"ip"`
}

type DiskPlugged struct {
	Event string `json:"event"`
	Uuid  string `json:"uuid"`
	Ip    string `json:"ip"`
}

type RaidRemove struct {
	Event     string   `json:"event"`
	Uuid      string   `json:"uuid"`
	RaidDisks []string `json:"raid_disks"`
	Ip        string   `json:"ip"`
}

type Disks struct {
	Disk string `json:"disk"`
}

type Machine struct {
	Uid     int    `orm:"pk"` // json:"uid"`
	Uuid    string // `json:"uuid"`
	Ip      string // `json:"ip"`
	Devtype string // `json:"ip"`
	Slotnr  int    // `json:"slotnr"`
	Status  bool   //  `json:"status"`
}
