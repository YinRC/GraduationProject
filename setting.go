package oj

import (
	"fmt"
	"encoding/json"
	"os"
	"flag"
)

const (
	AC = "AC"			// ok
	PE = "PE"			// 格式错误
	TLE = "TLE"			// 超时错误
	MLE = "MLE"			// 超过内存空间限制
	WA = "WA"			// 答案错误
	RE = "RE"			// 运行时错误
	OLE = "OLE"			// 输出超过限制
	CE = "CE"			// 编译错误
	SE = "SE"			// 系统错误
	PC = "PC"			// 特判分数
)

const (
	NormalMode = 0
	GrammarMode = 1		// 语法结构模式
	PointNumMode = 2	// 浮点数模式
	IncompleteMode = 3	// 完成比例模式
)

type Config struct {
	Lang string

	ProblemDir string	// 题目路径（定位到文件夹）
	CodePath string		// 代码路径（定位到文件）
	// WorkDir string		// 工作路径（测试样例、测试输出、标程输出）
	TmpDir string		// 临时文件路径（用户代码、用户程序）
}


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


// 从文件config.json中读取设置的参数并应用
func (cfg *Config) GetConfig() (err error) {
	filePtr, err := os.Open("config.json")
	defer filePtr.Close()
	if err != nil {
		return fmt.Errorf("cfg open config fail")
	}

	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(cfg)
	if err != nil {
		return fmt.Errorf("cfg decode config fail")
	}
	return nil
}

// 得到题目的设置参数
func (p_cfg *Problem) GetProblemConfig(path string) (err error) {
	filePtr, err := os.Open(path)
	defer filePtr.Close()
	if err != nil {
		return fmt.Errorf("p_cfg open config fail")
	}

	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(p_cfg)
	if err != nil {
		return fmt.Errorf("p_cfg decode config fail")
	}
	return nil
}

// [可选]是否从命令行设置config.json文件
func (cfg *Config) SetFromCmd() (err error) {
	var (
		lang string

		problem string
		code string
		//work string
		tmp string
	)

	flag.StringVar(&lang, "l", "cpp", "编译所用的语言(小写)-默认为C++语言" )

	flag.StringVar(&problem, "p", "problem/1", "题目路径")
	flag.StringVar(&code, "c", "code/usr.c", "代码路径")
	// flag.StringVar(&work, "w", "work", "工作路径")
	flag.StringVar(&tmp, "tmp", "tmp", "临时代码、程序存放路径")


	flag.Parse()
	fmt.Printf("编程语言：%v\n题目路径：%v\n代码路径：%v\n临时路径：%v\n", lang, problem, code, tmp)
	*cfg = Config {
		Lang:			lang,

		ProblemDir:		problem,
		CodePath:		code,
		//WorkDir:		work,
		TmpDir:			tmp,
	}
	err = cfg.setConfig()
	if err != nil {
		return err
	}
	return nil
}

// 将设置参数记录在config.json文件里
func (cfg *Config) setConfig() (err error) {
	_ = os.Remove("config.json")
	filePtr, err := os.OpenFile("config.json", os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cfg open/create config fail")
	}
	defer filePtr.Close()

	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(*cfg)
	if err != nil {
		return fmt.Errorf("cfg encode config fail")
	}
	return nil
}
