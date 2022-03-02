package storage

import (
	bt "bytes"
	"io/ioutil"
	path "path/filepath"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

// storage interface for saving files on Dropbox

// DropboxStorage struct
type DropboxStorage struct {
	path   string
	client files.Client
}

// NewDropboxStorage creates a new DropboxStorage
//
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

// Save saves a file with bytes
func (s *DropboxStorage) Save(filename string, bytes []byte) error {
	reader := ioutil.NopCloser(bt.NewReader(bytes))
	defer reader.Close()

	dst := path.Join(s.path, filename)

	arg := files.NewUploadArg(dst)
	arg.CommitInfo.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}} // overwrite
	_, err := s.client.Upload(arg, reader)

	return err
}
