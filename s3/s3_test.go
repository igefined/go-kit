package s3

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/igefined/go-kit/utils/uslice"
)

const testFilename = "test_filename"

func (s *Suite) TestStore() {
	s.Run("success", func() {
		err := s.client.Store(s.ctx, testFilename, s.randomBytes())
		s.Require().NoError(err)
	})
}

func (s *Suite) TestDelete() {
	var (
		number    = rand.Intn(10)
		filenames = make([]string, number)
	)

	s.Run("not equal files", func() {
		for i := 0; i < number; i++ {
			filename := fmt.Sprintf("%s_%d", testFilename, i+1)

			err := s.client.Store(s.ctx, filename, s.randomBytes())
			s.Require().NoError(err)

			filenames[i] = filename
		}

		err := s.client.Delete(s.ctx, filenames[:number-2])
		s.Require().NoError(err)
	})

	s.Run("success", func() {
		err := s.client.Store(s.ctx, testFilename, s.randomBytes())
		s.Require().NoError(err)

		filenames = append(filenames, testFilename)

		err = s.client.Delete(s.ctx, filenames)
		s.Require().NoError(err)
	})
}

func (s *Suite) TestList() {
	s.Run("success, empty list", func() {
		medias, err := s.client.List(s.ctx)
		s.Require().Error(err, ErrNoContents)
		s.Require().Empty(medias)
	})

	s.Run("success", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var (
			err          error
			number       = rand.Intn(10)
			filesToStore = make([]string, 0, number+1)
		)

		for i := 0; i < number; i++ {
			filename := strconv.Itoa(i + 1)
			content := s.randomBytes()
			filesToStore = append(filesToStore, filename)

			err = s.client.Store(context.Background(), filename, content)
			s.Require().NoError(err)
		}

		filesToStore = append(filesToStore, strconv.Itoa(number+1))

		medias, err := s.client.List(ctx)
		s.Require().NoError(err)

		for i := range medias {
			uslice.Contains(filesToStore, medias[i].Filename)
		}
	})
}

func (s *Suite) randomBytes() []byte {
	bytes := make([]byte, 128)
	_, err := rand.Read(bytes)
	s.Require().NoError(err)

	return bytes
}
