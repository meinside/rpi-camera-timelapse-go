package storage

import (
	bt "bytes"
	"io"
	"os"
	path "path/filepath"
)

// storage interface for saving files locally

// LocalStorage struct
type LocalStorage struct {
	path *string
}

// NewLocalStorage creates a new LocalStorage
func NewLocalStorage(path *string) *LocalStorage {
	if path == nil {
		panic("Parameter missing or invalid for Local")
	}

	return &LocalStorage{
		path: path,
	}
}

// Save saves a file with bytes
func (s *LocalStorage) Save(filename string, bytes []byte) error {
	src := bt.NewReader(bytes)
	dst := path.Join(*s.path, filename)

	df, err := os.Create(dst)

	if err == nil {
		defer df.Close()

		_, err := io.Copy(df, src)
		return err
	}

	return err
}
