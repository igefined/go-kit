package http

import "time"

type Opt func(client *ClientImpl)

func WithRetryThreshold(retryThreshold int) Opt {
	return func(client *ClientImpl) {
		client.retryThreshold = retryThreshold
	}
}

func WithTimeout(timeout time.Duration) Opt {
	return func(client *ClientImpl) {
		client.httpClient.Timeout = timeout
	}
}
