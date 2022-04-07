package oj

import (
	"fmt"
	"syscall"
	"os"
)

type Result struct {
	Flag int

	Time int
	Memory int

	SE_log string
}

func fork_cpp() (int, error) {
	pid, _, errMsg := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errMsg != 0 {
		return -1, fmt.Errorf("syscall fork fail")
	}
	return int(pid), nil
}

func RunBinFile(cfg Config, tmp tmpFilePath, rst *Result) error {
	_ = os.Remove(cfg.WorkDir + "/usr.out")
	_ = os.Remove(cfg.WorkDir + "/usr.err")
	var (
		pid int
		err error

		fd int
	)

	pid, err = fork_cpp()
	if err != nil {
		rst.Flag = SE
		rst.SE_log = err.Error()+"\n"
		return err
	}

	if pid == 0 {
		var child_err error
		// 输入重定向为 case in
		fd,child_err := getFD(cfg.CaseInPath, os.O_RDONLY, 0)
		if child_err != nil {
			//return fmt.Errorf("get case in fd fail")
			return child_err
		}
		child_err = syscall.Dup2(fd, syscall.Stdin)
		if child_err != nil {
			return fmt.Errorf("dup2 case in fail")
		}

		// 输出重定向为 WorkDir/usr.out
		fd, child_err = getFD(cfg.WorkDir+"/usr.out", os.O_WRONLY | os.O_CREATE, 0644)
		if child_err != nil {
			//return fmt.Errorf("get usr.out fd fail")
			return child_err
		}
		child_err = syscall.Dup2(fd, syscall.Stdout)
		if child_err != nil {
			return fmt.Errorf("dup2 usr.out fail")
		}

		// 错误重定向为 WorkDir/usr.err
		fd, child_err = getFD(cfg.WorkDir+"/usr.err", os.O_WRONLY | os.O_CREATE, 0644)
		if child_err != nil {
			return fmt.Errorf("get usr.err fd fail")
		}
		child_err = syscall.Dup2(fd, syscall.Stderr)
		if child_err != nil {
			return fmt.Errorf("dup2 usr.err fail")
		}

		// 设置程序运行环境 (time--memory--outputSize)
		child_err = setRunLimit(cfg.Time, cfg.Memory, cfg.OutputSize)
		if child_err != nil {
			return fmt.Errorf("setRunLimit fail")
		}

		// 执行 tmp 路径下的二进制文件
		child_err = syscall.Exec(tmp.tmpProgramPath, nil, nil)
		if child_err != nil {
			return fmt.Errorf("exec bin file fail")
		}
		return nil
	} else {
		// 通过信号收集信息进行初步判断 0
		err = analysis_0(cfg, pid, rst)
		if err != nil {
			rst.Flag = SE
			rst.SE_log += err.Error()+"\n"
			return err
		}

		err = syscall.Close(fd)
		if err != nil {
			return fmt.Errorf("fd close fail")
		}
	}
	return nil
}




























