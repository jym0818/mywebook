package events

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jym/mywebook/interactive/repository"
	"github.com/jym/mywebook/pkg/logger"
	"github.com/jym/mywebook/pkg/saramax"
	"time"
)

type BatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.Logger
}

func NewBatchConsumer(client sarama.Client, l logger.Logger, repo repository.InteractiveRepository) *BatchConsumer {
	return &BatchConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (b *BatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive-read", b.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"read_articles"}, saramax.NewBatchHandler[ReadEvent](b.l, b.Consume))
		if err != nil {
			//记录日志
		}
	}()
	return err
}

func (b *BatchConsumer) Consume(msgs []*sarama.ConsumerMessage, evts []ReadEvent) error {
	ids := make([]int64, 0, len(evts))
	bizs := make([]string, 0, len(evts))
	for _, evt := range evts {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := b.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		b.l.Error("批量阅读计数失败")
		//可以不返回err
	}
	return nil
}
