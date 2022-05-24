# open_judge
**A testing machine for green hand ！(Linux only)**

## 文件结构
```html
-main.go
-config.json	// for system
-open_judge
	-setting.go
	-compiler.go
	-executor.go
	-analysis_0.go
	-analysis_1.go
	-check.go
	-judge.go
	-judge_spj.go
	-utils.go
	-clean_up.go
-satori
-problem
	-testlib.h
	-1
		-ac.cpp
		-generator.cpp
		-makefile
		-config.json	// for problem
		-1.md
	-2
		-ac.cpp
		-generator.cpp
		-makefile
		-config.json
		-1.md
-code
	-1.cpp
	-2.cpp
-tmp
	-uuid.cpp
	-uuid
	..
	..
```

## 常量部分
```html
AC = iota	// ok
PE		// 格式错误
TLE		// 超时错误
MLE		// 超过内存空间限制
WA		// 答案错误
RE		// 运行时错误
OLE		// 输出超过限制
CE		// 编译错误
SE		// 系统错误
PC		// 特判分数
```
## 结构体部分

### 测评机配置
+ 对文件路径、代码编译过程、程序运行过程进行设置并持久化到config.json。
+ 设置config可以通过命令行也可以直接修改config.json文件
```go
type Config struct {

	Lang string

	ProblemDir string	// 题目路径（定位到文件夹）
	CodePath string		// 代码路径（定位到文件）
	TmpDir string		// 临时文件路径（用户代码、用户程序）
}
```

### 题目配置
```go
type Problem struct {

	AC_Path string		// 标程的相对路径，根为题目目录
	
	Time int		// 题目的时间限制
	Memory int		// 题目的空间限制
	OutputSize int		// 题目的输出限制
	
	Mode int		// 设定题目的判题模式
				// 0-普通 1-语法结构 2-浮点数 3-完成比例
	GrammarOptionMap map[string]int	// 设置语法结构模式的采分点 main: 1
					// -1 代表不该出现此语法结构，出现则返回错误
					// 0 代表不检测该语法结构
					// 1 代表检测该语法结构，给出部分分数，只加分不扣分
	
	CaseNum int		// 测试样例的数目 与 generator.cpp 的编写相匹配
}
```

### 编译器临时文件路径
后期会探索特殊判题的方式，所以将编译前后的文件存放在tmp文件夹中。
```go
type tmpFilePath struct {

	tmpCodePath string      // 用户提交代码的路径（临时）
	tmpProgramPath string   // 编译完成后用户程序的路径（临时）
}
```
### 测评机判题结果
```go
type Result struct {

	Flag string

	Time int
	Memory int
	
	Score int
	Hint string

	SE_log string
}
```

## 可用的包外调用函数

可以据此编写适用于个人的main.go
```go
// compiler.go	编译用户代码
func Compile(lang string, codePath string, tmpPath string) (tmp tmpFilePath, err error)
// executor.go	运行二进制文件（多次运行，次数为测试样例数）
func Run(p_cfg Problem, problemDir string, tmp tmpFilePath, rst *Result) (case_i int, case_rst []Result, err error)
// judge.go	对运行结果进行评判
func Judge(p_cfg Problem, case_i int, case_rst []Result, tmp tmpFilePath, problemDir string, rst *Result) (err error)
// utils.go	输出判题结果（Result结构体）
func (rst *Result) String(mode int) {
	if rst.Flag != SE {
		// 特判模式输出
		if mode != NormalMode {
			fmt.Printf("-----------\nRESULT: %s\nTIME: %v\tMEMORY: %v\n-----------\nSCORE: %d\n-----------\nHINT:\n%s-----------\n", rst.Flag, rst.Time, rst.Memory, rst.Score, rst.Hint)
		// 普通模式输出
		} else {
			fmt.Printf("-----------\nRESULT: %s\nTIME: %v\tMEMORY: %v\n-----------\nHINT: %s", rst.Flag, rst.Time, rst.Memory, rst.Hint)
		}
	} else {
		fmt.Println("-----------\nRESULT: %s\n%s\n-----------\n", SE, rst.SE_log)
	}
}
// clean_up.go
func CleanUp(problemDir string) error		// 执行 make clean 命令
// setting.go
func (cfg *Config) GetConfig() (err error)	// 从文件config.json中读取设置的参数并应用
func (cfg *Config) SetFromCmd() (err error)	// 从命令行设置config.json文件
```


