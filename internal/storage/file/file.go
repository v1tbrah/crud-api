package file

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"

	"refactoring/internal/model"
	storageErr "refactoring/internal/storage/errors"
)

type StorageFile struct {
	path string
	mu   sync.Mutex
}

func New(path string) (newStorage *StorageFile, err error) {
	log.Debug().Str("path", path).Msg("file.New START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("file.New END")
		} else {
			log.Debug().Msg("file.New END")
		}
	}()

	newStorage = &StorageFile{}

	newStorage.path = path

	return newStorage, nil
}

func (s *StorageFile) GetAllUsers() (allUsers *model.UserStore, err error) {
	log.Debug().Msg("file.GetAllUsers START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("file.GetAllUsers END")
		} else {
			log.Debug().Msg("file.GetAllUsers END")
		}
	}()

	dataUsers, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("opening file with users: %w", err)
	}

	if err = json.Unmarshal(dataUsers, &allUsers); err != nil {
		return nil, fmt.Errorf("unmarshalling file with users: %w", err)
	}

	return allUsers, nil

}

func (s *StorageFile) CreateUser(newUser *model.User) (id int64, err error) {
	log.Debug().Str("new user", newUser.String()).Msg("file.CreateUser START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("file.CreateUser END")
		} else {
			log.Debug().Msg("file.CreateUser END")
		}
	}()

	s.mu.Lock()
	defer s.mu.Unlock()

	allUsers, err := s.GetAllUsers()
	if err != nil {
		return 0, fmt.Errorf("reading file with users: %w", err)
	}

	if len(allUsers.List) == 0 {
		allUsers.List = make(model.UserList, 1)
	}

	allUsers.Increment++
	fmt.Println("increment: ", allUsers.Increment)

	idForList := strconv.FormatInt(allUsers.Increment, 10)
	allUsers.List[idForList] = *newUser

	dataForFile, err := json.Marshal(&allUsers)
	if err != nil {
		return 0, fmt.Errorf("marshalling users for file: %w", err)
	}

	err = os.WriteFile(s.path, dataForFile, fs.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("writing new user to file: %w", err)
	}

	return allUsers.Increment, nil
}

func (s *StorageFile) GetUser(id int64) (user *model.User, err error) {
	log.Debug().Msg("file.GetUser START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("file.GetUser END")
		} else {
			log.Debug().Msg("file.GetUser END")
		}
	}()

	allUsers, err := s.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("reading file with users: %w", err)
	}

	if len(allUsers.List) == 0 {
		return nil, storageErr.ErrUserIsNotFound
	}

	idForList := strconv.FormatInt(id, 10)
	userFromList, ok := allUsers.List[idForList]
	if !ok {
		return nil, storageErr.ErrUserIsNotFound
	}

	return &userFromList, nil
}

func (s *StorageFile) UpdateUser(id int64, newDisplayName string) (err error) {
	log.Debug().Msg("file.UpdateUser START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("file.UpdateUser END")
		} else {
			log.Debug().Msg("file.UpdateUser END")
		}
	}()

	allUsers, err := s.GetAllUsers()
	if err != nil {
		return fmt.Errorf("reading file with users: %w", err)
	}

	if len(allUsers.List) == 0 {
		return storageErr.ErrUserIsNotFound
	}

	idForList := strconv.FormatInt(id, 10)
	userFromList, ok := allUsers.List[idForList]
	if !ok {
		return storageErr.ErrUserIsNotFound
	}

	userFromList.DisplayName = newDisplayName

	allUsers.List[idForList] = userFromList

	dataForFile, err := json.Marshal(&allUsers)
	if err != nil {
		return fmt.Errorf("marshalling users for file: %w", err)
	}

	err = os.WriteFile(s.path, dataForFile, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("writing updated user to file: %w", err)
	}

	return nil
}

func (s *StorageFile) DeleteUser(id int64) (err error) {
	log.Debug().Msg("file.DeleteUser START")
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("file.DeleteUser END")
		} else {
			log.Debug().Msg("file.DeleteUser END")
		}
	}()

	allUsers, err := s.GetAllUsers()
	if err != nil {
		return fmt.Errorf("reading file with users: %w", err)
	}

	if len(allUsers.List) == 0 {
		return storageErr.ErrUserIsNotFound
	}

	idForList := strconv.FormatInt(id, 10)
	_, ok := allUsers.List[idForList]
	if !ok {
		return storageErr.ErrUserIsNotFound
	}

	delete(allUsers.List, idForList)

	dataForFile, err := json.Marshal(&allUsers)
	if err != nil {
		return fmt.Errorf("marshalling users for file: %w", err)
	}

	err = os.WriteFile(s.path, dataForFile, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("writing users to file: %w", err)
	}

	return nil
}
