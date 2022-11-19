package config

import "strings"

type storage struct {
	typeOf      string
	fileStorage *fileStorage
}

func (s *storage) setDefaultIfNotConfigured() {

	if s.typeOf == "" || s.typeOf == "file" {

		s.typeOf = "file"

		if s.fileStorage == nil {
			s.fileStorage = &fileStorage{}
		}
		s.fileStorage.setDefaultIfNotConfigured()

	}

}

func (s *storage) String() string {

	var representMngr strings.Builder

	representMngr.WriteString("Type: " + "\"" + s.typeOf + "\"")

	if s.fileStorage != nil {

		representMngr.WriteString(" File storage: " + s.fileStorage.String())

	} else {

		representMngr.WriteString(" File storage is empty")

	}

	return representMngr.String()
}

type fileStorage struct {
	path string
}

func (f *fileStorage) setDefaultIfNotConfigured() {

	if f.path == "" {
		f.path = "users.json"
	}

}

func (f *fileStorage) String() string {
	return "Path to file storage: " + "\"" + f.path + "\""
}
