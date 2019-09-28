package loader

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// VersionLoader godoc
type VersionLoader interface {
	GetVersion() string
}

type versionLoader struct {
	fs afero.Fs
}

// NewVersionLoader godoc
func NewVersionLoader(fs afero.Fs) VersionLoader {
	return &versionLoader{
		fs: fs,
	}
}

const defaultAppVersion = "dev version"

func (loader *versionLoader) GetVersion() string {
	versionFilePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	file, err := loader.fs.Open(filepath.Join(versionFilePath, "VERSION"))
	if err != nil {
		return defaultAppVersion
	}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text()
	}

	return defaultAppVersion
}
