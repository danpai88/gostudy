package gostudy

import (
	"strings"
	"time"
)

type CyTime struct {
	
}

//格式化时间
//Date("Y-m-d H:i:s", ...时间戳)
func (this CyTime) Date(format string, timesteamp ... int64) string {
	var tt time.Time
	if len(timesteamp) > 0 && timesteamp[0] > 0 {
		tt = time.Unix(timesteamp[0], 0)
	}else {
		tt = time.Now()
	}

	var defaultFormat = make(map[string]string, 0)
	defaultFormat["2006"] = "Y"
	defaultFormat["01"] = "m"
	defaultFormat["02"] = "d"
	defaultFormat["15"] = "H"
	defaultFormat["04"] = "i"
	defaultFormat["05"] = "s"

	for key, val := range defaultFormat {
		format = strings.Replace(format, val, key, -1)
	}
	return tt.Format(format)
}

//获取当前时间
func (this CyTime) CurrentDate() string {
	var t = time.Now()
	return this.Date("Y-m-d H:i:s", t.Unix())
}

//获取当前时间戳
func (this CyTime) Time() int64 {
	var t = time.Now()
	return t.Unix()
}