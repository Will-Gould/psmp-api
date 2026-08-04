package profiles

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"sort"
	"time"

	Uuid "github.com/google/uuid"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/cache"
	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/Will-Gould/psmp-api/internal/translation"
	"github.com/go-chi/chi"
)

const DATE_FORMAT = "2006-01-02"

type profileHandler struct {
	service    Service
	dataStore  *cache.DataStore
	gameLogger string
}

func NewHandler(service Service, ds *cache.DataStore, gameLogger string) *profileHandler {
	return &profileHandler{
		service:    service,
		dataStore:  ds,
		gameLogger: gameLogger,
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
	//TODO is this lock in correct place?
	ph.dataStore.Mu.RUnlock()
	// model player profile from data in cache
	profile := responsemodels.PlayerProfile{
		Player:        player,
		CombatRank:    ph.dataStore.CombatLeaderboard[uuid].Rank,
		CraftingRank:  ph.dataStore.CraftingLeaderboard[uuid].Rank,
		AdventureRank: ph.dataStore.AdventureLeaderboard[uuid].Rank,
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
	loggerUser, err := ph.findLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	firstJoin := ph.getFirstJoin(r.Context(), loggerUser.ID, psmpStatsPlayer.ID)

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

	loggerUser, err := ph.findLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	var blocksBroken []repo.GroupCountBlocksPlacedByUserRow
	switch ph.gameLogger {
	case "grieflogger":
		glBlocksBroken, err := ph.service.GriefLoggerGroupCountBlocksByUser(r.Context(), loggerUser.ID, ph.dataStore.MappingData.Actions["block-break"], ph.dataStore.MappingData.BannedBrokenObjects)
		if err != nil {
			json.Write(w, http.StatusInternalServerError, nil)
			return
		}
		blocksBroken = glBlocksBroken
	default:
		ledgerBlocksBroken, err := ph.service.LedgerGroupCountBlocksByUser(r.Context(), ph.dataStore.MappingData.Actions, loggerUser.ID, ph.dataStore.MappingData.Actions["block-break"], ph.dataStore.MappingData.BannedBrokenObjects)
		if err != nil {
			json.Write(w, http.StatusInternalServerError, nil)
			slog.Log(r.Context(), slog.LevelError, err.Error())
			return
		}
		blocksBroken = ledgerBlocksBroken
	}

	// transform into materials
	for _, b := range blocksBroken {
		name := ph.findObjectName(b.Type)
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

	loggerUser, err := ph.findLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	var blocksBroken []translation.Block
	var blocksPlaced []translation.Block
	switch ph.gameLogger {
	case "grieflogger":
		glBlocksBroken, err := ph.service.GriefLoggerListBlocksByUser(r.Context(), loggerUser.ID, mapping.BLOCK_BROKEN_ACTION, ph.dataStore.MappingData.BannedBrokenObjects)
		glBlocksPlaced, err := ph.service.GriefLoggerListBlocksByUser(r.Context(), loggerUser.ID, mapping.BLOCK_PLACED_ACTION, ph.dataStore.MappingData.BannedPlacedObjects)
		if err != nil {
			json.Write(w, http.StatusInternalServerError, nil)
		}
		blocksBroken = translation.GriefLoggerTranslateToBlocks(glBlocksBroken)
		blocksPlaced = translation.GriefLoggerTranslateToBlocks(glBlocksPlaced)
	default:
		ledgerBlocksBroken, err := ph.service.LedgerListBlocksByUser(r.Context(), ph.dataStore.MappingData.Actions, loggerUser.ID, ph.dataStore.MappingData.Actions["block-break"], ph.dataStore.MappingData.BannedBrokenObjects)
		ledgerBlocksPlaced, err := ph.service.LedgerListBlocksByUser(r.Context(), ph.dataStore.MappingData.Actions, loggerUser.ID, ph.dataStore.MappingData.Actions["block-place"], ph.dataStore.MappingData.BannedPlacedObjects)
		if err != nil {
			json.Write(w, http.StatusInternalServerError, nil)
		}
		blocksBroken = translation.LedgerTranslateToBlocks(ledgerBlocksBroken)
		blocksPlaced = translation.LedgerTranslateToBlocks(ledgerBlocksPlaced)
	}

	// combine block action into one slice & sort
	blocks := append(blocksBroken, blocksPlaced...)

	chart := GetTotalBlocksChartData(blocks, ph.dataStore.MappingData.Actions)

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) GetDeathsChart(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	psmpStatsPlayer, err := ph.service.FindPsmpstatsPlayerByUuid(r.Context(), uuid)
	loggerUser, err := ph.findLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	firstJoin := ph.getFirstJoin(r.Context(), loggerUser.ID, psmpStatsPlayer.ID)

	deaths, err := ph.service.ListDeathsByPlayer(r.Context(), psmpStatsPlayer.ID)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, err.Error())
		return
	}

	chart := GetDeathsChartData(deaths, firstJoin)

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) GetMostKilledMob(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	psmpStatsPlayer, err := ph.service.FindPsmpstatsPlayerByUuid(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	groupedMobKills, err := ph.service.GroupCountMobKillsByPlayer(r.Context(), psmpStatsPlayer.ID)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}
	if len(groupedMobKills) < 1 {
		json.Write(w, http.StatusOK, nil)
		return
	}

	sort.Slice(groupedMobKills, func(i, j int) bool {
		if groupedMobKills[i].TotalKilled > groupedMobKills[j].TotalKilled {
			return true
		}
		return false
	})

	name := ph.findMobName(groupedMobKills[0].Mob)
	mob := struct {
		Name   string
		Killed int64
	}{
		Name:   name,
		Killed: groupedMobKills[0].TotalKilled,
	}

	json.Write(w, http.StatusOK, mob)
}

func (ph *profileHandler) findObjectName(id int32) string {
	for _, o := range ph.dataStore.MappingData.Objects {
		if id == o.ID {
			return o.Name
		}
	}
	return ""
}

func (ph *profileHandler) findMobName(id int32) string {
	for _, m := range ph.dataStore.MappingData.MobMapping {
		if id == int32(m.ID) {
			return m.Name
		}
	}
	return ""
}

func (ph *profileHandler) findLoggerUser(ctx context.Context, uuid string) (translation.LoggerUser, error) {
	loggerUser := translation.LoggerUser{}
	switch ph.gameLogger {
	case "grieflogger":
		glUser, err := ph.service.FindGriefLoggerUser(ctx, uuid)
		if err != nil {
			slog.Log(ctx, slog.LevelDebug, "Unable to find Grieflogger user")
			return loggerUser, err
		}
		loggerUser.ID = glUser.ID
		loggerUser.Name = glUser.Name
		loggerUser.Uuid = glUser.Uuid
	default:
		parsedUuid, err := Uuid.Parse(uuid)
		if err != nil {
			slog.Log(ctx, slog.LevelDebug, "Error parsing uuid")
			return loggerUser, err
		}
		ledgerUser, err := ph.service.FindLedgerPlayer(ctx, parsedUuid[:])
		if err != nil {
			slog.Log(ctx, slog.LevelDebug, "Unable to find Ledger player")
			return loggerUser, err
		}
		loggerUser.ID = ledgerUser.ID
		loggerUser.Name = ledgerUser.PlayerName
		loggerUser.Uuid = uuid
	}

	return loggerUser, nil
}

func (ph *profileHandler) getFirstJoin(ctx context.Context, loggerId int32, psmpstatsId int32) translation.Session {
	var firstJoin translation.Session
	switch ph.gameLogger {
	case "grieflogger":
		glFirstJoin, err := ph.service.FindGriefLoggerFirstJoin(ctx, loggerId)
		if err != nil {
			slog.Log(ctx, slog.LevelError, "Error finding Grief Logger first join")
		} else {
			firstJoin = translation.Session{
				ID:     glFirstJoin.User,
				Time:   glFirstJoin.Time,
				Action: glFirstJoin.Action,
			}
		}
	default:
		pFirstJoin, err := ph.service.FindPsmpstatsFirstJoin(ctx, psmpstatsId)
		if err != nil {
			slog.Log(ctx, slog.LevelError, "Error finding PSMP Stats first join")
		} else {
			firstJoin = translation.Session{
				ID:     pFirstJoin.PlayerID,
				Time:   time.Unix(int64(pFirstJoin.Time), 0).UnixMilli(),
				Action: pFirstJoin.Action,
			}
		}
	}
	return firstJoin
}
