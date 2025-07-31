package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/sms/v1"
	"webook/sms/service"
)

type SmsServiceServer struct {
	smsv1.UnimplementedSmsServiceServer
	svc service.SmsService
}

func NewSmsServiceServer(svc service.SmsService) *SmsServiceServer {
	return &SmsServiceServer{svc: svc}
}

func (s *SmsServiceServer) Register(server *grpc.Server) {
	smsv1.RegisterSmsServiceServer(server, s)
}

func (s *SmsServiceServer) Send(ctx context.Context, request *smsv1.SmsSendRequest) (*smsv1.SmsSendResponse, error) {
	err := s.svc.Send(ctx, request.GetTplId(), request.GetArgs(), request.GetNumbers()...)
	return &smsv1.SmsSendResponse{}, err
}
