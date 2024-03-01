package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cossim/hipush/api/grpc/v1"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/internal/consts"
	"github.com/cossim/hipush/internal/factory"
	"github.com/cossim/hipush/internal/notify"
	"github.com/cossim/hipush/internal/push"
	"github.com/go-logr/logr"
	"github.com/golang/protobuf/jsonpb"
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
		logger:  logger.WithValues("server", "grpc"),
		factory: factory,
	}
}

func (h *Handler) Start(ctx context.Context) error {
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
		//h.cancel()
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
	h.logger.Info("Received push request", "platform", req.Platform, "tokens", req.Tokens, "req", req)

	service, err := h.factory.GetPushService(req.Platform)
	if err != nil {
		h.logger.Error(err, "failed to create push service")
		return resp, err
	}

	r, err := h.getPushRequest(req)
	if err != nil {
		h.logger.Error(err, "failed to get push request")
		return nil, err
	}

	if err := service.Send(ctx, r); err != nil {
		h.logger.Error(err, "failed to send push")
		return resp, err
	}

	h.logger.Info("Push request processed success")
	return resp, nil
}

func (h *Handler) getPushRequest(req *v1.PushRequest) (push.PushRequest, error) {
	badge := int(req.Badge)

	data := make(map[string]interface{})
	if req.Data != nil {
		jsonStr, err := (&jsonpb.Marshaler{}).MarshalToString(req.Data)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal([]byte(jsonStr), &data); err != nil {
			return nil, err
		}
	}

	alert := notify.Alert{}
	if req.Alert != nil {
		alert = notify.Alert{
			Action:       req.Alert.Action,
			ActionLocKey: req.Alert.ActionLocKey,
			Body:         req.Alert.Body,
			LaunchImage:  req.Alert.LaunchImage,
			LocArgs:      req.Alert.LocArgs,
			LocKey:       req.Alert.LocKey,
			Title:        req.Alert.Title,
			Subtitle:     req.Alert.Subtitle,
			TitleLocArgs: req.Alert.TitleLocArgs,
			TitleLocKey:  req.Alert.TitleLocKey,
		}
	}

	return &notify.ApnsPushNotification{
		ApnsID:           req.AppID,
		Tokens:           req.Tokens,
		Title:            req.Title,
		Message:          req.Message,
		Topic:            req.Topic,
		Category:         req.Category,
		Sound:            req.Sound,
		Alert:            alert,
		Badge:            &badge,
		ThreadID:         req.ThreadID,
		Data:             data,
		PushType:         req.PushType,
		Priority:         string(req.Priority),
		ContentAvailable: req.ContentAvailable,
		MutableContent:   req.MutableContent,
		Development:      req.Development,
	}, nil
}

func (h *Handler) validatePushRequest(req *v1.PushRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	if !consts.Platform(req.Platform).IsValid() {
		return errors.New("invalid platform")

	}

	if len(req.Tokens) == 0 {
		return errors.New("tokens are required")
	}

	// 检查其他必填字段
	if req.Title == "" {
		return errors.New("title is required")
	}

	if req.Message == "" {
		return errors.New("message is required")
	}

	// 检查 Alert 字段
	if req.Alert != nil {
		if err := h.validateAlert(req.Alert); err != nil {
			return err
		}
	}

	// 检查 Data 字段
	if req.Data == nil {
		//return errors.New("data is required")
	}

	return nil
}

func (h *Handler) validateAlert(alert *v1.Alert) error {
	// TODO 检查 Alert 字段的必填参数
	return nil
}

func (h *Handler) mustEmbedUnimplementedPushServiceServer() {
	//TODO implement me
	panic("implement me")
}
