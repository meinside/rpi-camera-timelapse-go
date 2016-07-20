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
	DropboxKey    *string `json:"dropbox_key,omitempty"`
	DropboxSecret *string `json:"dropbox_secret,omitempty"`
	DropboxToken  *string `json:"dropbox_token,omitempty"`
}

type Interface interface {
	Save(filepath *string) error
}
