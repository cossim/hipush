package http

import (
	"encoding/json"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/notify"
	"github.com/cossim/hipush/pkg/push"
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
		Data:        r.Data,
		Development: req.Option.Development,
		Badge: &notify.BadgeNotification{
			AddNum:     r.Badge.AddNum,
			SetNum:     r.Badge.SetNum,
			BadgeClass: r.Badge.Class,
		},
		ClickAction: &notify.HonorClickAction{
			Action:     r.ClickAction.Action,
			Activity:   r.ClickAction.Activity,
			Url:        r.ClickAction.Url,
			Parameters: r.ClickAction.Parameters,
		},
		NotifyId: r.NotifyId,
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
