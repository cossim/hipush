package http

import (
	"encoding/json"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/pkg/notify"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleAndroidPush(c *gin.Context, req *dto.PushRequest) error {
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
	var r dto.AndroidPushRequestData
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return err
	}

	h.logger.Info("Handling push request", "platform", req.Platform, "appID", req.AppID, "tokens", req.Token, "req", r)

	rr := &notify.FCMPushNotification{
		AppID:            req.AppID,
		AppName:          req.AppName,
		Tokens:           req.Token,
		Topic:            "",
		Priority:         "",
		Title:            r.Title,
		Message:          r.Content,
		Image:            "",
		Sound:            "",
		CollapseID:       "",
		Category:         "",
		Condition:        "",
		TTL:              nil,
		Retry:            0,
		Badge:            nil,
		ContentAvailable: false,
		MutableContent:   false,
		Data:             nil,
		Apns:             nil,
	}
	resp, err := service.Send(c, rr, &push.SendOptions{
		DryRun:        req.Option.DryRun,
		Retry:         req.Option.Retry,
		RetryInterval: req.Option.RetryInterval,
	})
	if err != nil {
		h.logger.Error(err, "Failed to send push notification")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return err
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Push notification send success", Data: resp})
	return nil
}
