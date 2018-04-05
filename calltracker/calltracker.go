package calltracker

type Phone struct {
	Id     int64
	Number string
}

type Audio struct {
	Buffer []byte
}

type Call struct {
	Id    int64
	Phone Phone
	Audio Audio
}

type CallStorage interface {
	Save(c *Call) (*int64, error)
	Dump(c *Call) (*string, error)
}

type AudioService interface {
	Process(*Call) (*Call, error)
	Convert(*Call) (*Call, error)
	Send(*Call) (string, error)
	GetDuration() (string, error)
	GetSTT() (string, error)
}
