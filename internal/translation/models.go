package translation

type LoggerUser struct {
	ID   int32
	Name string
	Uuid string
}

type Block struct {
	Time   int64
	User   int32
	World  int32
	X      int32
	Y      int32
	Z      int32
	Object int32
	Action int32
}

type Session struct {
	ID     int32
	Time   int64
	Action int32
}
