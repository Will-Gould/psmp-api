package players

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	ListPlayers(ctx context.Context) ([]repo.LuckpermsPlayer, error)
	FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error)
	FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s svc) ListPlayers(ctx context.Context) ([]repo.LuckpermsPlayer, error) {
	return s.repo.ListLuckpermsPlayers(ctx)
}

func (s svc) FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error) {
	return s.repo.FindLuckpermsPlayerByUuid(ctx, uuid)
}

func (s svc) FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error) {
	return s.repo.FindGriefLoggerUserByUuid(ctx, uuid)
}
