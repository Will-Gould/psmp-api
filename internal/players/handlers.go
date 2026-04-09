package players

import (
	"context"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"sort"

	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

type playerHandler struct {
	service             Service
	dataStore           *mapping.DataStore
	combatLeaderboard   map[string]responsemodels.CombatLeaderboardPlayer
	craftingLeaderboard map[string]responsemodels.CraftingLeaderboardPlayer
	storyLeaderboard    map[string]responsemodels.StoryLeaderboardPlayer
	statLeaderboards    StatLeaderboards
}

func NewHandler(service Service, ds *mapping.DataStore) *playerHandler {
	return &playerHandler{
		service:             service,
		dataStore:           ds,
		craftingLeaderboard: make(map[string]responsemodels.CraftingLeaderboardPlayer),
		combatLeaderboard:   make(map[string]responsemodels.CombatLeaderboardPlayer),
		storyLeaderboard:    make(map[string]responsemodels.StoryLeaderboardPlayer),
	}
}

func (ph *playerHandler) Load(ctx context.Context) {
	ph.loadServerLeaderboard(ctx, &ph.dataStore.MappingData)
	ph.loadStatLeaderboards()
	ph.formServerRanks()
	ph.formCombatRanks()
	ph.formCraftingRanks()
	ph.formStoryRanks()
	ph.formStatRanks()
}

func (ph *playerHandler) GetServerPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	ph.dataStore.Mu.RLock()
	sp, ok := ph.dataStore.Players[uuid]
	if ok {
		json.Write(w, http.StatusOK, sp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
	ph.dataStore.Mu.RUnlock()
}

func (ph *playerHandler) GetChampionPlayer(w http.ResponseWriter, r *http.Request) {
	ph.dataStore.Mu.RLock()
	for _, p := range ph.dataStore.Players {
		if p.ServerRank == 1 {
			json.Write(w, http.StatusOK, p)
			return
		}
	}
	ph.dataStore.Mu.RLock()
	json.Write(w, http.StatusInternalServerError, nil)
}

func (ph *playerHandler) GetMostDangerousPlayer(w http.ResponseWriter, r *http.Request) {
	ph.dataStore.Mu.RLock()
	kdRanking := slices.Collect(maps.Values(ph.dataStore.Players))
	ph.dataStore.Mu.RUnlock()
	sort.Slice(kdRanking, func(i, j int) bool {
		if kdRanking[i].CombatOverview.PvpKdRatio < kdRanking[j].CombatOverview.PvpKdRatio {
			return true
		}
		return false
	})

	json.Write(w, http.StatusOK, kdRanking[0])
}

func (ph *playerHandler) GetBiggestBuilder(w http.ResponseWriter, r *http.Request) {
	ph.dataStore.Mu.RLock()
	buildRanking := slices.Collect(maps.Values(ph.dataStore.Players))
	ph.dataStore.Mu.RUnlock()
	sort.Slice(buildRanking, func(i, j int) bool {
		if buildRanking[i].CraftingOverview.BlocksPlaced < buildRanking[j].CraftingOverview.BlocksPlaced {
			return true
		}
		return false
	})

	json.Write(w, http.StatusOK, buildRanking[0])
}

func (ph *playerHandler) loadServerLeaderboard(ctx context.Context, md *mapping.MappingData) {
	luckpermsPlayers, err := ph.service.ListLuckpermsPlayers(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to retrieve Luckperms players")
	}

	// build overviews
	for _, p := range luckpermsPlayers {
		glUser, err := ph.service.FindGriefLoggerUser(ctx, p.Uuid)
		if err != nil {
			continue
		}
		player := &responsemodels.ServerPlayer{
			Uuid:         p.Uuid,
			Name:         glUser.Name,
			GlId:         glUser.ID,
			PrimaryGroup: p.PrimaryGroup,
		}
		ph.getPlayerData(ctx, player, md)

		ph.dataStore.Players[player.Uuid] = *player
	}
}

func (ph *playerHandler) loadStatLeaderboards() {
	sl := StatLeaderboards{
		BlocksPlacedLeaderboard:  map[string]responsemodels.LeaderboardPlayer{},
		BlocksBrokenLeaderboard:  map[string]responsemodels.LeaderboardPlayer{},
		DiamondsMinedLeaderboard: map[string]responsemodels.LeaderboardPlayer{},
		TimePlayedLeaderboard:    map[string]responsemodels.LeaderboardPlayer{},
		PvpKillsLeaderboard:      map[string]responsemodels.LeaderboardPlayer{},
		DeathsLeaderboard:        map[string]responsemodels.LeaderboardPlayer{},
		MobKillsLeaderboard:      map[string]responsemodels.LeaderboardPlayer{},
		PvpKdRatioLeaderboard:    map[string]responsemodels.LeaderboardPlayer{},
	}
	ph.statLeaderboards = sl
}

func (ph *playerHandler) getPlayerData(ctx context.Context, player *responsemodels.ServerPlayer, md *mapping.MappingData) {

	// get combat overview
	combatOverview, err := ph.getCombatOverview(ctx, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	// get crafting overview
	craftingOverview, err := ph.getCraftingOverview(ctx, player.Uuid, player.GlId, md)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	// get story overview
	storyOverview, err := ph.getStoryOverview(ctx, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	player.CombatOverview = combatOverview
	player.CraftingOverview = craftingOverview
	player.StoryOverview = storyOverview

	// calculate total score
	player.Score = player.CombatOverview.CombatScore + player.CraftingOverview.CraftingScore + player.StoryOverview.StoryScore
}

func (ph *playerHandler) formServerRanks() {
	// get player slice and sort by score
	l := slices.Collect(maps.Values(ph.dataStore.Players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].Score < l[j].Score {
			return true
		}
		return false
	})
	for i, sp := range l {

		lp := ph.dataStore.Players[sp.Uuid]
		lp.ServerRank = int64(i) + 1
		ph.dataStore.Players[sp.Uuid] = lp
	}
}

func (ph *playerHandler) formCombatRanks() {
	// copy player to slice and sort by combat score
	l := slices.Collect(maps.Values(ph.dataStore.Players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].CombatOverview.CombatScore < l[j].CombatOverview.CombatScore {
			return true
		}
		return false
	})

	for i, sp := range l {
		lp := responsemodels.CombatLeaderboardPlayer{
			Uuid:         sp.Uuid,
			Name:         sp.Name,
			PrimaryGroup: sp.PrimaryGroup,
			Rank:         int64(i) + 1,
			PvpKdRatio:   sp.CombatOverview.PvpKdRatio,
			PvpKills:     sp.CombatOverview.PvpKills,
			MobKills:     sp.CombatOverview.MobKills,
			Deaths:       sp.CombatOverview.Deaths,
		}
		ph.combatLeaderboard[sp.Uuid] = lp
	}
}

func (ph *playerHandler) formCraftingRanks() {
	// copy player to slice and sort by crafting score
	l := slices.Collect(maps.Values(ph.dataStore.Players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].CraftingOverview.CraftingScore < l[j].CraftingOverview.CraftingScore {
			return true
		}
		return false
	})

	for i, sp := range l {
		lp := responsemodels.CraftingLeaderboardPlayer{
			Uuid:          sp.Uuid,
			Name:          sp.Name,
			PrimaryGroup:  sp.PrimaryGroup,
			Rank:          int64(i) + 1,
			BlocksPlaced:  sp.CraftingOverview.BlocksPlaced,
			BlocksBroken:  sp.CraftingOverview.BlocksBroken,
			DiamondsMined: sp.CraftingOverview.DiamondsMined,
		}
		ph.craftingLeaderboard[sp.Uuid] = lp
	}
}

func (ph *playerHandler) formStoryRanks() {
	// copy player to slice and sort by story score
	l := slices.Collect(maps.Values(ph.dataStore.Players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].StoryOverview.StoryScore < l[j].StoryOverview.StoryScore {
			return true
		}
		return false
	})

	for i, sp := range l {
		lp := responsemodels.StoryLeaderboardPlayer{
			Uuid:         sp.Uuid,
			Name:         sp.Name,
			PrimaryGroup: sp.PrimaryGroup,
			Rank:         int64(i) + 1,
		}
		ph.storyLeaderboard[sp.Uuid] = lp
	}
}

