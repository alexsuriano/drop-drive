package bucket

import (
	"errors"
	"io"
	"os"
	"reflect"
)

const (
	AwsProvider BucketType = iota
)

type BucketType int

type BucketInterface interface {
	Upload(io.Reader, string) error
	Download(string, string) (*os.File, error)
	Delete(string) error
}

type Bucket struct {
	Provider BucketInterface
}

func (bucket *Bucket) Upload(file io.Reader, key string) error {
	return bucket.Provider.Upload(file, key)
}

func (bucket *Bucket) Download(source string, destiny string) (file *os.File, err error) {
	return bucket.Provider.Download(source, destiny)
}

func (bucket *Bucket) Delete(source string) error {
	return bucket.Provider.Delete(source)
}

func New(bucketType BucketType, cfg any) (bucket *Bucket, err error) {
	config := reflect.TypeOf(cfg)

	switch bucketType {
	case AwsProvider:
		if config.Name() != "AwsConfig" {
			return nil, errors.New("config need's to be of type AwsConfig")
		}

		bucket.Provider = newAwsSession(cfg.(AwsConfig))
	default:
		return nil, errors.New("bucket type not implemented")
	}

	return
}
