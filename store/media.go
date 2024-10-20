package store

import (
	"github.com/editorpost/donq/res"
	"time"
)

// JoinTimeStampFolder joins folder and chunk time stamp
func JoinTimeStampFolder(folder string, chunkStamp time.Time) string {
	chunk := chunkStamp.UTC().Format(ChunkTimeFormat)
	return folder + "/" + chunk
}

func DropMediaStorage(folder string, bucket res.S3) error {

	// folder without chunk suffix
	storage, err := NewStorage(bucket, folder)
	if err != nil {
		return err
	}

	return storage.Reset()
}
