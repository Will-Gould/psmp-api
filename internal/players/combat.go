package players

import (
	"context"

	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
)

func calculateCombatScore(
	pvpKd float64,
	mobKd float64,
	pvpKills float64,
	mobKills float64,
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

func (ph *playerHandler) getCombatOverview(ctx context.Context, id int32) (responsemodels.CombatOverview, error) {
	// count deaths
	deaths, err := ph.service.CountDeathsByPlayer(ctx, id)
	if err != nil {
		deaths = 0
	}

	// count PvP kills
	pvpKills, err := ph.service.CountPvpKillsByPlayer(ctx, id)
	if err != nil {
		pvpKills = 0
	}

	// count mobs killed
	mobKills, err := ph.service.CountMobKillsByPlayer(ctx, id)
	if err != nil {
		mobKills = 0
	}

	// get pvp & mob K/D ratio
	var pvpKd float64 = 0
	var mobKd float64 = 0
	if deaths == 0 {
		pvpKd = float64(pvpKills)
		mobKd = float64(mobKills)
	} else {
		pvpKd = float64(pvpKills) / float64(deaths)
		mobKd = float64(mobKills) / float64(deaths)
	}

	// get notable mobs killed
	elderGuardiansKilled, err := ph.service.CountSpecificMobKillsByPlayer(ctx, mapping.ELDER_GUARDIAN, id)
	if err != nil {
		elderGuardiansKilled = 0
	}
	enderDragonsKilled, err := ph.service.CountSpecificMobKillsByPlayer(ctx, mapping.ENDER_DRAGON, id)
	if err != nil {
		enderDragonsKilled = 0
	}
	evokersKilled, err := ph.service.CountSpecificMobKillsByPlayer(ctx, mapping.EVOKER, id)
	if err != nil {
		evokersKilled = 0
	}
	piglinBrutesKilled, err := ph.service.CountSpecificMobKillsByPlayer(ctx, mapping.PIGLIN_BRUTE, id)
	if err != nil {
		piglinBrutesKilled = 0
	}
	wardensKilled, err := ph.service.CountSpecificMobKillsByPlayer(ctx, mapping.WARDEN, id)
	if err != nil {
		wardensKilled = 0
	}
	withersKilled, err := ph.service.CountSpecificMobKillsByPlayer(ctx, mapping.WITHER, id)
	if err != nil {
		withersKilled = 0
	}

	combatScore := calculateCombatScore(
		pvpKd, mobKd,
		float64(pvpKills),
		float64(mobKills),
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
		MobKills:             mobKills,
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
