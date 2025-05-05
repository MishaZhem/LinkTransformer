package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type ClickEvent struct {
	LinkKey   string    `json:"link_key"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Time      time.Time `json:"time"`
}

type KafkaProducer struct {
	writer *kafka.Writer
}

type Producer interface {
	SendClickEvent(ctx context.Context, key, ipAddress, userAgent string) error
	Close() error
}

func NewProducer(broker, topic string) Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) SendClickEvent(ctx context.Context, key, ipAddress, userAgent string) error {
	event := ClickEvent{
		LinkKey:   key,
		IP:        ipAddress,
		UserAgent: userAgent,
		Time:      time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{Value: data})
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
