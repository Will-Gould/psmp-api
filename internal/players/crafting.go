package players

import (
	"context"
	"sort"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/cache"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
)

func (ph *playerHandler) getCraftingOverview(ctx context.Context, uuid string, glId int32, md *cache.MappingData) (responsemodels.CraftingOverview, error) {
	// count blocks
	blocksBroken, err := ph.service.CountBlocksByUser(ctx, glId, mapping.BLOCK_BROKEN_ACTION, md.BannedBrokenMaterials)
	if err != nil {
		blocksBroken = 0
	}
	blocksPlaced, err := ph.service.CountBlocksByUser(ctx, glId, mapping.BLOCK_PLACED_ACTION, md.BannedPlacedMaterials)
	if err != nil {
		blocksPlaced = 0
	}

	// count time played
	var timePlayed int64
	sessions, err := ph.service.ListSessionDataByUser(ctx, glId)
	if err != nil {
		timePlayed = 0
	} else {
		timePlayed = countTimePlayed(sessions)
	}

	// get diamonds mined
	var diamondMined int64
	psmpstatsPlayer, err := ph.service.FindPsmpstatsPlayerByUuid(ctx, uuid)
	if err != nil {
		diamondMined = 0
	} else {
		diamondMined = int64(psmpstatsPlayer.DiamondsMined)
	}

	// calculate crafting score
	craftingScore := calculateCraftingScore(float64(blocksBroken), float64(blocksPlaced), float64(timePlayed), float64(diamondMined))

	return responsemodels.CraftingOverview{
		CraftingScore: craftingScore,
		BlocksPlaced:  blocksPlaced,
		BlocksBroken:  blocksBroken,
		DiamondsMined: diamondMined,
		TimePlayed:    timePlayed,
	}, nil
}

func countTimePlayed(sessions []repo.Session) int64 {
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
