package http

import (
	"encoding/json"
	"github.com/cossim/go-hms-push/push/model"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/notify"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleHuaweiPush(c *gin.Context, req *dto.PushRequest) error {
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
	var r dto.HuaweiPushRequestData
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return err
	}

	h.logger.Info("Handling push request", "platform", req.Platform, "appID", req.AppID, "tokens", req.Token, "req", r)

	rr := &notify.HMSPushNotification{
		AppID:       req.AppID,
		AppName:     req.AppName,
		Tokens:      req.Token,
		Development: r.Development,
		MessageRequest: &model.MessageRequest{
			Message: &model.Message{
				Notification: &model.Notification{
					Title: r.Title,
					Body:  r.Message,
				},
				Android: &model.AndroidConfig{
					Notification: &model.AndroidNotification{
						Title: r.Title,
						Body:  r.Message,
						Badge: &model.BadgeNotification{
							AddNum: r.Badge.AddNum,
							SetNum: r.Badge.SetNum,
							Class:  r.Badge.Class,
						},
					},
				},
				Token: req.Token,
			},
		},
	}
	resp, err := service.Send(c, rr, &push.SendOptions{
		DryRun:        req.Option.DryRun,
		Retry:         req.Option.Retry,
		RetryInterval: req.Option.RetryInterval,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return err
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Push notification send success", Data: resp})
	return nil
}
