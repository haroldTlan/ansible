package main

import (
	"cloud"
	"fmt"
	"time"
)

func RefreshOverViews(ip, event string) error {
	InsertJournals(event, ip)
	_, one, err := SelectMachine(ip)
	if err != nil {
		return err
	}
	if err := MulAttention(ip, event); err != nil {
		return err
	}

	if event == "ping.offline" {
		DelOutlineMachine(one.Uuid)

	} else if event == "ping.online" {
		time.Sleep(15 * time.Second)
		if err := RefreshStores(one.Uuid); err != nil {
			return err
		}

	} else {
		time.Sleep(15 * time.Second)
		if err := RefreshStores(one.Uuid); err != nil {
			return err
		}
		time.Sleep(4 * time.Second)
	}
	return nil
}

func MulAttention(ip, event string) error {
	ones := make([]Emergency, 0)
	if _, err := o.QueryTable("emergency").Filter("event", event).Filter("status", 0).All(&ones); err != nil || len(ones) < 1 {
		return err
	}
	RefreshMulAttention(ones[len(ones)-1].Uid, ip+" "+event)
	return nil
}

func RefreshMulAttention(uid int, message string) {
	go func() {
		_, err := SelectMulMails(uid, 1)
		if err != nil {
			AddLogtoChan("RefreshMulAttention_1", err)
		}
		//time.Sleep(time.Duration(ttl_1) * time.Second)
		SendMails(message, 1)

		status, err := SelectMulMails(uid, 2)
		if err != nil {
			AddLogtoChan("RefreshMulAttention_2", err)
		}
		if status {
			return
		} else {
			//	time.Sleep(time.Duration(ttl_2) * time.Second)
			SendMails(message, 2)
		}

		status, err = SelectMulMails(uid, 3)
		if err != nil {
			AddLogtoChan("RefreshMulAttention_3", err)
		}
		if status {
			return
		} else {
			//	time.Sleep(time.Duration(ttl_3) * time.Second)
			SendMails(message, 3)
			return
		}

	}()
}

func SendMails(message string, level int) {
	ones := make([]Mail, 0)
	mails := make([]string, 0)
	if _, err := o.QueryTable("mail").Filter("level", level).All(&ones); err != nil {
		AddLogtoChan("SendMails", err)
	}
	for _, val := range ones {
		mails = append(mails, val.Address)
	}
	cloud.Sendto(mails, message)
}

func AddLogtoChan(apiName string, err error) {
	var message string
	var log Log
	if err == nil {
		message = fmt.Sprintf("event success")
		log = Log{Level: "INFO", Message: message}
	} else {
		message = fmt.Sprintf("event %s, %s", apiName, err)
		log = Log{Level: "ERROR", Message: message}
	}

	ChanLogEvent <- log
	return
}
