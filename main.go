package main

import (
	"strconv"
	"sync"
	xlog "xlog/pkg/log"
)

func main() {
	xlog := xlog.NewXLog(nil, xlog.SetCodeDetail(true))
	wait := &sync.WaitGroup{}

	xlog.LogTrace("test1testtest")
	for i := 0; i < 2; i++ {
		wait.Add(1)
		go func(a int) {
			defer wait.Done()
			xlog.LogTrace("test1" + strconv.Itoa(a))

		}(i)
	}

	xlog.LogTrace("test1testtest1")
	wait.Wait()
}
