package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type ClickEvent struct {
	LinkKey   string    `json:"link_key"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Time      time.Time `json:"time"`
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

type Consumer interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

func NewConsumer(brokers []string, topic, groupID string) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &KafkaConsumer{reader: reader}
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
