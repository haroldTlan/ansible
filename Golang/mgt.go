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
					router.HandleFunc("/api/machines/{uuid}/disks", web.JsonResponse(getDisksOfMachine)).Methods("GET")
					router.HandleFunc("/api/machines", web.JsonResponse(addMachines)).Methods("POST")
					router.HandleFunc("/api/machines", web.JsonResponse(getMachines)).Methods("GET")
					router.HandleFunc("/api/ifaces", web.JsonResponse(getIfaces)).Methods("GET")
					router.HandleFunc("/api/systeminfo", web.JsonResponse(getSysteminfo)).Methods("GET")
					router.HandleFunc("/api/cloudset", web.JsonResponse(cloudSetting)).Methods("POST")
					router.HandleFunc("/api/cloudset", web.JsonResponse(getcloudSetting)).Methods("GET")
					router.HandleFunc("/api/cloudcheck", web.JsonResponse(cloudCheck)).Methods("POST")
					router.HandleFunc("/api/cloudstop", web.JsonResponse(cloudServiceStop)).Methods("POST")
					router.HandleFunc("/api/cloudtemp", web.JsonResponse(cloudTemp)).Methods("POST")
					router.HandleFunc("/api/journals", web.JsonResponse(getJournals)).Methods("GET")
					router.HandleFunc("/api/journalsdel", web.JsonResponse(delJournals)).Methods("POST")

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
}

func cloudSetting(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	settingtype := r.FormValue("settingtype")
	ip := r.FormValue("ip")

	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("python /root/code/new.py %s=%s", settingtype, ip)).Output()
	if err != nil {
		return string(out), err
	}

	results := stupidCmd(string(out))

	if (settingtype != "worker") && (results.Status == "True") {
		err = InsertCloudSetting(settingtype, ip, results.Status == "True")

		if err != nil {
			return string(out), err
		}
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
	addLogtoChan(ip, stoptype, "check", err, results.Status == "True")

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
	journals, err := Selectjournals()
	if err != nil {
		return journals, err
	}
	return journals, nil
}

func delJournals(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var result CmdRes
	err := ClearJournals()
	if err != nil {
		result.Status = "False"
		return result, err
	}
	result.Status = "True"
	return result, err
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

func addMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uuid := r.FormValue("uuid")
	ip := r.FormValue("ip")
	slotnr, _ := strconv.Atoi(r.FormValue("slotnr"))
	err := InsertMachine(uuid, ip, slotnr)
	return nil, err
}

func getMachines(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	machines, err := SelectAllMachines()
	if err != nil {
		return machines, err
	}
	return machines, nil
}

func getDisksOfMachine(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	machine, err := SelectMachine(uuid)
	if err != nil {
		return nil, nil
	}

	RefreshDisks(machine.Ip, uuid)
	//RefreshDisks("192.168.2.132", uuid)
	disks, _ := SelectDisksOfMachine(uuid)
	return disks, nil
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
		_ = InsertJournals(config.Level, fmt.Sprintf("Server %s %s %s ", ip, logtype, sertype, config.Result), fmt.Sprintf("服务器 %s %s %s %s", ip, config.ChLogType, sertype, config.ChResult))

	} else {
		_ = InsertJournals(config.Level, fmt.Sprintf("%s %s", logtype, err), fmt.Sprintf("%s %s 操作错误", config.ChLogType, err))
	}
}
