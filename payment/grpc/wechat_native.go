package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/payment/v1"
	"webook/payment/domain"
	"webook/payment/service/wechat"
)

type WechatPaymentServiceServer struct {
	// UnimplementedWechatPaymentServiceServer 提供了所有方法的默认实现（返回未实现错误）
	// 这样你只需覆盖需要实现的方法（如 NativePrePay）
	// 如果服务接口添加新方法，现有代码不会编译失败
	// 新方法会通过内嵌的 Unimplemented 结构体提供默认实现
	pmtv1.UnimplementedWechatPaymentServiceServer
	svc *wechat.NativePaymentService
}

func NewWechatPaymentServiceServer(svc *wechat.NativePaymentService) *WechatPaymentServiceServer {
	return &WechatPaymentServiceServer{svc: svc}
}

func (w *WechatPaymentServiceServer) Register(s *grpc.Server) {
	pmtv1.RegisterWechatPaymentServiceServer(s, w)
}

func (w *WechatPaymentServiceServer) GetPayment(ctx context.Context, request *pmtv1.GetPaymentRequest) (*pmtv1.GetPaymentResponse, error) {
	pmt, err := w.svc.GetPayment(ctx, request.GetBizTradeNo())
	if err != nil {
		return nil, err
	}
	return &pmtv1.GetPaymentResponse{
		Status: pmtv1.PaymentStatus(pmt.Status),
	}, nil
}

func (w *WechatPaymentServiceServer) NativePrePay(ctx context.Context, request *pmtv1.PrepayRequest) (*pmtv1.PrepayResponse, error) {
	CodeURL, err := w.svc.Prepay(ctx, domain.Payment{
		Description: request.GetDescription(),
		BizTradeNo:  request.GetBizTradeNo(),
		Amt: domain.Amount{
			Total:    request.GetAmt().GetTotal(),
			Currency: request.GetAmt().GetCurrency(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &pmtv1.PrepayResponse{
		CodeUrl: CodeURL,
	}, nil
}
