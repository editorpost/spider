package store

import (
	"fmt"
	"time"
)

func DefaultStoragePaths() Paths {
	return Paths{
		Collect: "spd/%s/collect",
		Media:   "spd/%s/media",
		Payload: "spd/%s/payload",
	}
}

func CheckStoragePaths() Paths {
	return Paths{
		Collect: "chk/%s/collect",
		Media:   "chk/%s/media",
		Payload: "chk/%s/payload",
	}
}

// Paths defines the path masks and filenames for storing spider data.
// It might be configured to store data in different locations.
//
// Changing the paths will affect the storage location,
// old data won't be moved automatically.
type Paths struct {
	Collect string `json:"collect"`
	Media   string `json:"media"`
	Payload string `json:"payload"`
}

func (paths Paths) joinChunkFolder(folder string) string {
	chunk := time.Now().UTC().Format(ChunkTimeFormat)
	return folder + "/" + chunk
}

func (paths Paths) MediaChunk(arg string) string {
	return paths.joinChunkFolder(paths.MediaRoot(arg))
}

func (paths Paths) MediaURL(arg, proxyURL string) string {
	// public url prefix for media files, e.g. http://my-proxy:8080
	// join public url with bucket folder, e.g. spider/%/media/123.jpg
	// to simplify further proxying the bucket, e.g. http://my-proxy:8080/spider/%/media/123.jpg
	return fmt.Sprintf("%s/%s", proxyURL, paths.MediaChunk(arg))
}

func (paths Paths) MediaRoot(arg string) string {
	return fmt.Sprintf(paths.Media, arg)
}

func (paths Paths) CollectRoot(arg string) string {
	return fmt.Sprintf(paths.Collect, arg)
}

func (paths Paths) PayloadRoot(arg string) string {
	return fmt.Sprintf(paths.Payload, arg)
}
