package http

import (
	"github.com/cossim/hipush/api/http/v1/dto"
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
	//
	//h.logger.Info("Received pushStat request", "platform", req.Platform, "method", req.Method)

	ps := &dto.PushStats{}
	ps.Total = status.StatStorage.GetTotalCount()
	ps.Success = status.StatStorage.GetSuccessCount()
	ps.Failed = status.StatStorage.GetFailedCount()

	ps.HTTP.Total = status.StatStorage.GetHttpTotal()
	ps.HTTP.Success = status.StatStorage.GetHttpSuccess()
	ps.HTTP.Failed = status.StatStorage.GetHttpFailed()

	ps.GRPC.Total = status.StatStorage.GetGrpcTotal()
	ps.GRPC.Success = status.StatStorage.GetGrpcSuccess()
	ps.GRPC.Failed = status.StatStorage.GetGrpcFailed()

	ps.IOS.Total = status.StatStorage.GetIosTotal()
	ps.IOS.Success = status.StatStorage.GetIosSuccess()
	ps.IOS.Failed = status.StatStorage.GetIosFailed()

	ps.Android.Total = status.StatStorage.GetAndroidTotal()
	ps.Android.Success = status.StatStorage.GetAndroidSuccess()
	ps.Android.Failed = status.StatStorage.GetAndroidFailed()

	ps.Huawei.Total = status.StatStorage.GetHuaweiTotal()
	ps.Huawei.Success = status.StatStorage.GetHuaweiSuccess()
	ps.Huawei.Failed = status.StatStorage.GetHuaweiFailed()

	ps.Xiaomi.Total = status.StatStorage.GetXiaomiTotal()
	ps.Xiaomi.Success = status.StatStorage.GetXiaomiSuccess()
	ps.Xiaomi.Failed = status.StatStorage.GetXiaomiFailed()

	ps.Vivo.Total = status.StatStorage.GetVivoTotal()
	ps.Vivo.Success = status.StatStorage.GetVivoSuccess()
	ps.Vivo.Failed = status.StatStorage.GetVivoFailed()

	ps.Oppo.Total = status.StatStorage.GetOppoTotal()
	ps.Oppo.Success = status.StatStorage.GetOppoSuccess()
	ps.Oppo.Failed = status.StatStorage.GetOppoFailed()

	ps.Meizu.Total = status.StatStorage.GetMeizuTotal()
	ps.Meizu.Success = status.StatStorage.GetMeizuSuccess()
	ps.Meizu.Failed = status.StatStorage.GetMeizuFailed()

	ps.Honor.Total = status.StatStorage.GetHonorTotal()
	ps.Honor.Success = status.StatStorage.GetHonorSuccess()
	ps.Honor.Failed = status.StatStorage.GetHonorFailed()

	c.JSON(http.StatusOK, Response{Code: http.StatusOK, Msg: "Get push stat success", Data: ps})
}
