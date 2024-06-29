//go:build e2e
// +build e2e

package store_test

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/editorpost/spider/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var bucket = store.Bucket{
	Name:   "ep-spider",
	URL:    "https://s3.ap-southeast-1.wasabisys.com",
	Region: "ap-southeast-1",
}

func TestMain(m *testing.M) {

	m.Run()
}

func TestNewS3Client(t *testing.T) {

	tests := []struct {
		name    string
		want    *s3.Client
		wantErr bool
	}{
		{
			name:    "test",
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.NewS3Client(bucket)
			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}

func TestBucketFolder_Save(t *testing.T) {

	client, err := store.NewS3Client(bucket)
	require.NoError(t, err)

	bf, err := store.NewBucketFolder("test", bucket.Name, client)
	require.NoError(t, err)

	data := []byte("test data")
	filename := "testfile.txt"
	_ = bf.Path(filename)

	_, err = bf.Save(data, filename)
	assert.NoError(t, err)
	// assert.Equal(t, fullPath, location)
}
