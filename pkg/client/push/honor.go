package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HonorPushClient struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	pushUrl      string
	authUrl      string
	httpClient   *http.Client
}

type SendMessageRequest struct {
	Data         string         `json:"data,omitempty"`
	Notification *Notification  `json:"notification,omitempty"`
	Android      *AndroidConfig `json:"android,omitempty"`
	Token        []string       `json:"token,omitempty"`
}

// Notification 结构体用于表示通知栏消息内容
type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image,omitempty"`
}

// AndroidConfig 结构体用于表示Android消息推送控制参数
type AndroidConfig struct {
	TTL            string               `json:"ttl,omitempty"`
	BiTag          string               `json:"biTag,omitempty"`
	Data           string               `json:"data,omitempty"`
	Notification   *AndroidNotification `json:"notification,omitempty"`
	TargetUserType int                  `json:"targetUserType,omitempty"`
}

// AndroidNotification 结构体用于表示Android通知栏消息结构体
type AndroidNotification struct {
	Title       string             `json:"title"`
	Body        string             `json:"body"`
	ClickAction *ClickAction       `json:"clickAction"`
	Image       string             `json:"image,omitempty"`
	Style       int                `json:"style,omitempty"`
	BigTitle    string             `json:"bigTitle,omitempty"`
	BigBody     string             `json:"bigBody,omitempty"`
	Importance  string             `json:"importance,omitempty"`
	When        string             `json:"when,omitempty"`
	Buttons     []*Button          `json:"buttons,omitempty"`
	Badge       *BadgeNotification `json:"badge,omitempty"`
	NotifyID    int                `json:"notifyId,omitempty"`
	Tag         string             `json:"tag,omitempty"`
	Group       string             `json:"group,omitempty"`
}

// ClickAction 结构体用于表示消息点击行为
type ClickAction struct {
	Type   int    `json:"type"`
	Intent string `json:"intent,omitempty"`
	URL    string `json:"url,omitempty"`
	Action string `json:"action,omitempty"`
}

// Button 结构体用于表示通知栏消息动作按钮
type Button struct {
	Name       string `json:"name"`
	ActionType int    `json:"actionType"`
	IntentType int    `json:"intentType,omitempty"`
	Intent     string `json:"intent,omitempty"`
	Data       string `json:"data,omitempty"`
}

// BadgeNotification 结构体用于表示Android通知消息角标控制
type BadgeNotification struct {
	AddNum     int    `json:"addNum,omitempty"`
	SetNum     int    `json:"setNum,omitempty"`
	BadgeClass string `json:"badgeClass"`
}

// SendMessageResponse 结构体用于表示SendMessage函数的响应
type SendMessageResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message,omitempty"`
	Data    SendMessageResponseData `json:"data,omitempty"`
}

// SendMessageResponseData 结构体用于表示SendMessage函数的响应数据
type SendMessageResponseData struct {
	SendResult   bool     `json:"sendResult"`
	RequestId    string   `json:"requestId,omitempty"`
	FailTokens   []string `json:"failTokens,omitempty"`
	ExpireTokens []string `json:"expireTokens,omitempty"`
}

type HonorPushOption func(*HonorPushClient)

func WithPushUrl(url string) HonorPushOption {
	return func(sdk *HonorPushClient) {
		sdk.pushUrl = url
	}
}

func WithAuthUrl(url string) HonorPushOption {
	return func(sdk *HonorPushClient) {
		sdk.authUrl = url
	}
}

// NewHonorPush 函数用于创建一个新的荣耀推送SDK实例
func NewHonorPush(clientID, clientSecret string, options ...HonorPushOption) *HonorPushClient {
	sdk := &HonorPushClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		pushUrl:      "https://push-api.cloud.hihonor.com/api/v1/",
		authUrl:      "https://iam.developer.hihonor.com/auth/token",
		httpClient:   &http.Client{},
	}

	for _, opt := range options {
		opt(sdk)
	}

	return sdk
}

// 获取Access Token
func (hpc *HonorPushClient) getAccessToken(ctx context.Context) error {
	requestBody := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", hpc.ClientID, hpc.ClientSecret)
	req, err := http.NewRequestWithContext(ctx, "POST", hpc.authUrl, bytes.NewBufferString(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hpc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	hpc.AccessToken = tokenResponse.AccessToken

	return nil
}

// SendMessage 方法用于推送消息
func (hpc *HonorPushClient) SendMessage(ctx context.Context, appID string, requestBody *SendMessageRequest) (*SendMessageResponse, error) {
	if hpc.AccessToken == "" {
		if err := hpc.getAccessToken(ctx); err != nil {
			return nil, err
		}
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	requestURL := hpc.pushUrl + appID + "/sendMessage"
	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hpc.AccessToken)
	req.Header.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	resp, err := hpc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应
	var response SendMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
