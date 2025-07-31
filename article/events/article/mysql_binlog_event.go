package article

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/article/domain"
	"webook/article/repository"
	"webook/article/repository/dao"
	"webook/pkg/canalx"
	"webook/pkg/logger"
	"webook/pkg/saramax"
)

// MysqlBinlogConsumer
// 这个例子我们主要是去更新pubart的缓存
type MysqlBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	// 这里利用binlog来更新缓存，是缓存策略，不是业务逻辑
	// 所以一般会绕开service，从repo进行操作
	// 而且这是具体的缓存策略，没必要在repo定义接口
	repo *repository.CachedArticleRepository
}

func NewMysqlBinlogConsumer(client sarama.Client, l logger.LoggerV1) *MysqlBinlogConsumer {
	return &MysqlBinlogConsumer{client: client, l: l, repo: &repository.CachedArticleRepository{}}
}

func (m *MysqlBinlogConsumer) Start() error {
	sg, err := sarama.NewConsumerGroupFromClient("pub_articles_cache", m.client)
	if err != nil {
		return err
	}
	go func() {
		err := sg.Consume(context.Background(), []string{"webook_binlog"},
			saramax.NewHandler[canalx.Message[dao.PublishedArticle]](m.l, m.Consume))
		if err != nil {
			m.l.Error("借助canal更新缓存的消费者退出循环", logger.Error(err))
		}
	}()
	return err
}

func (m *MysqlBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	// 这里因为我直接从repo操作，所以我可以去耦合dao的表
	val canalx.Message[dao.PublishedArticle]) error {
	if val.Table != "published_articles" {
		// 我不关心的数据
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, data := range val.Data {
		var err error
		switch data.Status {
		case domain.ArticleStatusPublished.ToUint8():
			// 发表
			err = m.repo.Cache().SetPub(ctx, m.repo.ToDomain(dao.Article(data)))
		case domain.ArticleStatusPrivate.ToUint8():
			err = m.repo.Cache().DelPub(ctx, data.Id)
		}
		if err != nil {
			// 你可以继续也可以中断
			return err
		}
	}
	return nil
}
