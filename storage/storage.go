package storage

// storage interface

type StorageType string

const (
	TypeLocal   StorageType = "local"
	TypeDropbox StorageType = "dropbox"
)

type Config struct {
	Type   StorageType `json:"type"`
	Path   *string     `json:"path"`
	Key    *string     `json:"key,omitempty"`
	Secret *string     `json:"secret,omitempty"`
	Token  *string     `json:"token,omitempty"`
}

type Interface interface {
	Save(filepath *string) error
}
