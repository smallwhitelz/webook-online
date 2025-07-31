package events

// PaymentEvent 也是最简设计
// 可以把支付详情也放进来，但是目前看来没必要
// 后续如果考虑接入大数据的话，可以增加其他的字段
type PaymentEvent struct {
	BizTradeNo string
	Status     uint8
}

func (PaymentEvent) Topic() string {
	return "payment_events"
}
