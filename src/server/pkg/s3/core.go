package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client interface {
	Upload(ctx context.Context, key string, fileBytes []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

func NewClient(
	s3 *s3.Client,
	bucket string,
) *ClientImpl {
	return &ClientImpl{
		s3:     s3,
		bucket: bucket,
	}
}

type ClientImpl struct {
	s3     *s3.Client
	bucket string
}

func (c ClientImpl) Upload(ctx context.Context, key string, fileBytes []byte) error {
	_, err := c.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileBytes),
	})
	if err != nil {
		return fmt.Errorf("s3 upload: %w", err)
	}
	return nil
}

func (c ClientImpl) Get(ctx context.Context, key string) ([]byte, error) {
	output, err := c.s3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 get: %w", err)
	}
	defer output.Body.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(output.Body); err != nil {
		return nil, fmt.Errorf("s3 read body: %w", err)
	}

	return buf.Bytes(), nil
}
