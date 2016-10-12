package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

/*
CREATE TABLE `disk` (
 `uid` INT(10) AUTO_INCREMENT,
 `uuid` VARCHAR(64),
 `location` VARCHAR(64),
 `machineId` VARCHAR(64),
 `health` VARCHAR(64),
 `role` VARCHAR(64),
 `cap_sector` INT(11),
 `raid` VARCHAR(64),
 `vendor` VARCHAR(64),
 `model` VARCHAR(64),
 `sn` VARCHAR(64),
 `created` DATETIME DEFAULT NULL,
 PRIMARY KEY (`uid`)
);
CREATE TABLE `machine` (
	`uid` INT(10) AUTO_INCREMENT,
	`uuid` VARCHAR(64),
	`ip` VARCHAR(64),
	`slotnr` INT(10),
	`created` DATETIME DEFAULT NULL,
	PRIMARY KEY (`uid`)
);
*/

type Machine struct {
	Uid     int `orm:"pk"`
	Uuid    string
	Ip      string
	Slotnr  int
	Created time.Time `orm:"index"`
}

type Disks struct { //the table is remote disk
	Uuid      string `orm:"pk"`
	Health    string
	Role      string
	Location  string
	Raid      string
	CapSector int64
	Vendor    string
	Model     string
	Sn        string
}

type Disk struct { //the table is local disk
	//Uid int `orm:"pk"`
	Disks
	MachineId string `orm:"column(machineId)"`
}

type Raids struct { //the table'name is disk
	Uuid   string `orm:"pk"`
	Health string
	Level  string
	Name   string
	Cap    int64
	Used   int64 `orm:"column(used_cap)"`
}

type Raid struct { //the table is local raid
	Raids
	MachineId string `orm:"column(machineId)"`
}

type Volumes struct { //the table'name is disk
	Uuid   string `orm:"pk"`
	Health string
	Name   string
	Cap    int64
	Used   int64
	Type   string `orm:"column(owner_type)"`
}

type Volume struct { //the table is local vol
	Volumes
	MachineId string `orm:"column(machineId)"`
}

type Xfs struct { //the table'name is xfs
	Uuid   string `orm:"pk"`
	Volume string
	Name   string
	Chunk  string `orm:"column(chunk_kb)"`
	Type   string
}

type Filesystems struct { //the table is local fs
	Xfs
	MachineId string `orm:"column(machineId)"`
}

type Initiators struct { //the table'name is initiators
	Wwn    string `orm:"pk"`
	Target string `orm:"column(target_wwn)"`
}

type Initiator struct { //the table is local initiator
	Initiators
	MachineId string `orm:"column(machineId)"`
}

type Setting struct {
	Uid         int    `orm:"pk"`
	Settingtype string `orm:"column(settingtype)"`
	Ip          string
	Status      bool
}

type Journals struct {
	Uid             int       `orm:"pk"`
	Created_at      time.Time `orm:"index"`
	Updated_at      time.Time `orm:"index"`
	Level           string
	Message         string
	Chinese_message string
}

type RaidVolumes struct {
	Id     int
	Raid   string
	Volume string
	Type   string
}

type RaidVolume struct {
	RaidVolumes
	MachineId string `orm:"column(machineId)"`
}

type InitiatorVolumes struct {
	Id        int
	Initiator string
	Volume    string
}

type InitiatorVolume struct {
	InitiatorVolumes
	MachineId string `orm:"column(machineId)"`
}

type NetworkInitiators struct {
	Id        int
	Initiator string
	Eth       string
	Port      int
}

type NetworkInitiator struct {
	NetworkInitiators
	MachineId string `orm:"column(machineId)"`
}

var o orm.Ormer

func Initdb() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:passwd@/speediodb?charset=utf8", 30)
	orm.RegisterModel(new(Machine), new(Disks), new(Disk), new(Raid), new(Raids), new(Volume), new(Volumes), new(Filesystems), new(Xfs), new(Initiator), new(Initiators), new(Setting), new(Journals), new(RaidVolumes), new(RaidVolume), new(InitiatorVolumes), new(InitiatorVolume), new(NetworkInitiators), new(NetworkInitiator))

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
			orm.RegisterDataBase(fmt.Sprintf("%s", name), "mysql", fmt.Sprintf("root:passwd@tcp(%s:3306)/speediodb?charset=utf8", ip), 30)
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
	err := o.Using(name)

	if err != nil {
		return err

	}
	err = InsertRemoteDisks(machine.Uuid)

	if err != nil {
		return err
	}
	return nil
}

