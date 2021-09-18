package file

import (
	"io"
	"os"
)

type core struct {
}

func (c *core) getReadCloser(fileName string) (io.ReadCloser, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}
