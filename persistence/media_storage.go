package persistence

import (
	"context"
	"os"

	"github.com/deluan/navidrome/model"
	"github.com/minio/minio-go/v7"
	s3 "github.com/minio/minio-go/v7"
	s3Creds "github.com/minio/minio-go/v7/pkg/credentials"
)

type CloudClient struct {
	*minio.Client
	bucket string
}

func (c *CloudClient) GetMedia(ctx context.Context, path string) (model.MediaData, error) {
	return c.Client.GetObject(ctx, c.bucket, path, minio.GetObjectOptions{})
}

type FilesystemStorage struct{}

func (fs *FilesystemStorage) GetMedia(ctx context.Context, path string) (model.MediaData, error) {
	return os.Open(path)
}

func newS3Client(endpoint, accessKeyID, secretAccessKey, bucket, path string) (*CloudClient, error) {
	useSSL := false
	creds := s3Creds.NewStaticV4(accessKeyID, secretAccessKey, "")
	client, err := s3.New(
		endpoint,
		&s3.Options{
			Creds:  creds,
			Secure: useSSL,
		},
	)
	return &CloudClient{client, bucket}, err
}
