package message

import (
	"xlog/pkg/log/level"
)

type IMessage interface {
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

func (m *XMessage) SetFields() {

}
