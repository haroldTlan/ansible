package main

import (
	"fmt"
	s "strings"
)

func main() {
	//a := s.Split("192.168.2.1", ".")
	//b := s.Join(a, "")
	b := s.Join(s.Split("192.168.2.1", "."), "")

	fmt.Println(b)

}
