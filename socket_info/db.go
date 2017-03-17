package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var o orm.Ormer

type Threshhold struct {
	Uid     int    `orm:"pk" json:"uid"`
	Type    string `json:"type"`
	Normal  int    `json:"normal"`
	Warning int    `json:"warning"`
}

func Initdb() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:passwd@/speediodb?charset=utf8&loc=Local", 30)
	orm.RegisterModel(new(Threshhold))
	o = orm.NewOrm()
}

func Detecting(info Stats) {
	var cpus, mems, cachesT, cachesU, sysT, sysA, fsT, fsA float64
	var num float64
	//exports = info.Exports

	for _, vals := range info.Storages {
		if len(vals.Dev) > 0 {
			fmt.Printf("%+v", vals)
			fmt.Printf("%+v", info.Storages)
			num = float64(len(info.Storages))
			//singleDet(vals)
			cpus += vals.Dev[0].Cpu
			mems += vals.Dev[0].Mem
			cachesT += vals.Dev[0].CacheT
			cachesU += vals.Dev[0].CacheU
			for _, df := range vals.Dev[0].Dfs {
				if df.Name == "system" {
					sysT += df.Total
					sysA += df.Available
				} else if df.Name == "filesystem" {
					fsT += df.Total
					fsA += df.Available
				}
			}
		}
	}

	if exist := o.QueryTable("threshhold").Filter("type", "cpu").Filter("warning__gt", cpus/num).Exist(); !exist {
		CpuNum += 1
		Publish("cpu", "All", cpus/num, CpuNum)
	} else {
		CpuNum = 0
	}
	if exist := o.QueryTable("threshhold").Filter("type", "mem").Filter("warning__gt", mems/num).Exist(); !exist {
		MemNum += 1
		Publish("mem", "All", mems/num, MemNum)
	} else {
		MemNum = 0
	}
	if exist := o.QueryTable("threshhold").Filter("type", "cache").Filter("warning__gt", cachesU/cachesT*100).Exist(); !exist {
		CacheNum += 1
		Publish("cache", "All", cachesU/cachesT*100, CacheNum)
	} else {
		CacheNum = 0
	}
	if exist := o.QueryTable("threshhold").Filter("type", "system").Filter("warning__gt", (sysT-sysA)/sysT*100).Exist(); !exist {
		SysNum += 1
		Publish("sys", "All", (sysT-sysA)/sysT*100, SysNum)
	} else {
		SysNum = 0
	}
	if exist := o.QueryTable("threshhold").Filter("type", "filesystem").Filter("warning__gt", (fsT-fsA)/fsT*100).Exist(); !exist {
		FsNum += 1
		Publish("fs", "All", (fsT-fsA)/fsT*100, FsNum)
	} else {
		FsNum = 0
	}

}

func Publish(sendtype, machine string, value float64, count int) {
	if count > 3 {
		return
	} else if count == 3 {
		NsqRequest("info.warning", sendtype, machine, "true", value)
	} else if count < 3 {
		NsqRequest("info.warning", sendtype, machine, "false", value)
	}

}
