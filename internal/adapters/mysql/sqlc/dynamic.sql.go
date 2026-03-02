package repo

import (
	"context"
	"strings"
)

const CountBlocksByUser = `-- name: CountBlocksBrokenByUser :many
SELECT
  count(*)
FROM
  blocks
WHERE
  user = ?
AND
  action = ?
`

func (q *Queries) CountBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) (int64, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(CountBlocksByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}

	row := q.db.QueryRowContext(ctx, query.String(), args...)
	var count int64
	err := row.Scan(&count)
	return count, err
}
