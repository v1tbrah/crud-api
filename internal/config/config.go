package config

import (
	"strings"
)

type Config struct {
	runAddress string
	storage    *storage
}

func New(options ...string) (*Config, error) {

	cfg := &Config{}

	storage := &storage{}
	cfg.storage = storage

	for _, opt := range options {
		switch opt {
		case WithFlag:
			cfg.parseFromFlag()
		case WithEnv:
			if err := cfg.parseFromEnv(); err != nil {
				return nil, err
			}
		}
	}

	cfg.setDefaultIfNotConfigured()

	return cfg, nil
}

func (c *Config) setDefaultIfNotConfigured() {

	if c.runAddress == "" {
		c.runAddress = ":3333"
	}
	if c.storage == nil {
		c.storage = &storage{}
	}
	c.storage.setDefaultIfNotConfigured()

}

func (c *Config) RunAddress() string {
	return c.runAddress
}

func (c *Config) TypeOfStorage() string {
	return c.storage.typeOf
}

func (c *Config) PathToFileStorage() string {
	if c.storage.fileStorage == nil {
		return ""
	}
	return c.storage.fileStorage.path
}

func (c *Config) String() string {

	var representMngr strings.Builder

	representMngr.WriteString("Run address: " + "\"" + c.runAddress + "\"")
	representMngr.WriteString("\n")

	if c.storage != nil {

		representMngr.WriteString("Storage: " + c.storage.String())

	} else {

		representMngr.WriteString("Storage is empty")

	}

	return representMngr.String()
}
