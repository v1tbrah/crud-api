package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

func (c *Config) parseFromFlag() {

	flag.StringVar(&c.runAddress, "a", "", "api server run address")

	var storageType string
	flag.StringVar(&storageType, "st", "", "type of data storage")

	var pathToFileStorage string
	flag.StringVar(&pathToFileStorage, "pf", "", "path to file data storage")

	flag.Parse()

	fmt.Println(storageType)
	fmt.Println(pathToFileStorage)

	if storageType != "" {

		if c.storage == nil {
			c.storage = &storage{}
		}

		c.storage.typeOf = storageType

		if storageType == "file" {

			c.storage.fileStorage = &fileStorage{}

			if pathToFileStorage != "" {
				c.storage.fileStorage.path = pathToFileStorage
			}
		}
	}

}

func (c *Config) parseFromEnv() error {

	envConfig := struct {
		RunAddress        string `env:"RUN_ADDRESS"`
		StorageType       string `env:"STORAGE_TYPE"`
		PathToFileStorage string `env:"PATH_FILE_STORAGE"`
	}{}

	if err := env.Parse(&envConfig); err != nil {
		return err
	}

	if envConfig.StorageType != "" {

		if c.storage == nil {
			c.storage = &storage{}
		}

		c.storage.typeOf = envConfig.StorageType

		if envConfig.StorageType == "file" {

			c.storage.fileStorage = &fileStorage{}

			if envConfig.PathToFileStorage != "" {
				c.storage.fileStorage.path = envConfig.PathToFileStorage
			}
		}
	}

	return nil
}
