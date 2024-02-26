package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// MiPushService 实现小米推送服务的具体逻辑
type MiPushService struct {
	AppSecret  string // 小米 App Secret
	Endpoint   string // 小米推送 API 端点
	APIVersion string // 小米推送 API 版本
}

// 推送消息结构体
type XiaomiMessage struct {
	RegistrationID string `json:"registration_id,omitempty"`
	Alias          string `json:"alias,omitempty"`
	UserAccount    string `json:"user_account,omitempty"`
	Topic          string `json:"topic,omitempty"`
	Topics         string `json:"topics,omitempty"`
	TopicOp        string `json:"topic_op,omitempty"`
	Payload        string `json:"payload"`
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	NotifyType     int    `json:"notify_type,omitempty"`
	TimeToLive     int64  `json:"time_to_live,omitempty"`
	NotifyID       int    `json:"notify_id,omitempty"`
	Extra          struct {
		SoundURI         string `json:"sound_uri,omitempty"`
		Ticker           string `json:"ticker,omitempty"`
		NotifyForeground string `json:"notify_foreground,omitempty"`
		NotifyEffect     string `json:"notify_effect,omitempty"`
		IntentURI        string `json:"intent_uri,omitempty"`
		WebURI           string `json:"web_uri,omitempty"`
	} `json:"extra,omitempty"`
}

// 发送单条消息
func (m *MiPushService) Send(message string, receiver string) error {
	return m.sendToEndpoint("/message/regid", message, receiver)
}

// 发送多条消息
func (m *MiPushService) MulticastSend(message string, receivers []string) error {
	return m.sendToEndpoint("/message/regid", message, strings.Join(receivers, ","))
}

// 订阅主题
func (m *MiPushService) Subscribe(topic string, receiver string) error {
	return m.sendToEndpoint("/topic/subscribe", topic, receiver)
}

// 取消订阅主题
func (m *MiPushService) Unsubscribe(topic string, receiver string) error {
	return m.sendToEndpoint("/topic/unsubscribe", topic, receiver)
}

// 向指定主题发送消息
func (m *MiPushService) SendToTopic(message string, topic string) error {
	return m.sendToEndpoint("/message/topic", message, topic)
}

// 向指定条件发送消息
func (m *MiPushService) SendToCondition(message string, condition string) error {
	return m.sendToEndpoint("/message/condition", message, condition)
}

// 检查设备可用性
func (m *MiPushService) CheckDevice(deviceToken string) bool {
	// 此处暂不实现，直接返回 true
	return true
}

// 获取推送服务名称
func (m *MiPushService) Name() string {
	return "MiPushService"
}

// 发送消息到指定端点
func (m *MiPushService) sendToEndpoint(endpoint, message, receiver string) error {
	url := m.Endpoint + m.APIVersion + endpoint

	xiaomiMessage := XiaomiMessage{
		RegistrationID: receiver,
		Payload:        message,
	}

	jsonStr, err := json.Marshal(xiaomiMessage)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("key=%s", m.AppSecret))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	return nil
}
