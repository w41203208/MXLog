package xlog

import (
	"os"
	"runtime"
	"sync"
	writerCore "xlog/pkg/log/WriterCore"
	"xlog/pkg/log/encoder"
	"xlog/pkg/log/level"
	"xlog/pkg/log/message"
	"xlog/pkg/log/pool"
)

type XLog struct {
	rwMu        sync.RWMutex
	closeChan   chan (struct{})
	level       level.XLevel
	formatFlag  XLogFormat
	WriterCores []writerCore.WriterCore
	enc         encoder.IEncoder

	messageList   []message.XMessage
	msgWriterPool *pool.Pool[*MessageWriter]
	newEncFn      func() encoder.IEncoder
}

type XLogFormat uint8

const (
	XLogDetail XLogFormat = 1 << iota
	XLogDate
	XLogTime
)

func NewXLog(newEncFn func() encoder.IEncoder, opts ...XOption) *XLog {

	xl := &XLog{
		level:      level.TraceLevel,
		formatFlag: XLogDate | XLogTime,
	}

	if fn := newEncFn; fn != nil {
		xl.enc = fn()
	} else {
		xl.enc = &encoder.XEncoder{}
	}

	xl.msgWriterPool = pool.NewPool[*MessageWriter](func() *MessageWriter {
		return &MessageWriter{
			enc: xl.enc,
		}
	})

	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(xl)
		}
	}

	xl.AddWriterCore(writerCore.NewLocalWriter(os.Stderr))

	// go xl.StartToWrite()
	return xl
}

func (xl *XLog) StartToWrite() {
	for {
		for len(xl.messageList) != 0 {
			xl.rwMu.RLock()

			xl.rwMu.RUnlock()
		}
	}
}

func (xl *XLog) AddWriterCore(wc writerCore.WriterCore) {
	// maybe need to use a struct to store wc
	xl.WriterCores = append(xl.WriterCores, wc)
}

func (xl *XLog) LogTrace(msg string, Fields ...interface{}) {
	if ent, m := xl.logCheck(level.TraceLevel, msg); ent != nil {
		ent.Write(m)
	}
}

func (xl *XLog) logCheck(level level.XLevel, msg string) (*MessageWriter, *message.XMessage) {
	if xl.level > level {
		return nil, nil
	}

	// step 1 new message

	//temp use this to take stacktrace
	var MsgCaller *message.MessageCaller
	if xl.formatFlag&XLogDetail != 0 {
		var pcs []uintptr = []uintptr{10}
		numFrames := runtime.Callers(3, pcs)
		pc := pcs[:numFrames]
		frames := runtime.CallersFrames(pc)
		frame, _ := frames.Next()

		MsgCaller = &message.MessageCaller{
			File: frame.File,
			Line: frame.Line,
		}
	}

	// add other attribute for message
	m := &message.XMessage{
		Level:   xl.level,
		Body:    msg,
		MCaller: MsgCaller,
	}

	// step 2 get MessageWriter
	// maybe need to use sync.Pool to store entry
	mWriter := xl.msgWriterPool.Get()

	// step 3 add core in MessageWriter
	if len(xl.WriterCores) > 0 {
		for _, core := range xl.WriterCores {
			if core.Check(m.Level) {
				mWriter.AddWriterCore(core)
			}
		}
	}

	return mWriter, m
}
