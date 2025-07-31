package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"webook/search/domain"
	"webook/search/repository/dao"
)

type userSearchRepository struct {
	userDao dao.UserSearchDAO
}

func NewUserSearchRepository(userDao dao.UserSearchDAO) UserSearchRepository {
	return &userSearchRepository{userDao: userDao}
}

func (s *userSearchRepository) SearchUser(ctx context.Context, keywords []string) ([]domain.User, error) {
	users, err := s.userDao.Search(ctx, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(users, func(idx int, src dao.User) domain.User {
		return domain.User{
			Id:       src.Id,
			Nickname: src.Nickname,
			Email:    src.Email,
			Phone:    src.Phone,
		}
	}), nil
}

func (s *userSearchRepository) InputUser(ctx context.Context, user domain.User) error {
	return s.userDao.InputUser(ctx, dao.User{
		Id:       user.Id,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
	})
}
