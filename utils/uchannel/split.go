package uchannel

// Split is a simple Fan-Out pattern
func Split[T any](source chan T, n int) []<-chan T {
	dest := make([]<-chan T, 0)

	for i := 0; i < n; i++ {
		ch := make(chan T)
		dest = append(dest, ch)

		go func() {
			defer close(ch)

			for val := range source {
				ch <- val
			}
		}()
	}

	return dest
}
