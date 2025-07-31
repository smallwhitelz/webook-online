package job

import (
	"context"
	"time"
	"webook/payment/service/wechat"
	"webook/pkg/logger"
)

type SyncWechatOrder struct {
	svc *wechat.NativePaymentService
	l   logger.LoggerV1
}

func NewSyncWechatOrder(svc *wechat.NativePaymentService, l logger.LoggerV1) *SyncWechatOrder {
	return &SyncWechatOrder{svc: svc, l: l}
}

func (s *SyncWechatOrder) Name() string {
	return "sync_wechat_order_job"
}

// Run 定时任务多久运行一次？
// 不必特别频繁，比如1分钟运行一次
func (s *SyncWechatOrder) Run() error {
	// 定时找到超时的wx支付订单，然后发起同步
	// 针对过期订单
	// 多一分钟的好处在于，临界点30分钟可能对账和用户支付会发生冲突
	// 比如用户刚好在30分钟的时候支付，而对账功能也在30分钟启用，用户先付账了，更新了订单状态
	// 然后对账系统又把用户支付了的订单状态改为未支付，这里会出现这样的冲突
	t := time.Now().Add(-time.Minute * 31)
	offset := 0
	const limit = 100
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		pmts, err := s.svc.FindExpiredPayment(ctx, offset, limit, t)
		cancel()
		if err != nil {
			// 建议直接中断，因为这里出错很有可能就是数据库出了问题
			return err
		}
		for _, pmt := range pmts {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			err := s.svc.SyncWechatInfo(ctx, pmt.BizTradeNo)
			cancel()
			if err != nil {
				s.l.Error("同步微信订单状态失败", logger.Error(err),
					logger.String("biz_trade_no", pmt.BizTradeNo))
			}
		}
		// 中断条件
		if len(pmts) < limit {
			return nil
		}
		offset = offset + len(pmts)
	}
}
