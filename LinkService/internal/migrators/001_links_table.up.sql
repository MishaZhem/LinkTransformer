-- migrate -path internal/migrators -database "postgres://postgres:123@localhost:5433/postgres?sslmode=disable" up
CREATE TABLE IF NOT EXISTS links (
    id serial not null,
    key TEXT not null,
    original_url TEXT not null
);