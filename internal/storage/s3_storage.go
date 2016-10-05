package storage

import (
	"bytes"
	"errors"
	"io"
	"time"

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

func S3Storage(keyID, key, regionName, bucketName string) (*s3Storage, error) {
	region, exists := aws.Regions[regionName]
	if !exists {
		return nil, errors.New("invalid region")
	}
	auth := aws.Auth{AccessKey: keyID, SecretKey: key, Token: ""}
	bucket := s3.New(auth, region).Bucket(bucketName)
	_, err := bucket.List("", "/", "", 1)
	return &s3Storage{bucket}, err
}

// Put provides a writer for saving the entry as an S3 file.
func (s3s *s3Storage) Put(key string) (io.WriteCloser, error) {
	var buffer bytes.Buffer
	return &writerWithCallback{
		&buffer,
		func() error {
			return s3s.bucket.Put(key, buffer.Bytes(), "binary/octet-stream", s3.Private)
		},
	}, nil
}

// Get provides a reader for reading the entry from an S3 file.
func (s3s *s3Storage) Get(key string) (io.ReadCloser, error) {
	return s3s.tryGet(key, 3)
}

// Delete removes an entry stored in an S3 file.
func (s3s *s3Storage) Delete(key string) error {
	return s3s.bucket.Del(key)
}

func (s3s *s3Storage) tryGet(key string, remaining int) (io.ReadCloser, error) {
	remaining--
	ret, err := s3s.bucket.GetReader(key)
	if err == nil || remaining == 0 {
		return ret, err
	}
	time.Sleep(3 * time.Second)
	return s3s.tryGet(key, remaining)
}
