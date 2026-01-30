package testutils

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type retryParams struct {
	maxTries int
	sleep    time.Duration
}

func NewRetry(maxTries int, sleep time.Duration) retryParams {
	p := retryParams{maxTries: maxTries, sleep: sleep}

	return p
}

func (p retryParams) Retry(t *testing.T, test func(*testing.T)) {
	for i := 0; i < p.maxTries; i++ {
		success := t.Run(fmt.Sprintf("attempt %v of %v\n", i+1, p.maxTries), func(t *testing.T) {
			test(t)
		})

		if success {
			return
		}

		if i == p.maxTries-1 {
			p := reflect.ValueOf(test).Pointer()
			f := runtime.FuncForPC(p)

			t.Errorf("test %v failed", f.Name())

			return
		}

		if i < p.maxTries-1 {
			pc := reflect.ValueOf(test).Pointer()
			f := runtime.FuncForPC(pc)

			t.Logf("test failed %v, retrying\n", f.Name())
			sleep := p.sleep * time.Duration(math.Pow(2, float64(i)))
			t.Logf("sleeping for %v\b", sleep)

			select {
			case <-time.After(sleep):
			case <-t.Context().Done():
				t.Fatalf("test %v timed out", f.Name())
				return
			}
		}
	}
}
