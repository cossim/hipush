package http

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/cossim/hipush/api/pb/v1"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/factory"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"net/http"
)

type Handler struct {
	cfg     *config.Config
	logger  logr.Logger
	factory *factory.PushServiceFactory
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewHandler(cfg *config.Config, logger logr.Logger, factory *factory.PushServiceFactory) *Handler {
	return &Handler{
		cfg:     cfg,
		logger:  logger.WithValues("server", "http"),
		factory: factory,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := gin.Default()
	r.POST("/api/v1/push", h.pushHandler)
	r.GET("/api/v1/push/stat", h.pushStatHandler)
	r.GET("/api/v1/message/stat", h.pushMessageStatHandler)

	srv := &http.Server{
		Addr:    h.cfg.HTTP.Addr(),
		Handler: r,
	}

	serverShutdown := make(chan struct{})
	go func() {
		<-ctx.Done()
		h.logger.Info("shutting down httpServer", "addr", h.cfg.HTTP.Addr())
		if err := srv.Shutdown(ctx); err != nil {
			h.logger.Error(err, "error shutting down httpServer")
		}
		close(serverShutdown)
	}()

	h.logger.Info("starting httpServer", "addr", h.cfg.HTTP.Addr())
	if err := srv.ListenAndServe(); err != nil {
		// Check if the error is not due to the server being closed intentionally
		if !errors.Is(err, http.ErrServerClosed) {
			// Log the error and return an error message
			h.logger.Error(err, fmt.Sprintf("Failed to start HTTP server: %v", err))
			return fmt.Errorf("failed to start HTTP server: %v", err)
		}
		// If the error is due to the server being closed intentionally, return nil
		return nil
	}

	<-serverShutdown
	return nil
}

func (h *Handler) pushHandler(c *gin.Context) {
	req := &v1.PushRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error(err, "failed to bind request")
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: err.Error(), Data: nil})
		return
	}

	status.StatStorage.AddHttpTotal(1)
	h.logger.Info("Received push request", "platform", req.Platform, "appid", req.AppID, "tokens", req.Token, "data", req.Data)

	var err error
	switch consts.Platform(req.Platform) {
	case consts.PlatformIOS:
		err = h.handleIOSPush(c, req)
	case consts.PlatformAndroid:
		err = h.handleAndroidPush(c, req)
	case consts.PlatformHuawei:
		err = h.handleHuaweiPush(c, req)
	case consts.PlatformVivo:
		err = h.handleVivoPush(c, req)
	case consts.PlatformOppo:
		err = h.handleOppoPush(c, req)
	case consts.PlatformXiaomi:
		err = h.handleXiaomiPush(c, req)
	case consts.PlatformMeizu:
		err = h.handleMeizuPush(c, req)
	case consts.PlatformHonor:
		err = h.handleHonorPush(c, req)

	default:
		c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Msg: "invalid platform", Data: nil})
		return
	}

	if err != nil {
		status.StatStorage.AddHttpFailed(1)
	} else {
		status.StatStorage.AddHttpSuccess(1)
	}
}
