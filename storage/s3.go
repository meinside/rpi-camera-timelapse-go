package storage

import (
	bt "bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"path/filepath"
	"fmt"
)

// storage interface for saving files in an S3 bucket

type S3Storage struct {
	bucket *string
	path   *string
	uploader *s3manager.Uploader
}

func NewS3Storage(bucket, path *string) *S3Storage {
	if bucket == nil {
		panic("Bucket parameter missing for S3 storage")
	}

	// use empty path if not specified in config, avoiding panic later on
	if path == nil {
		path = aws.String("")
	}

	// The session the S3 Uploader will use, credentials is read
	// from default places where AWS CLI usually finds them
	// e.g. environment variables or ~/.aws/credentials
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))

	uploader := s3manager.NewUploader(session)

	return &S3Storage{bucket, path, uploader}
}

func (instance *S3Storage) Save(filename string, bytes []byte) error {
	reader := ioutil.NopCloser(bt.NewReader(bytes))
	defer reader.Close()

	objectKey := filepath.Join(*instance.path, filename)

	_, err := instance.uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(*instance.bucket),
		Key:    &objectKey,
		Body:   reader,
	})

	return err
}

func (instance *S3Storage) String() string {
	return fmt.Sprintf("S3 bucket %s/%s", *instance.bucket, *instance.path)
}