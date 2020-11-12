package file_check

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
)

//GetPathSeparator 获取不同系统下的路径分隔符
func GetPathSeparator() string {
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

//ReadLineFile 从传入的文件路径读取每一行
func ReadLineFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file failed:%s\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
