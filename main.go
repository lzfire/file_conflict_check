package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	fc "file_conflict_check/file_check"
	rf "file_conflict_check/read_conf"
)

var (
	svnURL       string
	ssuName      string
	md5Str       string
	appversion   string
	dirPath      string
	conflictFile string
)

//解析命令行参数
func parseArg() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current dir path failed:%s", err)
	}
	flag.StringVar(&dirPath, "d", dir, "the dir path you want to search")
	flag.StringVar(&svnURL, "u", "", "the svn request url")
	flag.StringVar(&ssuName, "p", "", "the package you want to analay")
	flag.StringVar(&appversion, "a", "", "the file store appversion")
	flag.StringVar(&md5Str, "m", "9074170c68cdbb5dca8BA43226417741", "the conflict file md5 value")
	flag.StringVar(&conflictFile, "f", "README.md", "the conflict file name")
}

//主函数
func main() {
	//命令行处理
	parseArg()
	flag.Parse()

	if dirPath != "" {
		dirName := filepath.Dir(dirPath)
		//收集需要查找目录下的所有文件
		fileslist := fc.FileList{}
		if err := fileslist.GetAllFile(dirName); err != nil {
			log.Fatalf("getAllFile failed in dir:%s, err:%s", dirName, err)
		}
		isConflict := false
		//遍历每个文件的md5值，并做比较，找到冲突的md5则直接返回，否则提醒所查找的目录或包没有该冲突文件
		for _, file := range fileslist {
			if conflictFile != "" {
				pName := strings.LastIndex(file, "/")
				if file[pName+1:] == conflictFile {
					if isConflict = fc.CheckFileConflict(file, md5Str); isConflict {
						log.Printf("get the conflict file:%s", file)
						break
					}
				}
			}
		}
		if !isConflict {
			log.Printf("no file conflict in:%s with md5:%s", dirName, md5Str)
		}
	}
	//从appversion中读取每一个包名，并存放在切片中
	if appversion != "" {
		fc.ReadLineFile(appversion)
	}
	cfg, err2 := rf.Load("./read_conf/testdata/test.ini")
	if err2 != nil {
		log.Fatalf("rf Load failed:%s", err2)
	}
	log.Printf("rf.Load cfg:%#v", cfg)

	//根据前面分割好的包名，请求到改包的ssu包
	if svnURL != "" {
		result := fc.HTTPGet(svnURL)
		log.Printf("http get result:%s", string(result))
	}

}
