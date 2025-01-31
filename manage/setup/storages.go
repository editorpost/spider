package setup

import (
	"fmt"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/store"
	"log/slog"
)

func (s *Spider) withStorage(deps *config.Deps) error {

	if s.Deploy.Storage.Bucket == "" {
		return nil
	}

	if err := s.withVisitedHistory(deps); err != nil {
		return err
	}

	if err := WithFn(
		s.withExtractHistory,
		s.withExtractStore,
		s.withExtractIndex,
		s.withMedia,
	); err != nil {
		return err
	}

	return nil
}

func (s *Spider) withVisitedHistory(deps *config.Deps) error {

	// visit once,
	// stores collector history in S3 between runs
	if !s.Collect.VisitOnce {
		return nil
	}

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

	extractStore, err := store.NewExtractStorage(s.Deploy.Paths.PayloadRoot(s.ID), s.Deploy.Storage)
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

	payloads, err := store.NewSpiderPayloads(s.ID, s.Deploy.Database.DSN(), s.Deploy.Paths)
	if err != nil {
		return fmt.Errorf("failed to create extract index store: %w", err)
	}

	// provide save extractor func
	s.pipe.Finisher(payloads.Save)

	return nil
}

func (s *Spider) withExtractHistory() error {

	if !s.Extract.ExtractOnce {
		return nil
	}

	// load extracted urls from spider payloads database
	payloads, err := store.NewSpiderPayloads(s.ID, s.Deploy.Database.DSN(), s.Deploy.Paths)
	if err != nil {
		return fmt.Errorf("failed to create extract index store: %w", err)
	}

	// @todo: provide iterator, load from db in chunks
	s.pipe.History().Init(func() []string {
		urls, loadErr := payloads.URLs()
		if loadErr != nil {
			slog.Error("extractor history loading", "err", loadErr.Error())
			return nil
		}

		return urls
	})

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
