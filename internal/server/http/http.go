package http

import (
	"context"
	"github.com/cossim/hipush/api/http/v1/dto"
	"github.com/cossim/hipush/internal/consts"
	"github.com/cossim/hipush/internal/factory"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	factory *factory.PushServiceFactory
}

func NewHandler(factory *factory.PushServiceFactory) *Handler {
	return &Handler{factory: factory}
}

func (h *Handler) Start(ctx context.Context) error {
	r := gin.Default()
	r.POST("/api/v1/push", h.pushHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return srv.ListenAndServe()
}

func (h *Handler) pushHandler(c *gin.Context) {
	req := &dto.PushRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch consts.Platform(req.Platform) {
	case consts.PlatformIOS:
		h.handleIOSPush(c, req)
	case consts.PlatformHuawei:
		h.handleHuaweiPush(c, req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid platform"})
	}
}
