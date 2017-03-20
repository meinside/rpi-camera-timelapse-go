package storage

import (
	bt "bytes"
	"io/ioutil"
	path "path/filepath"

	dropbox "github.com/stacktic/dropbox"
)

// storage interface for saving files on Dropbox

type DropboxStorage struct {
	path    *string
	dropbox *dropbox.Dropbox
}

func NewDropboxStorage(key, secret, token, path *string) *DropboxStorage {
	if key == nil || secret == nil || token == nil || path == nil {
		panic("Parameter missing or invalid for Dropbox")
	}

	box := dropbox.NewDropbox()
	box.SetAppInfo(*key, *secret)
	box.SetAccessToken(*token)

	return &DropboxStorage{
		path:    path,
		dropbox: box,
	}
}

func (s *DropboxStorage) Save(filename string, bytes []byte) error {
	reader := ioutil.NopCloser(bt.NewReader(bytes))
	defer reader.Close()

	length := int64(len(bytes))
	dst := path.Join(*s.path, filename)

	_, err := s.dropbox.FilesPut(reader, length, dst, true, "")
	return err
}
