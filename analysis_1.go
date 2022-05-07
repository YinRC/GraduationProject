package oj

import (
	"fmt"
)


// 在得到初步的结果之后，对用户的输出与标程的输出进行对照，得到最终的结果
func analysis_1 (p_cfg Problem, case_i int, case_rst []Result, problemDir, tmpCodePath string, rst *Result) (err error) {
	// 是语法结构模式且有错误，特判
	if case_i == -1 {
		spj(GrammarMode, nil, p_cfg.GrammarOptionMap, tmpCodePath, rst)
		return nil
	// 初步来看是全部 AC 的，或者是完成比例模式，checker 做进一步的修正
	} else if case_i == 0 { 
		// 将所有用户输出与标程输出比较
		for i := 1; i <= p_cfg.CaseNum; i++ {
			// checker 看一看到底是不是 AC
			if case_rst[i].Flag == AC {
				// checker 更正(浮点数模式的特判在里面完成)
				err = checker(problemDir, p_cfg.Mode, i, rst)
				if err != nil {
					// SE: 文件不存在或读取文件内容失败
					return fmt.Errorf("check %v error: %v", i, err)
				} else if rst.Flag != AC {
					if p_cfg.Mode == GrammarMode {
						spj(GrammarMode, nil, p_cfg.GrammarOptionMap, tmpCodePath, rst)
						return nil
					// 普通模式和浮点数模式，更正出错后直接结束
					} else if p_cfg.Mode != IncompleteMode {
						rst.Hint += good(case_i, problemDir)
						return nil
					}
					case_rst[i] = *rst
				}
			// 检查 Run 函数阶段 case_i 的返回有没有 BUG (不是完成比例模式还允许出错的错误)
			} else {
				if p_cfg.Mode != IncompleteMode {
					return fmt.Errorf("***BUG***: wrong case_i: 0; check Run function")
				}
				// 可以考虑在这里加入 GrammarMode 取代 case_i = -1
			}
		}

		// 是完成比例模式的情况，结束
		if p_cfg.Mode == IncompleteMode {
			spj(IncompleteMode, case_rst, nil, "", rst)
			return nil
		}
	}
	// 有错误标号的情况，直接结束
	rst.Hint += good(case_i, problemDir)
	return nil
}
