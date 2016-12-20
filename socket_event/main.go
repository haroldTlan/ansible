package main

import (
	"cloud/logger"
	"cloud/topic"
	"github.com/googollee/go-socket.io"
	"net/http"
)

var eventTopic = topic.New()
var ChanLogEvent chan Log

func main() {
	ChanLogEvent = make(chan Log, 1)

	loggerChannel()
	Initdb()
	NsqConsumerInit()
	socket()
}

func socket() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})
	sio.Of("/event").On("connect", func(ns *socketio.NameSpace) {
		AddLogtoChan("connect", nil)
		go func(ns *socketio.NameSpace) {
			sub := eventTopic.Subscribe()
			defer eventTopic.Unsubscribe(sub)
			for {
				e := <-sub
				err := ns.Emit("event", e)
				if err != nil {
					AddLogtoChan("socket", err)
					return
				}
			}
		}(ns)
	})

	sio.Handle("/socket.io/", sio)
	http.ListenAndServe(":8012", sio)
}

func loggerChannel() {
	go func() {
		for {
			select {
			case v := <-ChanLogEvent:
				logger.OutputLogger(v.Level, v.Message)
			default:
			}
		}
	}()
}
