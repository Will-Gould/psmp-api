package statistics

import (
	"context"
	"sort"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type CraftingOverview struct {
	CraftingScore float64
	BlocksPlaced  int64
	BlocksBroken  int64
	DiamondsMined int64
	TimePlayed    int64
}

func GetCraftingOverview(ctx context.Context, sh StatisticsHandler, uuid string, glId int32) (CraftingOverview, error) {
	// count blocks
	blocksBroken, err := sh.Service.CountBlocksByUser(ctx, glId, BLOCK_BROKEN_ACTION, sh.BannedBrokenMaterials)
	if err != nil {
		blocksBroken = 0
	}
	blocksPlaced, err := sh.Service.CountBlocksByUser(ctx, glId, BLOCK_PLACED_ACTION, sh.BannedPlacedMaterials)
	if err != nil {
		blocksPlaced = 0
	}

	// count time played
	var timePlayed int64
	sessions, err := sh.Service.ListSessionDataByUser(ctx, glId)
	if err != nil {
		timePlayed = 0
	} else {
		timePlayed = countTimePlayed(sessions)
	}

	// get diamonds mined
	var diamondMined int64
	psmpstatsPlayer, err := sh.Service.FindPsmpstatsPlayerByUuid(ctx, uuid)
	if err != nil {
		diamondMined = 0
	} else {
		diamondMined = int64(psmpstatsPlayer.DiamondsMined)
	}

	// calculate crafting score
	craftingScore := calculateCraftingScore(float64(blocksBroken), float64(blocksPlaced), float64(timePlayed), float64(diamondMined))

	return CraftingOverview{
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
		if lastSession.Action == PLAYER_JOIN_ACTION && s.Action == PLAYER_LEAVE_ACTION {
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
