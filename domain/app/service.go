package app

import (
	"os"
	"path/filepath"
)

// FileReader read file
type FileReader interface {
	ReadString(path string) (string, error)
}

type appService struct {
	reader FileReader
}

// NewService return appService
func NewService(reader FileReader) UsecaseInteractor {
	return &appService{
		reader: reader,
	}
}

const _defaultVersion = "dev version"

func (srv *appService) Version() string {
	versionFilePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	version, err := srv.reader.ReadString(versionFilePath)
	if err != nil {
		return _defaultVersion
	}

	return version
}