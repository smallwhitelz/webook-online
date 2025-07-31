package service

import (
	"context"
	"webook/payment/domain"
)

// PaymentService 接入微信支付
type PaymentService interface {
	// Prepay 预支付，对应于微信创建订单的步骤
	Prepay(ctx context.Context, pmt domain.Payment) (string, error)
	GetPayment(ctx context.Context, bizTradeOn string) (domain.Payment, error)
}
