package pipe

import (
	"github.com/caarlos0/env/v10"
	"log/slog"
	"sync"
)

var (
	_job     = &Job{}
	_envOnce = &sync.Once{}
)

type Job struct {
	// JobProvider holds the provider configured for Windmill.
	// Must be set by provider
	Provider string `env:"SPIDER_JOB_PROVIDER" envDefault:"windmill"`
	// Username holds the username configured for Windmill.
	ID string `env:"SPIDER_JOB_ID"`
}

// JobFromEnv loads the job info from the environment.
func JobFromEnv() *Job {

	_envOnce.Do(func() {
		if err := env.Parse(_job); err != nil {
			_job.ID = ""
			_job.Provider = "unknown"
			slog.Error("failed to parse provider job from env", "err", err.Error())
		}
	})

	return _job
}

func JobMetadata(p *Payload) error {

	job := JobFromEnv()

	p.JobID = job.ID
	p.JobProvider = job.Provider

	return nil
}
