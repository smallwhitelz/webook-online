package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	"webook/oauth2/domain"
	"webook/pkg/logger"
	domain2 "webook/user/domain"
	"webook/user/repository/cache"
	"webook/user/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound  = dao.ErrRecordNotFound
)

//go:generate mockgen -source=./user.go -package=repomocks -destination=./mocks/user.mock.go UserRepository
type UserRepository interface {
	Create(ctx context.Context, user domain2.User) error
	FindByEmail(ctx context.Context, email string) (domain2.User, error)
	FindById(ctx context.Context, uid int64) (domain2.User, error)
	FindByPhone(ctx context.Context, phone string) (domain2.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain2.User) error
	FindByWechat(ctx context.Context, openId string) (domain2.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
	l     logger.LoggerV1
}

func NewCachedUserRepository(dao dao.UserDAO, c cache.UserCache, l logger.LoggerV1) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
		l:     l,
	}
}

func (repo *CachedUserRepository) FindByWechat(ctx context.Context, openId string) (domain2.User, error) {
	ue, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain2.User{}, err
	}
	return repo.toDomain(ue), nil
}

func (repo *CachedUserRepository) Create(ctx context.Context, user domain2.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain2.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain2.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain2.User {
	return domain2.User{
		Id:          u.Id,
		Email:       u.Email.String,
		Phone:       u.Phone.String,
		Password:    u.Password,
		Nickname:    u.Nickname,
		Birthday:    time.UnixMilli(u.Birthday),
		Description: u.Description,
		Ctime:       time.UnixMilli(u.Ctime),
		WechatInfo: domain.WechatInfo{
			OpenId:  u.WechatOpenId.String,
			UnionId: u.WechatUnionId.String,
		},
	}
}

func (repo *CachedUserRepository) toEntity(u domain2.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionId,
			Valid:  u.WechatInfo.UnionId != "",
		},
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenId,
			Valid:  u.WechatInfo.OpenId != "",
		},
		Nickname:    u.Nickname,
		Birthday:    u.Birthday.UnixMilli(),
		Description: u.Description,
	}
}

// FindById 建议用这种方式，虽然如果redis崩了，数据量会全打到数据库中
// 但是真正的高并发高流量场景把redis打崩的情况很少
func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain2.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err为 nil，就返回
	if err == nil {
		return du, nil
	}
	// 检测限流/熔断/降级标记位
	// 降级后就直接走快路径，也就是走redis不走数据库
	// 这里因为已经触发系统故障，所以为了保证就算大量请求，也不会打垮数据库
	if ctx.Value("downgrade") == "true" {
		return du, errors.New("触发降级，不再查询数据库")
	}
	// err 不为 nil，就要查询数据库
	// err 有两种可能
	// 1. key 不存在，说明redis是正常的
	// 2. 访问redis有问题。可能是网络有问题，也有可能是 redis 本身就崩溃了
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain2.User{}, err
	}
	du = repo.toDomain(u)
	// 异步
	// go中用异步可以，因为go开启一个异步非常方便
	// 单测中测试不到异步
	//go func() {
	//	err = repo.cache.Set(ctx, du)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()
	// 同步
	// java中用同步，因为java开一个线程很麻烦
	// 但是在这里同步和异步的效率优化并不大，所以都可以
	err = repo.cache.Set(ctx, du)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		repo.l.Error("FindById接口：回写user缓存失败", logger.Error(err), logger.Int64("uid", du.Id))
	}
	return du, nil
}

// FindByIdV1 考虑到redis崩溃以及网络不通的问题，保护住了数据库，避免瞬时流量过大打崩数据库，防止缓存击穿
func (repo *CachedUserRepository) FindByIdV1(ctx context.Context, uid int64) (domain2.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err为 nil，就返回
	switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExist:
		// 1. key 不存在，说明redis是正常的
		u, err := repo.dao.FindById(ctx, uid)
		if err != nil {
			return domain2.User{}, err
		}
		du = repo.toDomain(u)
		go func() {
			err = repo.cache.Set(ctx, du)
			if err != nil {
				log.Println(err)
			}
		}()
		return du, nil
	default:
		// 2. 访问redis有问题。可能是网络有问题，也有可能是 redis 本身就崩溃了
		// 接近降级的写法
		return domain2.User{}, err
	}
}

func (repo *CachedUserRepository) UpdateNonSensitiveInfo(ctx context.Context, user domain2.User) error {
	err := repo.dao.UpdateById(ctx, repo.toEntity(user))
	if err != nil {
		return err
	}
	// 延迟一秒
	time.AfterFunc(time.Second, func() {
		_ = repo.cache.Del(ctx, user.Id)
	})
	// 更新缓存,这里采用的是延迟双删避免db和缓存中的数据不一致问题
	return repo.cache.Del(ctx, user.Id)
}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain2.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain2.User{}, err
	}
	return repo.toDomain(u), nil
}
