package main

import (
	"encoding/json"
	"fmt"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	url    = "http://127.0.0.1:7070/api/v1/push"
	method = "POST"
)

func main() {
	payload := dto.PushRequest{
		AppID:    "xxx",
		AppName:  "cossim",
		Platform: "ios",
		Token: []string{
			"xxx",
		},
		Data: dto.APNsPushRequest{
			Title:            "cossim",
			Content:          "hello",
			Topic:            "com.hitosea.cossim",
			CollapseID:       "",
			ApnsID:           uuid.New().String(),
			Priority:         "normal",
			PushType:         "alert",
			URLArgs:          nil,
			TTL:              time.Now().Add(30 * time.Minute).Unix(),
			Badge:            1,
			Development:      false,
			MutableContent:   false,
			ContentAvailable: false,
			Sound: map[string]interface{}{
				"critical": 1,
				"volume":   4.5,
				"name":     "",
			},
			Data: nil,
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
