package app

import (
	"AnalyticsService/internal/adapters/repository"
	"AnalyticsService/internal/ports/kafka"
	"context"
	"encoding/json"
	"errors"
	"log"
)

type Program struct {
	repository repository.Repository
	consumer   kafka.Consumer
}

type App interface {
	GetStatistics(ctx context.Context, url string) ([]kafka.ClickEvent, error)
	GetTotalClicks(ctx context.Context, url string) (int64, error)
	RunConsumer(ctx context.Context) error
}

var ErrBadRequest = errors.New("bad request")
var ErrForbidden = errors.New("forbidden")

func NewApp(repository repository.Repository, consumer kafka.Consumer) App {
	return &Program{
		repository: repository,
		consumer:   consumer,
	}
}

func (r *Program) GetStatistics(ctx context.Context, url string) ([]kafka.ClickEvent, error) {
	events, err := r.repository.GetRows(ctx, url)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *Program) GetTotalClicks(ctx context.Context, url string) (int64, error) {
	clicks, err := r.repository.GetClicks(ctx, url)
	if err != nil {
		return 0, err
	}

	return clicks, nil
}

func (r *Program) RunConsumer(ctx context.Context) error {
	defer r.consumer.Close()
	log.Println("Kafka consumer loop starting...")
	for {
		m, err := r.consumer.ReadMessage(ctx)
		if err != nil {
			log.Println("Kafka consumer stopped by context")
			return err
		}

		var evt kafka.ClickEvent
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			log.Printf("bad event: %v", err)
			continue
		}

		if err := r.repository.SaveClickEvent(ctx, evt.LinkKey, evt.IP, evt.UserAgent, evt.Time); err != nil {
			log.Printf("save error: %v", err)
		}
	}
}
