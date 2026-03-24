-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_causes(
    id SERIAL PRIMARY KEY,
    name text NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_causes;
