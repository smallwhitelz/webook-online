package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
	"webook/payment/domain"
	"webook/payment/repository/dao"
)

type PaymentGORMRepository struct {
	dao dao.PaymentDAO
	db  *gorm.DB
}

func (p *PaymentGORMRepository) AddPayment(ctx context.Context, pmt domain.Payment) error {
	return p.dao.Insert(ctx, p.toEntity(pmt))
}

func (p *PaymentGORMRepository) GetPayment(ctx context.Context, bizTradeOn string) (domain.Payment, error) {
	pmt, err := p.dao.GetPayment(ctx, bizTradeOn)
	if err != nil {
		return domain.Payment{}, err
	}
	return p.toDomain(pmt), nil
}

func (p *PaymentGORMRepository) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error) {
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

func (p *PaymentGORMRepository) UpdatePayment(ctx context.Context, pmt domain.Payment) error {
	return p.dao.UpdateTxnIDAndStatus(ctx, pmt.BizTradeNo, pmt.TxnID, pmt.Status)
}

// Transaction 目前只会控制两个 dao，所以可以在 cb 里面直接传入两个 DAO
// 相当于聚合Payment和localMsg
func (p *PaymentGORMRepository) Transaction(ctx context.Context,
	cb func(pmt *PaymentGORMRepository, msg LocalMsgRepository) error) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return cb(NewPaymentGORMRepository(tx), NewLocalMsgGORMRepository(tx))
	})
}

func NewPaymentGORMRepository(db *gorm.DB) *PaymentGORMRepository {
	return &PaymentGORMRepository{
		dao: dao.NewPaymentGORMDAO(db),
		db:  db,
	}
}

func (p *PaymentGORMRepository) toEntity(pmt domain.Payment) dao.Payment {
	return dao.Payment{
		AmtTotal:    pmt.Amt.Total,
		Currency:    pmt.Amt.Currency,
		Description: pmt.Description,
		BizTradeNO:  pmt.BizTradeNo,
		Status:      pmt.Status.AsUint8(),
	}
}

func (p *PaymentGORMRepository) toDomain(pmt dao.Payment) domain.Payment {
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
