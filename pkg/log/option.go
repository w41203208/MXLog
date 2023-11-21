package xlog

import (
	"xlog/pkg/log/level"
)

type XOption interface {
	apply(x *XLog)
}

type XOptionFunc func(x *XLog)

func (xf XOptionFunc) apply(x *XLog) {
	xf(x)
}

func SetLevel(level level.XLevel) XOption {
	return XOptionFunc(func(x *XLog) {
		x.level = level
	})
}
func SetCodeDetail(t bool) XOption {
	return XOptionFunc(func(x *XLog) {
		if t {
			x.formatFlag |= XLogDetail
		} else {
			x.formatFlag &= ^XLogDetail
		}
	})
}

func SetLogDate(t bool) XOption {
	return XOptionFunc(func(x *XLog) {
		if t {
			x.formatFlag |= XLogDetail
		} else {
			x.formatFlag &= ^XLogDate
		}
	})
}

func SetLogTime(t bool) XOption {
	return XOptionFunc(func(x *XLog) {
		if t {
			x.formatFlag |= XLogTime
		} else {
			x.formatFlag &= ^XLogTime
		}
	})
}
