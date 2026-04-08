-- +goose Up
CREATE TABLE IF NOT EXISTS leaderboard_rankings(
    uuid varchar(36) PRIMARY KEY,
    previous_server_rank INT,
    previous_combat_rank INT,
    previous_crafting_rank INT,
    previous_pvp_kill_rank INT,
    previous_blocks_placed_rank INT,
    previous_block_broken_rank INT,
    previous_mobs_killed_rank INT,
    prevous_diamonds_mined_rank INT
);

-- +goose Down
SELECT 'down SQL query';
