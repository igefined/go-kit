package http

import "net/http"

type Client interface {
	DoRequest(*http.Request) (*http.Response, error)
}

type ClientImpl struct {
	retryThreshold int

	httpClient *http.Client
}

func New(opts ...Opt) Client {
	instance := &ClientImpl{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(instance)
	}

	return instance
}

func (c ClientImpl) DoRequest(r *http.Request) (*http.Response, error) {
	apiReq := func(r *http.Request) (*http.Response, error) {
		return c.httpClient.Do(r)
	}

	circuitBreaker := Breaker(apiReq, c.retryThreshold)

	response, err := circuitBreaker(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return response, err
}
