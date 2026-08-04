package cache

import (
	"sync"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
)

type DataStore struct {
	Mu                   sync.RWMutex
	MappingData          MappingData
	Players              map[string]responsemodels.ServerPlayer
	CombatLeaderboard    map[string]responsemodels.CombatLeaderboardPlayer
	CraftingLeaderboard  map[string]responsemodels.CraftingLeaderboardPlayer
	AdventureLeaderboard map[string]responsemodels.AdventureLeaderboardPlayer
	StatLeaderboards     StatLeaderboards
}

type MappingData struct {
	Objects             []Object
	Actions             map[string]int32
	BannedPlacedObjects []int32
	BannedBrokenObjects []int32
	CauseMapping        []repo.PsmpstatsCause
	MobMapping          []repo.PsmpstatsMob
	AdvancementMapping  []repo.PsmpstatsAdvancement
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

// An object in the case
type Object struct {
	ID   int32
	Name string
}
