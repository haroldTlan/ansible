package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	//	"strings"
	"reflect"
	"time"
)

type Machine struct {
	Uid     int       `orm:"pk"` // json:"uid"`
	Uuid    string    // `json:"uuid"`
	Ip      string    // `json:"ip"`
	Slotnr  int       // `json:"slotnr"`
	Created time.Time `orm:"index"` //json:"created"`
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

type Device struct {
	Uid     int       `orm:"pk" json:"uid"`
	Uuid    string    `json:"uuid"`
	Ip      string    `json:"ip"`
	Version string    `json:"version"`
	Devtype string    `json:"devtype"`
	Size    string    `json:"size"`
	Status  bool      `json:"status"`
	Export  string    `json:"export"`
	Cid     int       `json:"cid"` //cid:cluster
	Sid     int       `json:"sid"` //sid:host
	Slot    string    `json:"slot"`
	Expand  bool      `json:"expand"`
	Created time.Time `orm:"index"json:"created"`
}

//var o orm.Ormer

func init() {
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:passwd@/speediodb?charset=utf8", 30)
	name := "haha"
	ip := "192.168.2.190"

	err := orm.RegisterDataBase(fmt.Sprintf("%s", name), "mysql", fmt.Sprintf("root:passwd@tcp(%s:3306)/speediodb?charset=utf8", ip), 30)
	if err != nil {
		fmt.Println("\n\n\n")
		//strings.Contains((err).(string), "No route")
		for m := 0; m < reflect.TypeOf(err).NumMethod(); m++ {
			method := reflect.TypeOf(err).Method(m)
			fmt.Println(method.Type)         // func(*main.MyStruct) string
			fmt.Println(method.Name)         // GetName
			fmt.Println(method.Type.NumIn()) // 参数个数
			fmt.Println(method.Type)         // 参数类型
		}
		fmt.Println("\n\n\n")
	}

	//orm.RegisterModel(new(Machine), new(Disks), new(Disk), new(Raid), new(Raids), new(Volume), new(Volumes), new(Filesystems), new(Xfs), new(Initiator), new(Initiators), new(Setting), new(Journals), new(RaidVolumes), new(RaidVolume), new(InitiatorVolumes), new(InitiatorVolume), new(NetworkInitiators), new(NetworkInitiator), new(RozofsSetting), new(Device))
	orm.RegisterModel(new(Machine), new(Disks), new(Device))
	//orm.RunSyncdb("default", false, false)
	orm.Debug = true

}

func main() {
	o := orm.NewOrm()

	err := o.Using("haha")

	if err != nil {
		fmt.Println(err)

	}

	ones := make([]Disks, 0)
	if _, err := o.QueryTable("disks").All(&ones); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", ones)
}
