package responsemodels

type ServerPlayer struct {
	Uuid              string
	Name              string
	PrimaryGroup      string
	ServerRank        int64
	CraftingOverview  CraftingOverview
	CombatOverview    CombatOverview
	AdventureOverview AdventureOverview
	Score             float64
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
	MobKills             int64
	Deaths               int64
	ElderGuardiansKilled int64
	EnderDragonsKilled   int64
	EvokersKilled        int64
	PiglinBrutesKilled   int64
	WardensKilled        int64
	WithersKilled        int64
}

type AdventureOverview struct {
	AdventureScore float64
	Advancements   int32
	FishCaught     int64
}

type LeaderboardPlayer struct {
	Uuid string
	Rank int64
}

type CombatLeaderboardPlayer struct {
	Uuid         string
	Name         string
	PrimaryGroup string
	Rank         int64
	PvpKdRatio   float64
	PvpKills     int64
	MobKills     int64
	Deaths       int64
}

type CraftingLeaderboardPlayer struct {
	Uuid          string
	Name          string
	PrimaryGroup  string
	Rank          int64
	BlocksPlaced  int64
	BlocksBroken  int64
	DiamondsMined int64
}

type AdventureLeaderboardPlayer struct {
	Uuid         string
	Name         string
	PrimaryGroup string
	Rank         int64
	Advancements int64
	FishCaught   int64
}

type StatLeaderboardPlayer struct {
	Uuid              string
	BlocksPlacedRank  int64
	BlocksBrokenRank  int64
	DiamondsMinedRank int64
	TimePlayedRank    int64
	PvpKillsRank      int64
	DeathsRank        int64
	MobKillsRank      int64
	PvpKdRatioRank    int64
}
