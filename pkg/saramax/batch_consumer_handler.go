package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/jym/mywebook/pkg/logger"
	"time"
)

type BatchHandler[T any] struct {
	l            logger.Logger
	fn           func(msgs []*sarama.ConsumerMessage, ts []T) error
	batchSize    int
	batchTimeout time.Duration
}

func NewBatchHandler[T any](l logger.Logger, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{
		l:            l,
		fn:           fn,
		batchSize:    10,
		batchTimeout: time.Second,
	}
}

func (b BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()

	for {
		var last *sarama.ConsumerMessage
		ctx, cancel := context.WithTimeout(context.Background(), b.batchTimeout)
		done := false
		msgs := make([]*sarama.ConsumerMessage, 0, b.batchSize)
		ts := make([]T, 0, b.batchSize)
		for i := 0; i < b.batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					cancel()
					return nil
				}
				last = msg

				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败")
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)

			}

		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := b.fn(msgs, ts)
		if err != nil {
			b.l.Error("调用业务批量接口失败")
			//批量处理失败也要继续下去，不能停止

		}
		if last != nil {
			session.MarkMessage(last, "")
		}

	}

}
