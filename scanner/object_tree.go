package scanner

import (
	"context"
	"strings"

	s3 "github.com/minio/minio-go/v7"
	s3Creds "github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Client represents an object store client
type S3Client struct {
	client *s3.Client
	bucket string
}

func newS3Client(endpoint, accessKeyID, secretAccessKey, bucket, path string) (*S3Client, error) {
	useSSL := false
	creds := s3Creds.NewStaticV4(accessKeyID, secretAccessKey, "")
	client, err := s3.New(
		endpoint,
		&s3.Options{
			Creds:  creds,
			Secure: useSSL,
		},
	)
	return &S3Client{client, bucket}, err
}

// LoadDirTree lists objects in a B2 bucket and maps media found therein
func (c S3Client) LoadDirTree(ctx context.Context) (dirMap, error) {
	newMap := make(dirMap)
	err := c.loadMap(ctx, newMap)
	return newMap, err
}

func (c S3Client) loadMap(ctx context.Context, mapping dirMap) error {
	opts := s3.ListObjectsOptions{
		Prefix: "musix",
		Recursive: true,
	}
	for object := range c.client.ListObjects(ctx, "buckie", opts) {
		if object.Err != nil {
			return object.Err
		}
		m, path := getObjectPrefixMapping(object.Key, mapping)
		m.hasImages = m.hasImages || strings.HasPrefix(object.ContentType, "image/")
		m.hasPlaylist = m.hasPlaylist || strings.HasSuffix(object.Key, ".m3u8") || strings.HasSuffix(object.Key, ".m3u")
		m.hasAudioFiles = m.hasAudioFiles || strings.HasPrefix(object.ContentType, "audio/")
		mapping[path] = m
	}
	return nil
}

// getObjectPrefixMapping Uses an object's key to determine its mapping and returns its path prefix
func getObjectPrefixMapping(key string, mapping dirMap) (dirMapValue, string) {
	breadcrumbs := strings.Split(key, "/")
	fileName := breadcrumbs[len(breadcrumbs)-1]
	fileNameLength := len(fileName)
	pathPrefix := key[0 : len(key)-fileNameLength]
	m, ok := mapping[pathPrefix]
	if !ok {
		m = dirMapValue{}
	}
	return m, pathPrefix
}
