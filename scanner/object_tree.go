package scanner

import (
	"context"

	"github.com/deluan/navidrome/utils"
	b2 "github.com/kothar/go-backblaze"
)

// B2Client manages state and credentials of a b2 object store
type B2Client struct {
	bucket   *b2.Bucket
	rootPath string
}

func newB2Client(account, appkey, bucket, path string) (*B2Client, error) {
	creds := b2.Credentials{
		AccountID:      account,
		ApplicationKey: appkey,
	}
	client, err := b2.NewB2(creds)
	if err != nil {
		return nil, err
	}
	b, err := client.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	return &B2Client{b, path}, err
}

// LoadDirTree lists objects in a B2 bucket and maps media found therein
func (c B2Client) LoadDirTree(ctx context.Context) (dirMap, error) {
	newMap := make(dirMap)
	err := c.loadMap(ctx, c.rootPath, newMap)
	return newMap, err
}

func (c B2Client) loadMap(ctx context.Context, path string, mapping dirMap) error {
	resp := b2.ListFilesResponse{
		NextFileName: "",
	}
	for {
		resp, err := c.bucket.ListFileNamesWithPrefix(resp.NextFileName, 999, path, "")
		if err != nil {
			return err
		}
		for _, f := range resp.Files {
			path := ""
			m, ok := mapping[path]
			if !ok {
				m = dirMapValue{}
			}
			m.hasImages = m.hasImages || utils.IsImageFile(f.Name)
			m.hasPlaylist = m.hasPlaylist || utils.IsPlaylist(f.Name)
			m.hasAudioFiles = m.hasAudioFiles || utils.IsAudioFile(f.Name)
			mapping[path] = m
		}
		if resp.NextFileName == "" {
			break
		}
	}
	return nil
}
