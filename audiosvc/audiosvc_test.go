package audiosvc_test

import (
	"testing"
	"os/exec"
	"bytes"
	"fmt"
)

func TestExec(t *testing.T) {
	cmd := exec.Command("ls", "-la")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Dir: %q\n", out.String())
}
