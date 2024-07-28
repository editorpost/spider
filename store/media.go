package store

import (
	"fmt"
	"github.com/editorpost/donq/res"
	"time"
)

// NewMediaStorage is a wrapper for NewStorage to upload media in given folder.
func NewMediaStorage(spiderID string, bucket res.S3) (Storage, error) {
	return NewStorage(bucket, GetMediaStorageFolder(spiderID, time.Now()))
}

func GetMediaStorageFolder(spiderID string, chunkStamp time.Time) string {

	chunk := chunkStamp.UTC().Format(ChunkTimeFormat)
	return fmt.Sprintf(MediaFolder, spiderID) + "/" + chunk
}

func DropMediaStorage(spiderID string, bucket res.S3) error {

	// folder without chunk suffix
	folder := fmt.Sprintf(MediaFolder, spiderID)

	storage, err := NewStorage(bucket, folder)
	if err != nil {
		return err
	}

	return storage.Reset()
}
