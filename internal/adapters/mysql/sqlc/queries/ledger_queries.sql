-- name: ListLedgerPlayers :many
SELECT
  *
FROM
  players;

-- name: ListObjects :many
SELECT
  *
FROM
  ObjectIdentifiers;

-- name: ListActionsIdentifiers :many
SELECT
  *
FROM
  ActionIdentifiers;

-- name: FindLedgerPlayerByUuid :one
SELECT
  *
FROM 
  players
WHERE
  player_id = ?;

-- name: ListActions :many
SELECT
  *
FROM
  actions;