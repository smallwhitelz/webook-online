package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	accountv1 "webook/api/proto/gen/account/v1"
	pmtv1 "webook/api/proto/gen/payment/v1"
	"webook/pkg/logger"
	"webook/reward/domain"
	"webook/reward/repository"
)

type WechatNativeRewardService struct {
	client pmtv1.WechatPaymentServiceClient
	acli   accountv1.AccountServiceClient
	repo   repository.RewardRepository
	l      logger.LoggerV1
}

func NewWechatNativeRewardService(client pmtv1.WechatPaymentServiceClient,
	acli accountv1.AccountServiceClient, repo repository.RewardRepository,
	l logger.LoggerV1) RewardService {
	return &WechatNativeRewardService{client: client, acli: acli, repo: repo, l: l}
}

// GetReward 个人获取打赏记录
func (w *WechatNativeRewardService) GetReward(ctx context.Context, rid int64, uid int64) (domain.Reward, error) {
	// 快路径
	res, err := w.repo.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	// 确保是自己打赏的
	if res.Uid != uid {
		return domain.Reward{}, errors.New("非法访问别人的打赏记录")
	}
	// 降级或者限流的时候，不走慢路径
	if ctx.Value("limited") == "true" {
		return res, nil
	}
	// 慢路径，
	// 查询打赏结果的时候，可能存在一个问题，reward消费者还没有收到这个消息，或者正在处理中也就是UpdateReward中
	// 所以这里引入慢路径，就是我们在发现打赏结果还处于一种未知状态时，去payment查询一下
	// 这里依然可以引入降级或者限流的时候，我们就放弃慢路径
	if !res.IsComplete() {
		pmt, err := w.client.GetPayment(ctx, &pmtv1.GetPaymentRequest{
			BizTradeNo: w.bizTradeOn(rid),
		})
		if err != nil {
			w.l.Error("慢路径查询支付失败", logger.Error(err), logger.Int64("rid", rid))
			return res, nil
		}
		switch pmt.GetStatus() {
		case pmtv1.PaymentStatus_PaymentStatusFailed:
			res.Status = domain.RewardStatusFailed
		case pmtv1.PaymentStatus_PaymentStatusRefund:
			res.Status = domain.RewardStatusFailed
		case pmtv1.PaymentStatus_PaymentStatusSuccess:
			res.Status = domain.RewardStatusPayed
		case pmtv1.PaymentStatus_PaymentStatusInit:
			res.Status = domain.RewardStatusInit
		case pmtv1.PaymentStatus_PaymentStatusUnknown:
			res.Status = domain.RewardStatusUnknown
		}
		// 这里是幂等的，也就是说reward直接去查payment后，要及时更新自己保存的状态
		err = w.UpdateReward(ctx, w.bizTradeOn(rid), res.Status)
		if err != nil {
			w.l.Error("慢路径更新本地状态失败", logger.Error(err),
				logger.Int64("rid", rid))
		}
	}
	return res, nil
}

// UpdateReward 这个操作一定是幂等的，不管执行多少次，结果都是一样的
func (w *WechatNativeRewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	// 这里我们更新reward状态其实是根据id去更新的，所以这个要从bizTradeNO中提取出来rid
	// 也就是创建reward的时候返回的rid
	rid := w.toRid(bizTradeNO)
	err := w.repo.UpdateReward(ctx, rid, status)
	if err != nil {
		return err
	}
	// 完成了支付，准备入账
	if status == domain.RewardStatusPayed {
		// 获取到打赏的记录
		r, err := w.repo.GetReward(ctx, rid)
		if err != nil {
			return err
		}
		// webook抽成，也就是平台分到多少钱，这里按照10%计算
		weAmt := int64(float64(r.Amt) * 0.1)
		_, err = w.acli.Credit(ctx, &accountv1.CreditRequest{
			Biz:   "reward",
			BizId: rid,
			Items: []*accountv1.CreditItem{
				{
					AccountType: accountv1.AccountType_AccountTypeSystem,
					// 虽然可能为0，但是也要记录出来
					Amt:      weAmt,
					Currency: "CNY",
				},
				{
					Account:     r.Uid,
					Uid:         r.Uid,
					AccountType: accountv1.AccountType_AccountTypeReward,
					Amt:         r.Amt - weAmt,
					Currency:    "CNY",
				},
			},
		})
		if err != nil {
			w.l.Error("入账失败了，快来修数据啊！！！",
				logger.String("biz_trade_no", bizTradeNO),
				logger.Error(err))
			// 做好监控和告警，这里
			return err
		}
	}
	return nil
}

func (w *WechatNativeRewardService) toRid(bizTradeNO string) int64 {
	ridStr := strings.Split(bizTradeNO, "-")
	val, _ := strconv.ParseInt(ridStr[1], 10, 64)
	return val
}

func (w *WechatNativeRewardService) bizTradeOn(rid int64) string {
	return fmt.Sprintf("reward-%d", rid)
}

func (w *WechatNativeRewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	// 可以缓存一下CodeURL
	// 因为和第三方wx打交道走的是公网，公网会比较慢
	// 而且我们也不能保证wx一定不出问题，所以尽量少和wx打交道
	res, err := w.repo.GetCachedCodeURL(ctx, r)
	if err == nil {
		return res, nil
	}
	// 存一个打赏记录
	r.Status = domain.RewardStatusInit
	rid, err := w.repo.CreateReward(ctx, r)
	if err != nil {
		return domain.CodeURL{}, err
	}
	// 和wx打交道
	pmtResp, err := w.client.NativePrePay(ctx, &pmtv1.PrepayRequest{
		Amt: &pmtv1.Amount{
			Total:    r.Amt,
			Currency: "CNY",
		},
		// 一定要保证唯一性
		BizTradeNo:  fmt.Sprintf("reward-%d", rid),
		Description: fmt.Sprintf("打赏-%s", r.BizName),
	})
	if err != nil {
		return domain.CodeURL{}, err
	}
	cu := domain.CodeURL{
		Rid: rid,
		Url: pmtResp.CodeUrl,
	}
	// 设置二维码缓存
	err1 := w.repo.SetCachedCodeURL(ctx, cu, r)
	if err1 != nil {
		w.l.Error("缓存打赏二维码失败", logger.Error(err1),
			logger.Int64("rid", rid))
	}
	return cu, nil
}
