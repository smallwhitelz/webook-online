package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/oauth2/v1"
	"webook/oauth2/service"
)

type Oauth2ServiceServer struct {
	oauth2v1.UnimplementedOauth2ServiceServer
	svc service.Oauth2Service
}

func NewOauth2ServiceServer(svc service.Oauth2Service) *Oauth2ServiceServer {
	return &Oauth2ServiceServer{svc: svc}
}

func (o *Oauth2ServiceServer) Register(server *grpc.Server) {
	oauth2v1.RegisterOauth2ServiceServer(server, o)
}

func (o *Oauth2ServiceServer) AuthURL(ctx context.Context, request *oauth2v1.AuthURLRequest) (*oauth2v1.AuthURLResponse, error) {
	url, err := o.svc.AuthURL(ctx, request.GetState())
	if err != nil {
		return nil, err
	}
	return &oauth2v1.AuthURLResponse{
		Url: url,
	}, nil
}

func (o *Oauth2ServiceServer) VerifyCode(ctx context.Context, request *oauth2v1.VerifyCodeRequest) (*oauth2v1.VerifyCodeResponse, error) {
	wechatInfo, err := o.svc.VerifyCode(ctx, request.GetCode())
	if err != nil {
		return nil, err
	}
	return &oauth2v1.VerifyCodeResponse{
		WechatInfo: &oauth2v1.WechatInfo{
			OpenId:  wechatInfo.OpenId,
			UnionId: wechatInfo.UnionId,
		},
	}, nil
}
