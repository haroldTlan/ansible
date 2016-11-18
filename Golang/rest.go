package main

import (
	"fmt"
	"time"
)

type ResInitiators struct {
	Portals   []string `json:"portals"`
	Wwn       string   `json:"wwn"`
	Id        string   `json:"id"`
	Volumes   []string `json:"volumes"`
	MachineId string   `json:"machineId"`
	Ip        string   `json:"ip"`
}

type ResDisks struct {
	Uuid      string `json:"id"`
	Health    string `json:"health"`
	Role      string `json:"role"`
	Location  string `json:"location"`
	Raid      string `json:"raid"`
	CapSector int64  `json:"cap"`
	Vendor    string `json:"vendor"`
	Model     string `json:"model"`
	Sn        string `json:"sn"`
	Ip        string `json:"ip"`
}

type ResRaids struct {
	Uuid      string  `json:"id"`
	Health    string  `json:"health"`
	Level     string  `json:"level"`
	Name      string  `json:"name"`
	Cap       int64   `json:"cap"`
	Used      int64   `json:"used"`
	CapMb     float64 `json:"cap_mb"`
	UsedMb    float64 `json:"used_mb"`
	Ip        string  `json:"ip"`
	MachineId string  `json:"machineId"`
}

type ResVols struct {
	Uuid      string  `json:"id"`
	Health    string  `json:"health"`
	Name      string  `json:"name"`
	Cap       int64   `json:"cap"`
	CapMb     float64 `json:"cap_mb"`
	Type      string  `json:"type"`
	Owner     string  `json:"owner"`
	Deleted   int     `json:"deleted"`
	Ip        string  `json:"ip"`
	MachineId string  `json:"machineId"`
}

type ResFilesystems struct { //the table'name is xfs
	Uuid      string `json:"id""`
	Volume    string `json:"volume"`
	Name      string `json:"name"`
	Chunk     string `json:"chunk"`
	Type      string `json:"type"`
	Ip        string `json:"ip"`
	MachineId string `json:"machineId"`
}

type ResJournals struct {
	Message   string    `json:"message"`
	Created   time.Time `orm:"index" json:"created"`
	Unix      int64     `json:"unix"`
	Level     string    `json:"level"`
	Ip        string    `json:"ip"`
	MachineId string    `json:"machineId"`
	Devtype   string    `json:"devtype"`
}

func resDisks(uuid string) ([]ResDisks, error) {
	disks := make([]ResDisks, 0)

	var one Machine
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return nil, err
	}

	ones := make([]Disk, 0)
	_, err := o.QueryTable("disk").Filter("machineId__exact", uuid).Exclude("location__exact", "").All(&ones)
	if err != nil {
		return disks, err
	}

	for _, val := range ones {
		var disk ResDisks
		disk.Uuid = val.Uuid
		disk.Health = val.Health
		disk.Role = val.Role
		disk.Location = val.Location
		disk.Raid = val.Raid
		disk.CapSector = val.CapSector
		disk.Vendor = val.Vendor
		disk.Model = val.Model
		disk.Sn = val.Sn
		disk.Ip = one.Ip
		disks = append(disks, disk)
	}

	return disks, nil
}

func resRaids(uuid string) ([]ResRaids, error) {
	raids := make([]ResRaids, 0)

	var one Machine
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return nil, err
	}

	ones := make([]Raid, 0)
	_, err := o.QueryTable("raid").Filter("machineId__exact", uuid).Filter("deleted__exact", 0).All(&ones)
	if err != nil {
		return raids, err
	}

	for _, val := range ones {
		var raid ResRaids
		raid.Uuid = val.Uuid
		raid.Health = val.Health
		raid.Level = val.Level
		raid.Name = val.Name
		raid.Cap = int64(val.Cap) * 1024 * 1024
		raid.Used = int64(val.Used) * 1024 * 1024
		raid.CapMb = float64(val.Cap * 1024)
		raid.UsedMb = float64(val.Used) * 1024
		raid.MachineId = val.MachineId
		raid.Ip = one.Ip

		raids = append(raids, raid)
	}

	return raids, nil
}

