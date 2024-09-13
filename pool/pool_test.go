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
		PoolLimit:       5000,
		WorkerStopCount: 2,
		WorkerWaitTime:  time.Millisecond * 500,
	})

	go monitor()

	for i := 0; i < 30000; i++ {
		//fmt.Println("set: ", i)
		pool.Submit(func() any {
			sleep := rand.Intn(10)
			time.Sleep(time.Duration(sleep) * time.Second)
			//fmt.Println("running job ===> ", i, sleep)
			return nil
		})
		//fmt.Println("starting goroutines:", runtime.NumGoroutine())
		if i == 15000 {
			//fmt.Println("1010101010101010101010101010")
			time.Sleep(10 * time.Second)
			//fmt.Println("after sleep goroutines:", runtime.NumGoroutine())
		}
	}

	defer func() {
		//fmt.Println("closing goroutines:", runtime.NumGoroutine())
		pool.Close()
		//time.Sleep(time.Second)
		//fmt.Println("end goroutines:", runtime.NumGoroutine())
	}()
}

func monitor() {
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)
		fmt.Printf("Alloc: %v KB, NumGoroutine: %d\n", memStats.Alloc/1024, runtime.NumGoroutine())
		time.Sleep(time.Second * 1)
	}
}

func TestChannel(t *testing.T) {
	job := make(chan int, 10)
	fmt.Println(len(job), cap(job))
}
