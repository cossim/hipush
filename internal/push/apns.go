package push

import (
	"context"
	"errors"
	"github.com/cossim/hipush/internal/notify"
	"github.com/sideshow/apns2"
	"log"
	"net/http"
	"sync"
)

var (
	// MaxConcurrentIOSPushes pool to limit the number of concurrent iOS pushes
	MaxConcurrentIOSPushes chan struct{}
)

// APNsService 实现APNs推送
type APNsService struct {
	client *apns2.Client
}

func (a *APNsService) Send(ctx context.Context, request interface{}) error {
	req, ok := request.(*notify.PushNotification)
	if !ok {
		return errors.New("invalid request")
	}

	// 设置最大重试次数
	maxRetry := 5
	if req.Retry > 0 && req.Retry < maxRetry {
		maxRetry = req.Retry
	}

	// 重试计数
	retryCount := 0

	for {
		newTokens, err := a.sendNotifications(req)
		if err != nil {
			return err
		}
		// 如果有重试的 Token，并且未达到最大重试次数，则进行重试
		if len(newTokens) > 0 && retryCount < maxRetry {
			retryCount++
			req.Tokens = newTokens
		} else {
			break
		}
	}

	return nil
}

func (a *APNsService) sendNotifications(req *notify.PushNotification) ([]string, error) {
	var newTokens []string
	notification := notify.GetIOSNotification(req)
	var wg sync.WaitGroup

	for _, token := range req.Tokens {
		// occupy push slot
		MaxConcurrentIOSPushes <- struct{}{}
		wg.Add(1)
		go func(notification apns2.Notification, token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentIOSPushes
				wg.Done()
			}()
			notification.DeviceToken = token
			res, err := a.client.Push(&notification)
			if err != nil || (res != nil && res.StatusCode != http.StatusOK) {
				if err == nil {
					err = errors.New(res.Reason)
				}
				// 记录失败的 Token
				if res != nil && res.StatusCode >= http.StatusInternalServerError {
					newTokens = append(newTokens, token)
				}
				log.Printf("apns send error: %s", err)
			} else {
				log.Printf("apns send success: %s", res.Reason)
			}
		}(*notification, token)
	}
	wg.Wait()
	return newTokens, nil
}

func (a *APNsService) MulticastSend(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Subscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Unsubscribe(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) SendToTopic(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) SendToCondition(ctx context.Context, req interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) CheckDevice(ctx context.Context, req interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (a *APNsService) Name() string {
	//TODO implement me
	panic("implement me")
}
