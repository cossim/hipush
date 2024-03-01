package http

import (
	"encoding/json"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/internal/notify"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleIOSPush(c *gin.Context, req *dto.PushRequest) {
	service, err := h.factory.GetPushService(req.Platform)
	if err != nil {
		h.logger.Error(err, "Failed to get push service")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		h.logger.Error(err, "Failed to marshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return
	}
	var r dto.APNsPushRequest
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return
	}

	h.logger.Info("Handling push request", "platform", req.Platform, "appID", req.AppID, "tokens", req.Token, "req", r)

	rr := &notify.ApnsPushNotification{
		//Retry:      5,
		Tokens:           req.Token,
		Title:            r.Title,
		Topic:            req.AppID,
		Message:          r.Message,
		ApnsID:           req.AppID,
		Sound:            r.Sound,
		Production:       r.Production,
		Development:      r.Development,
		MutableContent:   r.MutableContent,
		CollapseID:       r.CollapseID,
		ContentAvailable: r.ContentAvailable,
		Priority:         r.Priority,
		Data:             r.Data,
		Expiration:       nil,
	}
	if err := service.Send(c, rr); err != nil {
		h.logger.Error(err, "Failed to send push notification")
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Push notification send success", Data: nil})
}
