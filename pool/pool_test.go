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

	pool := New(1)

	for i := 0; i < 11; i++ {
		fmt.Println("set: ", i)
		pool.Submit(func() any {
			sleep := rand.Intn(10)
			//time.Sleep(time.Duration(sleep) * time.Second)
			fmt.Println("running job ===> ", i, sleep)
			return nil
		})
		fmt.Println("starting goroutines:", runtime.NumGoroutine())
	}

	defer func() {
		fmt.Println("closing goroutines:", runtime.NumGoroutine())
		pool.Close()
		time.Sleep(time.Second)
		fmt.Println("end goroutines:", runtime.NumGoroutine())
	}()
}
