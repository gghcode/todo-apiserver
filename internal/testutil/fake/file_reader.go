package fake

import "github.com/stretchr/testify/mock"

// FileReader fake file reader
type FileReader struct {
	mock.Mock
}

// NewFileReader return fake file reader
func NewFileReader() *FileReader {
	return &FileReader{}
}

// ReadString godoc
func (reader *FileReader) ReadString(path string) (string, error) {
	args := reader.Called(path)
	return args.String(0), args.Error(1)
}
