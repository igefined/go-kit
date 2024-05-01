package uchannel

import (
	"fmt"
	"sync"
	"testing"
)

func TestSplit(t *testing.T) {
	const chSize = 5

	t.Run("int channel", func(t *testing.T) {
		source := make(chan int)
		destinations := Split[int](source, chSize)

		go func() {
			for i := 0; i < 10; i++ {
				source <- i
			}

			close(source)
		}()

		var wg sync.WaitGroup
		wg.Add(chSize)

		for i, ch := range destinations {
			go func(i int, ch <-chan int) {
				defer wg.Done()

				for val := range ch {
					fmt.Printf("#%d got %d\n", i, val)
				}
			}(i, ch)
		}

		wg.Wait()
	})

	t.Run("string channel", func(t *testing.T) {
		source := make(chan string)
		destinations := Split[string](source, chSize)

		go func() {
			for i := 0; i < 10; i++ {
				source <- fmt.Sprintf("string #%d", i+1)
			}

			close(source)
		}()

		var wg sync.WaitGroup
		wg.Add(chSize)

		for i, ch := range destinations {
			go func(i int, ch <-chan string) {
				defer wg.Done()

				for val := range ch {
					fmt.Printf("#%d got %s\n", i, val)
				}
			}(i, ch)
		}

		wg.Wait()
	})

	t.Run("custom struct channel", func(t *testing.T) {
		type User struct {
			ID   uint64
			Name string
		}

		source := make(chan User)
		destinations := Split[User](source, chSize)

		go func() {
			for i := 0; i < 10; i++ {
				source <- User{
					ID:   uint64(i + 1),
					Name: fmt.Sprintf("user #%d", i+1),
				}
			}

			close(source)
		}()

		var wg sync.WaitGroup
		wg.Add(chSize)

		for i, ch := range destinations {
			go func(i int, ch <-chan User) {
				defer wg.Done()

				for val := range ch {
					fmt.Printf("#%d got %+v\n", i, val)
				}
			}(i, ch)
		}

		wg.Wait()
	})
}
