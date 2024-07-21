package store

import "fmt"

// NewMediaStorage is a wrapper for NewStorage to upload media in given folder.
func NewMediaStorage(spiderID string, bucket Bucket) (Storage, error) {
	return NewStorage(bucket, GetMediaStorageFolder(spiderID))
}

func GetMediaStorageFolder(spiderID string) string {
	return fmt.Sprintf(MediaFolder, spiderID)
}
