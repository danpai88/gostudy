package gostudy

import (
	"github.com/danpai88/gostudy/lib"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

func Intval(str string, def ... int) int {
	tmp, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return int(tmp)
}

func Strtotime(formatTimeStr string) int {
	formatTime, err := time.Parse("2006-01-02 15:04:05",formatTimeStr)
	if err != nil {
		return 0
	}
	return int(formatTime.Unix())
}


func DelAllHtmlTag(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func FileGetContent(file string) string {
	var cyIo gostudy.CyIO
	return cyIo.FileGetContent(file)
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}