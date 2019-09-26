package loader

import (
	"bufio"

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
	file, err := loader.fs.Open("VERSION")
	if err != nil {
		return defaultAppVersion
	}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text()
	}

	return defaultAppVersion
}
