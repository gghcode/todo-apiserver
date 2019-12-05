package file

import (
	"bufio"

	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/spf13/afero"
)

type aferoReader struct {
	fs afero.Fs
}

// NewAferoFileReader return aferoFileReader
func NewAferoFileReader(fs afero.Fs) app.FileReader {
	return &aferoReader{
		fs: fs,
	}
}

func (reader *aferoReader) ReadString(path string) (string, error) {
	file, err := reader.fs.Open(path)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", nil
}
