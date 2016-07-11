package storage

import (
	"io"
	"os"
	path "path/filepath"
)

// storage interface for saving files locally

type LocalStorage struct {
	path *string
}

func NewLocalStorage(path *string) *LocalStorage {
	return &LocalStorage{
		path: path,
	}
}

func (s *LocalStorage) Save(filepath *string) error {
	if sf, err := os.Open(*filepath); err == nil {
		defer sf.Close()

		destination := path.Join(*s.path, path.Base(*filepath))
		if df, err := os.Create(destination); err == nil {
			defer df.Close()

			_, err := io.Copy(df, sf)
			return err
		} else {
			return err
		}
	} else {
		return err
	}
}
