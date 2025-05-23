package ioc

import (
	"github.com/IBM/sarama"
	"github.com/jym/mywebook/interactive/events"
	"github.com/jym/mywebook/pkg/saramax"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewConsumers(c1 *events.KafkaConsumer) []saramax.Consumer {
	return []saramax.Consumer{c1}
}
