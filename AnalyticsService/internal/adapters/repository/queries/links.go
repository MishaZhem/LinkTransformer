package queries

import (
	"AnalyticsService/internal/ports/kafka"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

const saveClickEventQuery = `
    INSERT INTO link_clicks (link_key, clicked_at, user_agent, ip_address)
    VALUES ($1, $2, $3, $4)`

func (q *Queries) SaveClickEvent(ctx context.Context, key, ip, ua string, t time.Time) error {

	if _, err := q.pool.Exec(ctx, saveClickEventQuery, key, t, ua, ip); err != nil {
		return err
	}
	return nil
}

const getClicksByKeyQuery = `SELECT COUNT(*) FROM link_clicks WHERE link_key = $1`

func (q *Queries) GetClicks(ctx context.Context, key string) (int64, error) {
	row := q.pool.QueryRow(ctx, getClicksByKeyQuery, key)

	var res int64
	if err := row.Scan(&res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrLinkNotFound
		}
		return 0, err
	}
	return res, nil
}

const getRowsByKeyQuery = `SELECT COUNT(*) FROM link_clicks WHERE link_key = $1`

func (q *Queries) GetRows(ctx context.Context, key string) ([]kafka.ClickEvent, error) {
	rows, err := q.pool.Query(ctx, getRowsByKeyQuery, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []kafka.ClickEvent
	for rows.Next() {
		var e kafka.ClickEvent
		if err := rows.Scan(&e.LinkKey, &e.IP, &e.UserAgent, &e.Time); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if len(events) == 0 {
		return nil, ErrLinkNotFound
	}
	return events, nil
}
