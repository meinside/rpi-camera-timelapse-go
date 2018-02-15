package storage

// storage interface

type StorageType string

const (
	TypeLocal   StorageType = "local"
	TypeSmtp    StorageType = "smtp"
	TypeDropbox StorageType = "dropbox"
	TypeS3      StorageType = "s3"
)

type Config struct {
	Type StorageType `json:"type"`

	// for local, dropbox, and S3
	Path *string `json:"path,omitempty"`

	// for SMTP
	SmtpRecipients *string `json:"smtp_recipients,omitempty"`
	SmtpEmail      *string `json:"smtp_email,omitempty"`
	SmtpPasswd     *string `json:"smtp_passwd,omitempty"`
	SmtpServer     *string `json:"smtp_server,omitempty"`

	// for dropbox
	DropboxToken *string `json:"dropbox_token,omitempty"`

	// for S3
	S3Bucket *string `json:"s3_bucket,omitempty"`
}

type Interface interface {
	Save(filename string, bytes []byte) error
}
