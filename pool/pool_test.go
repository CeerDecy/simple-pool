package pool

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	pool := New(Options{
		PoolLimit:       100,
		WorkerStopCount: 2,
		WorkerWaitTime:  time.Millisecond * 500,
	})

	for i := 0; i < 1000; i++ {
		fmt.Println("set: ", i)
		pool.Submit(func() any {
			sleep := rand.Intn(10)
			time.Sleep(time.Duration(sleep) * time.Second)
			fmt.Println("running job ===> ", i, sleep)
			return nil
		})
		fmt.Println("starting goroutines:", runtime.NumGoroutine())
		if i == 10 {
			fmt.Println("1010101010101010101010101010")
			time.Sleep(10 * time.Second)
			fmt.Println("after sleep goroutines:", runtime.NumGoroutine())
		}
	}

	defer func() {
		fmt.Println("closing goroutines:", runtime.NumGoroutine())
		pool.Close()
		//time.Sleep(time.Second)
		fmt.Println("end goroutines:", runtime.NumGoroutine())
	}()
}

func TestChannel(t *testing.T) {
	job := make(chan int, 10)
	fmt.Println(len(job), cap(job))
}
