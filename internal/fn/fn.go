package fn

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Execute(command string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(stdout.String())
	}
}
