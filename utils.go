package oj

import (
	"os"
	"os/exec"
	"fmt"
	"syscall"
	"bytes"
	"regexp"
	"io/ioutil"
	"strconv"
)


// 判断传入的路径上是否存在文件
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 根据给定的路径和设定的权限打开文件，返回文件描述符
func getFD(path string, flag int, perm uint32) (fd int, err error) {
	fd, err = syscall.Open(path, flag, perm)
	if err != nil {
		return -1, fmt.Errorf("open %s fail", path)
	}
	return fd, nil
}

// 去掉字节数组中的空白符
func rmBlank(text *[]byte) {
	r_blank := regexp.MustCompile(`[\s]+`)
	*text = r_blank.ReplaceAll(*text, []byte(""))
}


// 判断是否是空白字符
func isBlank(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\f' || b == '\r' || b == '\v'
}

// 将字节数组转换成小数
func byteToFloat64(bytes []byte) (float64, error) {
	return strconv.ParseFloat(string(bytes), 64)
}

// 生成标程输出和用户输出的文件名
func findua(i, suffix int) string {
	if suffix == 0 {
		return strconv.Itoa(i)+".ans"
	} else {
		return strconv.Itoa(i)+".usr"
	}
}

// 检查是不是直接打印了结果
func isCodeSizeFine(acPath string, tmpCodePath string) (bool, error) {
	// 标程代码
	ac, err := ioutil.ReadFile(acPath)
	if err != nil {
		return false, fmt.Errorf("read %s fail", acPath)
	}
	// 用户代码
	code, err := ioutil.ReadFile(tmpCodePath)
	if err != nil {
		return false, fmt.Errorf("read %s fail", tmpCodePath)
	}
	// 除去字符流中的空白符（用户代码已经过处理，不包含注释，【标程代码也不应该加入注释】）
	rmBlank(&ac)
	rmBlank(&code)
	// 用户代码非空白字符不能比标程代码的一半还小 或 比标程代码的二倍还大
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
	if case_i == 0 {
		return ""
	}
	in, _ := ioutil.ReadFile(problemDir+"/"+strconv.Itoa(case_i))
	ans, _ := ioutil.ReadFile(problemDir+"/"+findua(case_i, 0))
	out, _ := ioutil.ReadFile(problemDir+"/"+findua(case_i, 1))
	hint = fmt.Sprintf("\tcase in:\n%s\n\tans:\n%s\n\tout:\n%s\n", in, ans, out)
	return hint
}

// 比较两实数的大小，返回较小的那个
func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

// 只有特判时显示 Score 和 Hint
func (rst *Result) String(mode int) {
	if rst.Flag != SE {
		// 特判模式输出
		if mode != NormalMode {
			fmt.Printf("-----------\nRESULT: %s\nTIME: %v\tMEMORY: %v\n-----------\nSCORE: %d\n-----------\nHINT:\n%s-----------\n", rst.Flag, rst.Time, rst.Memory, rst.Score, rst.Hint)
		// 普通模式输出
		} else {
			fmt.Printf("-----------\nRESULT: %s\nTIME: %v\tMEMORY: %v\n-----------\nHINT: %s", rst.Flag, rst.Time, rst.Memory, rst.Hint)
		}
	} else {
		fmt.Println("-----------\nRESULT: %s\n%s\n-----------\n", SE, rst.SE_log)
	}
}

// 判断字节数组是否是数字
func isNum(s []byte) bool {
	str := string(s)
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}
