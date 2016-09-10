package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

func GetSerialNum() (string, error) {
	o, err := exec.Command("cmd", "/c", "wmic diskdrive get serialnumber").Output()
	if err != nil {
		return "", err
	}

	rSN := regexp.MustCompile(`SerialNumber\s*(\w+)`)
	match := rSN.FindStringSubmatch(string(o))
	if match == nil {
		return "", errors.New("no match sn")
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(match[1]))), nil
}
