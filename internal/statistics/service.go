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
	ListSessionDataByUser(ctx context.Context, id int32) ([]repo.Session, error)
	ListDeathsByPlayer(ctx context.Context, uuid string) ([]repo.PsmpstatsDeath, error)
	CountDeathsByPlayer(ctx context.Context, uuid string) (int64, error)
	CountPvpKillsByPlayer(ctx context.Context, uuid string) (int64, error)
	CountMobKillsByPlayer(ctx context.Context, uuid string) (int64, error)
	ListMobs(ctx context.Context) ([]repo.PsmpstatsMob, error)
	ListCausesOfDeath(ctx context.Context) ([]repo.PsmpstatsCause, error)
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

func (s svc) ListSessionDataByUser(ctx context.Context, id int32) ([]repo.Session, error) {
	return s.repo.ListSessionDataByUser(ctx, id)
}

func (s svc) ListDeathsByPlayer(ctx context.Context, uuid string) ([]repo.PsmpstatsDeath, error) {
	return s.repo.ListDeathsByUuid(ctx, uuid)
}

func (s svc) CountDeathsByPlayer(ctx context.Context, uuid string) (int64, error) {
	return s.repo.CountDeathsByUuid(ctx, uuid)
}

func (s svc) CountPvpKillsByPlayer(ctx context.Context, uuid string) (int64, error) {
	return s.repo.CountPvpKillsByUuid(ctx, uuid)
}

func (s svc) CountMobKillsByPlayer(ctx context.Context, uuid string) (int64, error) {
	return s.repo.CountMobsKilledByUuid(ctx, uuid)
}

func (s svc) ListMobs(ctx context.Context) ([]repo.PsmpstatsMob, error) {
	return s.repo.ListMobs(ctx)
}

func (s svc) ListCausesOfDeath(ctx context.Context) ([]repo.PsmpstatsCause, error) {
	return s.repo.ListCausesOfDeath(ctx)
}
