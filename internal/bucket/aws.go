package bucket

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsConfig struct {
	Config         *aws.Config
	BucketDownload string
	BucketUpload   string
}

type awsSession struct {
	session        *session.Session
	bucketDownload string
	bucketUpload   string
}

func newAwsSession(cfg AwsConfig) *awsSession {
	session := session.New(cfg.Config)

	return &awsSession{
		session:        session,
		bucketDownload: cfg.BucketDownload,
		bucketUpload:   cfg.BucketUpload,
	}
}

func (session *awsSession) Download(source string, destiny string) (file *os.File, err error) {
	file, err = os.Create(destiny)
	if err != nil {
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(session.session)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(session.bucketDownload),
			Key:    aws.String(source),
		})

	return
}

func (session *awsSession) Upload(file io.Reader, key string) error {
	uploader := s3manager.NewUploader(session.session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(session.bucketUpload),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

func (session *awsSession) Delete(source string) error {
	service := s3.New(session.session)
	_, err := service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(session.bucketDownload),
		Key:    aws.String(source),
	})

	if err != nil {
		return err
	}

	return service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(session.bucketDownload),
		Key:    aws.String(source),
	})
}
