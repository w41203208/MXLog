package writerCore

import (
	"io"
	"xlog/pkg/log/level"
)

type LocalWriter struct {
	level level.XLevel
	w     io.Writer
}

func NewLocalWriter(out io.Writer) *LocalWriter {
	return &LocalWriter{
		w: out,
	}
}

func (lw *LocalWriter) Check(lvl level.XLevel) bool {
	return lw.level.Enable(lvl)
}

func (lw *LocalWriter) Write(p []byte) {
	lw.write(p)
}

func (lw *LocalWriter) write(p []byte) {
	lw.w.Write(p)
}
