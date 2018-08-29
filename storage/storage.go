package storage

// storage interface

// Type type
type Type string

// types
const (
	TypeLocal   Type = "local"
	TypeSMTP    Type = "smtp"
	TypeDropbox Type = "dropbox"
	TypeS3      Type = "s3"
)

// Config struct
type Config struct {
	Type Type `json:"type"`

	// for local, dropbox, and S3
	Path *string `json:"path,omitempty"`

	// for SMTP
	SMTPRecipients *string `json:"smtp_recipients,omitempty"`
	SMTPEmail      *string `json:"smtp_email,omitempty"`
	SMTPPasswd     *string `json:"smtp_passwd,omitempty"`
	SMTPServer     *string `json:"smtp_server,omitempty"`

	// for dropbox
	DropboxToken *string `json:"dropbox_token,omitempty"`

	// for S3
	S3Bucket *string `json:"s3_bucket,omitempty"`
}

// Interface for interfacing
type Interface interface {
	Save(filename string, bytes []byte) error
}
