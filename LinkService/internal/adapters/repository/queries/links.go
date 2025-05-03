package queries

import (
	"context"
	"errors"

	"github.com/jackc/pgx"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

const saveLinkQuery = `
		INSERT INTO links (key, original_url)
		VALUES ($1, $2)
	`

func (q *Queries) SaveLink(ctx context.Context, key string, originalURL string) error {
	if _, err := q.pool.Exec(ctx, saveLinkQuery, key, originalURL); err != nil {
		return err
	}
	return nil
}

const getOriginalUrlByKeyQuery = `SELECT original_url
        FROM links
        WHERE key = $1`

func (q *Queries) GetOriginalURL(ctx context.Context, key string) (string, error) {
	row := q.pool.QueryRow(ctx, getOriginalUrlByKeyQuery, key)

	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", err
	}
	return originalURL, nil
}

const getKeyByOriginalUrlQuery = `SELECT key
        FROM links
        WHERE original_url = $1`

func (q *Queries) GetKey(ctx context.Context, originalURL string) (string, error) {
	row := q.pool.QueryRow(ctx, getKeyByOriginalUrlQuery, originalURL)

	var key string
	if err := row.Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", err
	}
	return key, nil
}