func MachineType(machine Machine) (string, string) {
	ip := machine.Ip
	temp := strings.Split(ip, ".")
	name := "remote" + temp[2] + temp[3]

	return name, ip
}

func InitLocal(dev []Machine) {
	for i := 0; i < len(dev); i++ {
		name, _ := MachineType(dev[i])
		err := o.Using(name)

		if err != nil {
			fmt.Println(err)
		}
		InsertRemoteDisks(dev[i].Uuid)
	}
}

func InsertMachine(uuid string, ip string, slotnr int) error {
	orm.Debug = true

	var one Machine
	one.Uuid = uuid
	one.Ip = ip
	one.Slotnr = slotnr
	one.Created = time.Now()
	if _, err := o.Insert(&one); err != nil {
		return err
	}
	/*err := InitSingleRemote(one)
	if err != nil {
		return err
	}*/

	RefreshViews(uuid)

	return nil
}

func SelectAllMachines() ([]Machine, error) {
	//get all machine
	var ones []Machine
	if _, err := o.QueryTable("machine").All(&ones); err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectMachine(uuid string) (Machine, error) {
	var one Machine
	return one, nil
}

func DeleteMachine(uuid string) error {
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
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).Delete(); err != nil {
		fmt.Println(err)
		fmt.Println("db")
		return err
	}

	return nil
}

func RefreshAllViews() {}

func RefreshViews(uuid string) error {
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
	return nil
}

func InsertRemoteDisks(id string) error {
	var ones []Disks
	_, err := o.QueryTable("disks").All(&ones)
	if err != nil {
		fmt.Println(err)
	}

	if mlen := len(ones); mlen > 0 {
		err = o.Using("default")
		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < mlen; i++ {
			one := Disk{Disks: ones[i], MachineId: id}
			if _, err = o.Insert(&one); err != nil {
				return err
			}

		}
	}
	return nil
}