func resVols(uuid string) ([]ResVols, error) {
	vols := make([]ResVols, 0)

	var one Machine
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return nil, err
	}

	ones := make([]Volume, 0)
	_, err := o.QueryTable("volume").Filter("machineId__exact", uuid).Filter("deleted__exact", 0).All(&ones)
	if err != nil {
		return vols, err
	}

	raidVols := make([]RaidVolume, 0)
	if _, err := o.QueryTable("raid_volume").Filter("machineId__exact", uuid).All(&raidVols); err != nil {
		return vols, err
	}

	for _, val := range ones {
		var vol ResVols
		vol.Uuid = val.Uuid
		vol.Health = val.Health
		vol.Name = val.Name
		vol.Cap = int64(val.Cap) * 1024 * 1024
		vol.CapMb = float64(val.Cap) * 1024
		vol.Type = val.Type
		vol.Deleted = val.Deleted
		vol.Ip = one.Ip
		vol.MachineId = val.MachineId

		for _, raidVol := range raidVols {
			if raidVol.RaidVolumes.Volume == val.Uuid {
				var raid Raid
				if _, err := o.QueryTable("raid").Filter("uuid", raidVol.RaidVolumes.Raid).All(&raid); err != nil {
					return nil, err
				}
				vol.Owner = raid.Raids.Name

			}
		}

		vols = append(vols, vol)
	}

	return vols, nil
}

func resFs(uuid string) ([]ResFilesystems, error) {
	fs := make([]ResFilesystems, 0)

	var one Machine
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return nil, err
	}

	ones := make([]Filesystems, 0)
	_, err := o.QueryTable("filesystems").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return fs, err
	}

	for _, val := range ones {
		var f ResFilesystems

		f.Uuid = val.Uuid
		f.Name = val.Name
		f.Chunk = val.Chunk
		f.Type = val.Type

		var vols Volume
		if _, err := o.QueryTable("volume").Filter("uuid", val.Volume).All(&vols); err != nil {
			return nil, err
		}
		f.Volume = vols.Volumes.Name

		f.Ip = one.Ip
		f.MachineId = uuid

		fs = append(fs, f)
	}

	return fs, nil
}

func resInits(uuid string) ([]ResInitiators, error) {
	inits := make([]ResInitiators, 0)

	var one Machine
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return nil, err
	}

	Nets := make([]NetworkInitiator, 0)
	if _, err := o.QueryTable("network_initiator").Filter("machineId__exact", uuid).All(&Nets); err != nil {
		return inits, err
	}

	Inits := make([]Initiator, 0)
	if _, err := o.QueryTable("initiator").Filter("machineId__exact", uuid).All(&Inits); err != nil {
		return inits, err
	}

	InitVols := make([]InitiatorVolume, 0)
	if _, err := o.QueryTable("initiator_volume").Filter("machineId__exact", uuid).All(&InitVols); err != nil {
		return inits, err
	}

	for _, val := range Inits {
		var slice ResInitiators
		vol_name := make([]string, 0)
		vol_port := make([]string, 0)

		for _, initvol := range InitVols {
			var vols []Volume

			if initvol.InitiatorVolumes.Initiator == val.Initiators.Wwn {
				if _, err := o.QueryTable("volume").Filter("uuid", initvol.InitiatorVolumes.Volume).All(&vols); err != nil {
					return nil, err
				}

				for _, name := range vols {
					vol_name = append(vol_name, name.Volumes.Name)
				}
			}
		}

		for _, net := range Nets {
			if net.NetworkInitiators.Initiator == val.Initiators.Wwn {
				vol_port = append(vol_port, net.NetworkInitiators.Eth)
			}

		}

		slice.Volumes = vol_name
		slice.Portals = vol_port

		slice.Wwn = val.Wwn
		slice.Id = val.Wwn
		slice.MachineId = uuid

		slice.Ip = one.Ip
		inits = append(inits, slice)
	}
	fmt.Println(inits)
	return inits, nil
}

func resJournals(uuid string) ([]ResJournals, error) {
	jours := make([]ResJournals, 0)

	var one Machine
	if _, err := o.QueryTable("machine").Filter("uuid", uuid).All(&one); err != nil {
		return nil, err
	}

	ones := make([]Journal, 0)
	_, err := o.QueryTable("Journal").Filter("machineId__exact", uuid).All(&ones)
	if err != nil {
		return jours, err
	}

	for _, val := range ones {
		var jour ResJournals
		jour.Message = val.Message
		jour.Created = val.Created_at
		jour.Unix = val.Created_at.Unix()
		jour.Level = val.Level
		jour.MachineId = uuid
		jour.Ip = one.Ip
		jours = append(jours, jour)
	}

	return jours, nil
}
