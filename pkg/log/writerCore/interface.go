package writerCore

import "xlog/pkg/log/level"

type WriterCore interface {
	Check(lvl level.XLevel) bool
	Write(p []byte)
}
