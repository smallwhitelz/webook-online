package service

import (
	"context"
	"webook/reward/domain"
)

//go:generate mockgen -source=./types.go -destination=./mocks/reward.mock.go -package=svcmocks RewardService
type RewardService interface {
	// PreReward 准备打赏，
	// 你也可以直接理解为对标到创建一个打赏的订单
	// 因为目前我们只支持微信扫码支付，所以实际上直接把接口定义成这个样子就可以了
	PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	// UpdateReward 从payment的消息队列中监听payment的支付结果
	// 这里相当于是payment -> reward
	// 如果你还有其他业务需要支付，那么也可以监听
	UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error
	GetReward(ctx context.Context, rid int64, uid int64) (domain.Reward, error)
}
