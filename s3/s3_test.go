package s3

import (
	"context"
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
	s.Run("success", func() {
		err := s.client.Store(s.ctx, testFilename, s.randomBytes())
		s.Require().NoError(err)

		err = s.client.Delete(s.ctx, []string{testFilename})
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
			filesToStore = make([]string, 0, number)
		)

		for i := 0; i < number; i++ {
			filename := strconv.Itoa(i + 1)
			content := s.randomBytes()
			filesToStore = append(filesToStore, filename)

			err = s.client.Store(context.Background(), filename, content)
			s.Require().NoError(err)
		}

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
