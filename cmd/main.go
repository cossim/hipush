package main

import (
	"context"
	"flag"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/adapter"
	"github.com/cossim/hipush/internal/consts"
	"github.com/cossim/hipush/internal/factory"
	"github.com/cossim/hipush/internal/push"
	g "github.com/cossim/hipush/internal/server/grpc"
	h "github.com/cossim/hipush/internal/server/http"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "", "Configuration file path.")
	flag.Parse()
}

func main() {
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLogger := zap.NewExample()

	pushServiceFactory := factory.NewPushServiceFactory()

	pushServiceFactory.Register(consts.PlatformIOS.String(), func() push.PushService {
		return adapter.NewPushServiceAdapter(push.NewAPNsService(cfg))
	})

	httpHandler := h.NewHandler(pushServiceFactory)

	grpcHandler := g.NewHandler(pushServiceFactory, zapr.NewLogger(zapLogger))

	go func() {
		if err := httpHandler.Start(context.Background()); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	go func() {
		if err := grpcHandler.Start(context.Background()); err != nil {
			log.Fatalf("failed to start GRPC server: %v", err)
		}
	}()

	// 等待信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	log.Println("shutting down...")
}
