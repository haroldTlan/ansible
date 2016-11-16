package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("sh", "-c", `rozo node status -E 192.168.2.107`)
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		fmt.Printf("Run returns: %s\n", err)
	}
	fmt.Printf("Stderr: %s\n", string(w.Bytes()))
}
