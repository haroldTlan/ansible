package main

import (
	//"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Emergency struct {
	Id              int
	Created_at      time.Time `orm:"index" json:"created"`
	Updated_at      time.Time `orm:"index" json:"updated"`
	Ip              string    `json:"ip"`
	Event           string    `json:"event"`
	Level           string    `json:"level"`
	Status          bool      `json:"status"`
	Message         string    `json:"message"`
	Chinese_message string    `json:"chinese_message"`
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
	message, chinese_message := messageTransform(event)

	one.Message = message
	one.Chinese_message = "服务器" + " " + machine + " " + chinese_message
	one.Event = event
	one.Ip = machine
	one.Created_at = time.Now()
	one.Updated_at = time.Now()
	if _, err := o.Insert(&one); err != nil {
		return err
	}
	return nil
}

func messageTransform(event string) (string, string) {
	switch event {
	case "ping.offline":
		message := "offline"
		chinese_message := "设备掉线"
		return message, chinese_message
	case "ping.online":
		message := "online"
		chinese_message := "设备上线"
		return message, chinese_message
	case "disk.unplugged":
		message := "disk unplugged"
		chinese_message := "磁盘拔出"
		return message, chinese_message
	case "disk.plugged":
		message := "disk plugged"
		chinese_message := "磁盘插入"
		return message, chinese_message
	case "raid.created":
		message := "raid created"
		chinese_message := "创建阵列"
		return message, chinese_message
	case "raid.removed":
		message := "raid removed"
		chinese_message := "删除阵列"
		return message, chinese_message
	case "volume.created":
		message := "volume created"
		chinese_message := "创建虚拟磁盘"
		return message, chinese_message
	case "volume.removed":
		message := "volume removed"
		chinese_message := "删除虚拟磁盘"
		return message, chinese_message
	case "raid.degraded":
		message := "raid.degraded"
		chinese_message := "阵列降级"
		return message, chinese_message
	case "raid.failed":
		message := "raid failed"
		chinese_message := "阵列损坏"
		return message, chinese_message
	case "volume.failed":
		message := "volume failed"
		chinese_message := "虚拟磁盘损坏"
		return message, chinese_message
		/*	case :
				message :=
				chinese_message :=
				return message ,chinese_message
			case :
				message :=
				chinese_message :=
				return message ,chinese_message*/
	default:
		return "", "未知"
	}
}
