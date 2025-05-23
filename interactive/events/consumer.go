package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jym/mywebook/interactive/repository"
	"github.com/jym/mywebook/pkg/logger"
	"github.com/jym/mywebook/pkg/saramax"
	"time"
)

type KafkaConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.Logger
}

func (k *KafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive-read", k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"read_articles"}, saramax.NewHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			//记录日志
		}
	}()
	return err
}
func (k *KafkaConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return k.repo.IncrReadCnt(ctx, "article", evt.Uid)
}

func NewKafkaConsumer(client sarama.Client, l logger.Logger, repo repository.InteractiveRepository) *KafkaConsumer {
	return &KafkaConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}
