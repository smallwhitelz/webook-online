package grpc

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
	"webook/account/domain"
	"webook/account/service"
	"webook/api/proto/gen/account/v1"
)

type AccountServiceServer struct {
	accountv1.UnimplementedAccountServiceServer
	svc service.AccountService
}

func NewAccountServiceServer(svc service.AccountService) *AccountServiceServer {
	return &AccountServiceServer{svc: svc}
}

func (a *AccountServiceServer) Register(s *grpc.Server) {
	accountv1.RegisterAccountServiceServer(s, a)
}

func (a *AccountServiceServer) Credit(ctx context.Context, request *accountv1.CreditRequest) (*accountv1.CreditResponse, error) {
	//res:=make([]domain.CreditItem,0,len(request.GetItems()))
	//for _, item := range request.Items {
	//	res = append(res,a.itemToDomain(item))
	//}
	err := a.svc.Credit(ctx, domain.Credit{
		Biz:   request.GetBiz(),
		BizId: request.GetBizId(),
		Items: slice.Map(request.GetItems(), func(idx int, src *accountv1.CreditItem) domain.CreditItem {
			return a.itemToDomain(src)
		}),
	})
	return &accountv1.CreditResponse{}, err
}

func (a *AccountServiceServer) itemToDomain(c *accountv1.CreditItem) domain.CreditItem {
	return domain.CreditItem{
		Account:     c.Account,
		AccountType: domain.AccountType(c.AccountType),
		Amt:         c.Amt,
		Currency:    c.Currency,
		Uid:         c.Uid,
	}
}