## config.json
```json
{
        "Lang": "cpp",
        "CodePath": "code/ac.cpp",
        "CaseInPath": "data/1.in",
        "CaseOutPath": "data/1.out",
        "WorkDir": "work",
        "Time": 1000,
        "Memory": 5000,
        "OutputSize": 104857600
}
```
### 命令行选项
```go
// setting.go
flag.StringVar(&lang, "l", "cpp", "编译所用的语言(小写)-默认为C++语言" )
flag.StringVar(&problem, "p", "problem/2", "题目路径")
flag.StringVar(&code, "c", "code/usr1.c", "代码路径")
flag.StringVar(&tmp, "tmp", "tmp", "临时代码、程序存放路径")
```

## 题目配置文件的设置
```go
package main

import (
	"fmt"
	"os"
	"encoding/json"
)

type Problem struct {
	AC_Path string
	Time int
	Memory int
	OutputSize int
	Mode int
	GrammarOptionMap map[string]int
	CaseNum int
}

// 设置并生成题目的 config.json
func(p_cfg *Problem) setProblem() (err error){
	_ = os.Remove("config.json")
	fp, err := os.OpenFile("config.json", os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open/create config fail")
	}
	defer fp.Close()

	encoder := json.NewEncoder(fp)
	err = encoder.Encode(p_cfg)
	if err != nil {
		return fmt.Errorf("encode config fail")
	}
	return nil
}

func (p_cfg *Problem) getProblem() (err error) {
	fp, err := os.Open("config.json")
	defer fp.Close()
	if err != nil {
		return fmt.Errorf("open config fail")
	}

	decoder := json.NewDecoder(fp)
	err = decoder.Decode(p_cfg)
	if err != nil {
		return fmt.Errorf("decode config fail")
	}
	return nil
}

func main(){
	_ = os.Remove("config.json")
	p_cfg := &Problem{
		AC_Path:	"ac.c",
		Time:		500,
		Memory:		500,
		OutputSize:	500,
		Mode:	1,
		GrammarOptionMap: map[string]int {
			// 采分点 1：头文件
			"stdio.h":	1,
			"string.h":	0,
			"ctype.h":	0,
			"math.h":	0,
			// 采分点 2：main 函数
			"main":		1,
			// 采分点 3：选择结构
			"if":		0,
			"switch":	0,
			// 采分点 4：循环结构
			"while":	1,
			"for":		1,
			// 采分点 5：基本数据类型
			"float":	0,
			"double":	0,
		},
		CaseNum:	10,	// 测试样例数目
	}
	_ = p_cfg.setProblem()

	var tmp Problem
	_ = tmp.getProblem()

	t, _ := json.MarshalIndent(tmp, "", "\t")
	fmt.Printf("%+v\n", string(t))
}
```

## 题目中 makefile 的编写
```makefile
CASE := 1 2 3 4 5 6 7 8 9 10

todo: gen ac
	./gen 1 1 
	@$(foreach var, $(CASE), \
		./ac < $(var) > $(var).ans; \
	)

gen: generator.cpp
	g++ -o gen $<

ac: ac.*
	g++ -o ac $<

clean:
	rm gen ac $(CASE) *.ans *.usr
```


## main.go
```go
package main

import (
	"fmt"
	"github.com/open_judge"
)

func main() {
	// 测评机设置
	var cfg oj.Config
	
	// 题目设置
	var p_cfg oj.Problem

	// 从命令行中读取参数更改config.json
	err := cfg.SetFromCmd()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	// 从config.json读取题目目录，在题目路径中读取config.json即题目设置
	err = p_cfg.GetProblemConfig(cfg.ProblemDir+"/"+"config.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 用选定的编译语言编译选定的代码，返回临时代码路径和用户程序路径（在指定的临时文件目录）
	tmpFilePath, err := oj.Compile(cfg.Lang, cfg.CodePath, cfg.TmpDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	// 运行各个测试样例，得到初步的结果
	var rst oj.Result
	case_i, case_rst, err := oj.Run(p_cfg, cfg.ProblemDir, tmpFilePath, &rst)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Printf("%+v\t%+v", case_i,case_rst)
	
	// 用户的输出与标程的输出进行对照，得到最终结果
	err = oj.Judge(p_cfg, case_i, case_rst, tmpFilePath, cfg.ProblemDir, &rst)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	// 清理测试用例、用户输出、标程输出文件
	err = oj.CleanUp(cfg.ProblemDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 输出测评结果
	rst.String(p_cfg.Mode)
}
```
