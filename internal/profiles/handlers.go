package profiles

import (
	"log/slog"
	"net/http"
	"slices"
	"sort"

	"github.com/Will-Gould/psmp-api/internal/cache"
	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

const DATE_FORMAT = "2006-01-02"

type profileHandler struct {
	service   Service
	dataStore *cache.DataStore
}

func NewHandler(service Service, ds *cache.DataStore) *profileHandler {
	return &profileHandler{
		service:   service,
		dataStore: ds,
	}
}

func (ph *profileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	ph.dataStore.Mu.RLock()
	player, ok := ph.dataStore.Players[uuid]
	if !ok {
		json.Write(w, http.StatusNotFound, nil)
		ph.dataStore.Mu.RUnlock()
		return
	}
	ph.dataStore.Mu.RUnlock()
	// model player profile from data in cache
	profile := responsemodels.PlayerProfile{
		Player:       player,
		CombatRank:   ph.dataStore.CombatLeaderboard[uuid].Rank,
		CraftingRank: ph.dataStore.CraftingLeaderboard[uuid].Rank,
		StoryRank:    ph.dataStore.StoryLeaderboard[uuid].Rank,
		StatLeaderboardPlayer: responsemodels.StatLeaderboardPlayer{
			Uuid:              uuid,
			BlocksPlacedRank:  ph.dataStore.StatLeaderboards.BlocksPlacedLeaderboard[uuid].Rank,
			BlocksBrokenRank:  ph.dataStore.StatLeaderboards.BlocksBrokenLeaderboard[uuid].Rank,
			DiamondsMinedRank: ph.dataStore.StatLeaderboards.DiamondsMinedLeaderboard[uuid].Rank,
			TimePlayedRank:    ph.dataStore.StatLeaderboards.TimePlayedLeaderboard[uuid].Rank,
			PvpKillsRank:      ph.dataStore.StatLeaderboards.PvpKillsLeaderboard[uuid].Rank,
			DeathsRank:        ph.dataStore.StatLeaderboards.DeathsLeaderboard[uuid].Rank,
			MobKillsRank:      ph.dataStore.StatLeaderboards.MobKillsLeaderboard[uuid].Rank,
			PvpKdRatioRank:    ph.dataStore.StatLeaderboards.PvpKdRatioLeaderboard[uuid].Rank,
		},
	}

	json.Write(w, http.StatusOK, profile)
}

func (ph *profileHandler) GetMobKillChart(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	psmpStatsPlayer, err := ph.service.FindPsmpstatsPlayerByUuid(r.Context(), uuid)
	glUser, err := ph.service.FindGriefLoggerUser(r.Context(), uuid)
	firstJoin, err := ph.service.FindFirstJoinByUser(r.Context(), glUser.ID)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	mobKills, err := ph.service.ListMobKillsByPlayer(r.Context(), psmpStatsPlayer.ID)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, err.Error())
		return
	}

	chart := GetMobKillChartData(mobKills, firstJoin)

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) GetBlocksBrokenPieChart(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	data := []responsemodels.BlockChartItem{}

	glUser, err := ph.service.FindGriefLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	blocksBroken, err := ph.service.GroupCountBlocksByUser(r.Context(), glUser.ID, mapping.BLOCK_BROKEN_ACTION, ph.dataStore.MappingData.BannedBrokenMaterials)

	// transform into materials
	for _, b := range blocksBroken {
		name := ph.findMaterialName(b.Type)
		data = append(data, responsemodels.BlockChartItem{
			Block: name,
			Value: b.TotalPlaced,
		})
	}

	// consolidate into 'others' category
	if len(data) > 10 {
		sort.Slice(data, func(i, j int) bool {
			if data[i].Value > data[j].Value {
				return true
			}
			return false
		})

		data[9].Block = "others"

		for i := 10; i < len(data); i++ {
			data[9].Value += data[i].Value
		}

		data = slices.Delete(data, 10, len(data))
	}

	blockChart := responsemodels.BlockChart{
		Title:     "Blocks Broken",
		ChartData: data,
	}

	json.Write(w, http.StatusOK, blockChart)
}

func (ph *profileHandler) GetTotalBlocksChart(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	glUser, err := ph.service.FindGriefLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	blocksBroken, err := ph.service.ListBlocksByUser(r.Context(), glUser.ID, mapping.BLOCK_BROKEN_ACTION, ph.dataStore.MappingData.BannedBrokenMaterials)
	if err != nil {
		json.Write(w, http.StatusInternalServerError, nil)
	}
	blocksPlaced, err := ph.service.ListBlocksByUser(r.Context(), glUser.ID, mapping.BLOCK_PLACED_ACTION, ph.dataStore.MappingData.BannedPlacedMaterials)
	if err != nil {
		json.Write(w, http.StatusInternalServerError, nil)
	}

	// combine block action into one slice & sort
	blocks := append(blocksBroken, blocksPlaced...)

	chart := GetTotalBlocksChartData(blocks)

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) GetDeathsChart(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	psmpStatsPlayer, err := ph.service.FindPsmpstatsPlayerByUuid(r.Context(), uuid)
	glUser, err := ph.service.FindGriefLoggerUser(r.Context(), uuid)
	firstJoin, err := ph.service.FindFirstJoinByUser(r.Context(), glUser.ID)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	deaths, err := ph.service.ListDeathsByPlayer(r.Context(), psmpStatsPlayer.ID)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, err.Error())
		return
	}

	chart := GetDeathsChartData(deaths, firstJoin)

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) findMaterialName(id int32) string {
	for _, m := range ph.dataStore.MappingData.Materials {
		if id == m.ID {
			return m.Name
		}
	}
	return ""
}
