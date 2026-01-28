package testutils

import (
	"math"
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

		t.Logf("test failed, retrying, attempt %v of %v\n", i+1, p.maxTries)
		if i < p.maxTries-1 {
			sleep := p.sleep * time.Duration(math.Pow(2, float64(i)))
			t.Logf("sleeping for %v\b", sleep)

			time.Sleep(sleep)
		}
	}
}
