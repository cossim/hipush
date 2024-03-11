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

func (h *Handler) handleHonorPush(c *gin.Context, req *dto.PushRequest) {
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
	var r dto.HonorPushRequestData
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return
	}

	h.logger.Info("Handling push request", "platform", req.Platform, "appID", req.AppID, "tokens", req.Token, "req", r)

	rr := &notify.HonorPushNotification{
		AppID:    req.AppID,
		Tokens:   req.Token,
		Title:    r.Title,
		Content:  r.Content,
		Image:    r.Icon,
		Priority: "",
		Category: "",
		//TTL:         strconv.Itoa(r.TTL),
		Data:        "",
		Development: true,
		ClickAction: &notify.HonorClickAction{
			Action:   r.ClickAction.Action,
			Activity: r.ClickAction.Activity,
			Url:      r.ClickAction.Url,
			//Parameters: r.ClickAction.Parameters,
		},
	}
	if err := service.Send(c, rr, &push.SendOptions{
		DryRun:        req.Option.DryRun,
		Retry:         req.Option.Retry,
		RetryInterval: req.Option.RetryInterval,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Push notification send success", Data: nil})
}
