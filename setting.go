package oj

import (
	"fmt"
	"encoding/json"
	"os"
	"flag"
)

const (
	AC = iota	// ok
	PE			// 格式错误
	TLE			// 超时错误
	MLE			// 超过内存空间限制
	WA			// 答案错误
	RE			// 运行时错误
	OLE			// 输出超过限制
	CE			// 编译错误
	SE			// 系统错误
)

type Config struct {
	Lang string

	CodePath string
	CaseInPath string
	CaseOutPath string
	WorkDir string

	Time int
	Memory int
	OutputSize int
}

// 从文件config.json中读取设置的参数并应用
func (cfg *Config) GetConfig() (err error) {
	filePtr, err := os.Open("config.json")
	defer filePtr.Close()
	if err != nil {
		return fmt.Errorf("open config fail")
	}

	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(cfg)
	if err != nil {
		return fmt.Errorf("decode config fail")
	}
	return nil
}

// [可选]是否从命令行设置config.json文件
func (cfg *Config) SetFromCmd() (err error) {
	var (
		lang string
		codePath string
		caseInPath string
		caseOutPath string
		workDir string

		time int
		memory int
		outputSize int
	)

	flag.StringVar(&lang, "l", "cpp", "编译所用的语言(小写)-默认为C++语言" )
	flag.StringVar(&codePath, "code", "code/ac.cpp", "代码的路径")
	flag.StringVar(&caseInPath, "in", "data/1.in", "测试用例路径")
	flag.StringVar(&caseOutPath, "out", "data/1.out", "测试用例路径")
	flag.StringVar(&workDir, "work", "work", "工作路径")
	flag.IntVar(&time, "t", 1000, "程序时间限制")
	flag.IntVar(&memory, "m", 5000, "程序空间限制")
	flag.IntVar(&outputSize, "s", 100*1024*1024, "程序输出大小限制")

	flag.Parse()
	fmt.Printf("编程语言：%v\n代码路径：%v\n用例输入：%v\n用例输出：%v\n工作目录：%v\n程序时间限制（ms）：%v\n程序空间限制（KB）：%v\n程序输出限制（B）：%v\n", lang, codePath, caseInPath, caseOutPath, workDir, time, memory, outputSize)
	*cfg = Config {
		Lang:			lang,

		CodePath:		codePath,
		CaseInPath:		caseInPath,
		CaseOutPath:	caseOutPath,
		WorkDir:		workDir,

		Time:			time,
		Memory:			memory,
		OutputSize:		outputSize,
	}
	err = cfg.setConfig()
	if err != nil {
		return err
	}
	return nil
}

// 将设置参数记录在config.json文件里
func (cfg *Config) setConfig() (err error) {
	_ = os.Remove("config/config.json")
	filePtr, err := os.OpenFile("config.json", os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open/create config fail")
	}
	defer filePtr.Close()

	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(*cfg)
	if err != nil {
		return fmt.Errorf("encode config fail")
	}
	return nil
}
