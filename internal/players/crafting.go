package players

import (
	"context"
	"log/slog"
	"sort"

	"github.com/Will-Gould/psmp-api/internal/cache"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/Will-Gould/psmp-api/internal/translation"
)

func (ph *playerHandler) getCraftingOverview(ctx context.Context, loggerId int32, psmpStatsId int32, md *cache.MappingData) (responsemodels.CraftingOverview, error) {
	// count blocks
	var blocksBroken int64
	var blocksPlaced int64
	switch ph.gameLogger {
	case "grieflogger":
		glBlocksBroken, err := ph.service.GriefLoggerCountBlocksByUser(ctx, loggerId, md.Actions["block-break"], md.BannedBrokenObjects)
		if err != nil {
			blocksBroken = 0
		}
		glBlocksPlaced, err := ph.service.GriefLoggerCountBlocksByUser(ctx, loggerId, md.Actions["block-place"], md.BannedPlacedObjects)
		if err != nil {
			blocksPlaced = 0
		}
		blocksBroken = glBlocksBroken
		blocksPlaced = glBlocksPlaced
	default:
		ledgerBlocksBroken, err := ph.service.LedgerCountBlocksByUser(ctx, md.Actions, loggerId, md.Actions["block-break"], md.BannedBrokenObjects)
		if err != nil {
			blocksBroken = 0
		}
		ledgerBlocksPlaced, err := ph.service.LedgerCountBlocksByUser(ctx, md.Actions, loggerId, md.Actions["block-place"], md.BannedPlacedObjects)
		if err != nil {
			blocksPlaced = 0
		}
		blocksBroken = ledgerBlocksBroken
		blocksPlaced = ledgerBlocksPlaced
	}

	// count time played
	var sessions []translation.Session
	var timePlayed int64
	switch ph.gameLogger {
	case "grieflogger":
		glSessions, err := ph.service.ListGriefLoggerSessions(ctx, loggerId)
		if err != nil {
			slog.Log(ctx, slog.LevelError, "Unable to list Grief Logger sessions")
		}
		sessions = translation.GriefLoggerTranslateSessions(glSessions)
	default:
		pSessions, err := ph.service.ListPsmpstatsSessions(ctx, psmpStatsId)
		if err != nil {
			slog.Log(ctx, slog.LevelError, "Unable to list PSMP Stats sessions")
		}
		sessions = translation.PsmpstatsTranslateSessions(pSessions)
	}
	if len(sessions) > 0 {
		timePlayed = countTimePlayed(sessions)
	} else {
		timePlayed = 0
	}

	// get diamonds mined
	diamondsMined, err := ph.service.CountDiamondsMinedByPlayer(ctx, psmpStatsId)
	if err != nil {
		diamondsMined = 0
	}

	// calculate crafting score
	craftingScore := calculateCraftingScore(float64(blocksBroken), float64(blocksPlaced), float64(timePlayed), float64(diamondsMined))

	return responsemodels.CraftingOverview{
		CraftingScore: craftingScore,
		BlocksPlaced:  blocksPlaced,
		BlocksBroken:  blocksBroken,
		DiamondsMined: diamondsMined,
		TimePlayed:    timePlayed,
	}, nil
}

func countTimePlayed(sessions []translation.Session) int64 {
	// check if there are no sessions
	if len(sessions) < 1 {
		return 0
	}

	// make sure slice is in chronological order
	sort.Slice(sessions, func(i, j int) bool {
		if sessions[i].Time < sessions[j].Time {
			return true
		}
		return false
	})

	var timePlayed int64
	var lastSession = sessions[0]
	for _, s := range sessions {
		if lastSession.Action == mapping.PLAYER_JOIN_ACTION && s.Action == mapping.PLAYER_LEAVE_ACTION {
			timePlayed += (s.Time - lastSession.Time)
		}
		lastSession = s
	}

	return timePlayed / 1000
}

func calculateCraftingScore(blocksBroken float64, blocksPlaced float64, timePlayed float64, diamondsMined float64) float64 {
	var craftingScore float64 = 0

	//1/3600 Blocks broken
	craftingScore += (blocksBroken / 3600)
	//1/2100 Blocks Placed
	craftingScore += (blocksPlaced / 2100)
	//+1 For every 3 hours of play time
	craftingScore += (timePlayed / 10800)
	//2% of diamonds mined
	craftingScore += (diamondsMined * 0.02)

	return craftingScore
}
