package oj

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"io/ioutil"
	"github.com/satori/go.uuid"
	"bytes"
)

const (
	GNU_C = "gcc %s -o %s"
	GNU_CPP = "g++ %s -o %s"
)

type tmpFilePath struct {
	tmpCodePath string
	tmpProgramPath string
}

func Compile(lang string, codePath string, tmpPath string) (tmp tmpFilePath, err error) {
	var ext string
	var runStr string

	switch lang {
		case "c":
			ext = ".c"
			runStr = GNU_C
		case "c++":
			ext = ".cpp"
			runStr = GNU_CPP
		case "cpp":
			ext = ".cpp"
			runStr = GNU_CPP
	}

	tmpCodeName := fmt.Sprintf("%s%s", uuid.Must(uuid.NewV4()), ext)
	tmpProgramName := fmt.Sprintf("%s%s", uuid.Must(uuid.NewV4()), "")

	// all in tmp

	tmpCodePath := tmpPath + "/" +tmpCodeName
	tmpProgramPath := tmpPath + "/" +tmpProgramName
	tmp.tmpCodePath = tmpCodePath
	tmp.tmpProgramPath = tmpProgramPath


	fp0, err := os.OpenFile(tmpCodePath, os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		return tmp, fmt.Errorf("open/create tmpCodeFile fail")
	}

	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		return tmp, fmt.Errorf("get code content fail")
	}

	_, err = fp0.WriteString(string(code))
	if err != nil {
		return tmp, fmt.Errorf("write string fail")
	}
	fp0.Close()

	runStr = fmt.Sprintf(runStr, tmpCodePath, tmpProgramPath)
	// test command
	fmt.Println(runStr)
	args := strings.Split(runStr, " ")
	cmd := exec.Command(args[0], args[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return tmp, fmt.Errorf(err.Error()+": "+stderr.String())
		//return err
	}
	return tmp, nil
}
