package oj

import (
	"io/ioutil"
	"fmt"
)

func standBy4Check(problemDir string, i int, rst *Result) (usr , ans []byte, err error) {
	// i.ans 存在  SE
	ansFile := problemDir+"/"+findua(i, 0)
	if !isFileExist(ansFile) {
		rst.Flag = SE
		return usr, ans, fmt.Errorf("%d.ans may not exist", i)
	}
	// i.usr 存在
	usrFile := problemDir+"/"+findua(i, 1)
	if !isFileExist(usrFile) {
		rst.Flag = SE
		return usr, ans, fmt.Errorf("%d.usr may not exist", i)
	}
	// ans 接收标程输出 buffer
	ans, err = ioutil.ReadFile(ansFile)
	if err != nil {
		rst.Flag = SE
		return usr, ans, fmt.Errorf("read %d.ans fail", i)
	}
	// usr 接收用户输出 buffer
	usr, err = ioutil.ReadFile(usrFile)
	if err != nil {
		rst.Flag = SE
		return usr, ans, fmt.Errorf("read %d.usr fail", i)
	}
	return usr, ans, nil
}


// 对照用户输出和标程输出，接收 Run 运行结果: []Result
func checker(problemDir string, mode int, i int, rst *Result) (err error) {
	usr, ans, err := standBy4Check(problemDir, i, rst)
	if err != nil {
		return fmt.Errorf("in function standBy4Check: %s", err.Error())
	}
	// 首先检验内容大小（初步检测有没有直接打印答案的可能 PPT）
	usrLen, ansLen := len(usr), len(ans)
	if usrLen == 0 && ansLen == 0 {
		rst.Flag = AC
		rst.Hint += "answer size: zero; output size: zero;\n"
		return nil
	} else if usrLen > 0 && ansLen > 0 {
		// 对输出文件大小限制的补充（一般来讲此限制应该以定为 2 倍的答案大小为宜）
		if usrLen / 2 > ansLen {
			rst.Flag = OLE
			rst.Hint += "OLE: larger than 2 times of answer size\n"
			return nil
		}
	} else { // 有一方为 0 的情况
		rst.Flag = WA
		if usrLen == 0 {
			rst.Hint += "WA: output size is zero\n"
			return nil
		} else if ansLen == 0 {
			rst.Hint += "WA: answer size is zero but your output size not\n"
			return nil
		}
	}
	if mode == PointNumMode {
		spj_point(problemDir, i, rst)
	} else {
		// 不计空白符的检查
		flag, hint := diffUtil_0(usr, ans, mode)
		rst.Flag = flag
		rst.Hint += hint
		if flag != AC {
			rst.Hint += good(i, problemDir)
		}
		// 通过了不计空白符的检测再严格检查，查出可能的 PE
		if flag == AC {
			if !diffUtil_1(usr, ans, mode) {
				rst.Flag = PE
				rst.Hint += "PE: check your blank char output\n"
				return nil
			} else {
				rst.Flag = AC
				return nil
			}
		}
	}
	
	return nil
}

// 将用户输出与标程输出进行对照（忽略空白字符）同时进行浮点数特判
func diffUtil_0(usr, ans []byte, mode int) (flag string, hint string) {
	var (
		usrLen, ansLen = len(usr), len(ans)
		left, right = 0, 0
	)

	for left <= usrLen && right <= ansLen {
		//fmt.Println(left, usrLen, right, ansLen)
		// usr 跳过遇到的空白
		for left <usrLen && isBlank(usr[left]) {
			left++
		}
		// ans 跳过遇到的空白
		for right <ansLen && isBlank(ans[right]) {
			right++
		}
		// 都没到头
		if left < usrLen && right < ansLen {
			if usr[left] != ans[right] {
				return WA, "wrong answer\n"
			} else {
				left++
				right++
				continue
			}
		// left 读 usr 越界
		} else if left == usrLen && right < ansLen {
			for right < ansLen {
				if !isBlank(ans[right]) {
					return WA, "your output is shorter than answer\n"
				}
				right++
			}
		// right 读 ans 越界
		} else if right == ansLen && left < usrLen {
			for left < usrLen {
				if !isBlank(usr[left]) {
					return WA, "your output is longer than answer\n"
				}
				left++
			}
		// 同时到达边界
		} else {
			break
		}
	}
	return AC, ""
}

// 在不考虑空白符号的检查之后进行更严格的检查，看看是不是 PE
func diffUtil_1(usr, ans []byte, mode int) bool {
	var (
		usrLen = len(usr)
		ansLen = len(ans)
	)
	if usrLen != ansLen {
		return false
	}
	for i := 0; i < ansLen; i++ {
		if usr[i] != ans[i] {
			return false
		}
	}
	return true
}
