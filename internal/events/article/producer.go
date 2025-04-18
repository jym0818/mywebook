package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	//aricleService会产生很多event事件ProduceReadEvent
	ProduceReadEvent(ctx context.Context, evt ReadEvent) error
}

type ReadEvent struct {
	Uid int64
	Aid int64
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: producer,
	}
}
func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, evt ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_articles",
		Value: sarama.ByteEncoder(data),
	})
	return err
}
