package main

import (
	"cloud"
	"fmt"
	"runtime"
	"time"
)

func RefreshOverViews(ip, event string) error {
	if err := InsertJournals(event, ip); err != nil {
		return err
	}

	_, one, err := SelectMachine(ip)
	if err != nil {
		return err
	}
	if err := MulAttention(ip, event); err != nil { //mul email sending
		return err
	}

	switch event {

	case "ping.offline":
		if err := DelOutlineMachine(one.Uuid); err != nil {
			return err
		}

	default:
		time.Sleep(19 * time.Second)
		if err := RefreshStores(one.Uuid); err != nil {
			return err
		}
	}

	/*
		if event == "ping.offline" {
			DelOutlineMachine(one.Uuid)

				} else if event == "ping.online" {
				time.Sleep(15 * time.Second)
				if err := RefreshStores(one.Uuid); err != nil {
					AddLogtoChan(err)
					return err
				}

		} else {
			time.Sleep(15 * time.Second)
			if err := RefreshStores(one.Uuid); err != nil {
				AddLogtoChan(err)
				return err
			}
			time.Sleep(4 * time.Second)
		}*/
	return nil
}

func MulAttention(ip, event string) error {
	ones := make([]Emergency, 0)
	if _, err := o.QueryTable("emergency").Filter("event", event).Filter("status", 0).All(&ones); err != nil || len(ones) < 1 {
		return err
	}
	_, message := messageTransform(event)
	RefreshMulAttention(ones[len(ones)-1].Uid, ip+" "+message) //the lastest attention!
	return nil
}

func RefreshMulAttention(uid int, message string) {
	go func() {
		if _, err := SelectMulMails(uid, 1); err != nil {
			AddLogtoChan(err)
		}
		SendMails(message, 1)

		status, err := SelectMulMails(uid, 2)
		if err != nil {
			AddLogtoChan(err)
		}
		if status {
			return
		} else {
			SendMails(message, 2)
		}

		status, err = SelectMulMails(uid, 3)
		if err != nil {
			AddLogtoChan(err)
		}
		if status {
			return
		} else {
			SendMails(message, 3)
			return
		}
	}()
}

func SendMails(message string, level int) {
	ones := make([]Mail, 0)
	mails := make([]string, 0)
	if _, err := o.QueryTable("mail").Filter("level", level).All(&ones); err != nil {
		AddLogtoChan(err)
	}
	for _, val := range ones {
		mails = append(mails, val.Address)
	}
	cloud.Sendto(mails, message)
}

func AddLogtoChan(err error) {
	var message string
	var log Log
	if err == nil {
		message = fmt.Sprintf("[EVENT]event success")
		log = Log{Level: "INFO", Message: message}
	} else {
		pc, fn, line, _ := runtime.Caller(1)
		message = fmt.Sprintf("[EVENT][%s %s:%d] %s, %s", runtime.FuncForPC(pc).Name(), fn, line, err)
		log = Log{Level: "ERROR", Message: message}
	}

	ChanLogEvent <- log
	return
}
