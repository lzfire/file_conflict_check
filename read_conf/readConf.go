package readConf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"unicode"
)

type section struct {
	keyValue map[string]string
}

type config struct {
	sectionList []string            //存key
	sections    map[string]*section //实际存储数据的地方（形成key和value的映射）
}

func (cfg *config) newSection(name string) error {
	if len(name) == 0 { //如果名字长度为0
		return errors.New("section name is empty")
	}
	if cfg.sections[name] != nil { //如果sections中某一项已经有数据
		return errors.New("already has same section name")

	}
	cfg.sectionList = append(cfg.sectionList, name) //扩展
	cfg.sections[name] = &section{
		keyValue: make(map[string]string),
	}
	return nil
}

func (cfg *config) init() {
	cfg.sections = make(map[string]*section) //初始化map
}

//ListenFunc 监听回调函数
type ListenFunc func(string)

//Listener 监听接口
type Listener interface {
	listen(inifile string)
}

func listen() {

}

//Watch 对外提供的借口
func Watch(filename string, listener Listener) (*config, error) {
	var cfg *config
	listener.listen(filename)
	var e error = nil
	cfg, e = load(filename)
	return cfg, e
}

func load(filename string) (*config, error) {
	f, err := os.Open(filename) //打开文件
	if err != nil {             //如果打开失败--遇到错误
		return nil, err
	}
	defer f.Close() //在函数执行结束后关闭文件

	buf := bufio.NewReader(f) //从文件里读取数据
	cfg, err := parse(buf)    //执行parse函数，详情见parse.go
	return cfg, err
}

func getLine(buf *bufio.Reader, isEOF *bool) ([]byte, error) {
	line, err := buf.ReadBytes('\n') //按行读取数据
	if err == io.EOF {
		*isEOF = true
		err = nil
	} else if err != nil { //如果读取失败返回错误
		return nil, err
	}
	line = bytes.TrimLeftFunc(line, unicode.IsSpace)
	return line, err

}

//解析section的名字
func parseSecName(line []byte, cfg *config) (string, error) {
	close := bytes.LastIndexByte(line, ']') //获得最后‘]’的下标
	if close == -1 {                        //如果缺少‘]’，报错
		return "", fmt.Errorf("unclosed section: %s", line)
	}
	secName := string(line[1:close]) //获取section的名字
	err := cfg.newSection(secName)   //开辟新空间并获得错误信息
	if err != nil {                  //如果存在错误，返回错误信息
		return "", err
	}
	return secName, nil
}

//解析key的名字
func parseKeyName(line string) (string, int, error) {
	end := strings.IndexAny(line, keyValueDelim) //确定key-alue的分割符的下标
	if end < 0 {                                 //如果下表<0，即存在错误
		return "", -1, fmt.Errorf("delimiter(%s) not found", keyValueDelim)
	}
	return strings.TrimSpace(line[0:end]), end + 1, nil //这里返回的下标为分割符下标+1
}

func (sec *section) newKeyValue(keyName string, value string) error {
	if _, ok := sec.keyValue[keyName]; ok {
		return fmt.Errorf("key(%v) already exists", keyName)
	}
	sec.keyValue[keyName] = value
	return nil
}

//解析value
func parseValue(line string) (string, error) {
	line = strings.TrimSpace(line) //去空格
	if len(line) == 0 {            //如果长度为0
		return "", nil
	}
	i := strings.IndexAny(line, "#;") //处理一种情况‘key = value # this is comment’也就是注释在语句后方
	if i > -1 {
		line = strings.TrimSpace(line[:i])

	}
	return line, nil //由于之前直接返回的是key-value分割符下标+1，就是value开始的地方
}

func parse(reader *bufio.Reader) (*config, error) {
	var cfg config
	cfg.init()     //初始化
	isEOF := false //文件是否关闭
	secName := defaultSection
	for !isEOF { //如果文件没有关闭则执行循环
		line, err := getLine(reader, &isEOF)
		if err != nil {
			return nil, err
		}
		if len(line) == 0 { //跳过空行
			continue
		}
		if line[0] == notes[0] { //跳过注释
			continue
		}
		if line[0] == '[' { //如果第0位是‘[’，即secname的前面一位
			secName, err = parseSecName(line, &cfg) //解析secname
			if err != nil {                         //如果解析错误返回错误
				return nil, err
			}
			continue
		}
		if len(cfg.sectionList) == 0 { //没有空间则开辟新空间
			err = cfg.newSection(secName)
			if err != nil {
				return nil, err
			}
		}
		keyName, offset, err := parseKeyName(string(line)) //解析keyname
		if err != nil {
			return nil, err
		}
		value, err := parseValue(string(line[offset:])) //解析value
		if err != nil {
			return nil, err
		}
		err = cfg.sections[secName].newKeyValue(keyName, value) //存key和value
		if err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

func init() { //根据系统来确定注释行符号
	switch runtime.GOOS {
	case "windows": //如果是windows系统
		notes = ";"
		lineBreaker = "\r\n"
	case "linux": //如果是linux系统
		notes = "#"
		lineBreaker = "\n"
	}
}
