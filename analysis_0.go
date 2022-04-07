package oj

import (
	"syscall"
	"fmt"
)

func analysis_0(cfg Config, pid int, rst *Result) (err error) {
	var (
		status syscall.WaitStatus
		rusage syscall.Rusage
	)

	_, err = syscall.Wait4(pid, &status, syscall.WUNTRACED, &rusage)
	if err != nil {
		return fmt.Errorf("wait4 fail")
	}

	// MS
	uTimeUsed := int(rusage.Utime.Sec*1000+int64(rusage.Utime.Usec/1000))
	sTimeUsed := int(rusage.Stime.Sec*1000+int64(rusage.Stime.Usec/1000))
	rst.Time = uTimeUsed + sTimeUsed

	// KB
	rst.Memory = int(rusage.Minflt*int64(syscall.Getpagesize()/1024))

	if status.Signaled() {
		signal := status.Signal()

		if signal == syscall.SIGSEGV {
			if rst.Memory > cfg.Memory {
				rst.Flag = MLE
			} else {
				rst.Flag = RE
			}
		} else if signal == syscall.SIGXFSZ {
			rst.Flag = OLE
		} else if signal == syscall.SIGXCPU {
			rst.Flag = TLE
		} else if signal == syscall.SIGKILL {
			if rst.Time > cfg.Time {
				rst.Flag = TLE
			} else if rst.Memory > cfg.Memory {
				rst.Flag = MLE
			} else {
				rst.Flag = AC
			}
		}
	}

	if rst.Flag == TLE {
		rst.Time = 0
	} else if rst.Flag == MLE {
		rst.Memory = 0
	}
	return nil
}

