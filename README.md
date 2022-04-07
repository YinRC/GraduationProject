# GraduationProject
**A testing machine for green hand ！(Linux only)**

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
```
## 结构体部分

### 测评机设置
+ 对文件路径、代码编译过程、程序运行过程进行设置并持久化到config.json。
+ 设置config可以通过命令行也可以直接修改config.json文件
```go
// setting.go
type Config struct {

	Lang string         // 编译所用的语言

	CodePath string     // 用户提交代码所在的路径
	CaseInPath string   // 测试用例输入文件的路径（一次提问）
	CaseOutPath string  // 测试用例输出文件的路径（一次回答）
	WorkDir string      // 工作目录（存放用户代码的运行结果和运行错误）

	Time int            // 用户代码运行时程序的时间限制（ms）
	Memory int          // 用户代码运行时程序的空间限制（KB）
	OutputSize int      // 用户代码运行时程序的输出限制（B）
}
```
### 编译器临时文件路径
后期会探索特殊判题的方式，所以将编译前后的文件存放在tmp文件夹中。
```go
// compiler.go
type tmpFilePath struct {

	tmpCodePath string      // 用户提交代码的路径（临时）
	tmpProgramPath string   // 编译完成后用户程序的路径（临时）
}
```
### 测评机判题结果
```go
type Result struct {

	Flag int        // 判题结果（AC WA PE...）

	Time int        // 判题过程中程序所用的时间（ms）
	Memory int      // 判题过程中程序所用的空间（KB）

	SE_log string   // 系统错误日志信息
}
```
### 可用的包外调用函数
可以据此编写适用于个人的main.go
```go
// compiler.go	编译用户代码
func Compile(lang string, codePath string, tmpPath string) (tmp tmpFilePath, err error)
// executor.go	运行二进制文件
func RunBinFile(cfg Config, tmp tmpFilePath, rst *Result) error
// utils.go	输出判题结果（Result结构体）
func String(result Result) {
	b, _ := json.Marshal(result)
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	fmt.Printf("%+v\n", result)
}
// setting.go
func (cfg *Config) GetConfig() (err error)	// 从文件config.json中读取设置的参数并应用
func (cfg *Config) SetFromCmd() (err error)	// [可选]是否从命令行设置config.json文件
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
## 命令行选项
```go
// setting.go
flag.StringVar(&lang, "l", "cpp", "编译所用的语言(小写)-默认为C++语言" )
flag.StringVar(&codePath, "code", "code/ac.cpp", "代码的路径")
flag.StringVar(&caseInPath, "in", "data/1.in", "测试用例路径")
flag.StringVar(&caseOutPath, "out", "data/1.out", "测试用例路径")
flag.StringVar(&workDir, "work", "work", "工作路径")
flag.IntVar(&time, "t", 1000, "程序时间限制(ms)")
flag.IntVar(&memory, "m", 5000, "程序空间限制(KB)")
flag.IntVar(&outputSize, "s", 100*1024*1024, "程序输出大小限制(B)")
```

## main.go
```go
package main

import (
	"fmt"
	"github.com/GraduationProject"
)

func main() {
	// [可选]从命令行中读取参数更改config.json
	//err := cfg.SetFromCmd()
	//if err != nil {
		//fmt.Println(err.Error())
	//}
	
	// 读取并应用config.json，变量cfg接收设置信息
	var cfg oj.Config
	err := cfg.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	
	// 编译用户代码
	tmpFilePath, err := oj.Compile(cfg.Lang, cfg.CodePath, "tmp")
	if err != nil {
		fmt.Println(err.Error())
	}
	
	// 运行用户程序，变量rst收集测评结果
	var rst oj.Result
	err = oj.RunBinFile(cfg, tmpFilePath, &rst)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 终端输出测评结果
	oj.String(rst)
}
```
