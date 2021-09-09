package file

import "io"

// GetReadCloser returns an io.ReadCloser if the file name is valid.
func (m *Module) GetReadCloser(fileName string) (io.ReadCloser, error) {
	return m.c.getReadCloser(fileName)
}
