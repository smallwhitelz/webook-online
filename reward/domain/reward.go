package domain

type Reward struct {
	Id int64
	// 代表我被打赏的东西
	// 业务
	Biz string
	// 业务id
	BizId int64
	// 给用户看的，让用户明白自己打赏的是什么东西
	BizName string
	// 被打赏人
	TargetUid int64
	// 打赏人
	Uid int64
	// 金额
	Amt int64
	// 打赏状态
	Status RewardStatus
}

// IsComplete 是否已经完成
// 目前来说，也就是是否处理了支付回调
func (r Reward) IsComplete() bool {
	return r.Status == RewardStatusPayed || r.Status == RewardStatusFailed
}

type RewardStatus uint8

func (r RewardStatus) AsUint8() uint8 {
	return uint8(r)
}

const (
	RewardStatusUnknown = iota
	RewardStatusInit
	RewardStatusPayed
	RewardStatusFailed
)

type CodeURL struct {
	Url string
	Rid int64
}
