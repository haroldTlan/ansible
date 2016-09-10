package web

import (
	"encoding/json"
	"net/http"
)

func NewResponse(status string, detail interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	o["status"] = status
	if detail != nil {
		switch v := detail.(type) {
		case error:
			o["errcode"] = 1
			o["description"] = v.Error()
		default:
			o["detail"] = v
		}
	}
	return o
}

func JsonResponse(f func(w http.ResponseWriter, r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		o, err := f(w, r)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			rsp := NewResponse("error", err)
			bytes, _ := json.Marshal(rsp)
			w.Write(bytes)
		} else {
			rsp := NewResponse("success", o)
			bytes, _ := json.Marshal(rsp)
			w.Write(bytes)
		}
	}
}
