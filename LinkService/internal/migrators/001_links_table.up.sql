-- migrate -path internal/migrators -database "postgres://postgres:123@localhost:5432/postgres?sslmode=disable" up
CREATE TABLE IF NOT EXISTS links (
    id serial not null,
    key VARCHAR(25) not null,
    original_url VARCHAR(25) not null
);