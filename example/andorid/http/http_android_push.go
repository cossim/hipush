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
    "platform": "android",
    "token": [
        "xxx"
    ],
    "app_name": "cossim",
    "data": {
        "title": "cossim",
        "content": "hello"
    },
    "option": {
        "dry_run": true,
		"development": false,
        "retry": 1,
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
