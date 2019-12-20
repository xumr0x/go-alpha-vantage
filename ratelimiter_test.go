package av

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRateLimiter_Do(t *testing.T) {
	tests := []struct {
		desc   string
		perSec int
		perDay int
		calls  int
		err    error
	}{
		{
			desc:   "1 call per second, no daily limit",
			perSec: 1,
			perDay: 0,
			calls:  10,
		},
		{
			desc:   "10 calls per second, no daily limit",
			perSec: 10,
			perDay: 0,
			calls:  50,
		},
		{
			desc:   "1 call per second, over daily limit",
			perSec: 1,
			perDay: 1,
			calls:  3,
			err:    ErrDailyLimitReached,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			rl := NewRateLimiter(tt.perDay, tt.perSec)
			ticker := time.NewTicker(time.Second)
			endTicker := time.NewTicker(time.Duration(tt.calls/tt.perSec) * time.Second)
			go func() {
				for {
					select {
					case <-ticker.C:
						if count := rl.secCount; count > int32(tt.perSec) {
							t.Fatalf("too many calls: %+v", count)
						}
					case <-endTicker.C:
						ticker.Stop()
						endTicker.Stop()
					}
				}
			}()

			for i := 0; i < tt.calls; i++ {
				if err, _ := rl.Do(func() (*http.Response, error) { return nil, nil }); err != nil {
					if !reflect.DeepEqual(err, tt.err) {
						t.Errorf("unexpected error: %+v", err)
						return
					}
				}
			}
		})
	}
}
