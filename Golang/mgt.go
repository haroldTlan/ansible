package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"math/big"
	"net/http"
	"snmpserver/cfg"
	"snmpserver/web"
	"strconv"

	"net"
	"os/exec"
	"strings"
)

/*type LsiCmd struct {
	cli string
	enclId string
}*/

type CmdRes struct {
	Status string
	Info   string
}

type Session struct {
	Id int32 `json:"login_id"`
}

func FromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("bad number: " + base10)
	}
	return i
}

func Serve() {
	Initdb()
	c := cfg.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/api/version", web.JsonResponse(replyVersion)).Methods("GET")
	router.HandleFunc("/api/sn", web.JsonResponse(replySN)).Methods("GET")

	//if c.License != "" {
	if true {
		fmt.Println("step1\n")
		sn, err := GetSerialNum()
		_ = err //not used
		//if err == nil {
		if true {
			hash := md5.New()
			io.WriteString(hash, sn)
			hashed := hash.Sum(nil)

			var h crypto.Hash
			pubKey := &rsa.PublicKey{
				N: FromBase10("126038038516492034489881010707522756455005310820723628794048567491219653586876002712941473403005276243429681350407059668213363248724006391092540187693872519570891047411229657493659432418029829008660673664620025809544514419347167680091518538641680141780633312725341167771832755283446081256635145120586638842379"),
				E: 65537}

			var sig []byte
			//sig := make([]byte, len(c.License))
			_, err := fmt.Sscanf(c.License, "%x", &sig)
			_ = err //not used
			//if err == nil {
			if true {
				_ = pubKey //not used
				_ = h      //not used
				_ = hashed //not used
				//err := rsa.VerifyPKCS1v15(pubKey, h, hashed, sig)
				//if err == nil {
				if true {
					router.HandleFunc("/api/sessions", web.JsonResponse(createSession)).Methods("POST")
					router.HandleFunc("/api/ifaces", web.JsonResponse(getIfaces)).Methods("GET")
					router.HandleFunc("/api/systeminfo", web.JsonResponse(getSysteminfo)).Methods("GET")

					router.HandleFunc("/api/machines", web.JsonResponse(addMachines)).Methods("POST")
					router.HandleFunc("/api/machines", web.JsonResponse(getMachines)).Methods("GET")
					router.HandleFunc("/api/machines/{uuid}", web.JsonResponse(delMachines)).Methods("DELETE")
					router.HandleFunc("/api/machine/{uuid}", web.JsonResponse(restMachines)).Methods("GET")

					router.HandleFunc("/api/storeviews", web.JsonResponse(getAllviews)).Methods("GET")
					router.HandleFunc("/api/cloudset", web.JsonResponse(cloudSetting)).Methods("POST")
					router.HandleFunc("/api/cloudset", web.JsonResponse(getcloudSetting)).Methods("GET")
					router.HandleFunc("/api/cloudcheck", web.JsonResponse(cloudCheck)).Methods("POST")
					router.HandleFunc("/api/cloudstop", web.JsonResponse(cloudServiceStop)).Methods("POST")
					router.HandleFunc("/api/cloudtemp", web.JsonResponse(cloudTemp)).Methods("POST")
					router.HandleFunc("/api/journals", web.JsonResponse(getJournals)).Methods("GET")
					router.HandleFunc("/api/journalsdel", web.JsonResponse(delJournals)).Methods("POST")

					router.HandleFunc("/api/devices", web.JsonResponse(addDevices)).Methods("POST")
					router.HandleFunc("/api/devices", web.JsonResponse(getDevices)).Methods("GET")
					router.HandleFunc("/api/devices/{uuid}", web.JsonResponse(delDevices)).Methods("DELETE")
					router.HandleFunc("/api/export", web.JsonResponse(setExports)).Methods("POST")
					router.HandleFunc("/api/storage", web.JsonResponse(setStorages)).Methods("POST")
					router.HandleFunc("/api/client", web.JsonResponse(setClients)).Methods("POST")
					router.HandleFunc("/api/client/{uuid}", web.JsonResponse(delClients)).Methods("DELETE")
					router.HandleFunc("/api/rozostop", web.JsonResponse(rozoServiceStop)).Methods("POST")
					router.HandleFunc("/api/rozocheck", web.JsonResponse(rozoCheck)).Methods("POST")

					router.HandleFunc("/api/machines/{uuid}/disks", web.JsonResponse(getDisksOfMachine)).Methods("GET")
					router.HandleFunc("/api/machines/{uuid}/raids", web.JsonResponse(getRaidsOfMachine)).Methods("GET")
					router.HandleFunc("/api/machines/{uuid}/volumes", web.JsonResponse(getVolumesOfMachine)).Methods("GET")
					router.HandleFunc("/api/machines/{uuid}/filesystems", web.JsonResponse(getFilesystemsOfMachine)).Methods("GET")
					router.HandleFunc("/api/machines/{uuid}/initiators", web.JsonResponse(getInitiatorsOfMachine)).Methods("GET")

					fmt.Println("step2\n")
				}
			}
		}
	}

	TrapServer()
	Rundb()

	sio := NewSocketIOServer()
	sio.Handle("/", router)
	http.ListenAndServe(":8008", sio)

	Refreshing()
}

