package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
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

}
