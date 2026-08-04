package repo

import (
	"context"
	"strings"
)

const GriefLoggerCountBlocksByUser = `-- name: GriefLoggerCountBlocksByUser :one
SELECT
  count(*)
FROM
  blocks
WHERE
  user = ?
AND
  action = ?
`

func (q *Queries) GriefLoggerCountBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) (int64, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(GriefLoggerCountBlocksByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}

	row := q.db.QueryRowContext(ctx, query.String(), args...)

	var count int64
	err := row.Scan(&count)
	return count, err
}

const LedgerCountBlocksByUser = `-- name: LedgerCountBlocksByUser :one
SELECT
  count(*)
FROM
  actions
WHERE
  player_id = ?
AND
  action_id = ?
AND
  rolled_back = 0
`

func (q *Queries) LedgerCountBlocksByUser(ctx context.Context, actions map[string]int32, user int32, action int32, banned []int32) (int64, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(LedgerCountBlocksByUser)
	for range banned {
		// Ledger store block broken in `old_object_id` and block placed in `object_id`
		switch action {
		case actions["block-break"]:
			query.WriteString(" AND old_object_id != ?")
		default:
			query.WriteString(" AND object_id != ?")
		}
	}

	row := q.db.QueryRowContext(ctx, query.String(), args...)

	var count int64
	err := row.Scan(&count)
	return count, err
}

const GriefLoggerGroupCountBlocksByUser = `-- name: GriefLoggerGroupCountBlocksByUser :many
SELECT
  type, COUNT(*) AS total_placed
FROM
  blocks
WHERE
  user = ?
AND
  action = ?
`

func (q *Queries) GriefLoggerGroupCountBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) ([]GroupCountBlocksPlacedByUserRow, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(GriefLoggerGroupCountBlocksByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}
	query.WriteString(" GROUP BY type")

	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupCountBlocksPlacedByUserRow
	for rows.Next() {
		var i GroupCountBlocksPlacedByUserRow
		if err := rows.Scan(&i.Type, &i.TotalPlaced); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const LedgerGroupCountBlocksByUser = `-- name: LedgerGroupCountBlocksByUser :many
SELECT
  object_id, COUNT(*) AS total_placed
FROM
  actions
WHERE
  player_id = ?
AND
  action_id = ?
`

const LedgerGroupCountBlocksBrokenByUser = `-- name: LedgerGroupCountBlocksByUser :many
SELECT
  old_object_id, COUNT(*) AS total_placed
FROM
  actions
WHERE
  player_id = ?
AND
  action_id = ?
`

func (q *Queries) LedgerGroupCountBlocksByUser(ctx context.Context, actions map[string]int32, user int32, action int32, banned []int32) ([]GroupCountBlocksPlacedByUserRow, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	switch action {
	case actions["block-break"]:
		query.WriteString(LedgerGroupCountBlocksBrokenByUser)
		for range banned {
			query.WriteString(" AND old_object_id != ?")
		}
		query.WriteString(" GROUP BY old_object_id")
	default:
		query.WriteString(LedgerGroupCountBlocksByUser)
		for range banned {
			query.WriteString(" AND object_id != ?")
		}
		query.WriteString(" GROUP BY object_id")
	}

	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupCountBlocksPlacedByUserRow
	for rows.Next() {
		var i GroupCountBlocksPlacedByUserRow
		if err := rows.Scan(&i.Type, &i.TotalPlaced); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GriefLoggerListBlocksByUser = `-- name: GriefLoggerListBlocksByUser :many
SELECT
  time, user, level, x, y, z, type, action
FROM
  blocks
WHERE
  user = ?
AND
  action = ?
`

func (q *Queries) GriefLoggerListBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) ([]Block, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(GriefLoggerListBlocksByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}

	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Block
	for rows.Next() {
		var i Block
		if err := rows.Scan(
			&i.Time,
			&i.User,
			&i.Level,
			&i.X,
			&i.Y,
			&i.Z,
			&i.Type,
			&i.Action,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const LedgerListBlocksByUser = `-- name: ListActions :many
SELECT
  id, action_id, time, x, y, z, world_id, object_id, old_object_id, block_state, old_block_state, source, player_id, extra_data, rolled_back
FROM
  actions
WHERE
  player_id = ?
AND
  action_id = ?
AND
  rolled_back = 0
`

func (q *Queries) LedgerListBlocksByUser(ctx context.Context, actions map[string]int32, user int32, action int32, banned []int32) ([]Action, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(LedgerListBlocksByUser)
	for range banned {
		switch action {
		case actions["block-break"]:
			query.WriteString(" AND old_object_id != ?")
		default:
			query.WriteString(" AND object_id != ?")
		}
	}

	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Action
	for rows.Next() {
		var i Action
		if err := rows.Scan(
			&i.ID,
			&i.ActionID,
			&i.Time,
			&i.X,
			&i.Y,
			&i.Z,
			&i.WorldID,
			&i.ObjectID,
			&i.OldObjectID,
			&i.BlockState,
			&i.OldBlockState,
			&i.Source,
			&i.PlayerID,
			&i.ExtraData,
			&i.RolledBack,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
