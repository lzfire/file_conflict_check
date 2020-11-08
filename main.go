package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//FileList 收集需要查找目录下的所有文件
type FileList []string

//SEP seprator
var SEP string

//GetPathSeprator 获取不同系统下的路径分隔符
func GetPathSeprator() string {
	sysType := runtime.GOOS
	var sep string

	switch sysType {
	case "windows":
		sep = "\\"
	case "linux":
		fallthrough
	case "drawin":
		fallthrough
	default:
		sep = "/"
	}
	return sep
}

func (f *FileList) getAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			if err = f.getAllFile(pathname + SEP + fi.Name()); err != nil {
				return err
			}
		} else {
			*f = append(*f, pathname+SEP+fi.Name())
		}
	}
	return nil
}

func main() {
	//自测变量设置
	md5Str := "9074170c68cdbb5dca8BA43226417741"
	fileName := "README.md"
	pathStr, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Printf("filepath.Abs error:%s\n", fileName)
		os.Exit(1)
	}
	SEP = GetPathSeprator()
	dirName := filepath.Dir(pathStr)
	//收集需要查找目录下的所有文件
	fileslist := FileList{}
	if err = fileslist.getAllFile(dirName); err != nil {
		log.Fatalf("getAllFile failed in dir:%s, err:%s\n", dirName, err)
	}
	//遍历每个文件的md5值，并做比较，找到冲突的md5则直接返回，否则提醒所查找的目录或包没有该冲突文件
	for _, file := range fileslist {
		fileP, err := os.Open(file)
		if err != nil {
			fmt.Printf("file: %s open failed with err:%s\n", file, err)
			continue
		}
		defer fileP.Close()
		md5Value := md5.New()
		_, err = io.Copy(md5Value, fileP)
		if err != nil {
			fmt.Printf("io.Copy error:%s\n", err)
			continue
		}
		fileMd5 := hex.EncodeToString(md5Value.Sum(nil))
		if strings.EqualFold(md5Str, fileMd5) {
			fmt.Printf("get the conflict file:%s\n", file)
			return
		}
	}
	fmt.Printf("no file conflict in:%s with md5:%s\n", dirName, md5Str)
}
