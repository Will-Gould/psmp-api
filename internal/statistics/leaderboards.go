package statistics

import (
	"context"
	"sort"
)

type LeaderboardRanks struct {
	ServerRank        int64
	CraftingRank      int64
	CombatRank        int64
	PvpKillRank       int64
	MobKillRank       int64
	BlocksPlacedRank  int64
	BlocksBrokenRank  int64
	DiamondsMinedRank int64
}

func (sh StatisticsHandler) UpdateLeaderboard(ctx context.Context) (map[string]LeaderboardRanks, error) {
	leaderboard := make(map[string]LeaderboardRanks)
	players, err := sh.getPlayers(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range players {
		leaderboard[p.Player.Uuid] = LeaderboardRanks{}
	}

	getServerRanks(players, &leaderboard)
	getCombatRanks(players, &leaderboard)
	getCraftingRanks(players, &leaderboard)
	getPvpKillRanks(players, &leaderboard)
	getMobKillRanks(players, &leaderboard)
	getBlocksPlacedRanks(players, &leaderboard)
	getBlocksBrokenRanks(players, &leaderboard)
	getDiamondsMinedRanks(players, &leaderboard)

	return leaderboard, nil
}

func getServerRanks(players []Overview, l *map[string]LeaderboardRanks) {
	// sort by total score
	sort.Slice(players, func(i, j int) bool {
		if players[i].Score < players[j].Score {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.ServerRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}

func getCombatRanks(players []Overview, l *map[string]LeaderboardRanks) {
	sort.Slice(players, func(i, j int) bool {
		if players[i].CombatOverview.CombatScore < players[j].CombatOverview.CombatScore {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.CombatRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}

func getDiamondsMinedRanks(players []Overview, l *map[string]LeaderboardRanks) {
	sort.Slice(players, func(i, j int) bool {
		if players[i].CraftingOverview.DiamondsMined < players[j].CraftingOverview.DiamondsMined {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.DiamondsMinedRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}

func getBlocksBrokenRanks(players []Overview, l *map[string]LeaderboardRanks) {
	sort.Slice(players, func(i, j int) bool {
		if players[i].CraftingOverview.BlocksBroken < players[j].CraftingOverview.BlocksBroken {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.BlocksBrokenRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}

func getBlocksPlacedRanks(players []Overview, l *map[string]LeaderboardRanks) {
	sort.Slice(players, func(i, j int) bool {
		if players[i].CraftingOverview.BlocksPlaced < players[j].CraftingOverview.BlocksPlaced {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.BlocksPlacedRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}

func getMobKillRanks(players []Overview, l *map[string]LeaderboardRanks) {
	sort.Slice(players, func(i, j int) bool {
		if players[i].CombatOverview.MobsKilled < players[j].CombatOverview.MobsKilled {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.MobKillRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}

func getPvpKillRanks(players []Overview, l *map[string]LeaderboardRanks) {
	sort.Slice(players, func(i, j int) bool {
		if players[i].CombatOverview.PvpKills < players[j].CombatOverview.PvpKills {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.PvpKillRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}

func getCraftingRanks(players []Overview, l *map[string]LeaderboardRanks) {
	sort.Slice(players, func(i, j int) bool {
		if players[i].CraftingOverview.CraftingScore < players[j].CraftingOverview.CraftingScore {
			return true
		}
		return false
	})

	for i, p := range players {
		lp := (*l)[p.Player.Uuid]
		lp.CraftingRank = int64(i) + 1
		(*l)[p.Player.Uuid] = lp
	}
}
