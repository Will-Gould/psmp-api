package statistics

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	ListMaterials(ctx context.Context) ([]repo.Material, error)
	FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error)
	CountBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) (int64, error)
	ListBlocksBrokenByUser(ctx context.Context, id int32) ([]repo.Block, error)
	ListBlocksPlacedByUser(ctx context.Context, id int32) ([]repo.Block, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s svc) ListMaterials(ctx context.Context) ([]repo.Material, error) {
	return s.repo.ListMaterials(ctx)
}

func (s svc) FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error) {
	return s.repo.FindGriefLoggerUserByUuid(ctx, uuid)
}

func (s svc) ListBlocksBrokenByUser(ctx context.Context, id int32) ([]repo.Block, error) {
	return s.repo.ListBlocksBrokenByUser(ctx, id)
}

func (s svc) ListBlocksPlacedByUser(ctx context.Context, id int32) ([]repo.Block, error) {
	return s.repo.ListBlocksPlacedByUser(ctx, id)
}

func (s svc) CountBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) (int64, error) {
	return s.repo.CountBlocksByUser(ctx, id, action, banned)
}
