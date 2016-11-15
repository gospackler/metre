package metre

import (
	"encoding/json"
)

type msgType int

const (
	Status msgType = iota
	Request
	Debug
	Error
)

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
