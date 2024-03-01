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

	if !cfg.HTTP.Enabled && !cfg.GRPC.Enabled {
		log.Fatalf("Neither HTTP nor GRPC server is enabled. Please enable at least one server.")
	}

	zapLogger := zap.NewExample()

	pushServiceFactory := factory.NewPushServiceFactory()
	pushServiceFactory.Register(consts.PlatformIOS.String(), func() push.PushService {
		svc, err := push.NewAPNsService(cfg)
		if err != nil {
			panic(err)
		}
		return adapter.NewPushServiceAdapter(svc)
	})

	pushServiceFactory.Register(consts.PlatformAndroid.String(), func() push.PushService {
		svc, err := push.NewFCMService(cfg)
		if err != nil {
			panic(err)
		}
		return adapter.NewPushServiceAdapter(svc)
	})

	pushServiceFactory.Register(consts.PlatformHuawei.String(), func() push.PushService {
		svc, err := push.NewHMSService(cfg)
		if err != nil {
			panic(err)
		}
		return adapter.NewPushServiceAdapter(svc)
	})

	if cfg.HTTP.Enabled {
		go func() {
			httpHandler := h.NewHandler(cfg, zapr.NewLogger(zapLogger), pushServiceFactory)
			if err := httpHandler.Start(context.Background()); err != nil {
				log.Fatalf("failed to start HTTP server: %v", err)
			}
		}()
	}

	if cfg.GRPC.Enabled {
		go func() {
			grpcHandler := g.NewHandler(cfg, zapr.NewLogger(zapLogger), pushServiceFactory)
			if err := grpcHandler.Start(context.Background()); err != nil {
				log.Fatalf("failed to start GRPC server: %v", err)
			}
		}()
	}

	// 等待信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	log.Println("shutting down...")
}
