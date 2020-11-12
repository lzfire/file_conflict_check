package file_check

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//SEP separator
var SEP string

//FileList 收集需要查找目录下的所有文件
type FileList []string

/*
//FileCheck 定义专属结构体
type FileCheck struct {
	FileList
}
*/

//GetAllFile 根据路径获取路径下的所有文件
func (f *FileList) GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			if err = f.GetAllFile(pathname + SEP + fi.Name()); err != nil {
				return err
			}
		} else {
			*f = append(*f, pathname+SEP+fi.Name())
		}
	}
	return nil
}

//CheckFileConflict 检查某文件的md5值是否匹配
//这样单独抽出来能够防止文件指针打开太多而没释放
func CheckFileConflict(filepath string, md5Str string) bool {
	fileP, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("file: %s open failed with err:%s\n", filepath, err)
		return false
	}
	defer fileP.Close()
	md5Value := md5.New()
	_, err = io.Copy(md5Value, fileP)
	if err != nil {
		fmt.Printf("io.Copy error:%s\n", err)
		return false
	}
	fileMd5 := hex.EncodeToString(md5Value.Sum(nil))
	if strings.EqualFold(md5Str, fileMd5) {
		return true
	}
	return false
}

func init() {
	SEP = GetPathSeparator()
}
