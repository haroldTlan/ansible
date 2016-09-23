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

	return sio
}
