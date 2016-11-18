package main

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"time"
)

type Devices struct {
	Exports  []Export
	Storages []Storage
	Clients  []Client
}

type Machine struct {
	Uid     int       `orm:"pk"` // json:"uid"`
	Uuid    string    // `json:"uuid"`
	Ip      string    // `json:"ip"`
	Devtype string    // `json:"ip"`
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

type Disk struct { //the table is local disk
	Disks
	MachineId string `orm:"column(machineId)"` // json:"machineId"`
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

type Initiators struct { //the table'name is initiators
	Wwn    string `orm:"pk"`                 //json:"wwn"`
	Target string `orm:"column(target_wwn)"` //json:"target"`
}

type Initiator struct { //the table is local initiator
	Initiators
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
	Level      string    `json:"level"`
	Message    string    `json:"message"`
}

type Journal struct {
	Journals
	MachineId string `orm:"column(machineId)" json:"machineId"`
}

type Setting struct {
	Uid         int    `orm:"pk"`                  //json:"uid"`
	Settingtype string `orm:"column(settingtype)"` //json:"settingtype"`
	Ip          string //`json:"ip"`
	Status      bool   //`json:"status"`
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

type RozofsSetting struct {
	Uid         int    `orm:"pk"                  json:"uid"`
	Settingtype string `orm:"column(settingtype)" json:"settingtype"`
	Ip          string `json:"ip"`
	Expand      string `json:"expand"`
	Status      bool   `json:"status"`
}

type LocalLog struct {
	Uid             int       `orm:"pk"`    //json:"uid"`
	Created_at      time.Time `orm:"index"` //json:"created"`
	Updated_at      time.Time `orm:"index"` //json:"updated"`
	Level           string    //`json:"level"`
	Message         string    //`json:"message"`
	Chinese_message string    //`json:"chinese_message"`
}

type Export struct {
	Uid     int       `orm:"pk" json:"uid"`
	Uuid    string    `json:"uuid"`
	Ip      string    `json:"ip"`
	Version string    `json:"version"`
	Size    string    `json:"size"`
	Cid     int       `json:"cid"`
	Status  bool      `json:"status"`
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
	Created time.Time `orm:"index" json:"created"`
}

type Client struct {
	Uid     int       `orm:"pk" json:"uid"`
	Uuid    string    `json:"uuid"`
	Ip      string    `json:"ip"`
	Version string    `json:"version"`
	Size    string    `json:"size"`
	Status  bool      `json:"status"`
	Created time.Time `orm:"index" json:"created"`
}

var o orm.Ormer

func Initdb() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:passwd@/speediodb?charset=utf8", 30)
	orm.RegisterModel(new(Machine), new(Disks), new(Disk), new(Raid), new(Raids), new(Volume), new(Volumes), new(Filesystems), new(Xfs), new(Initiator), new(Initiators), new(Setting), new(Journal), new(Journals), new(RaidVolumes), new(RaidVolume), new(InitiatorVolumes), new(InitiatorVolume), new(NetworkInitiators), new(NetworkInitiator), new(RozofsSetting), new(Device), new(LocalLog), new(Export), new(Storage), new(Client))

	//orm.RunSyncdb("default", false, false)
	orm.Debug = true
	o = orm.NewOrm()

	InitRemote()
}

func InitRemote() error {
	machines, err := SelectAllMachines()
	if err != nil {
		return err
	}
	if mlen := len(machines); mlen > 0 {
		for i := 0; i < mlen; i++ {
			name, ip := MachineType(machines[i])
			/*
				err := o.Using(name)
				fmt.Println(name)
				if err != nil {
					err := orm.RegisterDataBase(fmt.Sprintf("%s", name), "mysql", fmt.Sprintf("root:passwd@tcp(%s:3306)/speediodb?charset=utf8", ip), 30)
					fmt.Println(err)
					if err != nil {
						DelMachine(machines[i].Uuid)
						return err
					}

				}*/
			err := orm.RegisterDataBase(fmt.Sprintf("%s", name), "mysql", fmt.Sprintf("root:passwd@tcp(%s:3306)/speediodb?charset=utf8", ip), 30)
			if err != nil {
				DelMachine(machines[i].Uuid)
				return err
			}

		}

		//InitLocal(machines)
	} else {
		//TODO!!!!!
	}
	return nil

}

func InitSingleRemote(machine Machine) error {
	name, ip := MachineType(machine)
	orm.RegisterDataBase(fmt.Sprintf("%s", name), "mysql", fmt.Sprintf("root:passwd@tcp(%s:3306)/speediodb?charset=utf8", ip), 30)
	return nil
}

func MachineType(machine Machine) (string, string) {
	ip := machine.Ip
	tempIp := strings.Join(strings.Split(ip, "."), "")
	name := "remote" + tempIp

	return name, ip
}

func InitLocal(dev []Machine) {
	for i := 0; i < len(dev); i++ {
		name, _ := MachineType(dev[i])
		err := o.Using(name)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func Urandom() string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func InsertMachine(ip string, slotnr int, devtype string) error {
	orm.Debug = true

	var one Machine

	uran := Urandom()
	uuid := uran + "zip" + strings.Join(strings.Split(ip, "."), "")

	one.Uuid = uuid
	one.Ip = ip
	one.Devtype = devtype
	one.Slotnr = slotnr
	one.Created = time.Now()
	if _, err := o.Insert(&one); err != nil {
		return err
	}

	err := InitSingleRemote(one)

	if err != nil {
		return err
	}
	RefreshStores(uuid)

	return nil
}

func SelectAllMachines() ([]Machine, error) {
	//get all machine
	ones := make([]Machine, 0)
	if _, err := o.QueryTable("machine").All(&ones); err != nil {
		return ones, err
	}
	return ones, nil
}

func DelMachine(uuid string) error {
	if _, err := o.QueryTable("disk").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("raid").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("volume").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("filesystems").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("initiator").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("initiator_volume").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("network_initiator").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("raid_volume").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("journal").Filter("machineId", uuid).Delete(); err != nil {
		return err
	}
	RefreshStatRemove(uuid)

	if _, err := o.QueryTable("machine").Filter("uuid", uuid).Delete(); err != nil {
		return err
	}

	return nil
}

func SelectDisks() ([]Disk, int64, error) {
	ones := make([]Disk, 0)
	num, err := o.QueryTable("disk").Exclude("location__exact", "").All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertDisksOfMachine(redisks []Disks, uuid string) error {
	if mlen := len(redisks); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}

		for i := 0; i < mlen; i++ {
			var loc Disk
			num, err := o.QueryTable("disk").Filter("uuid__exact", redisks[i].Uuid).Filter("machineId__exact", uuid).All(&loc) //decide update or not
			if err != nil {
				return err
			}
			one := Disk{Disks: redisks[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil

}

func RefreshReDisks(uuid string) error {
	o.QueryTable("disk").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ones := make([]Disks, 0)
	_, err = o.QueryTable("disks").All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("%+v", ones)

	InsertDisksOfMachine(ones, uuid)

	return nil
}

func SelectVolumes() ([]Volume, int64, error) {
	ones := make([]Volume, 0)
	num, err := o.QueryTable("volume").Filter("deleted__exact", 0).All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertVolumesOfMachine(revols []Volumes, uuid string) error {
	if mlen := len(revols); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}

		for i := 0; i < mlen; i++ {
			var loc Volume
			num, err := o.QueryTable("volume").Filter("uuid__exact", revols[i].Uuid).All(&loc) //decide update or not
			if err != nil {
				return err
			}
			one := Volume{Volumes: revols[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshReVolumes(uuid string) error {
	o.QueryTable("volume").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ones := make([]Volumes, 0)
	_, err = o.QueryTable("volumes").All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	InsertVolumesOfMachine(ones, uuid)

	return nil
}

func SelectRaids() ([]Raid, int64, error) {
	ones := make([]Raid, 0)
	num, err := o.QueryTable("raid").Filter("deleted__exact", 0).All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertRaidsOfMachine(reraids []Raids, uuid string) error {
	if mlen := len(reraids); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}

		for i := 0; i < mlen; i++ {
			var loc Raid
			num, err := o.QueryTable("raid").Filter("uuid__exact", reraids[i].Uuid).All(&loc) //decide update or not
			if err != nil {
				return err
			}
			one := Raid{Raids: reraids[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do

		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshReRaids(uuid string) error {
	o.QueryTable("raid").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ones := make([]Raids, 0)
	_, err = o.QueryTable("raids").All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("%+v", ones)
	InsertRaidsOfMachine(ones, uuid)

	return nil
}

func SelectFilesystems() ([]Filesystems, int64, error) {
	ones := make([]Filesystems, 0)
	num, err := o.QueryTable("filesystems").All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertFilesystemsOfMachine(refs []Xfs, uuid string) error {
	if mlen := len(refs); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}

		for i := 0; i < mlen; i++ {
			var loc Filesystems
			num, err := o.QueryTable("filesystems").Filter("uuid__exact", refs[i].Uuid).All(&loc) //decide update or not
			if err != nil {
				return err
			}
			one := Filesystems{Xfs: refs[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshReFilesystems(uuid string) error {
	o.QueryTable("filesystems").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ones := make([]Xfs, 0)
	_, err = o.QueryTable("xfs").All(&ones)
	if err != nil {
		return err
	}
	InsertFilesystemsOfMachine(ones, uuid)

	return nil
}

func SelectInitiators() ([]Initiator, int64, error) {
	ones := make([]Initiator, 0)
	num, err := o.QueryTable("initiator").All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertInitiatorsOfMachine(refs []Initiators, uuid string) error {
	if mlen := len(refs); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}
		for i := 0; i < mlen; i++ {
			var loc Initiator
			num, err := o.QueryTable("initiator").Filter("wwn__exact", refs[i].Wwn).All(&loc) //decide update or not
			if err != nil {
				fmt.Println(err)
				return err
			}
			one := Initiator{Initiators: refs[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshReInitiators(uuid string) error {
	o.QueryTable("initiator").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		return err
	}

	ones := make([]Initiators, 0)
	_, err = o.QueryTable("initiators").All(&ones)
	if err != nil {
		return err
	}
	InsertInitiatorsOfMachine(ones, uuid)

	return nil
}

func SelectRaidVolumes() ([]RaidVolume, int64, error) {
	ones := make([]RaidVolume, 0)
	num, err := o.QueryTable("raid_volume").All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertRaidVolumesOfMachine(remote []RaidVolumes, uuid string) error {
	if mlen := len(remote); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}

		for i := 0; i < mlen; i++ {
			var loc RaidVolume
			num, err := o.QueryTable("raid_volume").Filter("volume__exact", remote[i].Volume).All(&loc) //decide update or not
			if err != nil {
				return err
			}
			one := RaidVolume{RaidVolumes: remote[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshReRaidVolumes(uuid string) error {
	o.QueryTable("raid_volume").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ones := make([]RaidVolumes, 0)
	_, err = o.QueryTable("raid_volumes").All(&ones)
	if err != nil {
		return err
	}
	InsertRaidVolumesOfMachine(ones, uuid)

	return nil
}

func SelectInitVolumes() ([]InitiatorVolume, int64, error) {
	ones := make([]InitiatorVolume, 0)
	num, err := o.QueryTable("initiator_volume").All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertInitVolumesOfMachine(remote []InitiatorVolumes, uuid string) error {
	if mlen := len(remote); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}
		for i := 0; i < mlen; i++ {
			var loc InitiatorVolume
			num, err := o.QueryTable("initiator_volume").Filter("volume__exact", remote[i].Initiator).All(&loc) //decide update or not     !!!!!!!!!!!!!!!!!key is not initiator
			if err != nil {
				return err
			}
			one := InitiatorVolume{InitiatorVolumes: remote[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshReInitVolumes(uuid string) error {
	o.QueryTable("initiator_volume").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ones := make([]InitiatorVolumes, 0)
	_, err = o.QueryTable("initiator_volumes").All(&ones)
	if err != nil {
		return err
	}
	InsertInitVolumesOfMachine(ones, uuid)

	return nil
}

func SelectNetInits() ([]NetworkInitiator, int64, error) {
	ones := make([]NetworkInitiator, 0)
	num, err := o.QueryTable("network_initiator").All(&ones)
	if err != nil {
		return ones, 0, err
	}
	return ones, num, nil
}

func InsertNetInitsOfMachine(remote []NetworkInitiators, uuid string) error {
	if mlen := len(remote); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}

		for i := 0; i < mlen; i++ {
			var loc NetworkInitiator
			num, err := o.QueryTable("network_initiator").Filter("eth__exact", remote[i].Eth).Filter("initiator__exact", remote[i].Initiator).All(&loc) //decide update or not
			if err != nil {
				return err
			}
			one := NetworkInitiator{NetworkInitiators: remote[i], MachineId: uuid}

			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshReNetInits(uuid string) error {
	o.QueryTable("network_initiator").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ones := make([]NetworkInitiators, 0)
	_, err = o.QueryTable("network_initiators").All(&ones)
	if err != nil {
		return err
	}
	InsertNetInitsOfMachine(ones, uuid)

	return nil
}

func SelectJournals() ([]ResJournals, int64, error) {
	ones := make([]Journal, 0)
	jours := make([]ResJournals, 0)
	num, err := o.QueryTable("journal").Filter("level", "warning").All(&ones)
	if err != nil {
		return jours, 0, err
	}

	for _, val := range ones {
		var one Machine
		o.QueryTable("machine").Filter("uuid", val.MachineId).All(&one)

		var jour ResJournals
		jour.Message = val.Message
		jour.Created = val.Created_at
		jour.Unix = val.Created_at.Unix()
		jour.Level = val.Level
		jour.MachineId = val.MachineId
		jour.Ip = one.Ip
		jour.Devtype = one.Devtype
		jours = append(jours, jour)
	}

	return jours, num, nil
}

func InsertJournalsOfMachine(rejournals []Journals, uuid string) error {
	if mlen := len(rejournals); mlen > 0 {
		err := o.Using("default")
		if err != nil {
			return err
		}

		for i := 0; i < mlen; i++ {
			var loc Journal
			num, err := o.QueryTable("journal").Filter("id__exact", rejournals[i].Id).Filter("machineId__exact", uuid).All(&loc) //decide update or not
			if err != nil {
				return err
			}
			one := Journal{Journals: rejournals[i], MachineId: uuid}
			if num == 0 {
				_, err = o.Insert(&one)
				if err != nil {
					return err
				}
			} else {
				_, err = o.Update(&one)
				if err != nil {
					return err
				}
			}

		}
	} else {
		//TODO !!!!!!!! but if len == 0, nothing you can do
		err := o.Using("default")
		if err != nil {
			return err
		}
	}
	return nil

}

func RefreshReJournals(uuid string) error {
	o.QueryTable("journal").Filter("machineId", uuid).Delete()
	name := "remote" + strings.Split(uuid, "zip")[1]

	err := o.Using(name)
	if err != nil {
		return err
	}

	ones := make([]Journals, 0)
	_, err = o.QueryTable("journals").All(&ones)
	if err != nil {
		return err
	}

	InsertJournalsOfMachine(ones, uuid)

	return nil
}

func InsertDevice(ip string, version string, size string, devtype string) error {
	if devtype == "export" {
		var one Export
		num, err := o.QueryTable("export").Filter("ip", ip).All(&one)
		if err != nil {
			return err
		}
		if num == 0 {
			uran := Urandom()
			uuid := uran + "zip" + strings.Join(strings.Split(ip, "."), "")
			one.Uuid = uuid
			one.Ip = ip
			one.Version = version
			one.Size = size
			one.Status = false
			one.Created = time.Now()
			if _, err := o.Insert(&one); err != nil {
				return err
			}
		} else {
			return errors.New("Ip address already exits")

		}
	} else if devtype == "storage" {
		var one Storage
		num, err := o.QueryTable("storage").Filter("ip", ip).All(&one)
		if err != nil {
			return err
		}
		if num == 0 {
			uran := Urandom()
			uuid := uran + "zip" + strings.Join(strings.Split(ip, "."), "")
			one.Uuid = uuid
			one.Ip = ip
			one.Version = version
			one.Size = size
			one.Status = false
			one.Created = time.Now()
			if _, err := o.Insert(&one); err != nil {
				return err
			}
		} else {
			return errors.New("Ip address already exits")

		}
	} else if devtype == "client" {
		var one Client
		num, err := o.QueryTable("client").Filter("ip", ip).All(&one)
		if err != nil {
			return err
		}
		if num == 0 {
			uran := Urandom()
			uuid := uran + "zip" + strings.Join(strings.Split(ip, "."), "")
			one.Uuid = uuid
			one.Ip = ip
			one.Version = version
			one.Size = size
			one.Status = false
			one.Created = time.Now()
			if _, err := o.Insert(&one); err != nil {
				return err
			}
		} else {
			return errors.New("Ip address already exits")

		}
	}

	return nil
}

func SelectAllDevices() (Devices, error) {
	//get all machine
	var ones Devices
	if _, err := o.QueryTable("export").All(&ones.Exports); err != nil {
		return ones, err
	}
	if _, err := o.QueryTable("storage").All(&ones.Storages); err != nil {
		return ones, err
	}
	if _, err := o.QueryTable("client").All(&ones.Clients); err != nil {
		return ones, err
	}
	return ones, nil
}

func DelDevice(uuid string) error {
	if _, err := o.QueryTable("export").Filter("uuid", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("storage").Filter("uuid", uuid).Delete(); err != nil {
		return err
	}
	if _, err := o.QueryTable("client").Filter("uuid", uuid).Delete(); err != nil {
		return err
	}
	return nil
}

func InsertExports(ip string, status bool) error {
	var one Export
	num, err := o.QueryTable("export").Filter("ip", ip).All(&one)
	if err != nil {
		return err
	}

	if num == 0 {
		uran := Urandom()
		uuid := uran + "zip" + strings.Join(strings.Split(ip, "."), "")
		one.Uuid = uuid
		one.Ip = ip
		one.Version = "ZS2000"
		one.Size = "4U"
		one.Status = status
		one.Created = time.Now()
		if _, err := o.Insert(&one); err != nil {
			return err
		}
	} else {
		one.Status = status
		_, err = o.Update(&one)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertStorages(ip string, status bool, expands string) error {
	var one Storage

	one.Status = status
	if _, err := o.Update(&one); err != nil {
		return err
	}
	return nil
}

func InsertRozofsSetting(settingtype string, ip string, status bool, export string, cid int, sid int, slot string) error {
	var one Device

	num, err := o.QueryTable("device").Filter("devtype", settingtype).Filter("ip", ip).All(&one) //decide update or not
	if err != nil {
		return err
	}
	if num == 0 {
		uran := Urandom()
		uuid := uran + "zip" + strings.Join(strings.Split(ip, "."), "")

		one.Uuid = uuid

		one.Status = false
		one.Devtype = settingtype
		one.Ip = ip
		one.Status = status
		one.Version = "ZS2000"
		one.Size = "4U"
		one.Created = time.Now()

		one.Export = export
		one.Cid = cid
		one.Sid = sid
		one.Slot = slot
		_, err = o.Insert(&one)
		if err != nil {
			return err
		}
	} else {
		one.Status = status
		one.Export = export
		one.Cid = cid
		one.Sid = sid
		one.Slot = slot
		_, err = o.Update(&one)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateRozofsSetting(settingtype string, ip string, status bool) error {
	var ones []Device
	var export Device
	_, err := o.QueryTable("device").Filter("devtype", "storage").Filter("export", ip).Exclude("expand", 1).All(&ones) //decide update or not
	if err != nil {
		return err
	}

	for _, one := range ones {
		one.Status = status
		one.Export = ip
		one.Expand = status

		_, err = o.Update(&one)
		if err != nil {
			return err
		}
	}

	_, err = o.QueryTable("device").Filter("devtype", settingtype).Filter("ip", ip).All(&export) //decide update or not
	export.Status = status
	_, err = o.Update(&export)
	if err != nil {
		return err
	}

	return nil
}

func SelectUnexpanded(ip string, cid int) ([]Device, error) {
	//get all data
	ones := make([]Device, 0)
	_, err := o.QueryTable("device").Filter("export", ip).Filter("cid", cid).All(&ones) //decide update or not
	if err != nil {
		return ones, err
	}

	return ones, nil
}

func StopRozofsSevices(uuid string) error {
	var one Device
	if _, err := o.QueryTable("device").Filter("uuid", uuid).All(&one); err != nil {
		return err
	}
	one.Status = false
	one.Export = ""
	one.Cid = 0
	one.Sid = 0
	one.Slot = ""
	one.Expand = false

	if _, err := o.Update(&one); err != nil {
		return err
	}

	return nil
}

func StopRozofsSetting(stoptype string, ip string) error {
	var one Device
	if num, err := o.QueryTable("device").Filter("devtype", stoptype).Filter("ip", ip).All(&one); num > 0 { //decide update or not
		if err != nil {
			return err
		}
		one.Status = false
		one.Export = ""
		one.Cid = 0
		one.Sid = 0
		one.Slot = ""
		one.Expand = false

		_, err = o.Update(&one)
		if err != nil {
			return err
		}

	}

	if stoptype == "export" {

		var ones []Device
		_, err := o.QueryTable("device").Filter("devtype", "storage").Filter("export", ip).All(&ones) //decide update or not
		if err != nil {
			return err
		}
		fmt.Println("%+v", ones)
		for _, one := range ones {
			one.Expand = false
			_, err = o.Update(&one)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InsertCloudSetting(settingtype string, ip string, status bool) error {
	var one Setting

	num, err := o.QueryTable("setting").Filter("Settingtype", settingtype).Filter("Ip", ip).All(&one) //decide update or not
	if err != nil {
		return err
	}
	if num == 0 {
		one.Settingtype = settingtype
		one.Ip = ip
		one.Status = status
		_, err = o.Insert(&one)
		if err != nil {
			return err
		}
	} else {
		one.Status = status

		_, err = o.Update(&one)
		if err != nil {
			return err
		}
	}
	return nil
}

func OnlyCloudSetting(settingtype string, ip string) (Setting, error) {
	var one Setting
	return one, nil
}

func SelectCloudSetting() ([]Setting, error) {
	//get all data
	buks := make([]Setting, 0)
	if _, err := o.QueryTable("setting").All(&buks); err != nil {
		return buks, err
	}
	return buks, nil
}

func ClearCloudSetting(stoptype string, ip string) error {
	_, err := o.QueryTable("setting").Filter("Ip", ip).Filter("Settingtype", stoptype).Delete()
	if err != nil {
		return err
	}
	return nil
}

func SelectJournal() ([]LocalLog, error) {
	//get all data
	buks := make([]LocalLog, 0)
	if _, err := o.QueryTable("local_log").All(&buks); err != nil {
		return buks, err
	}
	return buks, nil
}

func InsertJournal(level string, message string, chinese_message string) error {
	var buk LocalLog
	buk.Created_at = time.Now()
	buk.Updated_at = time.Now()
	buk.Level = level
	buk.Message = message
	buk.Chinese_message = chinese_message
	if _, err := o.Insert(&buk); err != nil {
		return err
	}
	return nil
}

func ClearJournal() error {
	_, err := o.QueryTable("local_log").Filter("uid__isnull", false).Delete()

	if err != nil {
		return err
	}
	return nil
}
