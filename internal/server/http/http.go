package http

import (
	"context"
	"errors"
	"github.com/cossim/hipush/api/v1/http/dto"
	"github.com/cossim/hipush/internal/factory"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	factory factory.PushServiceFactory
}

func (h *Handler) Start(ctx context.Context) error {
	r := gin.Default()
	r.GET("/api/v1/push", h.pushHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return nil
}

func (h *Handler) pushHandler(c *gin.Context) {
	req := &dto.PushRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Printf("error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.Platform {
	case dto.PlatformIOS:
		h.handleIOSPush(c, req)
	case dto.PlatformHuawei:
		h.handleHuaweiPush(c, req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid platform"})
	}
}
