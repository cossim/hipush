package http

import (
	"encoding/json"
	"fmt"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/internal/consts"
	"github.com/cossim/hipush/internal/notify"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleIOSPush(c *gin.Context, req *dto.PushRequest) {
	fmt.Println("consts.Platform(req.Platform).String() => ", consts.Platform(req.Platform).String())
	service, err := h.factory.GetPushService(consts.Platform(req.Platform).String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将 req.Data 转换为 JSON 字节
	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to marshal data"})
		return
	}

	// 将 JSON 字节解码到 APNsPushRequest 结构体中
	var r dto.APNsPushRequest
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Push notification sent successfully"})
}
