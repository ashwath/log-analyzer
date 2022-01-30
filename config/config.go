package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var cfg Config

func init() {
	err := envconfig.Process("log-ananlyzer-api", &cfg)
	if err != nil {
		log.Fatalf("Unable to load API config, error: %v", err)
	}
}

// Get returns the config struct
func Get() Config {
	return cfg
}

// Config holds all configuration coming from env
// In admiral deployments these values can be controlled
// in the configuration panel that is populated by clipper using the configuration-schema.yaml
type Config struct {
	// config schema values
	NumberOfLogLines int `default:"10"`
}
