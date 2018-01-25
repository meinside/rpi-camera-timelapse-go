package storage

// storage interface

type StorageType string

const (
	TypeLocal   StorageType = "local"
	TypeSmtp    StorageType = "smtp"
	TypeDropbox StorageType = "dropbox"
)

type Config struct {
	Type StorageType `json:"type"`

	// for local & dropbox
	Path *string `json:"path,omitempty"`

	// for SMTP
	SmtpRecipients *string `json:"smtp_recipients,omitempty"`
	SmtpEmail      *string `json:"smtp_email,omitempty"`
	SmtpPasswd     *string `json:"smtp_passwd,omitempty"`
	SmtpServer     *string `json:"smtp_server,omitempty"`

	// for dropbox
	DropboxToken *string `json:"dropbox_token,omitempty"`
}

type Interface interface {
	Save(filename string, bytes []byte) error
}
