package setup

import (
	"fmt"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/store"
)

func (s *Spider) withStorage(deps *config.Deps) error {

	if s.Deploy.Storage.Bucket == "" {
		return nil
	}

	if err := s.withCollectStore(deps); err != nil {
		return err
	}

	if err := WithFn(
		s.withExtractStore,
		s.withExtractIndex,
		s.withMedia,
	); err != nil {
		return err
	}

	return nil
}

func (s *Spider) withCollectStore(deps *config.Deps) error {

	folder := fmt.Sprintf(s.Deploy.Paths.Collect, s.ID)
	storage, upload, err := store.NewCollectStorage(folder, s.Deploy.Storage)
	if err != nil {
		return err
	}

	// upload visited urls to S3
	s.onShutdown(upload)
	deps.Storage = storage

	return err
}

func (s *Spider) withExtractStore() error {

	extractStore, err := store.NewExtractStorage(s.Deploy.Paths.Payload, s.Deploy.Storage)
	if err != nil {
		return fmt.Errorf("failed to create extract S3 storage: %w", err)
	}

	// provide save extractor func
	s.pipe.Finisher(extractStore.Save)

	return nil
}

func (s *Spider) withExtractIndex() error {

	if len(s.Deploy.Database.Host) == 0 {
		return nil
	}

	extractIndex, err := store.NewExtractIndex(s.ID, s.Deploy.Database.DSN())
	if err != nil {
		return fmt.Errorf("failed to create extract index store: %w", err)
	}

	// provide save extractor func
	s.pipe.Finisher(extractIndex.Save)

	return nil
}

func (s *Spider) withMedia() error {

	if !s.Extract.Media.Enabled {
		return nil
	}

	folder := s.Deploy.Paths.MediaChunk(s.ID)

	storage, err := store.NewStorage(s.Deploy.Media.S3, folder)
	if err != nil {
		return err
	}

	proxyURL := s.Deploy.Paths.MediaURL(s.ID, s.Deploy.Media.PublicURL)
	uploader := media.NewMedia(proxyURL, media.NewLoader(storage))

	s.pipe.Starter(uploader.Claims)
	s.pipe.Finisher(uploader.Upload)

	return nil
}

func WithFn(fns ...func() error) error {

	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// WithDepsFn runs functions with dependencies from one to one passing result to the next.
func WithDepsFn(deps *config.Deps, fns ...func(*config.Deps) error) error {

	for _, fn := range fns {
		if err := fn(deps); err != nil {
			return err
		}
	}

	return nil
}
