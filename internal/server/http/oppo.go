package http

import (
	"encoding/json"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/internal/consts"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/push"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleOppoPush(c *gin.Context, req *dto.PushRequest) {
	service, err := h.factory.GetPushService(consts.Platform(req.Platform).String())
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		h.logger.Error(err, "Failed to marshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return
	}
	var r dto.OppoPushRequestData
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return
	}

	rr := &notify.OppoPushNotification{
		AppID:       req.AppID,
		Tokens:      req.Token,
		Title:       r.Title,
		Subtitle:    r.Subtitle,
		Message:     r.Message,
		Data:        r.Data,
		ClickAction: &r.ClickAction,
		TTL:         0,
		Option: notify.PushOption{
			DryRun: req.Option.DryRun,
			Retry:  req.Option.Retry,
		},
	}
	if err := service.Send(c, rr, &push.SendOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Push notification send success", Data: nil})
}
