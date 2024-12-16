package breaker

import (
	"github.com/sony/gobreaker"
	"time"
)

var DefaultToTrip = func(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 6 && failureRatio >= 0.6
}

func NewCircuitBreaker(name string, MaxRequests uint32, Interval time.Duration, Timeout time.Duration) *gobreaker.CircuitBreaker {
	var cbSettings gobreaker.Settings
	cbSettings.Name = name
	cbSettings.MaxRequests = MaxRequests
	cbSettings.Interval = Interval
	cbSettings.Timeout = Timeout
	cbSettings.ReadyToTrip = DefaultToTrip
	return gobreaker.NewCircuitBreaker(cbSettings)
}

func NewCircuitBreakerWithSettings(settings gobreaker.Settings) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(settings)
}

/*
demo
this.Breaker.Execute(this.Job)
**/
