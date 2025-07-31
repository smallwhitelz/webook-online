package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"webook/oauth2/domain"
	"webook/pkg/logger"
)

const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	appID     string
	appSecret string
	client    *http.Client
	l         logger.LoggerV1
}

func NewService(appID string, appSecret string, l logger.LoggerV1) Oauth2Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
		l:         l,
	}
}

// VerifyCode 给微信发一个http请求
func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const baseURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	// 这是另外一种写法
	queryParams := url.Values{}
	queryParams.Set("appid", s.appID)
	queryParams.Set("secret", s.appSecret)
	queryParams.Set("code", code)
	queryParams.Set("grant_type", "authorization_code")
	accessTokenURL := baseURL + "?" + queryParams.Encode()
	req, err := http.NewRequest("GET", accessTokenURL, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	// 把最上层传进来的 ctx 关联到请求里，以便超时、取消能下传到底层 HTTP 调用。
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	defer resp.Body.Close()
	var res Result

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, errors.New("换取 access_token 失败")
	}
	return domain.WechatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionId,
	}, nil
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}

type Result struct {
	// 接口调用凭证
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	OpenId string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`

	// 错误返回

	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
