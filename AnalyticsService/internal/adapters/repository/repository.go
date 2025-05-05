package repository

import (
	"LinkTransformer/internal/adapters/repository/queries"
	"context"

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
	SaveLink(ctx context.Context, key string, originalURL string) error
	GetOriginalURL(ctx context.Context, key string) (string, error)
	GetKey(ctx context.Context, originalURL string) (string, error)
}
