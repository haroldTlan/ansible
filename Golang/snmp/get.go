// get
package snmp

import (
	"fmt"
	"github.com/cdevr/WapSNMP"
	"time"
)

func Get(target string, oidstr string) (string, error){
	community := "public"
	version := wapsnmp.SNMPv2c

	oid := wapsnmp.MustParseOid(oidstr)

	wsnmp, err := wapsnmp.NewWapSNMP(target, community, version, 2*time.Second, 5)
	defer wsnmp.Close()
	if err != nil {
		fmt.Printf("Error creating wsnmp => %v\n", wsnmp)
		return "", nil
	}

	val, err := wsnmp.Get(oid)
	fmt.Printf("Getting %v\n", oid)
	if err != nil {
		fmt.Printf("Get error => %v\n", err)
		return "", nil
	}
	fmt.Printf("Get(%v, %v, %v, %v) => %v\n", target, community, version, oid, val)

	out := fmt.Sprintf("%v", val)
	return out, err
}