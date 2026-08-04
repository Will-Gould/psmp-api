-- +goose Up
CREATE TABLE IF NOT EXISTS actions(
    id int(11) PRIMARY KEY NOT NULL,
    action_id int(11) KEY UNIQUE NOT NULL,
    time timestamp(6) NOT NULL,
    x int(11) NOT NULL,
    y int(11) NOT NULL,
    z int(11) NOT NULL,
    world_id int(11) NOT NULL,
    object_id int(11) NOT NULL,
    old_object_id int(11) NOT NULL,
    block_state text,
    old_block_state text,
    source int(11) NOT NULL,
    player_id int(11),
    extra_data text,
    rolled_back tinyint(1) NOT NULL
);