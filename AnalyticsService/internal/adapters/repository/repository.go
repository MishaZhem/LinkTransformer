package repository

import (
	"AnalyticsService/internal/adapters/repository/queries"
	"AnalyticsService/internal/ports/kafka"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type repo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger logrus.FieldLogger
}

func NewRepository(pgxPool *pgxpool.Pool, logger logrus.FieldLogger) Repository {
	return &repo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
		logger:  logger,
	}
}

type Repository interface {
	SaveClickEvent(ctx context.Context, key, ip, ua string, t time.Time) error
	GetClicks(ctx context.Context, key string) (int64, error)
	GetRows(ctx context.Context, key string) ([]kafka.ClickEvent, error)
}
