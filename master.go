// TODO : Split broker and Master into two.
package metre

import (
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gospackler/metre/logging"
	"github.com/gospackler/metre/transport"
	"github.com/robfig/cron"
	uuid "github.com/satori/go.uuid"
)

var metreMasterZap = zap.String("metre.type", "master")

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
			logging.Logger.Error("request connection error in master",
				zap.Error(err),
			)
		}
		reqConnChan <- reqConn
	}
	go func() {
		for {
			logging.Logger.Debug("master's scheduler listening for messages")
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
		logging.Logger.Warn("attempted to add two tasks with the same ID",
			zap.String("id", t.ID),
		)
	}

	m.TaskMap[id] = t
	logging.Logger.Debug("adding task to master",
		zap.String("id", t.ID),
	)
	m.Cron.AddFunc(t.Interval, func() {
		go m.ScheduleFromId(id)
	})
}

func (m *Master) ScheduleFromId(ID string) error {
	logging.Logger.Debug("scheduling task",
		zap.String("id", ID),
	)
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
	logging.Logger.Info("metre master scheduling info",
		metreMasterZap,
		zap.String("message", message),
		zap.String("task_id", payload.TaskId),
		zap.String("payload_uid", payload.UID),
		zap.String("job_guid", payload.JobGUID),
	)
}

func (m *Master) logScheduleErr(message string, payload *MetreMessage) {
	logging.Logger.Error("metre master scheduling error",
		metreMasterZap,
		zap.String("message", message),
		zap.String("task_id", payload.TaskId),
		zap.String("payload_uid", payload.UID),
		zap.String("job_guid", payload.JobGUID),
	)
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
