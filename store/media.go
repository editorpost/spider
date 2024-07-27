package store

import (
	"fmt"
	"github.com/editorpost/donq/res"
)

// NewMediaStorage is a wrapper for NewStorage to upload media in given folder.
func NewMediaStorage(spiderID string, bucket *res.S3) (Storage, error) {
	return NewStorage(bucket, GetMediaStorageFolder(spiderID))
}

func GetMediaStorageFolder(spiderID string) string {
	return fmt.Sprintf(PayloadFolder, spiderID)
}
