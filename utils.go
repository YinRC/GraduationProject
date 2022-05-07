package oj

import (
	"os"
	"os/exec"
	"fmt"
	"syscall"
	"bytes"
	"encoding/binary"
	"regexp"
	"math"
	"io/ioutil"
	"strconv"
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

func rmBlank(text *[]byte) {
	r_blank := regexp.MustCompile(`[\s]+`)
	*text = r_blank.ReplaceAll(*text, []byte(""))
}


// 判断是否是空白字符
func isBlank(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\f' || b == '\r' || b == '\v'
}

// 将字节序列转换成小数
func byteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

// 生成标程输出和用户输出的文件名
func findua(i, suffix int) string {
	if suffix == 0 {
		return strconv.Itoa(i)+".a"
	} else {
		return strconv.Itoa(i)+".u"
	}
}

// 检查是不是直接打印了结果
func isCodeSizeFine(acPath string, tmpCodePath string) (bool, error) {
	ac, err := ioutil.ReadFile(acPath)
	if err != nil {
		return false, fmt.Errorf("read %s fail", acPath)
	}
	code, err := ioutil.ReadFile(tmpCodePath)
	if err != nil {
		return false, fmt.Errorf("read %s fail", tmpCodePath)
	}
	rmBlank(&ac)
	rmBlank(&code)
	if len(ac)/2 > len(code) || len(code)/2 > len(ac) {
		return false, nil
	}
	return true, nil
}

// 执行 make clean 清除 1 1.a 1.u
func makeClean(problemDir string) error {
	cmd := exec.Command("make", "clean", "-s", "-C", problemDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(err.Error()+": "+stderr.String())
	}
	return nil
}

// 打印指定的测试样例、标程输出、用户输出
func good(case_i int, problemDir string) (hint string) {
	in, _ := ioutil.ReadFile(problemDir+strconv.Itoa(case_i))
	ans, _ := ioutil.ReadFile(problemDir+findua(case_i, 0))
	out, _ := ioutil.ReadFile(problemDir+findua(case_i, 1))
	hint = fmt.Sprintf("\tcase in:\n%v\n\tans:\n%v\n\tout:\n%v\n", in, ans, out)
	return hint
}

func (rst *Result)String(mode int) {
	if mode != NormalMode {
		fmt.Printf("-----------\nRESULT: %s\nTIME: %v\tMEMORY: %v\n-----------\nSCORE: %d\n-----------\nHINT:\n%s-----------\n", rst.Flag, rst.Time, rst.Memory, rst.Score, rst.Hint)
	} else {
		fmt.Printf("-----------\nRESULT: %s\nTIME: %v\tMEMORY: %v\n-----------\n", rst.Flag, rst.Time, rst.Memory)
	}
}

//func IsNum(s string) bool {
	//_, err := strconv.ParseFloat(s, 64)
	//return err == nil
//}
