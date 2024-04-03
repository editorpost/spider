package config

import (
	"github.com/caarlos0/env/v10"
)

type Env struct {
	TaskID        string `env:"CRAWLAB_TASK_ID"`
	DbAuthSource  string `env:"CRAWLAB_MONGO_AUTHSOURCE"`
	DbUser        string `env:"CRAWLAB_MONGO_USERNAME"`
	DbPass        string `env:"CRAWLAB_MONGO_PASSWORD"`
	DbHost        string `env:"CRAWLAB_MONGO_HOST"`
	DbPort        string `env:"CRAWLAB_MONGO_PORT"`
	DbName        string `env:"CRAWLAB_MONGO_DB"`
	DbCollection  string `env:"CRAWLAB_COLLECTION"`
	DbUniqueField string `env:"CRAWLAB_COLLECTION_UNIQUE"`
	GrpcAuthKey   string `env:"CRAWLAB_GRPC_AUTH_KEY"`
	GrpcAddress   string `env:"CRAWLAB_GRPC_ADDRESS"`
	DockerNode    string `env:"CRAWLAB_DOCKER"`
	MasterNode    string `env:"CRAWLAB_NODE_MASTER"`
}

func GetEnv() Env {
	cfg := Env{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
