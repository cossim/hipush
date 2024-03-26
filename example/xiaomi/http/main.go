package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	url    = "http://127.0.0.1:7070/api/v1/push"
	method = "POST"
)

func main() {
	payload := []byte(`{
    "platform": "xiaomi",
    "token": [
        "xxx"
    ],
    "app_id": "2882303761520159644",
    "app_name": "cossim",
    "data": {
        "title": "测试标题",
        "subtitle": "测试子标题",
        "content": "测试内容",
        "foreground": true
    },
    "option": {
        "dry_run": true,
		"development": false,
        "retry": 0,
        "retry_interval": 0
    }
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
