package http

import (
	"encoding/json"
	"fmt"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/internal/notify"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) handleIOSPush(c *gin.Context, req *dto.PushRequest) {
	fmt.Println("consts.Platform(req.Platform).String() => ", req.Platform)
	service, err := h.factory.GetPushService(req.Platform)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
		return
	}
	var r dto.APNsPushRequest
	if err := json.Unmarshal(dataBytes, &r); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid data", Data: nil})
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
		c.JSON(http.StatusInternalServerError, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "push notification send success", Data: nil})
}
