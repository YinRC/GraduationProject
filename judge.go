package oj

import (
	"fmt"
)


func Judge(p_cfg Problem, case_i int, case_rst []Result, tmp tmpFilePath, problemDir string, rst *Result) (err error){
	acPath := problemDir+"/"+p_cfg.AC_Path
	tmpCodePath := tmp.tmpCodePath
	// 标程是否存在
	if !isFileExist(acPath) {
		return fmt.Errorf("%s may not exsit\n", acPath)
	}

	// isCodeSizeFine：比较标程代码和用户代码的字符数，判断是否是直接打印结果
	ok, err := isCodeSizeFine(acPath, tmpCodePath)
	if err != nil {
		return fmt.Errorf("in function isCodeSizeFine: %s", err.Error())
	}
	if !ok {
		return fmt.Errorf("Error: your code size is not fine.\n")
	}

	// analysis_1()：通过比较用户输出和标程输出给出进一步的判题结果
	err = analysis_1(p_cfg, case_i, case_rst, problemDir, tmpCodePath, rst)
	if err != nil {
		return fmt.Errorf("in function analysis_1: %s", err.Error())
	}
	
	return nil
}
