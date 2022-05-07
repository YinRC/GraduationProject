package oj

import (
	"io/ioutil"
	"regexp"
	"fmt"
)


// 根据题目设置的特判模式进行相关特判
func spj(mode int, case_rst []Result, p_cfg_map map[string]int, tmpCodePath string, rst *Result) (err error){
	if mode == GrammarMode {
		spj_grammar(tmpCodePath, p_cfg_map, rst)
	} else if mode == IncompleteMode {
		spj_incomplete(case_rst, rst)
	}
	return nil
}


// 特判 1 ：语法练习特判
// 特殊判题：1、~AC 给出部分分数 2、AC 判断是否是直接给出答案
func spj_grammar(tmpCodePath string, p_cfg_map map[string]int, rst *Result) error {
	// 读取用户代码，之后在此代码上进行各种修改和匹配
	code, err := ioutil.ReadFile(tmpCodePath)
	if err != nil {
		return fmt.Errorf("Grammar mode can not open tmpCodePath: %v", err)
	}
	// define 只能定义数字
	r_macro_define := regexp.MustCompile(`#[ \t]*?define.*[0-9]`)
	if ok, _ := regexp.Match(`#[ \t]*?define.*[0-9]`, code); ok {
		code = r_macro_define.ReplaceAll(code, []byte(""));
		rst.Score = 0
		rst.Hint += "WARMING: special judge only supply #define X number"+"\n"
		return nil
	}
	// 采分点（1）：匹配头文件
	//（stdio *stdlib<system> string ctype<isdigit>）
	//（*assert *limits *time *unistd<exec>）
	//（*float **math *error *locale *setjmp *signal *stdarg）
	// stdio.h
	if p_cfg_map["stdio.h"] == 1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<stdio.h>`, code); ok {
			rst.Score += 5
			rst.Hint += "#include <stdio.h>"+"\n"
		}
	} else if p_cfg_map["stdio.h"] == -1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<stdio.h>`, code); ok {
			return fmt.Errorf("violation: stdio.h is banned")
		}
	}
	// string.h
	if p_cfg_map["string.h"] == 1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<string.h>`, code); ok {
			rst.Score += 5
			rst.Hint += "#include <string.h>"+"\n"
		}
	} else if p_cfg_map["string.h"] == -1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<string.h>`, code); ok {
			return fmt.Errorf("violation: string.h is banned")
		}
	}
	// ctype.h
	if p_cfg_map["ctype.h"] == 1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<ctype.h>`, code); ok {
			rst.Score += 5
			rst.Hint += "#include <ctype.h>"+"\n"
		}
	} else if p_cfg_map["ctype.h"] == -1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<ctype.h>`, code); ok {
			return fmt.Errorf("violation: ctype.h is banned")
		}
	}
	// math.h
	if p_cfg_map["math.h"] == 1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<math.h>`, code); ok {
			rst.Score += 5
			rst.Hint += "#include <math.h>"+"\n"
		}
	} else if p_cfg_map["math.h"] == -1 {
		if ok, _ := regexp.Match(`#[ \t]*?include[ \t]*?<math.h>`, code); ok {
			return fmt.Errorf("violation: math.h is banned")
		}
	}
	
	// 采分点（2）：main函数 int main([void]?){return 0;}
	if p_cfg_map["main"] == 1 {
		if ok, _ := regexp.Match(`int[\s]*?main[\s]*?\((void)?\)[\s]*?\{[\s\S]*?return[\s]*?0;[\s]*?\}`, code); ok {
			rst.Score += 5
			rst.Hint += "int main(){return 0;}"+"\n"
		}
	}
	// 采分点（3）：选择结构 if switch
	// if
	if p_cfg_map["if"] == 1 {
		if ok, _ := regexp.Match(`if[\s]*?\([\s\S]*?\)`, code); ok {
			rst.Score += 5
			rst.Hint += "if()"+"\n"
		}
	}
	// switch
	if p_cfg_map["switch"] == 1 {
		if ok, _ := regexp.Match(`switch[\s]*?\([\s\S]*?\)[\s]*?\{[\s]*?case[\s\S]*?:[\s\S]*?break[\s]*?;[\s\S]*?\}`, code); ok {
			rst.Score += 5
			rst.Hint += "switch(){case x: line; break;}"+"\n"
		}
	}
	// 采分点（4）：循环结构 while for
	// while
	if p_cfg_map["while"] == 1 {
		if ok, _ := regexp.Match(`while[\s]*?\([\s\S]*?\)`, code); ok {
			rst.Score += 5
			rst.Hint += "while()"+"\n"
		}
	}
	// for
	if p_cfg_map["for"] == 1 {
		if ok, _ := regexp.Match(`for[\s]*?\([\s\S]*?;[\s\S]*?;[\s\S]*?\)`, code); ok {
			rst.Score += 5
			rst.Hint += "for(;;)"+"\n"
		}
	}
	// 采分点（5）：基本数据类型 float double
	// float
	if p_cfg_map["float"] == 1 {
		if ok, _ := regexp.Match(`float[\s\S]*?;`, code); ok {
			rst.Score += 5
			rst.Hint += "float var;"+"\n"
		}
	}
	// double
	if p_cfg_map["double"] == 1 {
		if ok, _ := regexp.Match(`double[\s\S]*?;`, code); ok {
			rst.Score += 5
			rst.Hint += "double var;"+"\n"
		}
	}
	return nil
}


// 特判 3 ：通过比例特判
// 特殊判题：保留 tmpProgramPath 跑多个 case in 得到的多个 case out 与 标程输出比较 
func spj_incomplete(case_rst []Result, rst *Result) error {
	count := 0
	time := 0
	memory := 0
	caseNum := len(case_rst)

	for i := 1; i <= caseNum; i++ {
		if case_rst[i].Flag == AC {
			count++
			time += case_rst[i].Time
			memory += case_rst[i].Memory
		}
	}
	rst.Score = (count / caseNum) * 100
	if count != 0 {
		rst.Time = time / count
		rst.Memory = memory / count
		// s_rst.Hint = for case_rst[i]
	} else {
		rst.Time = 0
		rst.Memory = 0
		rst.Hint = ""
	}
	return nil
}
