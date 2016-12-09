package main

import (
	//"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Emergency struct {
	Id         int
	Created_at time.Time `orm:"index" json:"created"`
	Updated_at time.Time `orm:"index" json:"updated"`
	Ip         string    `json:"ip"`
	Event      string    `json:"event"`
	Level      string    `json:"level"`
	Status     bool      `json:"status"`
	//	Message    string    `json:"message"`
}

var o orm.Ormer

func Initdb() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:passwd@/speediodb?charset=utf8", 30)
	orm.RegisterModel(new(Emergency), new(Machine))

	o = orm.NewOrm()
}

func InsertJournals(event, machine string) error {
	var one Emergency
	switch event {
	case "ping.offline", "disk.unplugged", "raid.degraded", "volume.failed", "raid.failed":
		one.Level = "warning"
		one.Status = false
	default:
		one.Level = "info"
		one.Status = true
	}

	one.Event = event
	one.Ip = machine
	one.Created_at = time.Now()
	one.Updated_at = time.Now()
	if _, err := o.Insert(&one); err != nil {
		return err
	}
	return nil
}
