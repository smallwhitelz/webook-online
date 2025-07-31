package service

import (
	"context"
	"webook/oauth2/domain"
)

type Oauth2Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}
