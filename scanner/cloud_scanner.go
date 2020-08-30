package scanner

import (
	"context"
	"time"

	"github.com/deluan/navidrome/model"
)

//CloudScanner allows scanning of online resources
type CloudScanner interface {
	FolderScanner
	Login(ID, Key string) error
}

//S3CloudScanner scans S3 storage
type S3CloudScanner struct{}

//NewS3CloudScanner returns an initialized S3CloudScanner
func NewS3CloudScanner(rootFolder string, ds model.DataStore) CloudScanner {
	return &S3CloudScanner{}
}

//Scan scans a media folder in the cloud
func (s *S3CloudScanner) Scan(ctx context.Context, lastModifiedSince time.Time) error {
	return nil
}

//Login logs into the cloud storage
func (s *S3CloudScanner) Login(ID, Key string) error {
	return nil
}
