package retry

import (
	"github.com/avast/retry-go"
	"time"
)

type Retry struct {
	job   func() error
	times int
	delay time.Duration
}

func NewRetry(job func() error) *Retry {
	return &Retry{job: job, times: 3, delay: time.Second * 1}
}
func (r *Retry) Do() error {
	return retry.Do(r.job, retry.Attempts(uint(r.times)), retry.Delay(r.delay), retry.LastErrorOnly(true), retry.DelayType(retry.BackOffDelay))
}
func (r *Retry) SetTimes(times int) *Retry {
	r.times = times
	return r
}
func (r *Retry) SetDelay(delay time.Duration) *Retry {
	r.delay = delay
	return r
}
