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
		AppID:    "xxx",
		AppName:  "cossim",
		Platform: "huawei",
		Token: []string{
			"xxx",
		},
		Data: dto.HuaweiPushRequestData{
			Title:       "cossim",
			Content:     "hello",
			Category:    "IM",
			Priority:    "normal",
			TTL:         "86400s",
			Icon:        "",
			Img:         "",
			Sound:       "",
			Foreground:  true,
			Development: true,
			ClickAction: dto.ClickAction{
				Action:     3,
				Activity:   "",
				Url:        "",
				Parameters: nil,
			},
			Badge: dto.BadgeNotification{
				AddNum: 1,
				//SetNum: 1,
				//Class:  "1",
			},
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
