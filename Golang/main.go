package main

import (
	"errors"
	"fmt"
	"os"
	"snmpserver/cfg"

	"bitbucket.org/kardianos/service"

	lcsapi "snmpserver/license/api"
)

func register() error {
	c := cfg.Parse()
	if c.License != "" {
		return nil
	}

	sn, err := GetSerialNum()
	if err != nil {
		return err
	}

	if c.Server == nil {
		return errors.New("no server")
	}

	license, err := lcsapi.Register(c.Server.Ipaddr, c.Server.Port, sn)
	if err != nil {
		return err
	}

	c.License = license
	c.Server = nil
	c.Save()
	return nil
}

func main() {
	var name = "ZexaBox Disks Management Service"
	var displayName = "ZexaBox Disks Management"
	var desc = "This service provides support for Disks management."

	var s, err = service.NewService(name, displayName, desc)
	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = register()
			if err != nil {
				fmt.Printf("Failed to register: %s\n", err)
				return
			}

			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			Serve()

		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		go Serve()
		return nil
	}, func() error {
		// stop
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}
	return
}
