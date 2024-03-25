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

type Queue[T any] struct {
	rwMu sync.RWMutex
	list []T
}

func (q *Queue[T]) Push(x T) {
	q.rwMu.Lock()
	q.list = append(q.list, x)
	q.rwMu.Unlock()
}

func (q *Queue[T]) Pop() T {
	len := len(q.list)
	var el T
	q.rwMu.RLock()
	el, q.list = q.list[0], q.list[1:len]
	q.rwMu.RUnlock()
	return el
}

func (q *Queue[T]) Length() int {
	return len(q.list)
}

type XLog struct {
	closeChan   chan (struct{})
	level       level.XLevel
	formatFlag  XLogFormat
	WriterCores []writerCore.WriterCore
	enc         encoder.IEncoder

	messageQueue  *Queue[message.IMessage]
	msgWriterPool *pool.Pool[*MessageWriter]
	newEncFn      func() encoder.IEncoder
}

type XLogFormat uint8

const (
	XLogDetail XLogFormat = 1 << iota
	XLogDate
	XLogTime
)

var mwIndex int = 0

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
	xl.messageQueue = &Queue[message.IMessage]{}

	xl.msgWriterPool = pool.NewPool[*MessageWriter](func() *MessageWriter {
		mwIndex++
		return &MessageWriter{
			added: false,
			index: mwIndex,
			enc:   xl.enc,
		}
	})

	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(xl)
		}
	}

	xl.AddWriterCore(writerCore.NewLocalWriter(os.Stderr))

	return xl
}

func (xl *XLog) Stop() {
	close(xl.closeChan)
}
func (xl *XLog) Start() {
	go xl.StartToWrite()

	<-xl.closeChan
}

func (xl *XLog) StartToWrite() {
	for {
		for xl.messageQueue.Length() != 0 {
			msg := xl.messageQueue.Pop()
			mWriter := xl.msgWriterPool.Get()

			// set core in writer
			if !mWriter.added {
				if len(xl.WriterCores) > 0 {
					for _, core := range xl.WriterCores {
						if core.Check(msg.GetLevel()) {
							mWriter.AddWriterCore(core)
						}
					}
				}
				mWriter.added = true
			}

			mWriter.Write(msg)

			// complete to use, need to put back to pool
			xl.msgWriterPool.Put(mWriter)
		}
	}
}

func (xl *XLog) AddWriterCore(wc writerCore.WriterCore) {
	// maybe need to use a struct to store wc
	xl.WriterCores = append(xl.WriterCores, wc)
}

func (xl *XLog) LogTrace(msg string, Fields ...interface{}) {
	xl.logCheck(level.TraceLevel, msg)
}

// check it can create XMessage
func (xl *XLog) logCheck(level level.XLevel, msg string) {
	if xl.level > level {
		return
	}

	// maybe need to be a object to control
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

	xl.messageQueue.Push(m)
}
