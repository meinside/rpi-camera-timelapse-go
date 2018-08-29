package storage

import (
	bt "bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

const (
	maxLineLength = 500
)

// storage interface for saving files through SMTP

// SMTPStorage struct
type SMTPStorage struct {
	senderEmail     *string
	senderPasswd    *string
	senderServer    *string
	recipientEmails []string

	auth smtp.Auth
}

// NewSMTPStorage creates a new SMTPStorage
func NewSMTPStorage(senderEmail, senderServer, senderPasswd, recipientEmails *string) *SMTPStorage {
	if senderEmail == nil || senderServer == nil || senderPasswd == nil || recipientEmails == nil {
		panic("Parameter missing or invalid for SMTP")
	}

	auth := smtp.PlainAuth(
		"",
		*senderEmail,
		*senderPasswd,
		strings.Split(*senderServer, ":")[0], // XXX - without port number
	)

	return &SMTPStorage{
		senderEmail:     senderEmail,
		senderPasswd:    senderPasswd,
		senderServer:    senderServer,
		recipientEmails: strings.Split(*recipientEmails, ","),
		auth:            auth,
	}
}

// Save sends a file with bytes
//
// referenced: http://www.robertmulley.com/golang/sending-emails-with-attachments/
func (s *SMTPStorage) Save(filename string, bytes []byte) error {
	// captured time
	now := time.Now().Format("2006-01-02 15:04:05")

	// unique boundary for multipart
	md5edTime := md5.Sum([]byte(now))
	boundary := fmt.Sprintf("__%s__", base64.StdEncoding.EncodeToString(md5edTime[:]))

	// mail body part
	from := fmt.Sprintf("From: RPi Timelapse Camera <%s>", *s.senderEmail)
	subject := fmt.Sprintf("Subject: RPi Timelapse Image: %s", now)
	html := fmt.Sprintf("<html><body>Captured on <strong>%s</strong></body></html>", now) // HTML body
	body := fmt.Sprintf("%s\r\n%s\r\nMIME-version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n--%s\r\nContent-Type: text/html\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s\r\n--%s", from, subject, boundary, boundary, html, boundary)

	// attachment part
	encodedFile := base64.StdEncoding.EncodeToString(bytes)
	numLines := len(encodedFile) / maxLineLength
	var buf bt.Buffer
	for i := 0; i < numLines; i++ {
		buf.WriteString(encodedFile[i*maxLineLength:(i+1)*maxLineLength] + "\n")
	}
	buf.WriteString(encodedFile[numLines*maxLineLength:])
	attachment := fmt.Sprintf("Content-type: image/jpeg; name=\"%s\"\r\nContent-Transfer-Encoding:base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n%s\r\n--%s--", filename, filename, buf.String(), boundary)

	return smtp.SendMail(
		*s.senderServer,
		s.auth,
		*s.senderEmail,
		s.recipientEmails,
		[]byte(body+"\r\n"+attachment),
	)
}
