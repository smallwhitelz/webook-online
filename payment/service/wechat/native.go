package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"gorm.io/gorm"
	"time"
	"webook/payment/domain"
	"webook/payment/events"
	"webook/payment/repository"
	"webook/pkg/logger"
)

var errUnknownTransactionState = errors.New("未知的微信事务状态")

// NativePaymentService 对接微信支付：https://pay.weixin.qq.com/doc/v3/merchant/4012791877
type NativePaymentService struct {
	// APPID是微信开放平台(移动应用)或微信公众平台(小程序、公众号)为开发者的应用程序提供的唯一标识
	// 必填，每次都一样
	appID string
	// 是由微信支付系统生成并分配给每个商户的唯一标识符
	// 必填，单个商户每次传递都一样
	mchID string
	// 商户接收支付成功回调通知的地址，创单时传入
	// 必填
	notifyURL string
	// 自己的支付记录
	repo repository.PaymentRepository

	// 用于保证消息发送一定成功
	// 本地消息表的玩法
	// 没有这个保证的话 repov1和msgRepo可以去掉
	repov1  *repository.PaymentGORMRepository
	msgRepo repository.LocalMsgRepository

	// 构造向wxzf的请求，wxzf提供的方式
	svc *native.NativeApiService

	producer events.Producer

	l logger.LoggerV1

	// 在微信 native 里面，分别是
	// SUCCESS：支付成功
	// REFUND：转入退款
	// NOTPAY：未支付
	// CLOSED：已关闭
	// REVOKED：已撤销（付款码支付）
	// USERPAYING：用户支付中（付款码支付）
	// PAYERROR：支付失败(其他原因，如银行返回失败)
	nativeCBTypeToStatus map[string]domain.PaymentStatus
}

func NewNativePaymentService(appID string, mchID string, db *gorm.DB,
	repo repository.PaymentRepository, svc *native.NativeApiService, l logger.LoggerV1) *NativePaymentService {
	return &NativePaymentService{
		appID: appID, mchID: mchID,
		notifyURL: "http://wechat.zl.com/pay/callback",
		repo:      repo, svc: svc, l: l,
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"SUCCESS":  domain.PaymentStatusSuccess,
			"PAYERROR": domain.PaymentStatusFailed,
			"NOTPAY":   domain.PaymentStatusInit,
			"CLOSED":   domain.PaymentStatusFailed,
			"REVOKED":  domain.PaymentStatusFailed,
			"REFUND":   domain.PaymentStatusRefund,
			// 其它状态你都可以加
		},
		repov1:  repository.NewPaymentGORMRepository(db),
		msgRepo: repository.NewLocalMsgGORMRepository(db),
	}
}

func (n *NativePaymentService) Prepay(ctx context.Context, pmt domain.Payment) (string, error) {
	pmt.Status = domain.PaymentStatusInit
	// 自己记录一下wx的支付记录
	err := n.repo.AddPayment(ctx, pmt)
	if err != nil {
		return "", err
	}
	// 和wx打交道，返回一个支付二维码连接
	// 当下我们的需求只是打赏，所以极简设计
	// 如果其他需求，那就需要更多的字段
	resp, _, err := n.svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(n.appID),
		Mchid:       core.String(n.mchID),
		Description: core.String(pmt.Description),
		OutTradeNo:  core.String(pmt.BizTradeNo),
		// 支付的过期时间，30分钟后不可以支付
		TimeExpire: core.Time(time.Now().Add(time.Minute * 30)),
		Amount: &native.Amount{
			Total:    core.Int64(pmt.Amt.Total),
			Currency: core.String(pmt.Amt.Currency),
		},
	})
	if err != nil {
		return "", err
	}

	return *resp.CodeUrl, err
}

// SyncWechatInfo 兜底机制
// 利用 bizTradeNo 主动去wx查询订单状态，而后更新本地的信息
func (n *NativePaymentService) SyncWechatInfo(ctx context.Context, bizTradeNo string) error {
	// 对账
	txn, _, err := n.svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizTradeNo),
		Mchid:      core.String(n.mchID),
	})
	if err != nil {
		return err
	}
	return n.updateByTxn(ctx, txn)
}