func addDevices(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ip := r.FormValue("ip")
	version := r.FormValue("version")
	devType := r.FormValue("devtype")
	size := r.FormValue("size")

	err := InsertDevice(ip, version, size, devType)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func getDevices(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	devices, err := SelectAllDevices()
	if err != nil {
		return devices, err
	}
	return devices, nil
}

func delDevices(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if err := DelDevice(uuid); err != nil {
		return nil, err
	}

	return nil, nil
}

func setExports(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ip := r.FormValue("ip")
	expands := r.FormValue("expand")

	var one Device
	if _, err := o.QueryTable("device").Filter("ip", ip).All(&one); err != nil {
		return nil, err
	}

	out, err := exec.Command("python", "/root/code/rozofs/rozofs.py", "--settingtype=export", fmt.Sprintf("--ip=%s", ip), fmt.Sprintf("--expand=%s", expands), fmt.Sprintf("--status=%s", one.Status)).Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))
	if results.Status == "True" {
		if err := InsertExports(ip, true); err != nil {
			return nil, err
		}
		expand := strings.Split(expands, ",")
		for i := 0; i < len(expand); i++ {
			InsertStorages(ip, true, expand[i])
		}
	} //TODO

	addLogtoChan(ip, "export", "set", err, results.Status == "True")

	return results, err
}

func setStorages(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ip := r.FormValue("ip")

	var one Device
	if _, err := o.QueryTable("device").Filter("devtype", "export").Filter("ip", ip).All(&one); err != nil {
		return nil, err
	}

	out, err := exec.Command("python", "/root/code/rozofs/rozofs.py", "--settingtype=storage", fmt.Sprintf("--ip=%s", ip)).Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))

	if results.Status == "True" {

		/*err = InsertRozofsSetting("storage", ip, results.Status == "True", export, cid, sid, slot)

		if err != nil {
			return string(out), err
		}*/
	}

	addLogtoChan(ip, "storage", "set", err, results.Status == "True")
	return results, err

}

func delStorages(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	err := StopRozofsSevices(uuid)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func setClients(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ip := r.FormValue("ip")

	out, err := exec.Command("python", "/root/code/rozofs/rozofs.py", "--settingtype=client", fmt.Sprintf("--ip=%s", ip), "--slot=None").Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))

	if results.Status == "True" {
		err = InsertRozofsSetting("client", ip, results.Status == "True", ip, 0, 0, "0")
		if err != nil {
			return string(out), err
		}

	} else {
		err = InsertRozofsSetting("client", ip, results.Status == "False", ip, 0, 0, "0")
		if err != nil {
			return string(out), err
		}

	}

	//addLogtoChan(ip, "export", "set", err, results.Status == "True")
	return nil, nil
}

func delClients(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	return uuid, nil
}

func rozoServiceStop(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	stopType := r.FormValue("stoptype")
	ip := r.FormValue("ip")

	out, err := exec.Command("python", "/root/code/rozofs/stop.py", fmt.Sprintf("--stoptype=%s", stopType), fmt.Sprintf("--ip=%s", ip)).Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))
	devNum := strings.Split(ip, ",")

	if results.Status == "True" {
		if len(devNum) > 1 {
			for _, dev := range devNum {
				err = StopRozofsSetting(stopType, dev)
			}
		} else {
			err = StopRozofsSetting(stopType, ip)
		}
		if err != nil {
			return string(out), err
		}

	}
	addLogtoChan(ip, stopType, "unset", err, results.Status == "True")

	return results, err

}

func rozoCheck(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	checkType := r.FormValue("checktype")
	ip := r.FormValue("ip")

	result, err := refreshRozoCheck(checkType, ip)

	//if results.Status == "True" {
	//err = InsertRozofsSetting(checkType, ip, results.Status == "True")
	//}

	return result, err
}

func addMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ip := r.FormValue("ip")
	slotnr, _ := strconv.Atoi(r.FormValue("slotnr"))
	devType := r.FormValue("devtype")

	err := InsertMachine(ip, slotnr, devType)
	if err != nil {
		return nil, err
	}
	RefreshStatAdd(ip)

	return nil, err
}

func getMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	machines, err := SelectAllMachines()
	if err != nil {
		return machines, err
	}
	return machines, nil
}

func delMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	if err := DelMachine(uuid); err != nil {
		return nil, err
	}

	return nil, nil
}

func restMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if err := RefreshStores(uuid); err != nil {
		fmt.Println(err)
		return nil, err
	}

	rest, err := restApi(uuid)
	if err != nil {
		return rest, err
	}

	return rest, nil
}

func getAllviews(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	views, err := RefreshAllViews()
	if err != nil {
		return nil, err
	}

	return views, nil
}

func getDisksOfMachine(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if err := RefreshReDisks(uuid); err != nil {
		return nil, err
	}

	redisks, err := resDisks(uuid)
	if err != nil {
		return redisks, err
	}
	return redisks, nil
}

