package calltracker

type Phone struct {
	Id     int64
	Number string
}

type Audio struct {
	Buffer []byte
}

type Call struct {
	Id       int64
	Phone    Phone
	Audio    Audio
	RemoteId string
}

type Metric struct {
	Id   int64
	Name string
	Call Call
	Data []byte
}

type CallStorage interface {
	Save(c *Call) (*int64, error)
	Dump(c *Call) (*string, error)
	SaveMetric(m *Metric) (error)
	GetMetrics(string) ([]byte, error)
}

type AudioService interface {
	Process(*Call) (error)
	Convert(string) (*string, error)
	Send(string) (*string, error)
	GetDuration(Call) (error)
	GetSTT(Call) (error)
}
