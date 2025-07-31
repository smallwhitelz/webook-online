package service

import (
	"context"
	"fmt"
	"math/rand"
	smsv1 "webook/api/proto/gen/sms/v1"
	"webook/code/repository"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

//go:generate mockgen -source=./code.go -package=svcmocks -destination=./mocks/code.mock.go CodeService
type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type codeService struct {
	repo repository.CodeRepository
	sms  smsv1.SmsServiceClient
}

func NewCodeService(repo repository.CodeRepository, sms smsv1.SmsServiceClient) CodeService {
	return &codeService{
		repo: repo,
		sms:  sms,
	}
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	// 模仿验证码，随机生成
	code := svc.generate()
	// 将验证码放到redis中
	err := svc.repo.Set(ctx, biz, phone, code)
	// 这里就要开始发送验证码了
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	_, err = svc.sms.Send(ctx, &smsv1.SmsSendRequest{
		TplId:   codeTplId,
		Args:    []string{code},
		Numbers: []string{phone},
	})
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者这个不对
		return false, nil
	}
	return ok, err
}

func (svc *codeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
