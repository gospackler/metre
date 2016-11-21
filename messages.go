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
	MessageType msgType `json:"msg_type"`
	TaskId      string  `json:"task_id"`
	UID         string  `json:"uid"`
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

func CreateMsg(mt msgType, id, uid, msg string) string {
	msgObj := &MetreMessage{
		MessageType: mt,
		TaskId:      id,
		UID:         uid,
		Message:     msg,
	}
	byteMsg, _ := json.Marshal(msgObj)
	return string(byteMsg)
}

func CreateErrorMsg(err error, id string, uid string) string {
	return CreateMsg(Error, id, uid, err.Error())
}
