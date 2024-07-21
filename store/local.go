package store

import (
	"fmt"
	"os"
)

const LocalBucket = "local"

type LocalStorage struct {
	folder string
}

func IsLocalBucket(bucket Bucket) bool {
	return bucket.Name == LocalBucket
}

func NewLocalBucket(publicURL string) Bucket {
	return Bucket{
		Name:      LocalBucket,
		PublicURL: publicURL,
	}
}

func NewFolderStorage(bucket Bucket) *LocalStorage {
	return &LocalStorage{
		folder: bucket.Endpoint,
	}
}

func (f *LocalStorage) Save(data []byte, filename string) error {
	return os.WriteFile(f.path(filename), data, 0644)
}

func (f *LocalStorage) Load(filename string) ([]byte, error) {
	return os.ReadFile(f.path(filename))
}

func (f *LocalStorage) Delete(filename string) error {
	return os.Remove(f.path(filename))
}

func (f *LocalStorage) Reset() error {
	return os.RemoveAll(f.folder)
}

func (f *LocalStorage) path(filename string) string {

	if len(f.folder) == 0 {
		return filename
	}

	return fmt.Sprintf("%s/%s", f.folder, filename)
}
