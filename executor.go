package oj

import (
	"fmt"
	"strconv"
	"syscall"
	"os"
	"math"
	"bytes"
	"os/exec"
	"log"
)

type Result struct {
	Flag string

	Time int
	Memory int
	
	Score int
	Hint string

	SE_log string
}


func fork_cpp() (int, error) {
	pid, _, errMsg := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errMsg != 0 {
		return -1, fmt.Errorf("syscall fork fail")
	}
	return int(pid), nil
}


// 先把测试样例都运行一遍，得到初步的结果和用户输出
func Run(p_cfg Problem, problemDir string, tmp tmpFilePath, rst *Result) (case_i int, case_rst []Result, err error) {
	// 生成测试样例和答案
	cmd := exec.Command("make", "-s", "-C", problemDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return 0, nil, fmt.Errorf(err.Error()+": "+stderr.String())
	}
	// 初始化 case_rst 切片，收集用户程序运行各个测试样例的结果
	case_rst = make([]Result, p_cfg.CaseNum)
	tmpProgramPath := tmp.tmpProgramPath
	for i := 1; i <= p_cfg.CaseNum; i++ {
		// 测试样例路径
		casePath := problemDir+"/"+strconv.Itoa(i)
		// 题目设置、用户程序路径、测试样例路径、判题结果（对 rst 进行修改）
		err = runBinFile(&p_cfg, tmpProgramPath, casePath, rst)
		if err != nil {
			return 0, nil, fmt.Errorf("in case %v: %v", i, err)
		}
		// 完成比例模式不管怎样都要测试完全部的测试样例
		if p_cfg.Mode == IncompleteMode {
			case_rst[i-1] = *rst
		// 语法结构模式：有错误就直接特判（检测代码）
		} else if p_cfg.Mode == GrammarMode && rst.Flag != AC {
			return -1, nil, nil
		// 普通模式、浮点数模式：有错误返回未通过的样例
		} else if rst.Flag != AC {
			//fmt.Println("$$$")
			return i, nil, nil
		// AC 的情况，记录初步结果
		} else {
			case_rst[i-1] = *rst
		}
	}
	// 全部初步 AC 或完成比例模式
	return 0, case_rst, nil
}


// 题目设置、用户程序路径、测试样例路径、判题结果
func runBinFile(p_cfg *Problem, tmpProgramPath string, casePath string, rst *Result) error {
	var (
		pid int
		err error
	)

	// 编译生成的用户程序是否存在
	if !isFileExist(tmpProgramPath) {
		rst.Flag = CE
		return fmt.Errorf("compile fail")
	}

	pid, err = fork_cpp()
	if err != nil {
		rst.Flag = SE
		rst.SE_log = err.Error()+"\n"
		return err
	}
	
	// 用户程序是否存在
	
	// father
	if pid == 0 {
		// reapChildren()
		// 输入重定向为 $(no) [i]
		fd, child_err := getFD(casePath, os.O_RDONLY, 0)
		if child_err != nil {
			return fmt.Errorf("get %s fd fail: %s", casePath, child_err.Error())
		}
		child_err = syscall.Dup2(fd, syscall.Stdin)
		if child_err != nil {
			return fmt.Errorf("dup2 %s fail", casePath)
		}

		// 输出重定向为 $(no).usr [i.usr]
		fd, child_err = getFD(casePath+".usr", os.O_WRONLY | os.O_CREATE, 0644)
		if child_err != nil {
			return fmt.Errorf("get %s.usr fd fail: %s", casePath, child_err.Error())
		}
		child_err = syscall.Dup2(fd, syscall.Stdout)
		if child_err != nil {
			return fmt.Errorf("dup2 %s.usr fail", casePath)
		}
		
		// 设置程序运行环境 (time--memory--outputSize)
		child_err = setRunLimit(p_cfg.Time, p_cfg.Memory, p_cfg.OutputSize)
		if child_err != nil {
			return fmt.Errorf("setRunLimit fail")
		}
		// 执行 tmp 路径下的二进制文件
		child_err = syscall.Exec(tmpProgramPath, nil, nil)
		if child_err != nil {
			return fmt.Errorf("exec bin file fail")
		}
		return nil
	} else {
		// 通过信号收集信息进行初步判断，rst.Flag\rst.Time\rst.Memory
		err = analysis_0(p_cfg, pid, rst)
		if err != nil {
			rst.Flag = SE
			rst.SE_log += err.Error()+"\n"
			return err
		}
		//reapChildren()
		
	}
	return nil
}

// 收割僵尸进程
func reapChildren() {
	for {
		var wstatus syscall.WaitStatus
		wpid, err := syscall.Wait4(-1, &wstatus, syscall.WNOHANG, nil)
		if err != nil {
			log.Printf("syscall.Wait4 call failed: %v", err)
			break
		}

		if wpid != 0 {
			log.Printf("reap dead child: %d, wstatus: %#08x", wpid, wstatus)
		} else {
			break
		}
	}
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
