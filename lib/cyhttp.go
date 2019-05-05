package gostudy

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Cyhttp struct {

}

//发起http请求
func (this Cyhttp) Get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return string(html)
}

//发起post请求
func (this Cyhttp) Post(url string) {

}