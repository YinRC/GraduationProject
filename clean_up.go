package oj

import (
	"fmt"
	"bytes"
	"os/exec"
)

func CleanUp(problemDir string) error {
	cmd := exec.Command("make", "clean", "-s", "-C", problemDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(err.Error()+": "+stderr.String())
	}
	return nil
}
