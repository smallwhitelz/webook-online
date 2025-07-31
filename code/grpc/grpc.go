package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/code/v1"
	"webook/code/service"
)

type CodeServiceServer struct {
	codev1.UnimplementedCodeServiceServer
	svc service.CodeService
}

func NewCodeServiceServer(svc service.CodeService) *CodeServiceServer {
	return &CodeServiceServer{svc: svc}
}

func (c *CodeServiceServer) Register(server *grpc.Server) {
	codev1.RegisterCodeServiceServer(server, c)
}

func (c *CodeServiceServer) Send(ctx context.Context, request *codev1.CodeSendRequest) (*codev1.CodeSendResponse, error) {
	err := c.svc.Send(ctx, request.GetBiz(), request.GetPhone())
	return &codev1.CodeSendResponse{}, err
}

func (c *CodeServiceServer) Verify(ctx context.Context, request *codev1.CodeVerifyRequest) (*codev1.CodeVerifyResponse, error) {
	verify, err := c.svc.Verify(ctx, request.GetBiz(), request.GetPhone(), request.GetInputCode())
	if err != nil {
		return nil, err
	}
	return &codev1.CodeVerifyResponse{
		Answer: verify,
	}, nil
}
