package oj

import (
	"fmt"
)



func Judge(p_cfg Problem, case_i int, case_rst []Result, tmp tmpFilePath, problemDir string, rst *Result) (err error){
	acPath := problemDir+"/"+p_cfg.AC_Path
	tmpCodePath := tmp.tmpCodePath
	// 标程是否存在
	if !isFileExist(acPath) {
		return fmt.Errorf("%v may not exsit\n", acPath)
	}
	// isCodeSizeFine
	ok, err := isCodeSizeFine(acPath, tmpCodePath)
	if err != nil {
		return fmt.Errorf("in function isCodeSizeFine: "+err.Error())
	}
	if !ok {
		return fmt.Errorf("Error: your code size is not fine.\n")
	}
	
	// analysis_1()
	err = analysis_1(p_cfg, case_i, case_rst, problemDir, tmpCodePath, rst)
	if err != nil {
		return fmt.Errorf("in function analysis_1: "+err.Error())
	}
	
	return nil
}
