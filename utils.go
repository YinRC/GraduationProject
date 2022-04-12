package oj

import (
	"os"
	"fmt"
	"math"
	"syscall"
	"bytes"
	"encoding/json"
)


func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func getFD(path string, flag int, perm uint32) (fd int, err error) {
	fd, err = syscall.Open(path, flag, perm)
	if err != nil {
		return -1, fmt.Errorf("open %s fail", path)
	}
	return fd, nil
}

func setRunLimit(time int, memory int, outputSize int) (err error) {
	var rlimit syscall.Rlimit

	// set time limit (CPU time in seconds)
	rlimit.Cur = uint64(math.Ceil(float64(time)/1000.0))
	rlimit.Max = rlimit.Cur + 1
	err = syscall.Setrlimit(syscall.RLIMIT_CPU, &rlimit)
	if err != nil {
		return fmt.Errorf("set cpu time limit fail")
	}

	// set memory limit: data + heap
	rlimit.Cur = uint64(memory*1024)
	rlimit.Max = rlimit.Cur
	err = syscall.Setrlimit(syscall.RLIMIT_DATA, &rlimit)
	if err != nil {
		return fmt.Errorf("set memory[data] limit fail")
	}

	// set memory limit: stack
	rlimit.Cur = uint64(memory*1024)
	rlimit.Max = rlimit.Cur
	err = syscall.Setrlimit(syscall.RLIMIT_STACK, &rlimit)
	if err != nil {
		return fmt.Errorf("set memory[stack] limit fail")
	}

	// set outputSize size limit
	rlimit.Cur = uint64(outputSize)
	rlimit.Max = rlimit.Cur+1
	err = syscall.Setrlimit(syscall.RLIMIT_FSIZE, &rlimit)
	if err != nil {
		return fmt.Errorf("set outputSize limit fail")
	}
	return nil
}

func String(result Result) {
	b, _ := json.Marshal(result)
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	fmt.Printf("%+v\n", result)
}
