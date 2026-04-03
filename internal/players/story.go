package players

import (
	"context"
	"log/slog"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
)

var ONE_POINT_MILESTONES = []string{
	"minecraft:story/follow_ender_eye",
	"minecraft:nether/find_bastion",
	"minecraft:nether/find_fortress",
	"minecraft:adventure/minecraft_trials_edition",
	"minecraft:end/find_end_city",
}

var TWO_POINT_MILESTONES = []string{
	"minecraft:adventure/salvage_sherd",
	"minecraft:nether/explore_nether",
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

func (ph playerHandler) getStoryOverview(ctx context.Context, uuid string) (responsemodels.StoryOverview, error) {
	advancements, err := ph.service.ListAdvancementsByPlayer(ctx, uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	storyOverview := responsemodels.StoryOverview{
		StoryScore: calculateStoryScore(advancements),
	}

	return storyOverview, nil
}

func calculateStoryScore(advancements []repo.PsmpstatsAdvancement) float64 {
	var storyScore float64

	for _, m := range ONE_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				storyScore += 1
			}
		}
	}
	for _, m := range TWO_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				storyScore += 2
			}
		}
	}
	for _, m := range THREE_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				storyScore += 3
			}
		}
	}
	for _, m := range FOUR_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				storyScore += 4
			}
		}
	}
	for _, m := range FIVE_POINT_MILESTONES {
		for _, a := range advancements {
			if a.Name == m {
				storyScore += 5
			}
		}
	}

	return storyScore
}
