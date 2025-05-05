-- migrate -path internal/migrators -database "postgres://postgres:123@localhost:5433/postgres?sslmode=disable" up
CREATE TABLE IF NOT EXISTS link_clicks (
    id SERIAL PRIMARY KEY,
    link_key VARCHAR(255) NOT NULL,
    clicked_at TIMESTAMP DEFAULT NOW(),
    user_agent TEXT,
    ip_address TEXT
);