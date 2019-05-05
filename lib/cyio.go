package gostudy

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CyIO struct {
	CyTime CyTime
}

//向文件写数据
func (this CyIO) FilePutContent(file string, data string) bool {
	var filePath = file

	fl, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalln(err.Error())
		return false
	}

	_, err = fl.WriteString(data)

	fl.Close()

	if err != nil {
		log.Fatalln(err.Error())
		return false
	}
	return true
}

//判断文件或者文件夹是否存在
func (this CyIO) FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return false
}

//获取当前程序执行的路径
func (this CyIO) GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}


//生成nginx log
func (this CyIO) NginxLog(req *http.Request) string {
	var strMap = make([]string, 0)

	var ip string
	if t_ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		ip = t_ip
	}

	strMap = append(strMap, ip)
	strMap = append(strMap, "[" + this.CyTime.CurrentDate() + "]")
	strMap = append(strMap, req.URL.Scheme + req.Host)
	strMap = append(strMap, req.Method)
	strMap = append(strMap, req.RequestURI)
	strMap = append(strMap, req.Proto)
	strMap = append(strMap, strconv.Itoa(int(req.ContentLength)))
	strMap = append(strMap, req.Referer())
	strMap = append(strMap, req.UserAgent())
	return strings.Join(strMap, " ")
}

//读文件
func (this CyIO) FileGetContent(file string) string {
	path, _ := os.Getwd()
	var filePath = path + "/" + file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return string(content)
	}
	return "";
}