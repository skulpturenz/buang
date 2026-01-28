package testutils

import (
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
		test(t)

		if !t.Failed() {
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

			t.Logf("test failed %v, retrying, attempt %v of %v\n", f.Name(), i+1, p.maxTries)
			sleep := p.sleep * time.Duration(math.Pow(2, float64(i)))
			t.Logf("sleeping for %v\b", sleep)

			time.Sleep(sleep)
		}
	}
}
