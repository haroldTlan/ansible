package cfg

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
)

//<Config>
//	<Server ipaddr="127.0.0.1" port="8080"/>
//  <License>
//    abc
//	</License>
//</Config>

type Server struct {
	Ipaddr string `xml:"ipaddr,attr"`
	Port   int    `xml:"port,attr"`
}

type Config struct {
	XMLName xml.Name `xml:"Config"`
	Server  *Server  `xml:"Server,omitempty"`
	License string   `xml:"License,omitempty"`
}

func Parse() *Config {
	c := &Config{}
	path := filepath.Join(filepath.Dir(os.Args[0]), "config.xml")

	file, err := os.Open(path)
	if err != nil {
		return c
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return c
	}

	xml.Unmarshal(buf, c)
	return c
}

func (c *Config) Save() error {
	o, err := xml.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	path := filepath.Join(filepath.Dir(os.Args[0]), "config.xml")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(o)
	if err != nil {
		return err
	}

	return nil
}
