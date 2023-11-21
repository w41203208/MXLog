package xlog

import (
	writerCore "xlog/pkg/log/WriterCore"
	"xlog/pkg/log/encoder"
	message "xlog/pkg/log/message"
)

type MessageWriter struct {
	index       int
	added       bool
	writerCores []writerCore.WriterCore
	enc         encoder.IEncoder
}

func (me *MessageWriter) AddWriterCore(core writerCore.WriterCore) {
	me.writerCores = append(me.writerCores, core)
}

func (me *MessageWriter) Write(message message.IMessage) {
	buf := me.enc.Encode(message)

	// add ErrorReceiver
	for _, core := range me.writerCores {
		core.Write(buf.Bytes())
	}
}
