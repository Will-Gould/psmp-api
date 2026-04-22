package profiles

import (
	"log/slog"
	"net/http"
	"slices"
	"sort"
	"time"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
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

func (ph *profileHandler) GetMobKillChartData(w http.ResponseWriter, r *http.Request) {
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
	sort.Slice(mobKills, func(i, j int) bool {
		if mobKills[i].Time < mobKills[j].Time {
			return true
		}
		return false
	})

	data := getAllTimeSingleDailySlice(firstJoin)

	for _, k := range mobKills {
		killDate := time.Unix(int64(k.Time), 0).Local().Format(DATE_FORMAT)
		for i, d := range data {
			if d.Date == killDate {
				data[i].Value += 1
				continue
			}
		}
	}

	// calculate trend over last two intervals
	var trend float64 = 0
	if len(data) > 1 {
		cur := data[len(data)-1].Value
		base := data[len(data)-2].Value

		if base > 0 {
			trend = ((float64(cur) - float64(base)) / float64(base)) * 100
		}
	}

	chart := responsemodels.SingleDailyChart{
		Title:     "Mob Kills",
		ChartData: data,
		Trend:     trend,
	}

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) GetBlocksBrokenPieChartData(w http.ResponseWriter, r *http.Request) {
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
	data := []responsemodels.DoubleDailyChartItem{}

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

	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i].Time < blocks[j].Time {
			return true
		}
		return false
	})
	firstBlockTime := time.Unix(blocks[0].Time/1000, 0)
	earliestMidnight := firstBlockTime.Truncate(24 * time.Hour)
	nextDay := time.Now().AddDate(0, 0, 1).Local()
	nextMidnight := nextDay.Truncate(24 * time.Hour)
	// initialise value for each date between now & first block action
	for d := earliestMidnight; !d.After(nextMidnight); d = d.AddDate(0, 0, 1) {
		dateString := d.Local().Format(DATE_FORMAT)
		data = append(data, responsemodels.DoubleDailyChartItem{
			Date:   dateString,
			Value1: 0,
			Value2: 0,
		})
	}

	for i, b := range blocks {
		blockDate := time.Unix(blocks[i].Time/1000, 0).Format(DATE_FORMAT)
		for j, d := range data {

			if d.Date == blockDate {

				if b.Action == mapping.BLOCK_PLACED_ACTION {
					data[j].Value1 += 1
				}
				if b.Action == mapping.BLOCK_BROKEN_ACTION {
					data[j].Value2 += 1
				}
			}
		}
	}

	// combine date values so each day has a running total
	for i := range data {
		if i == 0 {
			continue
		}
		data[i].Value1 += data[i-1].Value1
		data[i].Value2 += data[i-1].Value2
	}

	chart := responsemodels.DoubleDailyChart{
		Title:     "Total Blocks Placed vs Broken",
		ChartData: data,
		Trend:     0,
	}

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
	sort.Slice(deaths, func(i, j int) bool {
		if deaths[i].Time < deaths[j].Time {
			return true
		}
		return false
	})

	data := getAllTimeSingleDailySlice(firstJoin)

	for _, k := range deaths {
		deathDate := time.Unix(int64(k.Time), 0).Local().Format(DATE_FORMAT)
		for i, d := range data {
			if d.Date == deathDate {
				data[i].Value += 1
				continue
			}
		}
	}

	// calculate trend over last two intervals
	var trend float64 = 0
	if len(data) > 1 {
		cur := data[len(data)-1].Value
		base := data[len(data)-2].Value

		if base > 0 {
			trend = ((float64(cur) - float64(base)) / float64(base)) * 100
		}
	}

	chart := responsemodels.SingleDailyChart{
		Title:     "Deaths",
		ChartData: data,
		Trend:     trend,
	}

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

func getAllTimeSingleDailySlice(firstJoin repo.Session) []responsemodels.SingleDailyChartItem {
	data := []responsemodels.SingleDailyChartItem{}

	firstTime := time.Unix(int64(firstJoin.Time/1000), 0).Local()
	earliestMidnight := firstTime.Truncate(24 * time.Hour)

	nextDay := time.Now().AddDate(0, 0, 1).Local()
	nextMidnight := nextDay.Truncate(24 * time.Hour)
	// initialise value for each date between now and first join
	for d := earliestMidnight; !d.After(nextMidnight); d = d.AddDate(0, 0, 1) {
		dateString := d.Local().Format(DATE_FORMAT)
		data = append(data, responsemodels.SingleDailyChartItem{
			Date:  dateString,
			Value: 0,
		})
	}

	return data
}
