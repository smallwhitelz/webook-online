package service

import (
	"context"
	"webook/account/domain"
	"webook/account/repository"
	"webook/pkg/logger"
)

type accountService struct {
	repo repository.AccountRepository
	l    logger.LoggerV1
}

func NewAccountService(repo repository.AccountRepository, l logger.LoggerV1) AccountService {
	return &accountService{repo: repo, l: l}
}

// Credit 记账幂等性，幂等归根结底一定是依赖唯一索引的，没有唯一索引完全谈不上幂等
// 保证同一个 biz+bizId，不会重复记账。
// 并且要考虑有什么优化方案，可以进一步提高性能。参考唯一索引、Redis 去重、布隆过滤器等。
// 但是其实布隆过滤器会出现假阳性的问题，之所以可以用是因为我们的兜底是唯一索引
func (a *accountService) Credit(ctx context.Context, cr domain.Credit) error {
	err := a.repo.CheckUnique(ctx, cr)
	if err != nil {
		return err
	}
	// 这里有唯一索引
	err = a.repo.AddCredit(ctx, cr)
	if err == nil {
		// 注意这些部分失败是没有什么问题的
		// 因为我们始终有一个兜底，就是唯一索引。
		err1 := a.repo.SetUnique(ctx, cr)
		if err1 != nil {
			// 记录日志即可
			a.l.Error("使用redis设置记账幂等性失败",
				logger.Error(err1),
				logger.String("biz", cr.Biz),
				logger.Int64("biz_id", cr.BizId))
		}
	}
	return err
}
