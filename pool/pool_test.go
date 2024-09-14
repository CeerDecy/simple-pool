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

	pool := New(5)

	go monitor()

	start := time.Now()
	for i := 0; i < 100; i++ {
		//fmt.Println("set: ", i)
		pool.Submit(func() any {
			sleep := rand.Intn(3)
			time.Sleep(time.Duration(sleep) * time.Second)
			//fmt.Println("running job ===> ", i, sleep)
			return nil
		})
	}

	defer func() {
		pool.Close()
		fmt.Println("use time: ", time.Now().Sub(start), runtime.NumGoroutine())
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
