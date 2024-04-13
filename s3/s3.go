package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

const urlMedia = "https://%s.s3.amazonaws.com/%s"

func (c *Client) List(ctx context.Context) ([]*Media, error) {
	objects, err := c.client.ListObjects(ctx, &s3.ListObjectsInput{Bucket: aws.String(c.bucketName)})
	if err != nil {
		return nil, err
	}

	if len(objects.Contents) == 0 {
		return nil, ErrNoContents
	}

	out := make([]*Media, 0, len(objects.Contents))

	for i := range objects.Contents {
		if len(objects.Contents) == 0 {
			continue
		}

		object := objects.Contents[i]

		if *object.Size != 0 {
			u, _ := url.Parse(fmt.Sprintf(urlMedia, c.bucketName, *object.Key))
			media := &Media{
				Filename:     *object.Key,
				Url:          u.String(),
				LastModified: *object.LastModified,
			}
			out = append(out, media)
		}
	}

	return out, nil
}

func (c *Client) Store(ctx context.Context, filename string, contentBytes []byte) error {
	result, err := c.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(contentBytes),
	})
	if err != nil {
		return err
	}

	if result.Location == "" {
		return errors.New("error store file")
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, filenames []string) error {
	toDelete := make([]types.ObjectIdentifier, len(filenames))
	for i := range filenames {
		toDelete[i] = types.ObjectIdentifier{
			Key: aws.String(filenames[i]),
		}
	}

	objects, err := c.client.DeleteObjects(ctx,
		&s3.DeleteObjectsInput{Bucket: aws.String(c.bucketName), Delete: &types.Delete{Objects: toDelete}})
	if err != nil {
		return err
	}

	if len(objects.Deleted) != len(filenames) {
		deleted := make([]string, len(objects.Deleted))
		for i := range objects.Deleted {
			if objects.Deleted[i].Key != nil {
				deleted[i] = *objects.Deleted[i].Key
			}
		}

		c.logger.Error("s3 deleted objects are not equal as expected",
			zap.Strings("deleted", deleted),
			zap.Strings("expected", filenames),
		)
	}

	return nil
}
