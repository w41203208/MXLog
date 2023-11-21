package MingLog

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	WARN  = "warn"
	INFO  = "info"
	DEBUG = "debug"
)

const (
	API = 1 << iota
	LOCAL
	FILE
	LstdWhere = LOCAL | FILE
)

const (
	Lshortfile = 1 << iota
	LstdFlags  = Lshortfile
)

type MingLog struct {
	mu         sync.Mutex
	prefix     string
	out        io.Writer
	formatFlag int
	outputFlag int
	buf        []byte
	dataChan   chan []byte
	closeChan  chan struct{}
	wait       sync.WaitGroup
}

func New(out io.Writer, prefix string, fFlag int, oFlag int) *MingLog {
	var log = &MingLog{out: out, prefix: prefix, formatFlag: fFlag, outputFlag: oFlag}
	log.closeChan = make(chan struct{})
	log.dataChan = make(chan []byte, 10)
	return log
}

func Receive() {
	mLog.Receive()
}

func (ml *MingLog) Receive() {
	go func() {
		for {
			select {
			case data := <-ml.dataChan:
				ml.execute(data)
				ml.mu.Unlock()
			case <-ml.closeChan:
				return
			}
		}
	}()
}

func Stop() {
	mLog.Stop()
}
func (ml *MingLog) Stop() {
	close(ml.closeChan)
}

var mLog = New(os.Stderr, "", LstdFlags, LstdWhere)

func Println(v ...interface{}) {
	mLog.output(2, fmt.Sprintln(v...))
}
func (ml *MingLog) execute(data []byte) error {

	if ml.outputFlag&LOCAL != 0 {
		go func(data []byte) {
			_ = ml.consoleLocal(data)
		}(data)
	}

	if ml.outputFlag&FILE != 0 {
		go func(data []byte) {
			_ = ml.consoleFile(data)
		}(data)
	}

	// if err != nil {
	// 	return err
	// }
	return nil
}

func (ml *MingLog) SetPrefix(prefix string) {
	ml.mu.Lock()
	defer ml.mu.Lock()
	ml.prefix = prefix
}

func (ml *MingLog) mItoa(buf *[]byte, n int, wid int) {
	var b [20]byte
	bPointer := len(b) - 1
	for n >= 10 || wid > 1 {
		wid--
		mod := n % 10
		b[bPointer] = byte('0' + mod)
		n = (n - mod) / 10
		bPointer--
	}
	b[bPointer] = byte('0' + n)
	*buf = append(*buf, b[bPointer:]...)
}

// prefix, time, file, line
func (ml *MingLog) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	// add prefix
	*buf = append(*buf, ml.prefix...)

	// add year / month / day
	year, month, day := t.Date()
	ml.mItoa(buf, year, 4)
	*buf = append(*buf, '/')
	ml.mItoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	ml.mItoa(buf, day, 2)
	*buf = append(*buf, ' ')

	// add hour / min / second
	hour, min, sec := t.Clock()
	// depend on os platform
	// It is on windows, hour doesn't need to be added additional 8 hours
	// It is on linux, hour need to be added additional 8 hours
	if runtime.GOOS == "windows" {
		hour += 0
	}
	if runtime.GOOS == "linux" {
		hour += 8
	}
	ml.mItoa(buf, hour, 2)
	*buf = append(*buf, ':')
	ml.mItoa(buf, min, 2)
	*buf = append(*buf, ':')
	ml.mItoa(buf, sec, 2)
	*buf = append(*buf, ' ')

	if ml.formatFlag&Lshortfile != 0 {
		log.Println("testtesttest")
		var arr = strings.Split(file, "/")
		var last = arr[len(arr)-1]

		*buf = append(*buf, last...)
		*buf = append(*buf, ':')
		ml.mItoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
	// add time
	// *buf = append(*buf, t.Local())
}

func (ml *MingLog) output(callerDepth int, s string) error {
	now := time.Now()
	var file string
	var line int
	ml.mu.Lock()

	ml.mu.Unlock()
	var ok bool
	_, file, line, ok = runtime.Caller(callerDepth)
	if !ok {
		file = "???"
		line = 0
	}
	ml.mu.Lock()

	ml.buf = ml.buf[:0]
	ml.formatHeader(&ml.buf, now, file, line)

	ml.buf = append(ml.buf, s...)

	fmt.Println("Before: ", string(ml.buf))
	// ml.consoleLocal(ml.buf)
	ml.dataChan <- ml.buf

	return nil
}

func (ml *MingLog) consoleLocal(b []byte) error {
	_, err := ml.out.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (ml *MingLog) consoleApi(b []byte) {
	ml.out.Write(b)
}

func (ml *MingLog) consoleFile(b []byte) error {
	file, err := os.OpenFile("./text.txt", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
<<<<<<< HEAD
	w := bufio.NewWriter(file)
	w.Write(b)
	w.Flush()
=======

	w := bufio.NewWriter(file)
	w.Write(b)
	w.Flush()
	// ml.out.Write(b)
>>>>>>> 8bdabcf0df840b4eab25fa11c6370e09eeb5a9b0

	return nil
}

func SetPrefix(prefix string) {
	mLog.SetPrefix(prefix)
}