func (ph *playerHandler) formStatRanks() {
	playersSlice := slices.Collect(maps.Values(ph.dataStore.Players))
	tempLeaderboard := map[string]responsemodels.LeaderboardPlayer{}
	for _, p := range ph.dataStore.Players {
		tempLeaderboard[p.Uuid] = responsemodels.LeaderboardPlayer{
			Uuid: p.Uuid,
			Rank: 0,
		}
	}

	// order for blocks placed ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CraftingOverview.BlocksPlaced < playersSlice[j].CraftingOverview.BlocksPlaced {
			return true
		}
		return false
	})

	// transfer ranks to temp leaderboard and clone to leaderboard
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.BlocksPlacedLeaderboard = maps.Clone(tempLeaderboard)

	// order for blocks broken ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CraftingOverview.BlocksBroken < playersSlice[j].CraftingOverview.BlocksBroken {
			return true
		}
		return false
	})
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.BlocksBrokenLeaderboard = maps.Clone(tempLeaderboard)

	// order for diamonds mined ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CraftingOverview.DiamondsMined < playersSlice[j].CraftingOverview.DiamondsMined {
			return true
		}
		return false
	})
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.DiamondsMinedLeaderboard = maps.Clone(tempLeaderboard)

	// order for time played ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CraftingOverview.TimePlayed < playersSlice[j].CraftingOverview.TimePlayed {
			return true
		}
		return false
	})
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.TimePlayedLeaderboard = maps.Clone(tempLeaderboard)

	// order for pvp kills ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CombatOverview.PvpKills < playersSlice[j].CombatOverview.PvpKills {
			return true
		}
		return false
	})
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.PvpKillsLeaderboard = maps.Clone(tempLeaderboard)

	// order for deaths ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CombatOverview.Deaths < playersSlice[j].CombatOverview.Deaths {
			return true
		}
		return false
	})
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.DeathsLeaderboard = maps.Clone(tempLeaderboard)

	// order for mob kills ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CombatOverview.MobKills < playersSlice[j].CombatOverview.MobKills {
			return true
		}
		return false
	})
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.MobKillsLeaderboard = maps.Clone(tempLeaderboard)

	// order for pvp kd ratio ranks
	sort.Slice(playersSlice, func(i, j int) bool {
		if playersSlice[i].CombatOverview.PvpKdRatio < playersSlice[j].CombatOverview.PvpKdRatio {
			return true
		}
		return false
	})
	transferRanks(playersSlice, tempLeaderboard)
	ph.statLeaderboards.PvpKdRatioLeaderboard = maps.Clone(tempLeaderboard)
}

func transferRanks(players []responsemodels.ServerPlayer, l map[string]responsemodels.LeaderboardPlayer) {
	for i, p := range players {
		lp := l[p.Uuid]
		lp.Rank = int64(i) + 1
		l[p.Uuid] = lp
	}
}
