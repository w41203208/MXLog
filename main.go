package main

import (
	"strconv"
	xlog "xlog/pkg/log"
)

func main() {
	xlog := xlog.NewXLog(nil, xlog.SetCodeDetail(true))

	xlog.LogTrace("test1testtest")
	for i := 0; i < 5; i++ {
		go func(a int) {
			xlog.LogTrace("test1" + strconv.Itoa(a))

		}(i)
	}

	xlog.LogTrace("test1testtest1")

	xlog.Start()
}
