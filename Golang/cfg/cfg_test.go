package cfg

import (
	"encoding/xml"
	"testing"
)

//<config>
//	<server ipaddr="127.0.0.1" port="8008"/>
//  <license>
//    abc
//	</license>
//</config>

func TestUnmarshalConfig(t *testing.T) {
	data := `
		<Config>
			<Server ipaddr="127.0.0.1" port="8089"/>
			<License>e1dfa8a0fdafa</License>
		</Config>
	`
	var v Config
	err := xml.Unmarshal([]byte(data), &v)

	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if v.Server.Ipaddr != "127.0.0.1" || v.Server.Port != 8089 || v.License != "e1dfa8a0fdafa" {
		t.Fatalf("unmarshal error: %v", v)
	}
}

func TestUnmarshalConfigEmptyLicense(t *testing.T) {
	data := `
		<Config>
			<Server ipaddr="127.0.0.1" port="8089"/>
		</Config>
	`
	var v Config
	err := xml.Unmarshal([]byte(data), &v)

	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if v.License != "" {
		t.Fatalf("unmarshal error: %v", v)
	}
}

func TestMarshalConfig(t *testing.T) {
	var v Config
	v.Server = &Server{}
	v.Server.Ipaddr = "127.0.0.1"
	v.Server.Port = 8008
	v.License = "123456"

	o, err := xml.MarshalIndent(&v, "", "    ")
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var v1 Config
	err = xml.Unmarshal(o, &v1)

	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if v.Server.Ipaddr != v1.Server.Ipaddr || v.Server.Port != v1.Server.Port || v.License != v1.License {
		t.Fatalf("unmarshal error: %v", v)
	}
}

func TestMarshalConfigEmptyServer(t *testing.T) {
	var v Config
	v.License = "123456"

	o, err := xml.MarshalIndent(&v, "", "    ")
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var v1 Config
	err = xml.Unmarshal(o, &v1)

	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if v1.Server != nil {
		t.Fatalf("unmarshal error: %v", v)
	}
}
