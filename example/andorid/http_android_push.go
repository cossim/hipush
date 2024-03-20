package main

import (
	"encoding/json"
	"fmt"
	"github.com/cossim/hipush/api/http/v1/dto"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	url    = "http://127.0.0.1:7070/api/v1/push"
	method = "POST"
)

func main() {
	payload := dto.PushRequest{
		AppID:    "com.dootask.task",
		AppName:  "cossim",
		Platform: "android",
		Token: []string{
			"cb7f8a974eec5fbb2e36762fcb78e51327bcef4822d600f17e8f9bd845af1e12",
		},
		Data: dto.AndroidPushRequestData{
			Title:      "cossim",
			Content:    "hello",
			TTL:        "10m",
			Topic:      "",
			Priority:   "normal",
			CollapseID: "",
			Condition:  "",
			Sound:      "",
			Image:      "",
			Data:       nil,
		},
		Option: dto.PushOption{
			DryRun:        false,
			Retry:         1,
			RetryInterval: 1,
		},
	}

	// Marshal the request object to JSON
	reqBody, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(reqBody)))
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
