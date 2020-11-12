package file_check

import (
	"io/ioutil"
	"log"
	"net/http"
)

//HTTPGet 使用http get请求拿文件或数据
func HTTPGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("http get failed,err:%s\n", err)
		return nil
	}
	defer resp.Body.Close()
	log.Printf("resp.status:%s\n", resp.Status)
	if resp.Status != "200 OK" {
		log.Printf("http get failed\n")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("iotuil readall from resp body failed,err:%s\n", err)
		return nil
	}
	return body
}
