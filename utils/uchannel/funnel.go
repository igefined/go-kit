package uchannel

import "sync"

// Funnel is a simple Fan-In pattern
func Funnel[T any](sources ...<-chan T) <-chan T {
	var (
		dest = make(chan T)
		wg   = sync.WaitGroup{}
	)

	wg.Add(len(sources))

	for _, ch := range sources {
		go func(ch <-chan T) {
			defer wg.Done()

			for item := range ch {
				dest <- item
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}
