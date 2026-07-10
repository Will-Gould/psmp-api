-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_advancements(
    id SERIAL PRIMARY KEY NOT NULL,
    name text NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_advancements;
