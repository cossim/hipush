package http

import (
	"github.com/cossim/hipush/api/http/v1/dto"
	api "github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) pushStatHandler(c *gin.Context) {
	//req := &dto.PushStatRequest{}
	//if err := c.ShouldBindJSON(req); err != nil {
	//	h.logger.Error(err, "failed to bind request")
	//	c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
	//	return
	//}
	//h.logger.Info("Received pushStat request", "platform", req.Platform, "method", req.Method)

	ps := &dto.PushStats{}
	ps.Total = status.StatStorage.GetTotalCount()
	ps.Success = status.StatStorage.GetSuccessCount()
	ps.Failed = status.StatStorage.GetFailedCount()
	ps.Send = status.StatStorage.GetSendCount()
	ps.Receive = status.StatStorage.GetReceiveCount()
	ps.Display = status.StatStorage.GetDisplayCount()
	ps.Click = status.StatStorage.GetClickCount()

	ps.HTTP.Total = status.StatStorage.GetHttpTotal()
	ps.HTTP.Success = status.StatStorage.GetHttpSuccess()
	ps.HTTP.Failed = status.StatStorage.GetHttpFailed()

	ps.GRPC.Total = status.StatStorage.GetGrpcTotal()
	ps.GRPC.Success = status.StatStorage.GetGrpcSuccess()
	ps.GRPC.Failed = status.StatStorage.GetGrpcFailed()

	ps.IOS.Total = status.StatStorage.GetIosTotal()
	ps.IOS.Success = status.StatStorage.GetIosSuccess()
	ps.IOS.Failed = status.StatStorage.GetIosFailed()
	ps.IOS.Send = status.StatStorage.GetIosSend()
	ps.IOS.Receive = status.StatStorage.GetIosReceive()
	ps.IOS.Display = status.StatStorage.GetIosDisplay()
	ps.IOS.Click = status.StatStorage.GetIosClick()

	ps.Android.Total = status.StatStorage.GetAndroidTotal()
	ps.Android.Success = status.StatStorage.GetAndroidSuccess()
	ps.Android.Failed = status.StatStorage.GetAndroidFailed()
	ps.Android.Send = status.StatStorage.GetAndroidSend()
	ps.Android.Receive = status.StatStorage.GetAndroidReceive()
	ps.Android.Display = status.StatStorage.GetAndroidDisplay()
	ps.Android.Click = status.StatStorage.GetAndroidClick()

	ps.Huawei.Total = status.StatStorage.GetHuaweiTotal()
	ps.Huawei.Success = status.StatStorage.GetHuaweiSuccess()
	ps.Huawei.Failed = status.StatStorage.GetHuaweiFailed()
	ps.Huawei.Send = status.StatStorage.GetHuaweiSend()
	ps.Huawei.Receive = status.StatStorage.GetHuaweiReceive()
	ps.Huawei.Display = status.StatStorage.GetHuaweiDisplay()
	ps.Huawei.Click = status.StatStorage.GetHuaweiClick()

	ps.Xiaomi.Total = status.StatStorage.GetXiaomiTotal()
	ps.Xiaomi.Success = status.StatStorage.GetXiaomiSuccess()
	ps.Xiaomi.Failed = status.StatStorage.GetXiaomiFailed()
	ps.Xiaomi.Send = status.StatStorage.GetXiaomiSend()
	ps.Xiaomi.Receive = status.StatStorage.GetXiaomiReceive()
	ps.Xiaomi.Display = status.StatStorage.GetXiaomiDisplay()
	ps.Xiaomi.Click = status.StatStorage.GetXiaomiClick()

	ps.Vivo.Total = status.StatStorage.GetVivoTotal()
	ps.Vivo.Success = status.StatStorage.GetVivoSuccess()
	ps.Vivo.Failed = status.StatStorage.GetVivoFailed()
	ps.Vivo.Send = status.StatStorage.GetVivoSend()
	ps.Vivo.Receive = status.StatStorage.GetVivoReceive()
	ps.Vivo.Display = status.StatStorage.GetVivoDisplay()
	ps.Vivo.Click = status.StatStorage.GetVivoClick()

	ps.Oppo.Total = status.StatStorage.GetOppoTotal()
	ps.Oppo.Success = status.StatStorage.GetOppoSuccess()
	ps.Oppo.Failed = status.StatStorage.GetOppoFailed()
	ps.Oppo.Send = status.StatStorage.GetOppoSend()
	ps.Oppo.Receive = status.StatStorage.GetOppoReceive()
	ps.Oppo.Display = status.StatStorage.GetOppoDisplay()
	ps.Oppo.Click = status.StatStorage.GetOppoClick()

	ps.Meizu.Total = status.StatStorage.GetMeizuTotal()
	ps.Meizu.Success = status.StatStorage.GetMeizuSuccess()
	ps.Meizu.Failed = status.StatStorage.GetMeizuFailed()
	ps.Meizu.Send = status.StatStorage.GetMeizuSend()
	ps.Meizu.Receive = status.StatStorage.GetMeizuReceive()
	ps.Meizu.Display = status.StatStorage.GetMeizuDisplay()
	ps.Meizu.Click = status.StatStorage.GetMeizuClick()

	ps.Honor.Total = status.StatStorage.GetHonorTotal()
	ps.Honor.Success = status.StatStorage.GetHonorSuccess()
	ps.Honor.Failed = status.StatStorage.GetHonorFailed()
	ps.Honor.Send = status.StatStorage.GetHonorSend()
	ps.Honor.Receive = status.StatStorage.GetHonorReceive()
	ps.Honor.Display = status.StatStorage.GetHonorDisplay()
	ps.Honor.Click = status.StatStorage.GetHonorClick()

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Get push stat success", Data: ps})
}

func (h *Handler) pushMessageStatHandler(c *gin.Context) {
	req := &dto.PushMessageStatRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error(err, "failed to bind request")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	h.logger.Info("Received pushMessageStat request", "appid", req.AppID)

	service, err := h.factory.GetPushService(consts.Platform(req.Platform).String())
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	vps := &api.PushMessageStatsList{}
	if err := service.GetTasksStatus(c, req.AppID, req.TaskID, vps); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Get push message stat success", Data: vps.Get()})
}
