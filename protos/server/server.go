package server

import (
	"context"
	"time"

	protos "github.com/ssubedir/goriffin/protos/service/protos"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (c *Service) HeartBeatStatus(ctx context.Context, req *protos.StatusRequest) (*protos.StatusResponse, error) {
	return &protos.StatusResponse{Status: true, Time: time.Now().Local().String()}, nil
}
