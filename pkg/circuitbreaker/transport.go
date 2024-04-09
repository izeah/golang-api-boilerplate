package circuitbreaker

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type (
	// Transport is the application http transport.
	Transport struct {
		tripper http.RoundTripper
		breaker Breaker
	}
)

func newTransport(cb Breaker) *Transport {
	t := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 90 * time.Second,
		}).DialContext,
	}

	return &Transport{
		tripper: t,
		breaker: cb,
	}
}

// RoundTrip ...
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := t.breaker.Execute(func() (interface{}, error) {
		res, err := t.tripper.RoundTrip(r)
		if err != nil {
			return nil, err
		}

		if res != nil && res.StatusCode >= http.StatusInternalServerError {
			return res, fmt.Errorf("http response error: %v", res.StatusCode)
		}

		return res, err
	})

	if err != nil {
		return nil, err
	}

	return res.(*http.Response), err
}
