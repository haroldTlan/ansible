package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

var o orm.Ormer

func Initdb() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:passwd@/speediodb?charset=utf8", 30)
	orm.RegisterModel(new(Machine), new(Disks), new(Disk), new(Raid), new(Raids), new(Volume), new(Volumes), new(Filesystems), new(Xfs), new(Initiator), new(Initiators), new(Emergency), new(RaidVolumes), new(RaidVolume), new(InitiatorVolumes), new(InitiatorVolume), new(NetworkInitiators), new(NetworkInitiator), new(Mail), new(Journal))

	o = orm.NewOrm()
	InitRemote()
}

func InitRemote() error {
	machines := make([]Machine, 0)
	if _, err := o.QueryTable("machine").All(&machines); err != nil {
		return err
	}
	if mlen := len(machines); mlen > 0 {
		for i := 0; i < mlen; i++ {
			name, ip := MachineType(machines[i])
			err := orm.RegisterDataBase(fmt.Sprintf("%s", name), "mysql", fmt.Sprintf("root:passwd@tcp(%s:3306)/speediodb?charset=utf8", ip), 30)
			if err != nil {
				//DelMachine(machines[i].Uuid)
				AddLogtoChan("InitRemote failed", err)
				return err
			}
		}
	} else {
		//TODO!!!!!
	}
	return nil
}

func InitSingleRemote(ip string) error {
	name := strings.Join(strings.Split(ip, "."), "")
	fmt.Println(ip, name)
	orm.RegisterDataBase(fmt.Sprintf("%s", name), "mysql", fmt.Sprintf("root:passwd@tcp(%s:3306)/speediodb?charset=utf8", ip), 30)
	return nil
}

func MachineType(machine Machine) (string, string) {
	ip := machine.Ip
	tempIp := strings.Join(strings.Split(ip, "."), "")
	name := "remote" + tempIp

	return name, ip
}

func SelectMachine(ip string) (int64, Machine, error) {
	var one Machine
	num, err := o.QueryTable("machine").Filter("ip", ip).All(&one)
	if err != nil {
		return 0, one, err
	}
	return num, one, nil
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
	_, err = o.QueryTable("disks").Exclude("location__exact", "").All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}

	InsertDisksOfMachine(ones, uuid)

	return nil
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
	_, err = o.QueryTable("raids").Filter("deleted__exact", 0).All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	InsertRaidsOfMachine(ones, uuid)

	return nil
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
	_, err = o.QueryTable("volumes").Filter("deleted__exact", 0).All(&ones)
	if err != nil {
		fmt.Println(err)
		return err
	}
	InsertVolumesOfMachine(ones, uuid)

	return nil
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

func SelectMulMails(uid int, level int) (bool, error) {
	var one Mail
	var two Emergency
	if _, err := o.QueryTable("mail").Filter("level", level).All(&one); err != nil {
		return two.Status, err
	}
	time.Sleep(time.Duration(one.Ttl) * time.Second)

	if _, err := o.QueryTable("emergency").Filter("uid", uid).All(&two); err != nil {
		return two.Status, err
	}
	return two.Status, nil
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
	one.ChineseMessage = "设备" + " " + machine + " " + chinese_message
	one.Event = event
	one.Ip = machine
	one.Created_at = time.Now()
	one.Updated_at = time.Now()
	if _, err := o.Insert(&one); err != nil {
		AddLogtoChan("InsertJournals", err)
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
		message := "raid degraded"
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
	case "raid.normal":
		message := "raid normal"
		chinese_message := "阵列恢复正常"
		return message, chinese_message
	case "volume.normal":
		message := "volume normal"
		chinese_message := "虚拟磁盘恢复正常"
		return message, chinese_message
	default:
		return "", "未知"
	}
}

func RefreshStores(uuid string) error {
	var one Machine

	err := o.Using("default")
	if err != nil {
		return err
	}

	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return err
	}
	if one.Devtype == "export" {
		return nil
	}

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
	if err := RefreshReRaidVolumes(uuid); err != nil {
		return err
	}
	if err := RefreshReInitVolumes(uuid); err != nil {
		return err
	}
	if err := RefreshReNetInits(uuid); err != nil {
		return err
	}

	return nil
}

func DelOutlineMachine(uuid string) error {
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
	return nil
}
