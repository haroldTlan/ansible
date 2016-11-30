package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"hwraid/topic"
	"net/http"
)

var statTopic = topic.New()

func main() {
	ansible()
	InfoStat()
	temp()

}

func temp() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})

	sio.Of("/statistics").On("connect", func(ns *socketio.NameSpace) {
		fmt.Println("connect...")
		go func(ns *socketio.NameSpace) {
			sub := statTopic.Subscribe()
			defer statTopic.Unsubscribe(sub)
			for {
				stat := <-sub
				fmt.Printf("%+v\n", stat)
				err := ns.Emit("statistics", stat)
				if err != nil {
					return
				}
			}
		}(ns)
	})

	sio.Handle("/socket.io/", sio)

	http.ListenAndServe(":8014", sio)

}
