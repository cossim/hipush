package http

import (
	"github.com/cossim/hipush/api/v1/http/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleIOSPush(c *gin.Context, req *dto.PushRequest) {
	//service, err := h.factory.CreatePushService(req.Platform.String())
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//r, ok := req.Data.(dto.HuaweiPushRequestData)
	//if !ok {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
	//	return
	//}
	//
	//if err := service.Send(c, rr); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//c.JSON(http.StatusOK, gin.H{"message": "Push notification sent successfully"})
}
