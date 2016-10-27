package main

import (
	"fmt"
	//"log"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var data = `
			blog: xiaorui.cc
			best_authors: ["fengyun","lee","park"]
			desc:
			  counter: 521
			    plist: [3, 4]
				`

type T struct {
	Blog    string
	Authors []string `yaml:"best_authors,flow"`
	Desc    struct {
		Counter int   `yaml:"Counter"`
		Plist   []int `yaml:",flow"`
	}
}

type S struct {
	Mysql      string `yaml:"mysql"`
	Mongo      string `yaml:"mongo"`
	Master     string `yaml:"master"`
	Worker     string `yaml:"worker"`
	Store      string `yaml:"store"`
	Key        string `yaml:"key"`
	Public     string `yaml:"public_service"`
	Inside     string `yaml:"inside_cloudstor"`
	Beanstalkd string `yaml:"beanstalkd"`
	Storeaim   string `yaml:"storeaim"`
	Push       string `yaml:"push"`
}

func main() {
	/*t := T{}
	//把yaml形式的字符串解析成struct类型
	err := yaml.Unmarshal([]byte(data), &t)
	//修改struct里面的记录
	t.Blog = "this is Blog"
	t.Authors = append(t.Authors, "myself")
	t.Desc.Counter = 99
	fmt.Printf("--- t:\n%v\n\n", t)
	//转换成yaml字符串类型
	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))*/

	path := "server.yml"
	str := read(path)
	var s S
	fmt.Println(str)
	yaml.Unmarshal([]byte(str), &s)

	s.Store = "192.168.2.102,192.168.2.103"

	e, _ := yaml.Marshal(&s)
	write(path, fmt.Sprintf("---\n%s\n", string(e)))

}

func read(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func write(path string, str string) {
	yaml := []byte(str)

	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	err = ioutil.WriteFile(path, yaml, 0666)
	if err != nil {
		panic(err)
	}

}
