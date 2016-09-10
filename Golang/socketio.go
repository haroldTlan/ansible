package main

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	
	"fmt"
)

func NewSocketIOServer() *socketio.SocketIOServer {
	sio := socketio.NewSocketIOServer(&socketio.Config{})

	sio.Of("/diskevent").On("connect", func(ns *socketio.NameSpace) {
		go func(ns *socketio.NameSpace) {
			sub := trapTopic.Subscribe()
			defer trapTopic.Unsubscribe(sub)
			for {
				
				e := <-sub

				bytes, err := json.Marshal(e)
				if err != nil {
					continue
				}

				err = ns.Emit("diskevent", string(bytes))
				//err := ns.Emit("event", e)
				fmt.Println(string(bytes), bytes, sub)
				if err != nil {
					return
				}
			}
		}(ns)
	})

	sio.Of("/statistics").On("connect", func(ns *socketio.NameSpace) {
		go func(ns *socketio.NameSpace) {
			sub := statTopic.Subscribe()
			defer statTopic.Unsubscribe(sub)
			for {
				stat := <-sub
				/*m := map[string]string{"status": "ok", "sample": "haha"}
				bytes, err := json.Marshal(m)
				if err != nil {
					continue
				}

				err = ns.Emit("statistics", bytes)*/
				//m := map[string]interface{}{"status": "ok", "sample": stat}

				//ms,err := json.Marshal(m)
				
				fmt.Println(stat, sub)
				err := ns.Emit("statistics", stat)
				if err != nil {
					return
				}/**/
			}
		}(ns)
	})

	return sio
}
