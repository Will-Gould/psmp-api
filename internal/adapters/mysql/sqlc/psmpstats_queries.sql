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

-- name: ListAdvancementsByUuid :many
SELECT
  *
FROM
  psmpstats_player_advancements
WHERE
  player_uuid = ?;

-- name: ListDeathsByUuid :many
SELECT
  *
FROM
  psmpstats_deaths
WHERE
  player_uuid = ?;

-- name: ListPvpKillsByUuid :many
SELECT
 *
FROM
  psmpstats_combat
WHERE
  player_uuid = ?;

-- name: ListMobKillsByUuid :many
SELECT
  *
FROM
  psmpstats_mob_kills
WHERE
  player_uuid = ?;

-- name: CountDeathsByUuid :one
SELECT
  count(*)
FROM
  psmpstats_deaths
WHERE
  player_uuid = ?;

-- name: CountPvpKillsByUuid :one
SELECT
  count(*)
FROM
  psmpstats_combat
WHERE
  player_uuid = ?;

-- name: CountMobsKilledByUUid :one
SELECT
  count(*)
FROM
  psmpstats_mob_kills
WHERE
  player_uuid = ?;

-- name: CountSpecificMobKillsByUuid :one
SELECT
  count(*)
FROM
  psmpstats_mob_kills
WHERE
  player_uuid = ?
AND
  mob = ?;