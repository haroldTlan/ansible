package main

import (
	//"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"net/http"
)

func main() {
	InfoStat()
	//ServeStat()
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
				fmt.Println(stat, sub)
				err := ns.Emit("statistics", stat)
				if err != nil {
					return
				}
			}
		}(ns)
	})

	sio.Handle("/socket.io/", sio)

	http.ListenAndServe(":5000", sio)

}
