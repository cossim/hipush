package http

import (
	"encoding/json"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/consts"
	"github.com/cossim/hipush/notify"
	"github.com/cossim/hipush/push"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleXiaomiPush(c *gin.Context, req *dto.PushRequest) error {
	service, err := h.factory.GetPushService(consts.Platform(req.Platform).String())
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return err
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		h.logger.Error(err, "Failed to marshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return err
	}
	var r dto.XiaomiPushRequestData
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return err
	}

	h.logger.Info("Handling push request", "platform", req.Platform, "appID", req.AppID, "tokens", req.Token, "req", r)

	rr := &notify.XiaomiPushNotification{
		AppID:         req.AppID,
		Tokens:        req.Token,
		Title:         r.Title,
		Content:       r.Content,
		NotifyType:    r.NotifyType,
		TTL:           int64(r.TTL),
		IsShowNotify:  r.Foreground,
		IsScheduled:   r.IsScheduled,
		ScheduledTime: r.ScheduledTime,
	}
	if err := service.Send(c, rr, &push.SendOptions{
		DryRun: false,
		Retry:  0,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return err
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Push notification send success", Data: nil})
	return nil
}
