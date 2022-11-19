package storage

import (
	"errors"

	"github.com/rs/zerolog/log"

	"refactoring/internal/model"
	"refactoring/internal/storage/file"
)

var (
	ErrEmptyConfig            = errors.New("config is empty")
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
)

type Storage interface {
	GetAllUsers() (allUsers *model.UserStore, err error)
	CreateUser(newUser *model.User) (id int64, err error)
	GetUser(id int64) (user *model.User, err error)
	UpdateUser(id int64, newDisplayName string) (err error)
	DeleteUser(id int64) (err error)
}

func New(config Config) (newStorage Storage, err error) {
	log.Debug().Str("config", config.String()).Msg("storage.New START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("storage.New END")
		} else {
			log.Debug().Msg("storage.New END")
		}
	}()

	if config == nil {
		return nil, ErrEmptyConfig
	}

	if config.TypeOfStorage() == "file" {
		return file.New(config.PathToFileStorage())
	}

	return nil, ErrUnsupportedStorageType
}
