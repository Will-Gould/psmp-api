package players

import (
	"maps"
	"net/http"
	"slices"
	"sort"

	"github.com/Will-Gould/psmp-api/internal/json"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

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

func (ph *playerHandler) ListServerLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.players))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].ServerRank < lb[j].ServerRank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) ListCombatLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.combatLeaderboard))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) ListCraftingLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.craftingLeaderboard))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) ListStoryLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.storyLeaderboard))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) GetCombatLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	lp, ok := ph.combatLeaderboard[uuid]
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetCraftingLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	lp, ok := ph.craftingLeaderboard[uuid]
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetStoryLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	lp, ok := ph.storyLeaderboard[uuid]
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) ListStatLeaderboard(w http.ResponseWriter, r *http.Request) {
	stat := chi.URLParam(r, "stat")
	switch stat {
	case "blocks-placed":
		json.Write(w, http.StatusOK, ph.statLeaderboards.BlocksPlacedLeaderboard)
	case "blocks-broken":
		json.Write(w, http.StatusOK, ph.statLeaderboards.BlocksBrokenLeaderboard)
	case "diamonds-mined":
		json.Write(w, http.StatusOK, ph.statLeaderboards.DiamondsMinedLeaderboard)
	case "time-played":
		json.Write(w, http.StatusOK, ph.statLeaderboards.TimePlayedLeaderboard)
	case "pvp-kills":
		json.Write(w, http.StatusOK, ph.statLeaderboards.PvpKillsLeaderboard)
	case "deaths":
		json.Write(w, http.StatusOK, ph.statLeaderboards.DeathsLeaderboard)
	case "mob-kills":
		json.Write(w, http.StatusOK, ph.statLeaderboards.MobKillsLeaderboard)
	case "pvp-kd-ratio":
		json.Write(w, http.StatusOK, ph.statLeaderboards.PvpKdRatioLeaderboard)
	default:
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetStatRank(w http.ResponseWriter, r *http.Request) {
	stat := chi.URLParam(r, "stat")
	uuid := chi.URLParam(r, "uuid")

	// check if player exists
	_, ok := ph.players[uuid]
	if !ok {
		json.Write(w, http.StatusNotFound, nil)
	}

	switch stat {
	case "blocks-placed":
		json.Write(w, http.StatusOK, ph.statLeaderboards.BlocksPlacedLeaderboard[uuid])
	case "blocks-broken":
		json.Write(w, http.StatusOK, ph.statLeaderboards.BlocksBrokenLeaderboard[uuid])
	case "diamonds-mined":
		json.Write(w, http.StatusOK, ph.statLeaderboards.DiamondsMinedLeaderboard[uuid])
	case "time-played":
		json.Write(w, http.StatusOK, ph.statLeaderboards.TimePlayedLeaderboard[uuid])
	case "pvp-kills":
		json.Write(w, http.StatusOK, ph.statLeaderboards.PvpKillsLeaderboard[uuid])
	case "deaths":
		json.Write(w, http.StatusOK, ph.statLeaderboards.DeathsLeaderboard[uuid])
	case "mob-kills":
		json.Write(w, http.StatusOK, ph.statLeaderboards.MobKillsLeaderboard[uuid])
	case "pvp-kd-ratio":
		json.Write(w, http.StatusOK, ph.statLeaderboards.PvpKdRatioLeaderboard[uuid])
	default:
		json.Write(w, http.StatusNotFound, nil)
	}
}
