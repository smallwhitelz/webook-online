package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"webook/oauth2/domain"
	domain2 "webook/user/domain"
	"webook/user/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户名或密码错误")
)

//go:generate mockgen -source=./user.go -package=svcmocks -destination=./mocks/user.mock.go UserService
type UserService interface {
	Signup(ctx context.Context, u domain2.User) error
	Login(ctx context.Context, email string, password string) (domain2.User, error)
	FindById(ctx context.Context, uid int64) (domain2.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain2.User) error
	FindOrCreate(ctx context.Context, phone string) (domain2.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain2.User, error)
}

type userService struct {
	repo repository.UserRepository
	//logger *zap.Logger
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		//logger: zap.L(),
	}
}

func (svc *userService) Signup(ctx context.Context, u domain2.User) error {
	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
	// 这里如果你接入了IM，需要调用IM的api将注册的新用户导入到OpenIM的用户中
	// 那么就可以在这里写一个生产者，用户注册成功后发送消息到kafka
	// 这里接住canal保证数据的一致性
	// 然后在IM哪里消费这条消息进行导入
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain2.User, error) {
	// 先通过邮箱找到用户
	u, err := svc.repo.FindByEmail(ctx, email)
	// 没有的话报错找不到用户
	if err == repository.ErrUserNotFound {
		return domain2.User{}, ErrInvalidUserOrPassword
	}
	// 这里报错可能是数据库崩溃之类的错误
	if err != nil {
		return domain2.User{}, err
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	// 校验密码失败，报错用户或者密码错误
	if err != nil {
		return domain2.User{}, ErrInvalidUserOrPassword
	}
	// 找到用户并返回
	return u, nil
}

func (svc *userService) FindById(ctx context.Context, uid int64) (domain2.User, error) {
	return svc.repo.FindById(ctx, uid)
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain2.User) error {
	return svc.repo.UpdateNonSensitiveInfo(ctx, user)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain2.User, error) {
	// 先找一下，我们认为大部分用户是已经存在的用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// 有两种情况
		// err == nil，说明用户已经存在
		// err != nil，系统错误
		return u, err
	}
	// 用户没找到
	err = svc.repo.Create(ctx, domain2.User{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突 (phone)
	// 一种是 err != nil，说明系统错误
	if err != nil && err != repository.ErrDuplicateUser {
		return domain2.User{}, err
	}
	// 要么err==nil，要么ErrDuplicateUser也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain2.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	if err != repository.ErrUserNotFound {
		return u, err
	}

	// 这边意味着是一个新用户
	// JSON 格式的 wechatInfo
	// 直接使用包变量
	zap.L().Info("这是一个新用户", zap.Any("wechatInfo", wechatInfo))
	// 使用注入的logger
	//svc.logger.Info("这是一个新用户", zap.Any("wechatInfo", wechatInfo))

	err = svc.repo.Create(ctx, domain2.User{
		WechatInfo: wechatInfo,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain2.User{}, err
	}
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
}
