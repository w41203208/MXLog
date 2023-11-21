package message

import (
	"xlog/pkg/log/level"
)

type IMessage interface {
	GetLevel() level.XLevel
	SetFields()
}

type MessageCaller struct {
	File string
	Line int
}

type XMessage struct {
	Body    string
	Level   level.XLevel
	MCaller *MessageCaller
}

func (m *XMessage) GetLevel() level.XLevel {
	return m.Level
}

func (m *XMessage) SetFields() {

}
