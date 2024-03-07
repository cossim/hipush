package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthRequest struct {
	AppID     int64  `json:"appId"`
	AppKey    string `json:"appKey"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type AuthResponse struct {
	Result    int    `json:"result"`
	Desc      string `json:"desc"`
	AuthToken string `json:"authToken,omitempty"`
}

type UnicastRequest struct {
	AppID          int64                  `json:"appId"`
	RegID          string                 `json:"regId,omitempty"`
	Alias          string                 `json:"alias,omitempty"`
	NotifyType     int                    `json:"notifyType"`
	Title          string                 `json:"title"`
	Content        string                 `json:"content"`
	TimeToLive     int                    `json:"timeToLive,omitempty"`
	SkipType       int                    `json:"skipType"`
	SkipContent    string                 `json:"skipContent,omitempty"`
	NetworkType    int                    `json:"networkType,omitempty"`
	ClientCustom   map[string]interface{} `json:"clientCustomMap,omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty"`
	RequestID      string                 `json:"requestId"`
	PushMode       int                    `json:"pushMode,omitempty"`
	ForegroundShow bool                   `json:"foregroundShow,omitempty"`
}

type UnicastResponse struct {
	Result      int    `json:"result"`
	Desc        string `json:"desc"`
	TaskID      string `json:"taskId,omitempty"`
	InvalidUser struct {
		Status int    `json:"status"`
		UserID string `json:"userid"`
	} `json:"invalidUser,omitempty"`
}

func main() {
	// 认证请求
	authToken, err := sendAuthRequest()
	if err != nil {
		fmt.Println("Error sending auth request:", err)
		return
	}

	// 发送单播消息请求
	err = sendUnicastRequest(authToken)
	if err != nil {
		fmt.Println("Error sending unicast request:", err)
		return
	}
}

func generateSign(appID int64, appKey string, timestamp int64, appSecret string) string {
	signParams := fmt.Sprintf("%d%s%d%s", appID, appKey, timestamp, appSecret)
	hasher := md5.New()
	hasher.Write([]byte(signParams))
	sign := hex.EncodeToString(hasher.Sum(nil))
	return sign
}

func sendAuthRequest() (string, error) {
	appID := int64(105562603)
	appKey := "ea6a826256eb2896a2f36743d859dbd0"
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	appSecret := "b14a60d5-bcb1-4836-8d88-15f2ac040eee"

	sign := generateSign(appID, appKey, timestamp, appSecret)

	requestBody := AuthRequest{
		AppID:     appID,
		AppKey:    appKey,
		Timestamp: timestamp,
		Sign:      sign,
	}

	jsonBody, _ := json.Marshal(requestBody)

	resp, err := http.Post("https://api-push.vivo.com.cn/message/auth", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return "", err
	}

	return authResponse.AuthToken, nil
}

func sendUnicastRequest(authToken string) error {
	requestBody := UnicastRequest{
		AppID:          105562603,
		RegID:          "v2-CQanxnrM-uZu6i_y_E3PpymvRvSJhhFjaOwQxbGE-jJ0BNLa0IUm",
		NotifyType:     1,
		Title:          "标题1",
		Content:        "内容1",
		SkipType:       2,
		SkipContent:    "http://www.vivo.com",
		RequestID:      "25509283-3767-4b9e-83fe-b6e55ac6b123",
		PushMode:       1,
		ForegroundShow: true,
	}

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "https://api-push.vivo.com.cn/message/send", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authToken", authToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var unicastResponse UnicastResponse
	err = json.NewDecoder(resp.Body).Decode(&unicastResponse)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %d, Desc: %s, TaskID: %s\n", unicastResponse.Result, unicastResponse.Desc, unicastResponse.TaskID)

	return nil
}
