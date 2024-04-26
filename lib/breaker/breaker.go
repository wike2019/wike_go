package breaker

import (
	"github.com/sony/gobreaker"
	"time"
)

var DefaultCb *gobreaker.CircuitBreaker

func init() {
	var cbSettings gobreaker.Settings
	cbSettings.Name = "default CircuitBreaker"
	cbSettings.MaxRequests = 3
	cbSettings.Interval = time.Minute * 1
	cbSettings.Timeout = time.Minute * 2
	cbSettings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 6 && failureRatio >= 0.6
	}
	DefaultCb = gobreaker.NewCircuitBreaker(cbSettings)
}

func NewCircuitBreaker(name string, MaxRequests uint32, Interval time.Duration, Timeout time.Duration, ReadyToTrip func(counts gobreaker.Counts) bool) *gobreaker.CircuitBreaker {
	var cbSettings gobreaker.Settings
	cbSettings.Name = name
	cbSettings.MaxRequests = MaxRequests
	cbSettings.Interval = Interval
	cbSettings.Timeout = Timeout
	cbSettings.ReadyToTrip = ReadyToTrip
	return gobreaker.NewCircuitBreaker(cbSettings)
}
