package metre

import (
	"encoding/json"
	"errors"
)

type msgType int

const (
	Status msgType = iota
	Request
	Debug
	Error
)

// The set of messages metre passes between master and slave.
type MetreMessage struct {
	// JobGUID is the low-level unique id that metre will track between master and slaves.
	JobGUID string `json:"job_guid"`
	// UID is the application level unique id that the app can track for it's own purpose.
	UID         string  `json:"uid"`
	MessageType msgType `json:"msg_type"`
	TaskId      string  `json:"task_id"`
	Message     string  `json:"message"`
}

func ParseMessage(msg string) (*MetreMessage, error) {
	metMsg := new(MetreMessage)
	err := json.Unmarshal([]byte(msg), metMsg)
	return metMsg, err
}

func CleanResponseMessage(msg string) (string, error) {
	metMsg, err := ParseMessage(msg)
	if err != nil {
		return "", err
	}
	var replyMsg string
	switch metMsg.MessageType {
	case Status:
		replyMsg = metMsg.Message
	case Debug:
		replyMsg = metMsg.Message
	case Error:
		err = errors.New(metMsg.Message)
	case Request:
		err = errors.New("Response Message should not contain a Request type.")
	default:
		err = errors.New("Some message type needed in the response.")
	}
	return replyMsg, err
}

func CreateMsg(mt msgType, id, uid, jobguid, msg string) *MetreMessage {
	return &MetreMessage{
		JobGUID:     jobguid,
		MessageType: mt,
		TaskId:      id,
		UID:         uid,
		Message:     msg,
	}
}

func SerializeMsg(msgObj *MetreMessage) string {
	byteMsg, _ := json.Marshal(msgObj)
	return string(byteMsg)
}

func CreateErrorMsg(err error, id, uid, jobguid string) string {
	return SerializeMsg(CreateMsg(Error, id, uid, jobguid, err.Error()))
}
