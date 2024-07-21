package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const LocalBucket = "local"

type LocalStorage struct {
	folder string
}

func IsLocalBucket(bucket Bucket) bool {
	return bucket.Name == LocalBucket
}

func NewFolderStorage(bucket Bucket) (*LocalStorage, error) {

	folder, err := filepath.Abs(bucket.Endpoint)
	if err != nil {
		return nil, err
	}

	return &LocalStorage{
		folder: folder,
	}, nil
}

// Save writes or overwrites file
func (f *LocalStorage) Save(data []byte, filename string) error {

	path, err := f.path(filename)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (f *LocalStorage) Load(filename string) ([]byte, error) {

	path, err := f.path(filename)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []byte("{}"), nil
	}
	return b, err
}

func (f *LocalStorage) Delete(filename string) error {

	path, err := f.path(filename)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (f *LocalStorage) Reset() error {
	return os.RemoveAll(f.folder)
}

func (f *LocalStorage) path(filename string) (string, error) {

	if len(f.folder) == 0 {
		return filename, nil
	}

	// get folder path from filename
	dir := filepath.Dir(filename)

	// create folders recursively if not exists
	err := os.MkdirAll(filepath.Join(f.folder, dir), 0755)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", f.folder, filename), nil
}
