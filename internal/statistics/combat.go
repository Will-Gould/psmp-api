package statistics

import (
	"context"
	"log/slog"
	"time"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

const DATE_FORMAT = "2006-01-02"

type CombatOverview struct {
	CombatScore          float64
	PvpKdRatio           float64
	PvpKills             int64
	MobsKilled           int64
	Deaths               int64
	ElderGuardiansKilled int64
	EnderDragonsKilled   int64
	EvokersKilled        int64
	PiglinBrutesKilled   int64
	WardensKilled        int64
	WithersKilled        int64
}

type SingleDailyChartData struct {
	Date  string
	Value int64
}

func (sh StatisticsHandler) GetCombatOverview(ctx context.Context, uuid string) (CombatOverview, error) {
	// count deaths
	deaths, err := sh.Service.CountDeathsByPlayer(ctx, uuid)
	if err != nil {
		deaths = 0
	}

	// count PvP kills
	pvpKills, err := sh.Service.CountPvpKillsByPlayer(ctx, uuid)
	if err != nil {
		pvpKills = 0
	}

	// count mobs killed
	mobsKilled, err := sh.Service.CountMobKillsByPlayer(ctx, uuid)
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
	elderGuardiansKilled, err := sh.Service.CountSpecificMobKillsByUuid(ctx, ELDER_GUARDIAN, uuid)
	if err != nil {
		elderGuardiansKilled = 0
	}
	enderDragonsKilled, err := sh.Service.CountSpecificMobKillsByUuid(ctx, ENDER_DRAGON, uuid)
	if err != nil {
		enderDragonsKilled = 0
	}
	evokersKilled, err := sh.Service.CountSpecificMobKillsByUuid(ctx, EVOKER, uuid)
	if err != nil {
		evokersKilled = 0
	}
	piglinBrutesKilled, err := sh.Service.CountSpecificMobKillsByUuid(ctx, PIGLIN_BRUTE, uuid)
	if err != nil {
		piglinBrutesKilled = 0
	}
	wardensKilled, err := sh.Service.CountSpecificMobKillsByUuid(ctx, WARDEN, uuid)
	if err != nil {
		wardensKilled = 0
	}
	withersKilled, err := sh.Service.CountSpecificMobKillsByUuid(ctx, WITHER, uuid)
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

	combatOverview := CombatOverview{
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

func GetMobKillChartData(ctx context.Context, sh StatisticsHandler, uuid string) []SingleDailyChartData {
	chartData := []SingleDailyChartData{}
	mobKills, err := sh.Service.ListMobKillsByUuid(ctx, uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return chartData
	}

	for _, m := range mobKills {
		insertMobKill(&chartData, m)
	}

	return chartData

}

func insertMobKill(cd *[]SingleDailyChartData, m repo.PsmpstatsMobKill) {
	t := time.Unix(int64(m.Time), 0)
	for i, d := range *cd {
		if d.Date == t.Local().Format(DATE_FORMAT) {
			(*cd)[i].Value += 1
			return
		}
	}

	// insert new data point for new date
	d := SingleDailyChartData{
		Date:  t.Local().Format(DATE_FORMAT),
		Value: 1,
	}
	*cd = append(*cd, d)
}

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
