package oj

import (
	"syscall"
	"fmt"
)

func analysis_0(p_cfg *Problem, pid int, rst *Result) (err error) {
	var (
		status syscall.WaitStatus
		rusage syscall.Rusage
	)
	_, errMsg := syscall.Wait4(pid, &status, syscall.WUNTRACED, &rusage)
	if errMsg != nil {
		return errMsg
	}
	// MS
	uTimeUsed := int(rusage.Utime.Sec*1000 + int64(rusage.Utime.Usec/1000))
	sTimeUsed := int(rusage.Stime.Sec*1000 + int64(rusage.Stime.Usec/1000))
	rst.Time = uTimeUsed + sTimeUsed
	// KB
	rst.Memory = int(rusage.Minflt * int64(syscall.Getpagesize()/1024))
	
	// 接收信号
	if status.Signaled() {
		signal := status.Signal()
		fmt.Println(signal)
		if signal == syscall.SIGFPE {
			rst.Flag = RE
		} else if signal == syscall.SIGSEGV {
			if rst.Memory > p_cfg.Memory {
				rst.Flag = MLE
			} else {
				rst.Flag = RE
			}
		} else if signal == syscall.SIGXFSZ {
			rst.Flag = OLE
		} else if signal == syscall.SIGXCPU {
			rst.Flag = TLE
		} else if signal == syscall.SIGKILL {
			if rst.Time > p_cfg.Time {
				rst.Flag = TLE
			} else if rst.Memory > p_cfg.Memory {
				rst.Flag = MLE
			} else {
				rst.Flag = AC
			}
		}
	} else {
		if rst.Time > p_cfg.Time {
			rst.Flag = TLE
		} else if rst.Memory > p_cfg.Memory {
			rst.Flag = MLE
		} else {
			rst.Flag = AC
		}
	}
	
	if rst.Flag == TLE {
		rst.Time = 0
		rst.Memory = 0
	} else if rst.Flag == MLE {
		rst.Time = 0
		rst.Memory = 0
	}
	// fmt.Println(rst.Flag)
	return nil
}

