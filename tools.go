package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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