func (n *NativePaymentService) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error) {
	return n.repo.FindExpiredPayment(ctx, offset, limit, t)
}

func (n *NativePaymentService) GetPayment(ctx context.Context, bizTradeOn string) (domain.Payment, error) {
	return n.repo.GetPayment(ctx, bizTradeOn)
}

func (n *NativePaymentService) HandleCallback(ctx context.Context, txn *payments.Transaction) error {
	return n.updateByTxn(ctx, txn)
}

func (n *NativePaymentService) updateByTxn(ctx context.Context, txn *payments.Transaction) error {
	// 将wx返回的状态码转为我自己的状态码
	status, ok := n.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		return fmt.Errorf("%w, 微信的状态是 %s", errUnknownTransactionState, *txn.TradeState)
	}
	// 很显然就是更新一下本地数据库payment的状态
	err := n.repo.UpdatePayment(ctx, domain.Payment{
		// wx过来的 Transaction id，因为有些地方需要
		TxnID: *txn.TransactionId,
		// 因为要找到更新那一条payment
		BizTradeNo: *txn.OutTradeNo,
		Status:     status,
	})
	if err != nil {
		return err
	}
	// 通知业务方，payment状态有变化，业务方肯定也要跟着变化
	// 有些人的系统里，会根据支付状态决定要不要发送消息
	// 我要是发消息失败了怎么办?
	// 站在业务的角度，是不是应该至少发成功一次
	// 思路：可以重试+监控+告警
	//      异步补偿机制
	err1 := n.producer.ProducePaymentEvent(ctx, events.PaymentEvent{
		BizTradeNo: *txn.OutTradeNo,
		Status:     status.AsUint8(),
	})
	if err1 != nil {
		n.l.Error("发送支付回调成功事件失败", logger.Error(err1), logger.String("biz_trade_no", *txn.OutTradeNo))
	}
	return nil
}

// updateByTxnV1
// 使用本地消息表的关键就是要把更新支付状态和插入代发送消息在一个数据库事务内完成操作
// 这一步我们下沉到了 repository，来规避在 service 上操作本地数据库事务
func (n *NativePaymentService) updateByTxnV1(ctx context.Context, txn *payments.Transaction) error {
	status, ok := n.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		return fmt.Errorf("%w, 微信的状态是 %s", errUnknownTransactionState, *txn.TradeState)
	}
	evt := events.PaymentEvent{
		BizTradeNo: *txn.OutTradeNo,
		Status:     status.AsUint8(),
	}
	var msgId int64
	// 在这种情况下，你无可避免会和底层耦合在一起
	err := n.repov1.Transaction(ctx, func(pmt *repository.PaymentGORMRepository, msg repository.LocalMsgRepository) error {
		// 这里其实就是我们更新支付状态一定是和我们写入本地消息是同一个事务的
		err1 := pmt.UpdatePayment(ctx, domain.Payment{
			// 微信过来的 transaction id
			TxnID:      *txn.TransactionId,
			BizTradeNo: *txn.OutTradeNo,
			Status:     status,
		})
		if err1 != nil {
			return err1
		}
		evtData, err1 := json.Marshal(evt)
		if err1 != nil {
			return err1
		}
		msgId, err1 = n.msgRepo.AddMsg(ctx, string(evtData))
		return err1
	})
	if err != nil {
		return err
	}
	err1 := n.producer.ProducePaymentEvent(ctx, evt)
	if err1 != nil {
		// 失败的时候，我们并没有将本地消息表标记为失败，是因为我们后面还想继续重试
		// msg 表里面有一条处于 StatusInit 的数据
		n.l.Error("发送支付事件失败", logger.Error(err1), logger.String("biz_trade_no", *txn.OutTradeNo))
		return nil
	}
	// 更新本地消息表状态
	// 这里我认为即便是本地消息表更新失败，也不是业务失败
	err1 = n.msgRepo.MarkSuccess(ctx, msgId)
	if err1 != nil {
		// 没有把 msg 标记为发送成功
		// 消息队列里面会有至少两条消息
		n.l.Error("将本地消息表标记为成功操作失败", logger.Error(err1),
			logger.String("biz_trade_no", *txn.OutTradeNo))
	}
	return nil
}
