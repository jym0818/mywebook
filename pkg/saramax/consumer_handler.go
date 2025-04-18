package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/jym/mywebook/pkg/logger"
)

type Handler[T any] struct {
	l  logger.Logger
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func NewHandler[T any](l logger.Logger, fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			//记录日志
			h.l.Error("反序列化消息失败", logger.String("msg", string(msg.Value)), logger.String("topic", msg.Topic))
			continue
		}
		//在这里执行重试 或者记录日志
		for i := 0; i < 3; i++ {
			err = h.fn(msg, t)
			if err == nil {
				break
			}
			//重试日志
			h.l.Error("处理消息失败", logger.String("msg", string(msg.Value)), logger.String("topic", msg.Topic))

		}
		if err != nil {
			//重试都失败了  继续打印日志
			h.l.Error("处理消息失败并且重试失败", logger.String("msg", string(msg.Value)), logger.String("topic", msg.Topic))
		} else {
			session.MarkMessage(msg, "")
		}

	}

	return nil
}
