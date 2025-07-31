package domain

type Payment struct {
	// 商品信息描述
	Description string
	// 代表业务，业务方决定生成
	// Payment不管，表达这一次支付的唯一索引凭证
	BizTradeNo string
	// 订单金额
	Amt Amount

	// 记录支付状态
	Status PaymentStatus

	// 第三方那边返回的ID
	TxnID string
}

type Amount struct {
	// 总金额
	Total int64
	// 货币类型
	Currency string
}

type PaymentStatus uint8

func (s PaymentStatus) AsUint8() uint8 {
	return uint8(s)
}

const (
	PaymentStatusUnknown = iota
	PaymentStatusInit
	PaymentStatusSuccess
	PaymentStatusFailed
	PaymentStatusRefund
)
