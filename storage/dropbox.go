package storage

import (
	bt "bytes"
	"io/ioutil"
	path "path/filepath"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

// storage interface for saving files on Dropbox

type DropboxStorage struct {
	path   string
	client files.Client
}

// `token` can be obtained/generated in:
// Dropbox Developers page > My apps > Your App > Settings
func NewDropboxStorage(token, path *string) *DropboxStorage {
	if token == nil || path == nil {
		panic("Parameter missing or invalid for Dropbox")
	}

	return &DropboxStorage{
		path: *path,
		client: files.New(dropbox.Config{
			Token: *token,
		}),
	}
}

func (s *DropboxStorage) Save(filename string, bytes []byte) error {
	reader := ioutil.NopCloser(bt.NewReader(bytes))
	defer reader.Close()

	dst := path.Join(s.path, filename)

	_, err := s.client.Upload(&files.CommitInfo{
		Path:       dst,
		Mode:       &files.WriteMode{Tagged: dropbox.Tagged{"overwrite"}}, // overwrite
		Autorename: false,
		Mute:       false,
	}, reader)

	return err
}
