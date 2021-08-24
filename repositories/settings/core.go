package settings

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/kroppt/voxels/log"
)

type fileMod interface {
	GetFileReader(string) io.Reader
}

type core struct {
	fileMod fileMod
	fovY    float32
	width   int32
	height  int32
}

func (c *core) setFOV(degY float32) {
	c.fovY = degY
}

func (c *core) getFOV() float32 {
	return c.fovY
}

func (c *core) setResolution(width, height int32) {
	c.width = width
	c.height = height
}

func (c *core) getResolution() (int32, int32) {
	return c.width, c.height
}

func (c *core) setFromReader(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		err := scanner.Err()
		if err != nil {
			return err
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		elements := strings.Split(line, "=")
		if len(elements) != 2 {
			return errors.New("malformed settings line: expected key=value")
		}
		key := strings.TrimSpace(elements[0])
		value := strings.Trim(elements[1], "\t ")
		switch key {
		case "fov":
			fov, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			c.setFOV(float32(fov))
		case "resolutionX":
			resX, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			c.setResolution(int32(resX), c.height)
		case "resolutionY":
			resY, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			c.setResolution(c.width, int32(resY))
		default:
			log.Warnf("invalid settings entry: %v=%v", key, value)
		}
	}
	return nil
}
