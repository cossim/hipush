package push

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidAppID = errors.New("invalid appid or appid push is not enabled")
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	Success = 200
	Fail    = 400
)

type SendFunc func(ctx context.Context, token string) (*Response, error)

func RetrySend(ctx context.Context, send SendFunc, tokens []string, retry int32, retryInterval int32, maxConcurrent int) (*Response, error) {
	var wg sync.WaitGroup
	var resp = &Response{}
	if retryInterval <= 0 {
		retryInterval = 1
	}
	if maxConcurrent <= 0 {
		maxConcurrent = 100
	}
	var MaxConcurrentPushes = make(chan struct{}, maxConcurrent)
	var es []error

	for _, token := range tokens {
		// occupy push slot
		MaxConcurrentPushes <- struct{}{}
		wg.Add(1)
		go func(token string) {
			defer func() {
				// free push slot
				<-MaxConcurrentPushes
				wg.Done()
			}()
			for i := 0; i <= int(retry); i++ {
				res, err := send(ctx, token)
				if err != nil || (res != nil && res.Code != 200) {
					if err == nil {
						err = errors.New(res.Msg)
					} else {
						es = append(es, err)
					}
					if i == 0 {
						continue
					}
					log.Printf("send error: %s (attempt %d)", err, i)
					time.Sleep(time.Duration(retryInterval) * time.Second)
				} else {
					log.Printf("send success: %s", res.Msg)
					resp.Data = res.Data
					break
				}
			}
		}(token)
	}
	wg.Wait()
	if len(es) > 0 {
		uniqueErrors := make(map[string]struct{})
		for _, err := range es {
			uniqueErrors[err.Error()] = struct{}{}
		}
		var uniqueErrorStrings []string
		for err := range uniqueErrors {
			uniqueErrorStrings = append(uniqueErrorStrings, err)
		}
		allErrorsString := strings.Join(uniqueErrorStrings, ", ")
		return nil, errors.New(allErrorsString)
	}

	return resp, nil
}
