package scanner_test

import (
	"github.com/deluan/navidrome/model"
	"github.com/deluan/navidrome/scanner"
	"github.com/go-kit/kit/endpoint"
)

type CloudScanner interface{
	Login(accessKey, secretKey, endpoint string)
}

func Test(t *testing.T) {
	testCases := []struct {
		desc	string
		MusicFolder string
	}{
		{
			desc: "S3 folder",
			MusicFolder: "s3://some_bucket/some/path/",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			s := scanner.NewTagScanner(tC.MusicFolder, model.DataStore{})
		})
	}
}