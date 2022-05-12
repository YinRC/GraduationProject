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
			// checker 看一看到底是不是 AC（对于不确实的 AC）
			if case_rst[i-1].Flag == AC {
				// checker 更正(浮点数模式的特判在里面完成)
				err = checker(problemDir, p_cfg.Mode, i, rst)
				if err != nil {
					// 文件不存在或读取文件内容失败
					return fmt.Errorf("check %d error: %s", i, err)
				} else if rst.Flag != AC {
					rst.Time = 0
					rst.Memory = 0
					if p_cfg.Mode == GrammarMode {
						spj(GrammarMode, nil, p_cfg.GrammarOptionMap, tmpCodePath, rst)
						if rst.Flag == WA || rst.Flag == PC {
							rst.Hint += good(i, problemDir)
						}
						return nil
					// 普通模式和浮点数模式，更正出错后直接结束
					} else if p_cfg.Mode != IncompleteMode {
						rst.Hint += fmt.Sprintf("case_%-2d fail\n", i)
						// ana_1 错误，且错误为 WA PE PC，返回 good （普通模式和浮点数模式）
						if rst.Flag == WA || rst.Flag == PE || rst.Flag == PC {
							rst.Hint += good(i, problemDir)
						}
						return nil
					}
					// 完成比例模式，完整记录所有测试用例的判题结果
					case_rst[i-1] = *rst
				}
			// case_i 有错误但没返回（BUG）
			} else {
				// 不是完成比例模式，出错的时候不返回 case_i
				if p_cfg.Mode != IncompleteMode {
					return fmt.Errorf("***BUG***: wrong case_i: 0; check Run function")
				}
				// 可以考虑在这里加入 GrammarMode 取代 case_i = -1
			}
		}
		// 更正完毕，是完成比例模式的情况，结束
		if p_cfg.Mode == IncompleteMode {
			spj(IncompleteMode, case_rst, nil, "", rst)
			return nil
		}
	// 有错误标号的情况，直接结束（case_i > 0）
	} else {
		// ana_0 错误不返回 good
		rst.Hint += fmt.Sprintf("case_%-2d fail\n", case_i)
		return nil
	}
	rst.Hint += "good job\n"
	return nil
}
