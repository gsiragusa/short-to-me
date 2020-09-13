package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	*logrus.Logger

	// MongoUri is the connection string to mongo
	MongoUri    string `split_words:"true" default:"mongodb://localhost:27017"`
	MongoDbName string `split_words:"true" default:"short-to-me"`

	// Port is the port to run the HTTP server on
	Port int `split_words:"true" default:"8081"`
}

func Configure() (*AppConfig, error) {
	conf := &AppConfig{}
	if err := load(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// load accepts a struct to load the environment configuration from
func load(config interface{}) error {
	return envconfig.Process("", config)
}
