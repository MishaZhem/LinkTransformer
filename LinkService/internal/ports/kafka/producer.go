package kafka

import (
	"context"
	"encoding/json"
	"log"
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
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
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

	log.Printf("Sending to Kafka: %s", string(data))
	err = p.writer.WriteMessages(ctx, kafka.Message{Value: data})
	log.Printf("Sent to Kafka, err: %v", err)
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
