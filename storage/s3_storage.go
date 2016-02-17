package storage

import (
	"bytes"
	"errors"
	"io"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type s3Storage struct {
	bucket *s3.Bucket
}

type writerWithCallback struct {
	io.Writer
	callback func() error
}

func (wc *writerWithCallback) Close() error {
	return wc.callback()
}

func S3Storage(keyID, key, regionName, bucket string) (*s3Storage, error) {
	region, exists := aws.Regions[regionName]
	if !exists {
		return nil, errors.New("invalid region")
	}
	auth := aws.Auth{AccessKey: keyID, SecretKey: key, Token: ""}
	return &s3Storage{s3.New(auth, region).Bucket(bucket)}, nil
}

// Put provides a writer for saving the entry as a local file.
func (s3s *s3Storage) Put(key string) (io.WriteCloser, error) {
	var buffer bytes.Buffer
	return &writerWithCallback{
		&buffer,
		func() error {
			return s3s.bucket.Put(key, buffer.Bytes(), "binary/octet-stream", s3.Private)
		},
	}, nil
}

// Get provides a reader for reading the entry from a local file.
func (s3s *s3Storage) Get(key string) (io.ReadCloser, error) {
	reader, err := s3s.bucket.GetReader(key)
	if err != nil {
		return nil, err
	}
	return &readCloserWithCallback{
		reader,
		func() error {
			return s3s.bucket.Del(key)
		},
	}, nil
}