func getRaidsOfMachine(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if err := RefreshReRaids(uuid); err != nil {
		return nil, err
	}

	reraids, err := resRaids(uuid)
	if err != nil {
		return reraids, err
	}
	return reraids, nil
}

func getVolumesOfMachine(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if err := RefreshReVolumes(uuid); err != nil {
		return nil, err
	}

	revolumes, err := resVols(uuid)
	if err != nil {
		return revolumes, err
	}
	return revolumes, nil
}

func getFilesystemsOfMachine(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if err := RefreshReFilesystems(uuid); err != nil {
		return nil, err
	}

	refs, err := resFs(uuid)
	if err != nil {
		return refs, err
	}
	return refs, nil
}

func getInitiatorsOfMachine(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if err := RefreshReInitiators(uuid); err != nil {
		return nil, err
	}
	RefreshReRaidVolumes(uuid)
	RefreshReInitVolumes(uuid)
	RefreshReNetInits(uuid)
	refs, err := resInits(uuid)
	if err != nil {
		return refs, err
	}
	return refs, nil
}

func getIfaces(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	info, _ := net.InterfaceAddrs()
	ifaces := make([]string, 0)
	for _, addr := range info {
		ifaces = append(ifaces, strings.Split(addr.String(), "/")[0])
	}

	//addLogtoChan("getIfaces", nil, false)
	return ifaces, nil
}

func getSysteminfo(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	feature := make([]string, 0)
	feature = append(feature, "xfs")

	systeminfo := make(map[string]interface{})
	systeminfo["gui version"] = "2.7.3"
	systeminfo["version"] = "2.2"
	systeminfo["feature"] = feature

	//addLogtoChan("getSysteminfo", nil, false)
	return systeminfo, nil
}

func createSession(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var sess Session
	sess.Id = 111

	return &sess, nil
}

func cloudSetting(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	settingtype := r.FormValue("settingtype")
	ip := r.FormValue("ip")

	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("python /root/code/new.py %s=%s", settingtype, ip)).Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))

	//if (settingtype != "worker") && (results.Status == "True") {
	err = InsertCloudSetting(settingtype, ip, results.Status == "True")

	if err != nil {
		return string(out), err
	}

	addLogtoChan(ip, settingtype, "set", err, results.Status == "True")
	return results, err
}

func getcloudSetting(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	settings, err := SelectCloudSetting()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

func cloudCheck(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	ip := r.FormValue("ip")
	checktype := r.FormValue("checktype")
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("python /root/code/check.py --ip=%s --checktype=%s", ip, checktype)).Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))
	if results.Status == "True" {
		err = InsertCloudSetting(checktype, ip, results.Status == "True")
	} else {
		err = InsertCloudSetting(checktype, ip, results.Status == "True")
	}
	//addLogtoChan(ip, checktype, "check", err, results.Status == "True")

	return results, err
}

func cloudServiceStop(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	stoptype := r.FormValue("stoptype")
	ip := r.FormValue("ip")
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("python /root/code/stop.py --stoptype=%s --ip=%s", stoptype, ip)).Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))
	if results.Status == "True" {
		err = ClearCloudSetting(stoptype, ip)
	}
	addLogtoChan(ip, stoptype, "unset", err, results.Status == "True")

	return results, err
}

func cloudTemp(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var result CmdRes
	stoptype := r.FormValue("stoptype")
	ip := r.FormValue("ip")

	err := ClearCloudSetting(stoptype, ip)
	if err != nil {
		result.Status = "False"
		return result, err
	}
	result.Status = "True"
	return result, err
}

func getJournals(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	journals, err := SelectJournal()
	if err != nil {
		return journals, err
	}
	return journals, nil
}

func delJournals(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var result CmdRes
	err := ClearJournal()
	if err != nil {
		result.Status = "False"
		return result, err
	}
	result.Status = "True"
	return result, err
}

func replyVersion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return map[string]string{"version": "1.0"}, nil
}

func replySN(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	sn, err := GetSerialNum()
	if err != nil {
		return nil, err
	} else {
		return map[string]string{"sn": sn}, nil
	}
}

func stupidCmd(out string) CmdRes {
	var result CmdRes

	pure := strings.Replace(string(out), "\n", "", -1)

	results := strings.Split(pure, "?")
	result.Status = results[1]
	result.Info = results[2]

	return result

}

func addLogtoChan(ip string, sertype string, logtype string, err error, result bool) {
	var isGUILog bool
	if isGUILog {
		//when goansible ,TODO
	}
	config := getLogConfig(logtype, result)

	if err == nil {
		_ = InsertJournal(config.Level, fmt.Sprintf("Server %s %s %s ", ip, logtype, sertype, config.Result), fmt.Sprintf("服务器 %s %s %s %s", ip, config.ChLogType, sertype, config.ChResult))

	} else {
		_ = InsertJournal(config.Level, fmt.Sprintf("%s %s", logtype, err), fmt.Sprintf("%s %s 操作错误", config.ChLogType, err))
	}
}
