// TODO : Split broker and Master into two.
package metre

import (
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gospackler/metre/transport"
	"github.com/robfig/cron"
	uuid "github.com/satori/go.uuid"
)

// Master is the representation of a master server.
// It can register tasks which gets scheduled periodically.
// Multiple master instances can connect to a broker.
// The schParallel defines the number of request connections to the broker.
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
			log.Debug("Masters scheduler listening for messsages")
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
	log.Debug("Adding cron task to Master" + t.ID)
	m.Cron.AddFunc(t.Interval, func() {
		log.Debug("Cron triggered")
		go m.ScheduleFromId(id)
	})
}

func (m *Master) ScheduleFromId(ID string) error {
	log.Debug("About to schedule ID " + ID)
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

func (m *Master) logScheduleInfo(message string, payload *MetreMessage) {
	log.Info(fmt.Sprintf("metre.master.Schedule: %s - taskId:%s, uid:%s, jobGUID:%s",
		message,
		payload.TaskId,
		payload.UID,
		payload.JobGUID))
}

func (m *Master) logScheduleErr(message string, payload *MetreMessage) {
	log.Error(fmt.Sprintf("metre.master.Schedule: %s - taskId:%s, uid:%s, jobGUID:%s",
		message,
		payload.TaskId,
		payload.UID,
		payload.JobGUID))
}

// FIXME : Scheduling happens here. Is this the right design ?
func (m *Master) Schedule(taskId string, uid string) (string, error) {
	// jobGUID is a low-level tracking guid at the metre level.  Tasks can also have their
	// own UID for tracking which is specific to the app.
	jobGUID := uuid.NewV4().String()

	msg := CreateMsg(Request, taskId, uid, jobGUID, "")
	req := NewScheduleInput(SerializeMsg(msg))
	m.logScheduleInfo("About to schedule job", msg)
	m.SchInpChan <- req
	m.logScheduleInfo("Job scheduled", msg)
	defer req.Close()
	select {
	case resp := <-req.RespChan:
		cleanResp, err := CleanResponseMessage(resp)
		if err != nil {
			m.logScheduleErr("Failed to clean response", msg)
			return "", err
		}
		m.logScheduleInfo("Received slave successful response", msg)
		return cleanResp, nil
	case err := <-req.ErrorChan:
		m.logScheduleErr("Received slave error response", msg)
		return "", err
	}
}
