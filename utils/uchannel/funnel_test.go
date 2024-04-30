package uchannel

import (
	"fmt"
	"testing"
	"time"
)

func TestFunnel(t *testing.T) {
	const (
		sourcesSize = 5
		elementSize = 10
	)

	t.Run("simple int sources", func(t *testing.T) {
		sources := make([]<-chan int, 0)

		for i := 0; i < sourcesSize; i++ {
			ch := make(chan int)
			sources = append(sources, ch)

			go func() {
				defer close(ch)

				for j := 0; j < elementSize; j++ {
					ch <- j
					time.Sleep(time.Millisecond * 100)
				}
			}()
		}

		dest := Funnel[int](sources...)
		for d := range dest {
			fmt.Printf("%d ", d)
		}
	})

	t.Run("simple string sources", func(t *testing.T) {
		sources := make([]<-chan string, 0)

		for i := 0; i < sourcesSize; i++ {
			ch := make(chan string)
			sources = append(sources, ch)

			go func() {
				defer close(ch)

				for j := 0; j < elementSize; j++ {
					ch <- fmt.Sprintf("message #%d", j)
					time.Sleep(time.Millisecond * 100)
				}
			}()
		}

		dest := Funnel[string](sources...)
		for d := range dest {
			fmt.Printf("%s ", d)
		}
	})

	t.Run("simple struct sources", func(t *testing.T) {
		type user struct {
			ID   uint64
			Name string
		}

		sources := make([]<-chan user, 0)

		for i := 0; i < sourcesSize; i++ {
			ch := make(chan user)
			sources = append(sources, ch)

			go func() {
				defer close(ch)

				for j := 0; j < elementSize; j++ {
					ch <- user{
						ID:   uint64(j + 1),
						Name: fmt.Sprintf("username%d", j),
					}
					time.Sleep(time.Millisecond * 100)
				}
			}()
		}

		dest := Funnel[user](sources...)
		for d := range dest {
			fmt.Printf("%+v ", d)
		}
	})
}
