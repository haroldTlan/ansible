package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Register(svrIp string, port int, sn string) (string, error) {
	rsp, err := http.PostForm(fmt.Sprintf("http://%v:%v/api/license", svrIp, port), url.Values{"sn": {sn}})
	if err != nil {
		return "", nil
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	type Signature struct {
		Signature string
	}
	type Response struct {
		Status string
		Detail Signature
	}

	var o Response
	err = json.Unmarshal(body, &o)
	if err != nil {
		return "", err
	}

	if o.Status != "success" {
		return "", errors.New("response error")
	}

	return o.Detail.Signature, nil
}
