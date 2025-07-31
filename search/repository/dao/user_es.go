package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const UserIndexName = "user_index"

type User struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UserESDAO struct {
	client *elastic.Client
}

func NewUserESDAO(client *elastic.Client) UserSearchDAO {
	return &UserESDAO{client: client}
}

func (u *UserESDAO) Search(ctx context.Context, keywords []string) ([]User, error) {
	queryString := strings.Join(keywords, " ")
	query := elastic.NewMatchQuery("nickname", queryString)
	resp, err := u.client.Search(UserIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]User, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var user User
		err := json.Unmarshal(hit.Source, &user)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}

func (u *UserESDAO) InputUser(ctx context.Context, user User) error {
	_, err := u.client.Index().Index(UserIndexName).
		// 这里保持insert or update 语义
		// 因为同一个用户数据过来我们肯定只希望有一份数据就好
		// 这里指定了文档ID，所以就肯定只有一份数据
		Id(strconv.FormatInt(user.Id, 10)).
		BodyJson(user).
		Do(ctx)
	return err
}
