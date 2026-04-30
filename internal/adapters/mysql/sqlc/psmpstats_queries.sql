-- name: FindPsmpstatsPlayerByUuid :one
SELECT
  *
FROM
  psmpstats_players
WHERE
  uuid = ?;

-- name: ListAdvancements :many
SELECT
  *
FROM
  psmpstats_advancements;

-- name: ListCausesOfDeath :many
SELECT
  *
FROM
  psmpstats_causes;

-- name: ListMobs :many
SELECT
  *
FROM
  psmpstats_mobs;

-- name: ListPlayerAdvancementsById :many
SELECT
  *
FROM
  psmpstats_player_advancements
WHERE
  player_id = ?;

-- name: ListDeathsById :many
SELECT
  *
FROM
  psmpstats_deaths
WHERE
  player_id = ?;

-- name: ListPvpKillsById :many
SELECT
 *
FROM
  psmpstats_combat
WHERE
  player_id = ?;

-- name: ListMobKillsById :many
SELECT
  *
FROM
  psmpstats_mob_kills
WHERE
  player_id = ?;

-- name: CountDeathsById :one
SELECT
  count(*)
FROM
  psmpstats_deaths
WHERE
  player_id = ?;

-- name: CountPvpKillsById :one
SELECT
  count(*)
FROM
  psmpstats_combat
WHERE
  player_id = ?;

-- name: CountMobsKilledById :one
SELECT
  count(*)
FROM
  psmpstats_mob_kills
WHERE
  player_id = ?;

-- name: CountSpecificMobKillsById :one
SELECT
  count(*)
FROM
  psmpstats_mob_kills
WHERE
  player_id = ?
AND
  mob = (
    SELECT
      id
    FROM
      psmpstats_mobs
    WHERE
      name = ?
  );

-- name: ListAdvancementsById :many
SELECT
  *
FROM
  psmpstats_advancements
LEFT JOIN
  psmpstats_player_advancements
ON
  psmpstats_advancements.id = psmpstats_player_advancements.advancement
WHERE
  psmpstats_player_advancements.player_id = ?;

-- name: CountDiamondsMinedById :one
SELECT
  COUNT(*)
FROM
  psmpstats_diamonds_mined
WHERE
  player_id = ?;

-- name: ListDiamondsMinedById :many
SELECT
  *
FROM
  psmpstats_diamonds_mined
WHERE
  player_id = ?;