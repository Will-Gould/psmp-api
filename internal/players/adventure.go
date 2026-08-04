package players

import (
	"context"
	"log/slog"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
)

var ONE_POINT_MILESTONES = []string{
	"minecraft:story/follow_ender_eye",
	"minecraft:story/enter_the_nether",
	"minecraft:nether/find_fortress",
	"minecraft:adventure/minecraft_trials_edition",
	"minecraft:story/mine_diamond",
	"minecraft:story/enchant_item",
}

var TWO_POINT_MILESTONES = []string{
	"minecraft:adventure/salvage_sherd",
	"minecraft:nether/explore_nether",
	"minecraft:end/find_end_city",
}

var THREE_POINT_MILESTONES = []string{
	"minecraft:end/kill_dragon",
	"minecraft:adventure/hero_of_the_village",
	"minecraft:nether/create_beacon",
	"minecraft:nether/loot_bastion",
}

var FOUR_POINT_MILESTONES = []string{
	"minecraft:adventure/adventuring_time",
	"minecraft:adventure/kill_all_mobs",
	"minecraft:nether/uneasy_alliance",
}

var FIVE_POINT_MILESTONES = []string{
	"minecraft:nether/all_effects",
}

func (ph *playerHandler) getAdventureOverview(ctx context.Context, psmpStatsId int32) (responsemodels.AdventureOverview, error) {
	advancements, err := ph.service.ListAdvancementsByPlayer(ctx, psmpStatsId)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	fishCaught, err := ph.service.CountFishCaughtByPlayer(ctx, psmpStatsId)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Error counting fish caught")
	}

	adventureOverview := responsemodels.AdventureOverview{
		AdventureScore: calculateAdventureScore(advancements, fishCaught),
		Advancements:   int32(len(advancements)),
		FishCaught:     fishCaught,
	}

	return adventureOverview, nil
}

func calculateAdventureScore(advancements []repo.PsmpstatsAdvancement, fishCaught int64) float64 {
	var adventureScore float64 = 0

	// 1/100 fish caught
	adventureScore += (float64(fishCaught) / 100)

	for _, m := range ONE_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				adventureScore += 1
			}
		}
	}
	for _, m := range TWO_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				adventureScore += 2
			}
		}
	}
	for _, m := range THREE_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				adventureScore += 3
			}
		}
	}
	for _, m := range FOUR_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				adventureScore += 4
			}
		}
	}
	for _, m := range FIVE_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				adventureScore += 5
			}
		}
	}

	return adventureScore
}
