package bbr

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/notes/crtos/pkg/ratelimit"
)

func TestBBR(t *testing.T) {
	cfg := &Config{
		Window:       time.Second * 5,
		WinBucket:    50,
		CPUThreshold: 100,
	}
	limiter := newLimiter(cfg)
	var wg sync.WaitGroup
	var drop int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 300; i++ {
				f, err := limiter.Allow(context.TODO())
				if err != nil {
					atomic.AddInt64(&drop, 1)
				} else {
					count := rand.Intn(100)
					time.Sleep(time.Microsecond * time.Duration(count))
					f(ratelimit.DoneInfo{Op: ratelimit.Success})
				}
			}
		}()
	}
	wg.Wait()
	fmt.Println("drop:", drop)
}
