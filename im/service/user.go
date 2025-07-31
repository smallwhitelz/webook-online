package service

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/net/httpx"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"webook/im/domain"
)

type UserService interface {
	Sync(ctx context.Context, user domain.User) error
}

type RESTUserService struct {
	// http请求的域名端口
	base string
	// 一旦将来你要换client，很容易就换掉
	client *http.Client
}

func (s *RESTUserService) Sync(ctx context.Context, user domain.User) error {
	var operationID string
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		operationID = spanCtx.TraceID().String()
	} else {
		operationID = uuid.New().String()
	}
	var resp response
	err := httpx.NewRequest(ctx, http.MethodPost, s.base+"/user/user_register").
		AddHeader("operationID", operationID).
		JSONBody(request{
			Users: []domain.User{user},
		}).Client(s.client).Do().JSONScan(&resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("同步用户数据失败 %v", resp)
	}
	return nil
}

type request struct {
	Users []domain.User `json:"users"`
}

type response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}
