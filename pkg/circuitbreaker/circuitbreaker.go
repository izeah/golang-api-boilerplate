package circuitbreaker

import (
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type (
	// Breaker is the http circuit breaker.
	Breaker interface {
		// Execute runs the given request if the circuit breaker is closed or half-open states.
		// An error is instantly returned when the circuit breaker is tripped.
		Execute(func() (interface{}, error)) (interface{}, error)
	}
)

// NewClient ...
func NewClient() *http.Client {
	return newClient()
}

func newClient() *http.Client {
	return &http.Client{
		Transport: newTransport(newCircuitBreaker()),
	}
}

func newCircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "HTTP Client",
		Timeout: time.Second * 45,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})
}
