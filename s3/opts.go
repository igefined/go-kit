package s3

import (
	"errors"
	"time"
)

var ErrNoContents = errors.New("no contents")

type Media struct {
	Filename     string
	Url          string
	LastModified time.Time
}

type Opt func(client *Client)

func WithBucketName(bucketName string) Opt {
	return func(client *Client) {
		client.bucketName = bucketName
	}
}
