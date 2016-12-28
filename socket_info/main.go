package main

import (
	"cloud/logger"
	"cloud/topic"
	"github.com/googollee/go-socket.io"
	"net/http"
)

var statTopic = topic.New()
var ChanLogInfo chan Log

func main() {
	ChanLogInfo = make(chan Log, 1)

	LoggerChannel()
	ansible()
	InfoStat()
	socket()

}

func socket() {
	sio := socketio.NewSocketIOServer(&socketio.Config{})
	sio.Of("/statistics").On("connect", func(ns *socketio.NameSpace) {
		AddLogtoChan("connect", nil)
		go func(ns *socketio.NameSpace) {
			sub := statTopic.Subscribe()
			defer statTopic.Unsubscribe(sub)
			for {
				stat := <-sub
				err := ns.Emit("statistics", stat)
				if err != nil {
					AddLogtoChan("socket", err)
					return
				}
			}
		}(ns)
	})

	sio.Handle("/socket.io/", sio)
	http.ListenAndServe(":5000", sio)
}

func LoggerChannel() {
	go func() {
		for {
			select {
			case v := <-ChanLogInfo:
				logger.OutputLogger(v.Level, v.Message)
			default:
			}
		}
	}()
}
