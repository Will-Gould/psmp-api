-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_mobs(
    id SERIAL PRIMARY KEY,
    name text NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_mobs;
