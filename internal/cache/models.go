package cache

import (
	"sync"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
)

type DataStore struct {
	Mu                  sync.RWMutex
	MappingData         MappingData
	Players             map[string]responsemodels.ServerPlayer
	CombatLeaderboard   map[string]responsemodels.CombatLeaderboardPlayer
	CraftingLeaderboard map[string]responsemodels.CraftingLeaderboardPlayer
	StoryLeaderboard    map[string]responsemodels.StoryLeaderboardPlayer
	StatLeaderboards    StatLeaderboards
}

type MappingData struct {
	Materials             []repo.Material
	BannedPlacedMaterials []int32
	BannedBrokenMaterials []int32
	CauseMapping          []repo.PsmpstatsCause
	MobMapping            []repo.PsmpstatsMob
	AdvancementMapping    []repo.PsmpstatsAdvancement
}

type StatLeaderboards struct {
	BlocksPlacedLeaderboard  map[string]responsemodels.LeaderboardPlayer
	BlocksBrokenLeaderboard  map[string]responsemodels.LeaderboardPlayer
	DiamondsMinedLeaderboard map[string]responsemodels.LeaderboardPlayer
	TimePlayedLeaderboard    map[string]responsemodels.LeaderboardPlayer
	PvpKillsLeaderboard      map[string]responsemodels.LeaderboardPlayer
	DeathsLeaderboard        map[string]responsemodels.LeaderboardPlayer
	MobKillsLeaderboard      map[string]responsemodels.LeaderboardPlayer
	PvpKdRatioLeaderboard    map[string]responsemodels.LeaderboardPlayer
}
