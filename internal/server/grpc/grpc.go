package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cossim/hipush/api/pb/v1"
	push2 "github.com/cossim/hipush/api/push"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/factory"
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/status"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"net"
)

type Handler struct {
	cfg     *config.Config
	logger  logr.Logger
	factory *factory.PushServiceFactory
	v1.UnimplementedPushServiceServer
}

func NewHandler(cfg *config.Config, logger logr.Logger, factory *factory.PushServiceFactory) *Handler {
	return &Handler{
		cfg:     cfg,
		logger:  logger.WithValues("server", "pb"),
		factory: factory,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lisAddr := fmt.Sprintf("%s", h.cfg.GRPC.Addr())
	lis, err := net.Listen("tcp", lisAddr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterPushServiceServer(server, h)

	serverShutdown := make(chan struct{})
	go func() {
		<-ctx.Done()
		h.logger.Info("Shutting down grpcServer", "addr", lisAddr)
		server.GracefulStop()
		close(serverShutdown)
	}()

	h.logger.Info("Starting  grpcServer", "addr", lisAddr)
	if err := server.Serve(lis); err != nil {
		if !errors.Is(err, grpc.ErrServerStopped) {
			h.logger.Error(err, "failed to start grpcServer")
			return err
		}
	}

	<-serverShutdown
	return nil
}

func (h *Handler) Push(ctx context.Context, req *v1.PushRequest) (*v1.PushResponse, error) {
	resp := &v1.PushResponse{}
	h.logger.Info("Received push request", "platform", req.Platform, "tokens", req.Token, "req", req)

	service, err := h.factory.GetPushService(req.Platform)
	if err != nil {
		h.logger.Error(err, "failed to create push service")
		return resp, err
	}

	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		h.logger.Error(err, "Failed to marshal data")
		return resp, err
	}

	fmt.Println("dataBytes => ", dataBytes)

	meta := &v1.Meta{
		AppID:   req.AppID,
		AppName: req.AppName,
		Token:   req.Token,
	}

	var r push2.SendRequest
	switch consts.Platform(req.Platform) {
	case consts.PlatformIOS:
		r = &v1.APNsPushRequest{Meta: meta}
	case consts.PlatformAndroid:
		r = &v1.AndroidPushRequestData{Meta: meta}
	case consts.PlatformHuawei:
		r = &v1.HuaweiPushRequestData{Meta: meta}
	case consts.PlatformXiaomi:
		r = &v1.XiaomiPushRequestData{Meta: meta}
	case consts.PlatformVivo:
		r = &v1.VivoPushRequestData{Meta: meta}
	case consts.PlatformOppo:
		r = &v1.OppoPushRequestData{Meta: meta}
	case consts.PlatformMeizu:
		r = &v1.MeizuPushRequestData{Meta: meta}
	case consts.PlatformHonor:
		r = &v1.HonorPushRequestData{Meta: meta}
	default:
		return nil, errors.New("platform not supported")
	}

	marshalJSON, err := req.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(marshalJSON, &r); err != nil {
		h.logger.Error(err, "Failed to unmarshal data")
		return resp, err
	}

	if err := h.validatePushRequest(r); err != nil {
		return nil, err
	}

	option := v1.PushOption{}
	if req.Option != nil {
		option.Development = req.Option.Development
		option.DryRun = req.Option.DryRun
		option.Retry = req.Option.Retry
		option.RetryInterval = req.Option.RetryInterval
	}

	fmt.Println("r => ", r)

	status.StatStorage.AddGrpcTotal(1)
	_, err = service.Send(ctx, r, &push2.SendOptions{
		DryRun:        option.DryRun,
		Retry:         option.Retry,
		RetryInterval: option.RetryInterval,
	})
	if err != nil {
		status.StatStorage.AddGrpcFailed(1)
		h.logger.Error(err, "failed to send push")
		return resp, err
	}
	status.StatStorage.AddGrpcSuccess(1)

	h.logger.Info("Push request processed success")
	return resp, nil
}

func (h *Handler) validatePushRequest(req push2.SendRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	//if !consts.Platform(req.Platform).IsValid() {
	//	return errors.New("invalid platform")
	//
	//}

	//if len(req.Token) == 0 {
	//	return errors.New("tokens are required")
	//}

	// 检查其他必填字段
	if req.GetTitle() == "" {
		return errors.New("title is required")
	}

	if req.GetContent() == "" {
		return errors.New("message is required")
	}

	// 检查 Data 字段
	//if req.Data == nil {
	//return errors.New("data is required")
	//}

	return nil
}
