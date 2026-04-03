package leaderboards

import (
	"context"

	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
)

func calculateCombatScore(
	pvpKd float64,
	mobKd float64,
	pvpKills float64,
	mobsKilled float64,
	deaths float64,
	elderGuardiansKilled float64,
	enderDragonsKilled float64,
	evokersKilled float64,
	piglinBrutesKilled float64,
	wardensKilled float64,
	withersKilled float64) float64 {

	var combatScore float64

	// pvp k/d ratio * 3
	combatScore += (pvpKd * 4)
	//3% of mobs killed/death ratio
	combatScore += (mobKd * 0.03)
	// (Raids won) * 3

	//(Bastions cleared + woodland mansions cleared + ocean monuments cleared) * 2
	combatScore += ((elderGuardiansKilled * 3) + (evokersKilled * 5) + (piglinBrutesKilled * 5)) * 2

	//Wardens killed * 4
	combatScore += (wardensKilled * 4)

	//Ender dragons killed
	combatScore += (enderDragonsKilled)
	//50% of withers killed
	combatScore += (withersKilled * 0.5)

	return combatScore
}

func (lh leaderboardHandler) getCombatOverview(ctx context.Context, uuid string) (responsemodels.CombatOverview, error) {
	// count deaths
	deaths, err := lh.service.CountDeathsByPlayer(ctx, uuid)
	if err != nil {
		deaths = 0
	}

	// count PvP kills
	pvpKills, err := lh.service.CountPvpKillsByPlayer(ctx, uuid)
	if err != nil {
		pvpKills = 0
	}

	// count mobs killed
	mobsKilled, err := lh.service.CountMobKillsByPlayer(ctx, uuid)
	if err != nil {
		mobsKilled = 0
	}

	// get pvp & mob K/D ratio
	var pvpKd float64 = 0
	var mobKd float64 = 0
	if deaths == 0 {
		pvpKd = float64(pvpKills)
		mobKd = float64(mobsKilled)
	} else {
		pvpKd = float64(pvpKills) / float64(deaths)
		mobKd = float64(mobsKilled) / float64(deaths)
	}

	// get notable mobs killed
	elderGuardiansKilled, err := lh.service.CountSpecificMobKillsByUuid(ctx, mapping.ELDER_GUARDIAN, uuid)
	if err != nil {
		elderGuardiansKilled = 0
	}
	enderDragonsKilled, err := lh.service.CountSpecificMobKillsByUuid(ctx, mapping.ENDER_DRAGON, uuid)
	if err != nil {
		enderDragonsKilled = 0
	}
	evokersKilled, err := lh.service.CountSpecificMobKillsByUuid(ctx, mapping.EVOKER, uuid)
	if err != nil {
		evokersKilled = 0
	}
	piglinBrutesKilled, err := lh.service.CountSpecificMobKillsByUuid(ctx, mapping.PIGLIN_BRUTE, uuid)
	if err != nil {
		piglinBrutesKilled = 0
	}
	wardensKilled, err := lh.service.CountSpecificMobKillsByUuid(ctx, mapping.WARDEN, uuid)
	if err != nil {
		wardensKilled = 0
	}
	withersKilled, err := lh.service.CountSpecificMobKillsByUuid(ctx, mapping.WITHER, uuid)
	if err != nil {
		withersKilled = 0
	}

	combatScore := calculateCombatScore(
		pvpKd, mobKd,
		float64(pvpKills),
		float64(mobsKilled),
		float64(deaths),
		float64(elderGuardiansKilled),
		float64(enderDragonsKilled),
		float64(evokersKilled),
		float64(piglinBrutesKilled),
		float64(wardensKilled),
		float64(withersKilled))

	combatOverview := responsemodels.CombatOverview{
		CombatScore:          combatScore,
		PvpKdRatio:           pvpKd,
		PvpKills:             pvpKills,
		MobsKilled:           mobsKilled,
		Deaths:               deaths,
		ElderGuardiansKilled: elderGuardiansKilled,
		EnderDragonsKilled:   enderDragonsKilled,
		EvokersKilled:        evokersKilled,
		PiglinBrutesKilled:   piglinBrutesKilled,
		WardensKilled:        wardensKilled,
		WithersKilled:        withersKilled,
	}

	return combatOverview, nil
}
