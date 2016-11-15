// TODO : Split broker and Master into two.
package metre

import (
	"errors"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gospackler/metre/transport"
	"github.com/robfig/cron"
)

type SchInput struct {
	Input     string
	RespChan  chan string
	ErrorChan chan error
}

func NewScheduleInput(inp string) *SchInput {
	return &SchInput{
		Input:     inp,
		RespChan:  make(chan string),
		ErrorChan: make(chan error),
	}
}

func (s *SchInput) Close() {
	close(s.RespChan)
	close(s.ErrorChan)
}

type Master struct {
	Cron       cron.Cron
	SchInpChan chan *SchInput
	ReqConn    *transport.ReqConn
	TaskMap    map[string]*Task
}

// Remove dealer from the picture in Master.
func NewMaster(routerUri string, schParallel int) (*Master, error) {
	if routerUri == "" {
		routerUri = LocalHost + ":" + RouterPort
	} else if strings.Index(routerUri, ":") == 0 {
		routerUri = LocalHost + ":" + RouterPort
	}

	reqConn, err := transport.NewReqConn(routerUri)
	if err != nil {
		return nil, err
	}
	cron := *cron.New()

	schInpChan := make(chan *SchInput, schParallel)

	master := &Master{
		Cron:       cron,
		ReqConn:    reqConn,
		SchInpChan: schInpChan,
		TaskMap:    make(map[string]*Task),
	}
	return master, nil
}

func (m *Master) Start() {
	m.Cron.Start()
	go func() {
		for {
			log.Info("Masters scheduler listening for messsages")
			schInp := <-m.SchInpChan
			resp, err := m.ReqConn.MakeReq(schInp.Input)
			if err != nil {
				schInp.ErrorChan <- err
			}
			schInp.RespChan <- resp
		}
	}()
}

func (m *Master) AddTask(t *Task) {
	id := t.GetID()
	if _, exists := m.TaskMap[id]; exists {
		log.Warn("attempted to add two tasks with the same ID [" + t.ID + "]")
	}

	m.TaskMap[id] = t
	log.Info("Adding cron task to Master" + t.ID)
	m.Cron.AddFunc(t.Interval, func() {
		log.Info("Cron triggered")
		go m.scheduleFromId(id)
	})
}

func (m *Master) scheduleFromId(ID string) error {
	log.Info("About to schedule ID " + ID)
	var t *Task
	var ok bool
	if t, ok = m.TaskMap[ID]; !ok {
		return errors.New("task [" + ID + "] not recognized")
	}

	// Making sure the next run is not affected by previous runs.
	t.Lock()
	t.Zero()
	t.SendMessage(t.ID + ": Scheduled")
	t.Schedule(m)
	t.SendMessage(t.ID + ": Time taken " + time.Since(t.StartTime).String())
	t.Unlock()
	return nil
}
