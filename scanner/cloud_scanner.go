package scanner

import (
	"context"
	"net/url"
	"time"

	"github.com/deluan/navidrome/conf"
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
	ListObjects(ctx context.Context, bucket string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
}

//S3CloudScanner scans S3 storage
type S3CloudScanner struct {
	model.DataStore
	s3Client       s3Client
	bucket, prefix string
}

//NewS3CloudScanner returns an initialized S3CloudScanner
func NewS3CloudScanner(rootFolder string, ds model.DataStore) *S3CloudScanner {
	s := S3CloudScanner{DataStore: ds}
	err := s.Login(conf.Server.ObjectStoreAccessID, conf.Server.ObjectStoreAccessKey)
	if err != nil {
		log.Error(err)
	}
	u, err := url.Parse(rootFolder)
	if err != nil {
		log.Error(err)
	}
	s.bucket = u.Host
	s.prefix = u.Path
	return &s
}

//Scan scans a media folder in the cloud
func (s *S3CloudScanner) Scan(ctx context.Context, lastModifiedSince time.Time) error {
	// if there are existing songs in the local db... scan them for changes or deletions
	media, err := s.DataStore.MediaFile(ctx).GetAll()
	if err != nil {
		return err
	}
	scanForChanges(media)
	// then scan the object storage, and if there are songs we don't have, add them to local db
	addNewMedia(s.s3Client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		WithVersions: false,
		WithMetadata: false,
		Prefix:       s.prefix,
		Recursive:    false,
		MaxKeys:      0,
		UseV1:        false,
	}))
	return nil
}

func scanForChanges(mfs model.MediaFiles) {

}
func addNewMedia(objects <-chan minio.ObjectInfo) {

}

//Login logs into the cloud storage
func (s *S3CloudScanner) Login(id, key string) (err error) {
	c, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(id, key, ""),
		Secure: false,
	})
	s.s3Client = c
	return
}
