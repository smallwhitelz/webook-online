package repository

import (
	"context"
	"time"
	"webook/payment/domain"
	"webook/payment/repository/dao"
)

type paymentRepository struct {
	dao dao.PaymentDAO
}

func NewPaymentRepository(dao dao.PaymentDAO) PaymentRepository {
	return &paymentRepository{dao: dao}
}

func (p *paymentRepository) GetPayment(ctx context.Context, bizTradeOn string) (domain.Payment, error) {
	pmt, err := p.dao.GetPayment(ctx, bizTradeOn)
	if err != nil {
		return domain.Payment{}, err
	}
	return p.toDomain(pmt), nil
}

func (p *paymentRepository) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error) {
	pmts, err := p.dao.FindExpiredPayment(ctx, offset, limit, t)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Payment, 0, len(pmts))
	for _, pmt := range pmts {
		res = append(res, p.toDomain(pmt))
	}
	return res, nil
}

func (p *paymentRepository) UpdatePayment(ctx context.Context, pmt domain.Payment) error {
	return p.dao.UpdateTxnIDAndStatus(ctx, pmt.BizTradeNo, pmt.TxnID, pmt.Status)
}

func (p *paymentRepository) AddPayment(ctx context.Context, pmt domain.Payment) error {
	return p.dao.Insert(ctx, p.toEntity(pmt))
}

func (p *paymentRepository) toEntity(pmt domain.Payment) dao.Payment {
	return dao.Payment{
		AmtTotal:    pmt.Amt.Total,
		Currency:    pmt.Amt.Currency,
		Description: pmt.Description,
		BizTradeNO:  pmt.BizTradeNo,
		Status:      pmt.Status.AsUint8(),
	}
}

func (p *paymentRepository) toDomain(pmt dao.Payment) domain.Payment {
	return domain.Payment{
		Description: pmt.Description,
		BizTradeNo:  pmt.BizTradeNO,
		Amt: domain.Amount{
			Total:    pmt.AmtTotal,
			Currency: pmt.Currency,
		},
		Status: domain.PaymentStatus(pmt.Status),
		TxnID:  pmt.TxnID.String,
	}
}
