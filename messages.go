package metre

import (
	"encoding/json"
)

type msgType int

const (
	Status msgType = iota
	Debug
	Error
)

type trackMessage struct {
	MessageType msgType `json:"msg_type"`
	TaskId      string  `json:"task_id"`
	UID         string  `json:"uid"`
	Message     string  `json:"message"`
}

func parseMessage(msg string) (*trackMessage, error) {
	trackMsg := new(trackMessage)
	err := json.Unmarshal([]byte(msg), trackMsg)
	return trackMsg, err
}

func createMsg(mt msgType, id, uid, msg string) string {
	msgObj := &trackMessage{
		MessageType: mt,
		TaskId:      id,
		UID:         uid,
		Message:     msg,
	}
	byteMsg, _ := json.Marshal(msgObj)
	return string(byteMsg)
}
