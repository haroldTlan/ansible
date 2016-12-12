package main

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"hwraid/topic"
	"net/http"
)

var statTopic = topic.New()

func main() {
	Initdb()
	NsqConsumerInit()
	socket()

}

func socket() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})

	sio.Of("/event").On("connect", func(ns *socketio.NameSpace) {
		fmt.Println("\nevent...\n")
		go func(ns *socketio.NameSpace) {
			sub := eventTopic.Subscribe()
			defer eventTopic.Unsubscribe(sub)
			for {
				e := <-sub
				/*			bytes, err := json.Marshal(e)
							if err != nil {
								continue
							}*/
				/*			var dat map[string]interface{}
							if err := json.Unmarshal([]byte(e), &dat); err != nil {
								panic(err)
							}
							fmt.Printf("%+v\n", dat)*/

				//			err = ns.Emit("diskevent", string(bytes))
				err := ns.Emit("event", e)
				if err != nil {
					return
				}
			}
		}(ns)
	})
	fmt.Println("socket")
	sio.Handle("/socket.io/", sio)

	http.ListenAndServe(":8012", sio)

}
