package queries

import (
	"context"
	"time"
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
