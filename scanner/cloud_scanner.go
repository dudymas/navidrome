package scanner

import (
	"context"
	"time"

	"github.com/deluan/navidrome/log"
	"github.com/deluan/navidrome/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//CloudScanner allows scanning of online resources
type CloudScanner interface {
	FolderScanner
	Login(ID, Key string) error
}

type s3Client interface {
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
}

//S3CloudScanner scans S3 storage
type S3CloudScanner struct {
	s3Client s3Client
}

//NewS3CloudScanner returns an initialized S3CloudScanner
func NewS3CloudScanner(rootFolder string, ds model.DataStore) *S3CloudScanner {
	s := S3CloudScanner{}
	err := s.Login("AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	if err != nil {
		log.Error(err)
	}
	return &s
}

//Scan scans a media folder in the cloud
func (s *S3CloudScanner) Scan(ctx context.Context, lastModifiedSince time.Time) error {
	return nil
}

//Login logs into the cloud storage
func (s *S3CloudScanner) Login(id, key string) (err error) {
	s.s3Client, err = minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(id, key, ""),
		Secure: false,
	})
	return
}
