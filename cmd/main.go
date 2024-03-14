package main

import (
	"context"
	"flag"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/factory"
	g "github.com/cossim/hipush/internal/server/grpc"
	h "github.com/cossim/hipush/internal/server/http"
	"github.com/cossim/hipush/pkg/push"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
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

	if err := status.InitAppStatus(cfg); err != nil {
		panic(err)
	}

	zapLogger := zap.NewExample()
	logger := zapr.NewLogger(zapLogger)
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	pushServiceFactory := factory.NewPushServiceFactory()
	if err := pushServiceFactory.Register(
		//pushServiceFactory.WithPushService(push.NewAPNsService(cfg, logger)),
		//pushServiceFactory.WithPushService(push.NewFCMService(cfg, logger)),
		//pushServiceFactory.WithPushService(push.NewHMSService(cfg, logger)),
		//pushServiceFactory.WithPushService(push.NewXiaomiService(cfg, logger)),
		//pushServiceFactory.WithPushService(push.NewOppoService(cfg, logger)),
		pushServiceFactory.WithPushService(push.NewVivoService(cfg, logger, scheduler)),
		//pushServiceFactory.WithPushService(push.NewMeizuService(cfg, logger)),
		//pushServiceFactory.WithPushService(push.NewHonorService(cfg, logger)),
	); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	if cfg.HTTP.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			httpHandler := h.NewHandler(cfg, logger, pushServiceFactory)
			if err := httpHandler.Start(ctx); err != nil {
				log.Fatalf("failed to start HTTP server: %v", err)
			}
		}()
	}

	if cfg.GRPC.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			grpcHandler := g.NewHandler(cfg, logger, pushServiceFactory)
			if err := grpcHandler.Start(ctx); err != nil {
				log.Fatalf("failed to start GRPC server: %v", err)
			}
		}()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case <-sig:
			log.Println("receive system signal, cancel context")
			status.StatStorage.Close()
			cancel()
		}
	}()
	wg.Wait()
}
