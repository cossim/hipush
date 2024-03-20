package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/pkg/notify"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleIOSPush(c *gin.Context, req *dto.PushRequest) error {
	service, err := h.factory.GetPushService(req.Platform)
	if err != nil {
		h.logger.Error(err, "Failed to get push service")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return err
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		h.logger.Error(err, "Failed to marshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return err
	}
	var r dto.APNsPushRequest
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return err
	}

	h.logger.Info("Handling push request", "platform", req.Platform, "appID", req.AppID, "tokens", req.Token, "req", r)

	var topic = r.Topic

	if req.AppID == "" && r.Topic == "" {
		msg := errors.New("one of AppID and Topic cannot be empty")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: msg.Error(), Data: nil})
		return msg
	}

	if r.Topic == "" {
		topic = req.AppID
	}

	fmt.Println("topic => ", topic)

	rr := &notify.ApnsPushNotification{
		AppID:            req.AppID,
		AppName:          req.AppName,
		Tokens:           req.Token,
		Title:            r.Title,
		Topic:            topic,
		Content:          r.Content,
		ApnsID:           r.ApnsID,
		Sound:            r.Sound,
		Development:      r.Development,
		MutableContent:   r.MutableContent,
		CollapseID:       r.CollapseID,
		ContentAvailable: r.ContentAvailable,
		Priority:         r.Priority,
		Data:             r.Data,
		Badge:            &r.Badge,
		Expiration:       &r.TTL,
	}
	resp, err := service.Send(c, rr, &push.SendOptions{
		DryRun:        req.Option.DryRun,
		Retry:         req.Option.Retry,
		RetryInterval: req.Option.RetryInterval,
	})
	if err != nil {
		h.logger.Error(err, "Failed to send push notification")
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return err
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Push notification send success", Data: resp})
	return nil
}
