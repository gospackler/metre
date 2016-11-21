// Package metre is used to schedule end execute cron jobs in a simplified fashion
package metre

const LocalHost string = "tcp://127.0.0.1" // Default host for cache and queue
const RouterPort string = "5555"           // Default port for queue
const DealerPort string = "5556"           // Default port for queue

type Metre struct {
	MessageChannel chan string
	MasterIns      *Master
	SlaveIns       *Slave
}

// New creates a new scheduler to manage task scheduling and states
func New(m *Master, s *Slave) (*Metre, error) {
	msgChan := make(chan string)
	return &Metre{msgChan, m, s}, nil
}

// Add adds a cron job task to schedule and process
func (m *Metre) Add(t *Task) {
	if m.MasterIns != nil {
		t.MessageChannel = m.MessageChannel
		m.MasterIns.AddTask(t)
	}
	if m.SlaveIns != nil {
		m.SlaveIns.AddTask(t)
	}
}
