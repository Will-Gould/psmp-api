package responsemodels

type ServerPlayer struct {
	Uuid             string
	Name             string
	GlId             int32
	PrimaryGroup     string
	ServerRank       int64
	CraftingOverview CraftingOverview
	CombatOverview   CombatOverview
	StoryOverview    StoryOverview
	Score            float64
}

type CraftingOverview struct {
	CraftingScore float64
	BlocksPlaced  int64
	BlocksBroken  int64
	DiamondsMined int64
	TimePlayed    int64
}

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

type StoryOverview struct {
	StoryScore float64
}

type LeaderboardPlayer struct {
	Uuid    string
	Ranking int64
	Value   float64
}
