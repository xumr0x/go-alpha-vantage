package av

import (
	"math"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

const (
	DefaultDayLimit    = math.MaxInt32
	DefaultSecondLimit = math.MaxInt32
)

var ErrDailyLimitReached = errors.New("daily API limit has been reached")

// RateLimiter limits the per-second and per-day execution counts.
//
// It delays execution to comply with API restrictions (i.e. 5 calls per second).
//
// Usage
// 	rl := NewRateLimiter(500, 5) // 500 calls per day, 5 calls per second
//	rl.Do(funcToExecute())
type RateLimiter struct {
	secLimit int32
	secCount int32
	dayLimit int32
	dayCount int32
}

func NewRateLimiter(dayLimit int, secLimit int) *RateLimiter {
	if dayLimit == 0 {
		dayLimit = DefaultDayLimit
	}
	if secLimit == 0 {
		dayLimit = DefaultSecondLimit
	}

	l := &RateLimiter{
		secLimit: int32(secLimit),
		secCount: 0,
		dayLimit: int32(dayLimit),
		dayCount: 0,
	}

	l.init()

	return l
}

func (l *RateLimiter) init() {
	secTicker := time.NewTicker(time.Second)
	dayTicker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-secTicker.C:
				// Reset the current per second count.
				atomic.StoreInt32(&l.secCount, 0)
			case <-dayTicker.C:
				// Reset the current per day count.
				atomic.StoreInt32(&l.dayCount, 0)
			}
		}
	}()
}

// Do executes the given function.
//
// It will delays execution by 50ms steps if the per-second
// limit has been reached.
func (l *RateLimiter) Do(f func() (*http.Response, error)) (*http.Response, error) {
	if atomic.LoadInt32(&l.dayCount) == l.dayLimit {
		return nil, ErrDailyLimitReached
	}

	// Delay until the count is reset.
	for atomic.LoadInt32(&l.secCount) == l.secLimit {
		time.Sleep(50 * time.Millisecond)
	}

	// Execute function and increment count.
	res, err := f()
	atomic.AddInt32(&l.secCount, 1)
	atomic.AddInt32(&l.dayCount, 1)

	return res, err
}
