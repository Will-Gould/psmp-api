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

-- name: ListPlayerAdvancementsByUuid :many
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

-- name: CountMobsKilledByUuid :one
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
  mob = (
    SELECT
      id
    FROM
      psmpstats_mobs
    WHERE
      name = ?
  );

-- name: ListAdvancementsByUuid :many
SELECT
  *
FROM
  psmpstats_advancements
LEFT JOIN
  psmpstats_player_advancements
ON
  psmpstats_advancements.id = psmpstats_player_advancements.advancement
WHERE
  psmpstats_player_advancements.player_uuid = ?