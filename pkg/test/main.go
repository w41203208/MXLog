package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	MingLog "test/pkg/log"
)

func main() {

	finalClosedChan := make(chan struct{})

	go func() {
		MingLog.Receive()
	}()
	// MingLog.Println("test1")

	// MingLog.Println("test1")
	// MingLog.Println("test1")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		if sig.String() == "interrupt" {
			log.Println("Program terminated")
			func() {
				MingLog.Stop()

				close(finalClosedChan)
			}()
		}
	}()

	for i := 0; i < 10; i++ {
		go func(a int) {
			MingLog.Println("test ", a)
		}(i)
	}

	// for c := 0; c < 5; c++ {
	// 	go func(a int) {
	// 		log.Println("test ", a)
	// 	}(c)
	// }

	<-finalClosedChan
}

// // func Foo() {
// // 	fmt.Printf("我是 %s, %s 在调用我!\n", printMyName(), printCallerName())
// // 	Bar()
// // }
// // func Bar() {
// // 	fmt.Printf("我是 %s, %s 又在调用我!\n", printMyName(), printCallerName())
// // 	trace()
// // }
// // func printMyName() string {
// // 	pc, file, line, ok := runtime.Caller(1)
// // 	fmt.Println(file)
// // 	fmt.Println(line)
// // 	fmt.Println(ok)
// // 	return runtime.FuncForPC(pc).Name()
// // }
// // func printCallerName() string {
// // 	pc, _, _, _ := runtime.Caller(2)
// // 	return runtime.FuncForPC(pc).Name()
// // }

// // func trace() {
// // 	pc := make([]uintptr, 10) // at least 1 entry needed
// // 	n := runtime.Callers(0, pc)
// // 	for i := 0; i < n; i++ {
// // 		f := runtime.FuncForPC(pc[i])
// // 		file, line := f.FileLine(pc[i])
// // 		fmt.Printf("%s:%d %s\n", file, line, f.Name())
// // 	}
// // }

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func main() {
// 	ch := make(chan []byte)
// 	var mutex sync.Mutex

// 	wait := &sync.WaitGroup{}
// 	// go func() {
// 	// 	wait.Add(1)

// 	// 	defer wait.Done()
// 	// 	for i := 0; i < 3; i++ {
// 	// 		data := []byte{byte(i)}
// 	// 		mutex.Lock()
// 	// 		ch <- data
// 	// 		mutex.Unlock()
// 	// 	}
// 	// }()
// 	for i := 0; i < 3; i++ {
// 		wait.Add(1)
// 		go func(a int) {
// 			defer func() {
// 				wait.Done()
// 				mutex.Unlock()
// 			}()
// 			data := []byte{byte(a)}
// 			mutex.Lock()
// 			ch <- data

// 		}(i)
// 	}

// 	go func() {
// 		wait.Add(1)

// 		defer wait.Done()
// 		for {
// 			timeChan := time.After(2 * time.Second)
// 			select {
// 			case data := <-ch:
// 				fmt.Println(data)
// 				_ = <-timeChan
// 			case <-timeChan:
// 				fmt.Println("每次請求的2秒已過")
// 				return
// 			}
// 		}
// 	}()
// 	wait.Wait()
// }
