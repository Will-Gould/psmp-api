package players

import (
	"maps"
	"net/http"
	"slices"
	"sort"
	"strconv"

	"github.com/Will-Gould/psmp-api/internal/json"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

func (ph *playerHandler) ListServerLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(chi.URLParam(r, "limit"))
	if err != nil {
		limit = -1
	}

	ph.dataStore.Mu.RLock()
	lb := slices.Collect(maps.Values(ph.dataStore.Players))
	ph.dataStore.Mu.RUnlock()
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].ServerRank < lb[j].ServerRank {
			return true
		}
		return false
	})

	if limit < 0 {
		json.Write(w, http.StatusOK, lb)
		return
	}

	if len(lb) <= limit {
		json.Write(w, http.StatusOK, lb)
	} else {
		json.Write(w, http.StatusOK, lb[:limit])
	}
}

func (ph *playerHandler) ListCombatLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(chi.URLParam(r, "limit"))
	if err != nil {
		limit = -1
	}

	ph.dataStore.Mu.RLock()
	lb := slices.Collect(maps.Values(ph.dataStore.CombatLeaderboard))
	ph.dataStore.Mu.RUnlock()
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})

	if limit < 0 {
		json.Write(w, http.StatusOK, lb)
		return
	}

	if len(lb) <= limit {
		json.Write(w, http.StatusOK, lb)
	} else {
		json.Write(w, http.StatusOK, lb[:limit])
	}
}

func (ph *playerHandler) ListCraftingLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(chi.URLParam(r, "limit"))
	if err != nil {
		limit = -1
	}

	ph.dataStore.Mu.RLock()
	lb := slices.Collect(maps.Values(ph.dataStore.CraftingLeaderboard))
	ph.dataStore.Mu.RUnlock()
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})

	if limit < 0 {
		json.Write(w, http.StatusOK, lb)
		return
	}

	if len(lb) <= limit {
		json.Write(w, http.StatusOK, lb)
	} else {
		json.Write(w, http.StatusOK, lb[:limit])
	}
}

func (ph *playerHandler) ListAdventureLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(chi.URLParam(r, "limit"))
	if err != nil {
		limit = -1
	}

	ph.dataStore.Mu.RLock()
	lb := slices.Collect(maps.Values(ph.dataStore.AdventureLeaderboard))
	ph.dataStore.Mu.RUnlock()
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})

	if limit < 0 {
		json.Write(w, http.StatusOK, lb)
		return
	}

	if len(lb) <= limit {
		json.Write(w, http.StatusOK, lb)
	} else {
		json.Write(w, http.StatusOK, lb[:limit])
	}
}

func (ph *playerHandler) GetCombatLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	ph.dataStore.Mu.RLock()
	lp, ok := ph.dataStore.CombatLeaderboard[uuid]
	ph.dataStore.Mu.RUnlock()
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetCraftingLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	ph.dataStore.Mu.RLock()
	lp, ok := ph.dataStore.CraftingLeaderboard[uuid]
	ph.dataStore.Mu.RUnlock()
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetAdventureLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	ph.dataStore.Mu.RLock()
	lp, ok := ph.dataStore.AdventureLeaderboard[uuid]
	ph.dataStore.Mu.RUnlock()
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) ListStatLeaderboard(w http.ResponseWriter, r *http.Request) {
	stat := chi.URLParam(r, "stat")
	ph.dataStore.Mu.RLock()
	switch stat {
	case "blocks-placed":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.BlocksPlacedLeaderboard)
	case "blocks-broken":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.BlocksBrokenLeaderboard)
	case "diamonds-mined":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.DiamondsMinedLeaderboard)
	case "time-played":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.TimePlayedLeaderboard)
	case "pvp-kills":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.PvpKillsLeaderboard)
	case "deaths":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.DeathsLeaderboard)
	case "mob-kills":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.MobKillsLeaderboard)
	case "pvp-kd-ratio":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.PvpKdRatioLeaderboard)
	default:
		json.Write(w, http.StatusNotFound, nil)
	}
	ph.dataStore.Mu.RUnlock()
}

func (ph *playerHandler) GetStatRanks(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	// check if player exists
	ph.dataStore.Mu.RLock()
	_, ok := ph.dataStore.Players[uuid]
	ph.dataStore.Mu.RUnlock()
	if !ok {
		json.Write(w, http.StatusNotFound, nil)
	}

	sr := responsemodels.StatLeaderboardPlayer{
		Uuid:              uuid,
		BlocksPlacedRank:  ph.dataStore.StatLeaderboards.BlocksPlacedLeaderboard[uuid].Rank,
		BlocksBrokenRank:  ph.dataStore.StatLeaderboards.BlocksBrokenLeaderboard[uuid].Rank,
		DiamondsMinedRank: ph.dataStore.StatLeaderboards.DiamondsMinedLeaderboard[uuid].Rank,
		TimePlayedRank:    ph.dataStore.StatLeaderboards.TimePlayedLeaderboard[uuid].Rank,
		PvpKillsRank:      ph.dataStore.StatLeaderboards.PvpKillsLeaderboard[uuid].Rank,
		DeathsRank:        ph.dataStore.StatLeaderboards.DeathsLeaderboard[uuid].Rank,
		MobKillsRank:      ph.dataStore.StatLeaderboards.MobKillsLeaderboard[uuid].Rank,
		PvpKdRatioRank:    ph.dataStore.StatLeaderboards.PvpKdRatioLeaderboard[uuid].Rank,
	}

	json.Write(w, http.StatusOK, sr)
}

func (ph *playerHandler) GetStatRank(w http.ResponseWriter, r *http.Request) {
	stat := chi.URLParam(r, "stat")
	uuid := chi.URLParam(r, "uuid")

	// check if player exists
	ph.dataStore.Mu.RLock()
	_, ok := ph.dataStore.Players[uuid]
	if !ok {
		json.Write(w, http.StatusNotFound, nil)
	}
	ph.dataStore.Mu.RUnlock()

	switch stat {
	case "blocks-placed":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.BlocksPlacedLeaderboard[uuid])
	case "blocks-broken":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.BlocksBrokenLeaderboard[uuid])
	case "diamonds-mined":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.DiamondsMinedLeaderboard[uuid])
	case "time-played":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.TimePlayedLeaderboard[uuid])
	case "pvp-kills":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.PvpKillsLeaderboard[uuid])
	case "deaths":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.DeathsLeaderboard[uuid])
	case "mob-kills":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.MobKillsLeaderboard[uuid])
	case "pvp-kd-ratio":
		json.Write(w, http.StatusOK, ph.dataStore.StatLeaderboards.PvpKdRatioLeaderboard[uuid])
	default:
		json.Write(w, http.StatusNotFound, nil)
	}
}
