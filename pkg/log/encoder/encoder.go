package encoder

import (
	"bytes"
	"strconv"
	"xlog/pkg/log/message"
)

type IEncoder interface {
	Encode(msg message.IMessage) *bytes.Buffer
}

type XEncoder struct {
}

func (enc *XEncoder) Encode(iMsg message.IMessage) *bytes.Buffer {
	buf := new(bytes.Buffer)

	switch msg := iMsg.(type) {
	case *message.XMessage:
		// write level header
		buf.WriteString("[")
		buf.Write([]byte(msg.Level.String()))
		buf.WriteString("]")
		buf.WriteString(" ")

		// write code detail
		if msg.MCaller != nil {
			buf.WriteString(msg.MCaller.File)
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(msg.MCaller.Line))
			buf.WriteString(" ")
		}

		buf.Write([]byte(msg.Body))
		buf.WriteString("\n")
		return buf
	default:
		return nil
	}
}
