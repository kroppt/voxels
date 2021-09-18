package file

import "io"

// GetReadCloser returns an io.ReadCloser of the specified file if the file name is valid.
func (m *Module) GetReadCloser(fileName string) (io.ReadCloser, error) {
	return m.c.getReadCloser(fileName)
}
