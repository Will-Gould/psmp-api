package statistics

import (
	"context"
	"sort"
)

type FeaturedPlayers struct {
	ChampionPlayer      Overview
	BiggestBuilder      Overview
	MostDangerousPlayer Overview
}

func GetChampionPlayer(sh StatisticsHandler, ctx context.Context, players []Overview) Overview {

	// sort by total score
	sort.Slice(players, func(i, j int) bool {
		if players[i].Score < players[j].Score {
			return true
		}
		return false
	})

	return players[0]
}

func GetBiggestBuilder(sh StatisticsHandler, ctx context.Context, players []Overview) Overview {

	// sort by blocks placed
	sort.Slice(players, func(i, j int) bool {
		if players[i].CraftingOverview.BlocksPlaced < players[j].CraftingOverview.BlocksPlaced {
			return true
		}
		return false
	})

	return players[0]
}

func GetMostDangerousPlayer(sh StatisticsHandler, ctx context.Context, players []Overview) Overview {

	// sort by blocks placed
	sort.Slice(players, func(i, j int) bool {
		if players[i].CombatOverview.PvpKdRatio < players[j].CombatOverview.PvpKdRatio {
			return true
		}
		return false
	})

	return players[0]
}
