package oj

import (
	"io/ioutil"
	"regexp"
	"fmt"
	"math"
)


// 根据题目设置的特判模式进行相关特判 AC PC WA (*PE for PointMode)
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
// 宽容判题，只加分不扣分
func spj_grammar(tmpCodePath string, p_cfg_map map[string]int, rst *Result) error {
	// 读取用户代码，之后在此代码上进行各种修改和匹配
	code, err := ioutil.ReadFile(tmpCodePath)
	if err != nil {
		return fmt.Errorf("Grammar mode can not open tmpCodePath: %v", err)
	}
	// define 只能定义数字
	r_macro_define := regexp.MustCompile(`#[ \t]*?define.*[a-zA-Z]`)
	if ok, _ := regexp.Match(`#[ \t]*?define.*[a-zA-Z]`, code); ok {
		code = r_macro_define.ReplaceAll(code, []byte(""));
		rst.Hint += "WARMING: special judge only supply #define X number."+"\n"
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
			rst.Flag = WA
			rst.Score = 0
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
			rst.Flag = WA
			rst.Score = 0
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
			rst.Flag = WA
			rst.Score = 0
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
			rst.Flag = WA
			rst.Score = 0
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
	// 特判结果：如果特判分数为零判 WA 不为零判 PC
	if rst.Score == 0 {
		rst.Flag = WA
	} else if rst.Score > 0 {
		rst.Flag = PC
	}
	return nil
}


// 特判 2 ：浮点数特判
// 特殊判题：替换 diffUtil_0() diffUtil_1() 的字符比较，改为截取的数字比较(和空白符比较)
// 答案与用户输出对照，误差全部 <= 1e5：100 
// 答案与用户输出对照，存在的最大误差 >= 1e5 && <= 1.0：60
// 答案与用户输出对照，存在误差 > 1.0：0
func spj_point(problemDir string, i int, rst *Result) error {
	// 先通过字节流长短的比较来看是否存在问题（AC OLE WA）
	usr, ans, err := standBy4Check(problemDir, i, rst)
	if err != nil {
		return fmt.Errorf("in function standBy4Check: %s", err.Error())
	}
	// 将答案和用户输出的字节流切分成单词
	delim := regexp.MustCompile(`[\S]+`)
	usr_words := delim.FindAll(usr, -1)
	ans_words := delim.FindAll(ans, -1)
	// 两者单词数目不同则直接 WA
	usr_total := len(usr_words)
	ans_total := len(ans_words)
	if usr_total != ans_total {
		rst.Flag = WA
		rst.Hint += "PointMode: split result length different"
		return nil
	}
	var score int = 100
	// 单词数目相同则判断单词是否为数字
	for i:=0; i<ans_total; i++ {
		// 都是数字则将单词转换为数字，根据两者之差作出特判
		if isNum(ans_words[i]) && isNum(usr_words[i]) {
			ans_word, err := byteToFloat64(ans_words[i])
			if err != nil {
				return fmt.Errorf("byteToFloat64 ans_word fail: %s", err.Error())
			}
			usr_word, err := byteToFloat64(usr_words[i])
			if err != nil {
				return fmt.Errorf("byteToFloat64 usr_word fail: %s", err.Error())
			}
			abs_val := math.Abs(ans_word - usr_word)
			if abs_val <= 1e-5 {
				score = min(score, 100)
			} else if abs_val <= 1.0 {
				score = min(score, 60)
			} else {
				score = min(score, 0)
			}
		// 有一方不是数字则返回 WA
		} else if isNum(ans_words[i]) || isNum(usr_words[i]) {
			rst.Flag = WA
			rst.Hint += "print somthing not number\n"
			return nil
		}
	}
	rst.Score = score
	if score == 100 {
		rst.Flag = AC
	} else if score == 60 {
		rst.Flag = PC
		rst.Hint += "1e-5 < abs(ans_word - usr_word) > 1.0\n"
	} else {
		rst.Flag = WA
		rst.Hint += "abs(ans_word - usr_word) > 1.0\n"
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
	// 遍历结果数组，记录 AC 的测试样例数目，统计时间和空间信息
	for i := 0; i < caseNum; i++ {
		if case_rst[i].Flag == AC {
			count++
			time += case_rst[i].Time
			memory += case_rst[i].Memory
		}
	}
	// 根据通过的测试样例占全部测试样例的比例给出分数
	rst.Score = (count / caseNum) * 100
	if count == len(case_rst) {
		rst.Flag = AC
		rst.Time = time / count
		rst.Memory = memory / count
		rst.Hint = "good job\n"
	// 如果部分正确：给出平均的时间空间数据
	} else if count != 0 {
		rst.Flag = PC
		rst.Score = int(float64(count) / float64(len(case_rst)) * 100)
		rst.Time = time / count
		rst.Memory = memory / count
		rst.Hint = ""
		// 格式输出各测试样例的结果
		for i := 1; i <= caseNum; i++ {
			if i % 4 != 0 {
				rst.Hint += fmt.Sprintf("case_%-2d: %s\t", i, case_rst[i-1].Flag)
			} else {
				rst.Hint += fmt.Sprintf("case_%-2d: %s\n", i, case_rst[i-1].Flag)
			}
		}
		if caseNum % 4 != 0 {
			rst.Hint += "\n"
		}
	// 如果全部错误：时间和空间数据归零
	} else {
		rst.Flag = WA
		rst.Time = 0
		rst.Memory = 0
		rst.Hint = "IncompleteMode: all fail\n"
	}
	return nil
}
