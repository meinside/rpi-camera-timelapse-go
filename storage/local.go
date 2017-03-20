package storage

import (
	bt "bytes"
	"io"
	"os"
	path "path/filepath"
)

// storage interface for saving files locally

type LocalStorage struct {
	path *string
}

func NewLocalStorage(path *string) *LocalStorage {
	if path == nil {
		panic("Parameter missing or invalid for Local")
	}

	return &LocalStorage{
		path: path,
	}
}

func (s *LocalStorage) Save(filename string, bytes []byte) error {
	src := bt.NewReader(bytes)
	dst := path.Join(*s.path, filename)

	if df, err := os.Create(dst); err == nil {
		defer df.Close()

		_, err := io.Copy(df, src)
		return err
	} else {
		return err
	}
}
