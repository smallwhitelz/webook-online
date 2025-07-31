package domain

type Credit struct {
	Biz   string
	BizId int64
	Items []CreditItem
}

type CreditItem struct {
	Account     int64
	AccountType AccountType
	Amt         int64
	Currency    string
	Uid         int64
}

type AccountType uint8

func (a AccountType) AsUint8() uint8 {
	return uint8(a)
}

const (
	AccountTypeUnknown = iota
	// AccountTypeReward 个人赞赏账号
	AccountTypeReward
	// AccountTypeSystem 平台分成账号
	AccountTypeSystem
)