func SelectDisksOfMachine(uuid string) ([]Disk, error) {
	var ones []Disk
	_, err := o.QueryTable("disk").Filter("machineId__exact", uuid).Exclude("location__exact", "").All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectDisks() ([]Disk, int64, error) {
	var ones []Disk
	num, err := o.QueryTable("disk").Exclude("location__exact", "").Exclude("location__exact", "").All(&ones)
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
			num, err := o.QueryTable("disk").Filter("uuid__exact", redisks[i].Uuid).All(&loc) //decide update or not
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var ones []Disks
	_, err = o.QueryTable("disks").All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	InsertDisksOfMachine(ones, uuid)

	return nil
}

func SelectVolumesOfMachine(uuid string) ([]Volume, error) {
	var ones []Volume
	_, err := o.QueryTable("volume").Filter("machineId__exact", uuid).Filter("used__exact", 1).All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectVolumes() ([]Volume, int64, error) {
	var ones []Volume
	num, err := o.QueryTable("volume").Filter("used__exact", 1).All(&ones)
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var ones []Volumes
	_, err = o.QueryTable("volumes").All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	InsertVolumesOfMachine(ones, uuid)

	return nil
}

func SelectRaidsOfMachine(uuid string) ([]Raid, error) {
	var ones []Raid
	_, err := o.QueryTable("raid").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectRaids() ([]Raid, int64, error) {
	var ones []Raid
	num, err := o.QueryTable("raid").All(&ones)
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var ones []Raids
	_, err = o.QueryTable("raids").All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	InsertRaidsOfMachine(ones, uuid)

	return nil
}

func SelectFilesystemsOfMachine(uuid string) ([]Filesystems, error) {
	var ones []Filesystems
	_, err := o.QueryTable("filesystems").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectFilesystems() ([]Filesystems, int64, error) {
	var ones []Filesystems
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var ones []Xfs
	_, err = o.QueryTable("xfs").All(&ones)
	if err != nil {
		return err
	}
	InsertFilesystemsOfMachine(ones, uuid)

	return nil
}

func SelectInitiatorsOfMachine(uuid string) ([]Initiator, error) {
	var ones []Initiator
	_, err := o.QueryTable("initiator").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectInitiators() ([]Initiator, int64, error) {
	var ones []Initiator
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		return err
	}

	var ones []Initiators
	_, err = o.QueryTable("initiators").All(&ones)
	if err != nil {
		return err
	}
	InsertInitiatorsOfMachine(ones, uuid)

	return nil
}

func SelectRaidVolumesOfMachine(uuid string) ([]RaidVolume, error) {
	var ones []RaidVolume
	_, err := o.QueryTable("raid_volume").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectRaidVolumes() ([]RaidVolume, int64, error) {
	var ones []RaidVolume
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
			num, err := o.QueryTable("raid_volume").Filter("raid__exact", remote[i].Raid).All(&loc) //decide update or not
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var ones []RaidVolumes
	_, err = o.QueryTable("raid_volumes").All(&ones)
	if err != nil {
		return err
	}
	InsertRaidVolumesOfMachine(ones, uuid)

	return nil
}

func SelectInitVolumesOfMachine(uuid string) ([]InitiatorVolume, error) {
	var ones []InitiatorVolume
	_, err := o.QueryTable("initiator_volume").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectInitVolumes() ([]InitiatorVolume, int64, error) {
	var ones []InitiatorVolume
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
			num, err := o.QueryTable("initiator_volume").Filter("initiator__exact", remote[i].Initiator).All(&loc) //decide update or not
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var ones []InitiatorVolumes
	_, err = o.QueryTable("initiator_volumes").All(&ones)
	if err != nil {
		return err
	}
	InsertInitVolumesOfMachine(ones, uuid)

	return nil
}

func SelectNetInitsOfMachine(uuid string) ([]NetworkInitiator, error) {
	var ones []NetworkInitiator
	_, err := o.QueryTable("network_initiator").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectNetInits() ([]NetworkInitiator, int64, error) {
	var ones []NetworkInitiator
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
			num, err := o.QueryTable("network_initiator").Filter("initiator__exact", remote[i].Initiator).All(&loc) //decide update or not
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
	name := "remote" + strings.Split(uuid, "192168")[1]

	err := o.Using(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var ones []NetworkInitiators
	_, err = o.QueryTable("network_initiators").All(&ones)
	if err != nil {
		return err
	}
	InsertNetInitsOfMachine(ones, uuid)

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
	var buks []Setting
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

func Selectjournals() ([]Journals, error) {
	//get all data
	var buks []Journals
	if _, err := o.QueryTable("journals").All(&buks); err != nil {
		return buks, err
	}
	return buks, nil
}

func InsertJournals(level string, message string, chinese_message string) error {
	var buk Journals
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

func ClearJournals() error {
	_, err := o.QueryTable("journals").Filter("uid__isnull", false).Delete()

	if err != nil {
		return err
	}
	return nil
}

func RemoteRaids() ([]Raids, error) {
	var ones []Raids
	fmt.Println(ones)
	return ones, nil

}

func RemoteVolumes() ([]Volumes, error) {
	var ones []Volumes
	fmt.Println(ones)
	return ones, nil

}

func RemoteFs() ([]Initiators, error) {
	var ones []Initiators

	fmt.Println(ones)
	return ones, nil

}

func SelectAllDisks() ([]Disk, error) {
	//get all data
	var ones []Disk
	return ones, nil
}

func SelectDisk(uuid string) (Disk, error) {
	var one Disk
	return one, nil
}

/*
func UpdateDisk(uuid string, location string, machineId string, status string, role string, raid string, size string) error {
	// //update data
	saveone, _ := SelectDisk(uuid)
	saveone.Uuid = uuid
	saveone.Location = location
	saveone.MachineId = machineId
	saveone.Created = time.Now()
	saveone.Status = status
	saveone.Role = role
	saveone.Raid = raid
	saveone.Size = size
	if err := orm.Save(&saveone); err != nil {
		return err
	}
	return nil
}*/

func DeleteDisk(uuid string) error {
	// // //delete one data
	fmt.Printf("delete disk finished\n")
	return nil
}

func DeleteAllDisks() error {
	// //delete all data

	return nil
}

func UpdateMachine(uuid string, ip string, slotnr int) error {
	// //update data
	saveone, _ := SelectMachine(uuid)
	saveone.Uuid = uuid
	saveone.Ip = ip
	saveone.Slotnr = slotnr
	saveone.Created = time.Now()
	return nil
}
