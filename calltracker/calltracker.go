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

type CallService interface {
	Save(c *Call) (int, error)
}
