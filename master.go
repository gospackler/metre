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

type Master struct {
	Cron        cron.Cron
	SchInpChan  chan *SchInput
	TaskMap     map[string]*Task
	schParallel int
	routerUri   string
}

// Remove dealer from the picture in Master.
func NewMaster(routerUri string, schParallel int) (*Master, error) {
	if routerUri == "" {
		routerUri = LocalHost + ":" + RouterPort
	} else if strings.Index(routerUri, ":") == 0 {
		routerUri = LocalHost + ":" + RouterPort
	}

	cron := *cron.New()

	schInpChan := make(chan *SchInput)

	master := &Master{
		Cron:        cron,
		routerUri:   routerUri,
		SchInpChan:  schInpChan,
		TaskMap:     make(map[string]*Task),
		schParallel: schParallel,
	}
	return master, nil
}

func (m *Master) Start() {
	m.Cron.Start()
	reqConnChan := make(chan *transport.ReqConn, m.schParallel)
	for i := 0; i < m.schParallel; i++ {
		reqConn, err := transport.NewReqConn(m.routerUri)
		if err != nil {
			log.Error("Request connection error " + err.Error())
		}
		reqConnChan <- reqConn
	}
	go func() {
		for {
			log.Info("Masters scheduler listening for messsages")
			s := <-m.SchInpChan
			// Add go routine for parallel
			r := <-reqConnChan
			go func(schInp *SchInput, reqConn *transport.ReqConn) {
				resp, err := reqConn.MakeReq(schInp.Input)
				if err != nil {
					schInp.ErrorChan <- err
				}
				schInp.RespChan <- resp
				reqConnChan <- reqConn
			}(s, r)
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

func (m *Master) Schedule(taskId string, uid string) (string, error) {
	msg := CreateMsg(Request, taskId, uid, "")
	req := NewScheduleInput(msg)
	m.SchInpChan <- req
	defer req.Close()
	select {
	case resp := <-req.RespChan:
		return resp, nil
	case err := <-req.ErrorChan:
		return "", err
	}
	return "", nil
}
