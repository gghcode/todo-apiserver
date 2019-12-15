package app

import (
	"os"
	"path/filepath"
)

type (
	// FileReader read file
	FileReader interface {
		ReadString(path string) (string, error)
	}

	appService struct {
		reader FileReader
	}
)

// NewService return appService
func NewService(reader FileReader) UsecaseInteractor {
	return &appService{
		reader: reader,
	}
}

const _defaultVersion = "dev version"

func (srv *appService) Version() string {
	versionFileDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	versionFilePath := filepath.Join(versionFileDir, "VERSION")

	version, err := srv.reader.ReadString(versionFilePath)
	if err != nil {
		return _defaultVersion
	}

	return version
}
