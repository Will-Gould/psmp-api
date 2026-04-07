package profiles

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error)
	ListBlocksBrokenByUser(ctx context.Context, id int32) ([]repo.Block, error)
	ListBlocksPlacedByUser(ctx context.Context, id int32) ([]repo.Block, error)
	ListSessionDataByUser(ctx context.Context, id int32) ([]repo.Session, error)
	ListDeathsByPlayer(ctx context.Context, uuid string) ([]repo.PsmpstatsDeath, error)
	CountSpecificMobKillsByUuid(ctx context.Context, name string, uuid string) (int64, error)
	FindPsmpstatsPlayerByUuid(ctx context.Context, uuid string) (repo.PsmpstatsPlayer, error)
	FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error)
	ListAdvancementsByPlayer(ctx context.Context, uuid string) ([]repo.PsmpstatsAdvancement, error)
	ListMobKillsByUuid(ctx context.Context, uuid string) ([]repo.PsmpstatsMobKill, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

// CountSpecificMobKillsByUuid implements [Service].
func (s *svc) CountSpecificMobKillsByUuid(ctx context.Context, name string, uuid string) (int64, error) {
	return s.repo.CountSpecificMobKillsByUuid(ctx, repo.CountSpecificMobKillsByUuidParams{
		PlayerUuid: uuid,
		Name:       name,
	})
}

// FindGriefLoggerUser implements [Service].
func (s *svc) FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error) {
	return s.repo.FindGriefLoggerUserByUuid(ctx, uuid)
}

// FindLuckpermsPlayer implements [Service].
func (s *svc) FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error) {
	return s.repo.FindLuckpermsPlayerByUuid(ctx, uuid)
}

// FindPsmpstatsPlayerByUuid implements [Service].
func (s *svc) FindPsmpstatsPlayerByUuid(ctx context.Context, uuid string) (repo.PsmpstatsPlayer, error) {
	return s.repo.FindPsmpstatsPlayerByUuid(ctx, uuid)
}

// ListAdvancementsByPlayer implements [Service].
func (s *svc) ListAdvancementsByPlayer(ctx context.Context, uuid string) ([]repo.PsmpstatsAdvancement, error) {
	rows, err := s.repo.ListAdvancementsByUuid(ctx, uuid)
	advancements := []repo.PsmpstatsAdvancement{}
	if err != nil {
		return advancements, err
	}
	// map result to PsmpstatsAdvancement struct
	for _, a := range rows {
		advancements = append(advancements, repo.PsmpstatsAdvancement{
			ID:   a.ID,
			Name: a.Name,
		})
	}
	return advancements, nil
}

// ListBlocksBrokenByUser implements [Service].
func (s *svc) ListBlocksBrokenByUser(ctx context.Context, id int32) ([]repo.Block, error) {
	return s.repo.ListBlocksBrokenByUser(ctx, id)
}

// ListBlocksPlacedByUser implements [Service].
func (s *svc) ListBlocksPlacedByUser(ctx context.Context, id int32) ([]repo.Block, error) {
	return s.repo.ListBlocksPlacedByUser(ctx, id)
}

// ListDeathsByPlayer implements [Service].
func (s *svc) ListDeathsByPlayer(ctx context.Context, uuid string) ([]repo.PsmpstatsDeath, error) {
	return s.repo.ListDeathsByUuid(ctx, uuid)
}

// ListMobKillsByUuid implements [Service].
func (s *svc) ListMobKillsByUuid(ctx context.Context, uuid string) ([]repo.PsmpstatsMobKill, error) {
	return s.repo.ListMobKillsByUuid(ctx, uuid)
}

// ListSessionDataByUser implements [Service].
func (s *svc) ListSessionDataByUser(ctx context.Context, id int32) ([]repo.Session, error) {
	return s.repo.ListSessionDataByUser(ctx, id)
}
