package file

import "io"

type Interface interface {
	GetReadCloser(fileName string) (io.ReadCloser, error)
}

// GetReadCloser returns an io.ReadCloser of the specified file if the file name is valid.
func (m *Module) GetReadCloser(fileName string) (io.ReadCloser, error) {
	return m.c.getReadCloser(fileName)
}

type FnModule struct {
	FnGetReadCloser func(fileName string) (io.ReadCloser, error)
}

func (fn *FnModule) GetReadCloser(fileName string) (io.ReadCloser, error) {
	if fn.FnGetReadCloser != nil {
		return fn.FnGetReadCloser(fileName)
	}
	return nil, nil
}
