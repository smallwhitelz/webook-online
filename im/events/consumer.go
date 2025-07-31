package events

import (
	"context"
	"github.com/IBM/sarama"
	"strconv"
	"time"
	"webook/im/domain"
	"webook/im/service"
	"webook/pkg/canalx"
	"webook/pkg/logger"
	"webook/pkg/saramax"
)

type MysqlBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.UserService
}

func (m *MysqlBinlogConsumer) Start() error {
	sg, err := sarama.NewConsumerGroupFromClient("openim_sync", m.client)
	if err != nil {
		return err
	}
	go func() {
		err := sg.Consume(context.Background(), []string{"webook_binlog"}, saramax.NewHandler[canalx.Message[User]](m.l, m.Consume))
		if err != nil {
			m.l.Error("退出循环消费", logger.Error(err))
		}
	}()
	return err
}

func (m *MysqlBinlogConsumer) Consume(msg *sarama.ConsumerMessage, val canalx.Message[User]) error {
	if val.Table != "users" || val.Type != "INSERT" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for _, data := range val.Data {
		err := m.svc.Sync(ctx, domain.User{
			UserID:   strconv.FormatInt(data.Id, 10),
			Nickname: data.Nickname,
		})
		if err != nil {
			// 记录日志
			continue
		}
	}
	return nil
}

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Nickname string `json:"nickname"`
	Birthday int64  `json:"birthday"`
	AboutMe  string `json:"about_me"`
	Phone    string `json:"phone"`

	WechatOpenId  string `json:"wechat_open_id"`
	WechatUnionId string `json:"wechat_union_id"`

	Ctime int64 `json:"ctime"`
	Utime int64 `json:"utime"`
}
