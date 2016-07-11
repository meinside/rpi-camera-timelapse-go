package storage

import (
	path "path/filepath"

	dropbox "github.com/stacktic/dropbox"
)

// storage interface for saving files on Dropbox

type DropboxStorage struct {
	path    *string
	dropbox *dropbox.Dropbox
}

func NewDropboxStorage(key, secret, token, path *string) *DropboxStorage {
	box := dropbox.NewDropbox()
	box.SetAppInfo(*key, *secret)
	box.SetAccessToken(*token)

	return &DropboxStorage{
		path:    path,
		dropbox: box,
	}
}

func (s *DropboxStorage) Save(filepath *string) error {
	_, err := s.dropbox.UploadFile(*filepath, path.Join(*s.path, path.Base(*filepath)), true, "")
	return err
}
