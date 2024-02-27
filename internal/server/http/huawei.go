package http

import (
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/internal/consts"
	"github.com/cossim/hipush/internal/notify"
	"github.com/gin-gonic/gin"
	"github.com/msalihkarakasli/go-hms-push/push/model"
	"net/http"
)

func (h *Handler) handleHuaweiPush(c *gin.Context, req *dto.PushRequest) {
	service, err := h.factory.GetPushService(consts.Platform(req.Platform).String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, ok := req.Data.(dto.HuaweiPushRequestData)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	_ = &notify.HMSPushNotification{
		AppID:     req.AppID,
		AppSecret: req.AppSecret,
		Tokens:    req.Token,
		MessageRequest: &model.MessageRequest{
			Message: &model.Message{
				Notification: &model.Notification{
					Title: r.Title,
					Body:  r.Body,
				},
				Android: nil,
				Token:   req.Token,
			},
		},
	}

	if err := service.Send(c, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Push notification sent successfully"})
}
