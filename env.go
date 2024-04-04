package spider

import (
	"github.com/caarlos0/env/v10"
	"os"
	"sync"
)

var (
	_env     = &Windmill{}
	_envOnce = &sync.Once{}
)

// Env loads the Windmill config instance from the environment.
// Variables loaded with sync.Once. You can extend Windmill with setters.
func Env() *Windmill {

	_envOnce.Do(func() {
		if err := env.Parse(_env); err != nil {
			panic(err)
		}
	})

	return _env
}

// Windmill represents configuration settings for a Windmill application.
type Windmill struct {
	// Email holds the email address configured for Windmill.
	Email string `env:"WM_EMAIL"`
	// Username holds the username configured for Windmill.
	Username string `env:"WM_USERNAME"`
	// Workspace holds the workspace configured for Windmill.
	Workspace string `env:"WM_WORKSPACE"`
	// PermissionAs holds the permissioned user for Windmill.
	PermissionAs string `env:"WM_PERMISSIONED_AS"`
	// SHLVL represents the shell level, indicating the nesting depth of the current shell session.
	SHLVL string `env:"SHLVL"`
	// JobID holds the job ID configured for Windmill.
	JobID string `env:"WM_JOB_ID"`
	// FlowJobID holds the flow job ID configured for Windmill.
	FlowJobID string `env:"WM_FLOW_JOB_ID"`
	// RootFlowJobID holds the root flow job ID configured for Windmill.
	RootFlowJobID string `env:"WM_ROOT_FLOW_JOB_ID"`
	// FlowStepID holds the flow step ID configured for Windmill.
	FlowStepID string `env:"WM_FLOW_STEP_ID"`
	// Pwd holds the present working directory configured for Windmill.
	Pwd string `env:"PWD"`
	// HomePath holds the home directory path configured for Windmill.
	HomePath string `env:"HOME"`
	// StatePath holds the state path configured for Windmill.
	StatePath string `env:"WM_STATE_PATH"`
	// StatePathNew holds the new state path configured for Windmill.
	StatePathNew string `env:"WM_STATE_PATH_NEW"`
	// ObjectPath holds the object path configured for Windmill.
	ObjectPath string `env:"WM_OBJECT_PATH"`
	// JobPath holds the job path configured for Windmill.
	JobPath string `env:"WM_JOB_PATH"`
	// SchedulePath holds the schedule path configured for Windmill.
	SchedulePath string `env:"WM_SCHEDULE_PATH"`
	// FlowPath holds the flow path configured for Windmill.
	FlowPath string `env:"WM_FLOW_PATH"`
	// BaseURL holds the base URL configured for Windmill.
	BaseURL string `env:"WM_BASE_URL"`
	// BaseInternalURL holds the base internal URL configured for Windmill.
	BaseInternalURL string `env:"BASE_INTERNAL_URL"`
	// Token holds the token configured for Windmill.
	Token string `env:"WM_TOKEN"`
	// OidcJWT holds the OIDC JWT configured for Windmill.
	OidcJWT string `env:"WM_OIDC_JWT"`
}

// GetEmail returns the email address configured for Windmill.
func (w *Windmill) GetEmail() string {
	return os.Getenv("WM_EMAIL")
}

// GetUsername returns the username configured for Windmill.
func (w *Windmill) GetUsername() string {
	return os.Getenv("WM_USERNAME")
}

// GetWorkspace returns the workspace configured for Windmill.
func (w *Windmill) GetWorkspace() string {
	return os.Getenv("WM_WORKSPACE")
}

// GetPermissionAs returns the permissioned user for Windmill.
func (w *Windmill) GetPermissionAs() string {
	return os.Getenv("WM_PERMISSIONED_AS")
}

// GetSHLVL returns the shell level, indicating the nesting depth of the current shell session.
func (w *Windmill) GetSHLVL() string {
	return os.Getenv("SHLVL")
}

// GetJobID returns the job ID configured for Windmill.
func (w *Windmill) GetJobID() string {
	return os.Getenv("WM_JOB_ID")
}

// GetFlowJobID returns the flow job ID configured for Windmill.
func (w *Windmill) GetFlowJobID() string {
	return os.Getenv("WM_FLOW_JOB_ID")
}

// GetRootFlowJobID returns the root flow job ID configured for Windmill.
func (w *Windmill) GetRootFlowJobID() string {
	return os.Getenv("WM_ROOT_FLOW_JOB_ID")
}

// GetFlowStepID returns the flow step ID configured for Windmill.
func (w *Windmill) GetFlowStepID() string {
	return os.Getenv("WM_FLOW_STEP_ID")
}

// GetPwd returns the present working directory configured for Windmill.
func (w *Windmill) GetPwd() string {
	return os.Getenv("PWD")
}

// GetHomePath returns the home directory path configured for Windmill.
func (w *Windmill) GetHomePath() string {
	return os.Getenv("HOME")
}

// GetStatePath returns the state path configured for Windmill.
func (w *Windmill) GetStatePath() string {
	return os.Getenv("WM_STATE_PATH")
}

// GetStatePathNew returns the new state path configured for Windmill.
func (w *Windmill) GetStatePathNew() string {
	return os.Getenv("WM_STATE_PATH_NEW")
}

// GetObjectPath returns the object path configured for Windmill.
func (w *Windmill) GetObjectPath() string {
	return os.Getenv("WM_OBJECT_PATH")
}

// GetJobPath returns the job path configured for Windmill.
func (w *Windmill) GetJobPath() string {
	return os.Getenv("WM_JOB_PATH")
}

// GetSchedulePath returns the schedule path configured for Windmill.
func (w *Windmill) GetSchedulePath() string {
	return os.Getenv("WM_SCHEDULE_PATH")
}

// GetFlowPath returns the flow path configured for Windmill.
func (w *Windmill) GetFlowPath() string {
	return os.Getenv("WM_FLOW_PATH")
}

// GetBaseURL returns the base URL configured for Windmill.
func (w *Windmill) GetBaseURL() string {
	return os.Getenv("WM_BASE_URL")
}

// GetBaseInternalURL returns the base internal URL configured for Windmill.
func (w *Windmill) GetBaseInternalURL() string {
	return os.Getenv("BASE_INTERNAL_URL")
}

// GetToken returns the token configured for Windmill.
func (w *Windmill) GetToken() string {
	return os.Getenv("WM_TOKEN")
}

// GetOidcJWT returns the OIDC JWT configured for Windmill.
func (w *Windmill) GetOidcJWT() string {
	return os.Getenv("WM_OIDC_JWT")
}
