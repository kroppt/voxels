package settings

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/kroppt/voxels/log"
)

type core struct {
	fovY   float64
	width  uint32
	height uint32
}

func (c *core) setFOV(degY float64) {
	c.fovY = degY
}

func (c *core) getFOV() float64 {
	return c.fovY
}

func (c *core) setResolution(width, height uint32) {
	c.width = width
	c.height = height
}

func (c *core) getResolution() (uint32, uint32) {
	return c.width, c.height
}

func (c *core) setFromReader(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
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
			return &ErrParse{
				Line: lineNumber,
				Err:  ErrParseSyntax,
			}
		}
		key := strings.TrimSpace(elements[0])
		value := strings.TrimSpace(elements[1])
		switch key {
		case "fov":
			fov, err := strconv.Atoi(value)
			if err != nil {
				return &ErrParse{
					Line: lineNumber,
					Err:  ErrParseValue,
				}
			}
			c.setFOV(float64(fov))
		case "resolutionX":
			resX, err := strconv.Atoi(value)
			if err != nil || resX < 0 {
				return &ErrParse{
					Line: lineNumber,
					Err:  ErrParseValue,
				}
			}
			c.setResolution(uint32(resX), c.height)
		case "resolutionY":
			resY, err := strconv.Atoi(value)
			if err != nil || resY < 0 {
				return &ErrParse{
					Line: lineNumber,
					Err:  ErrParseValue,
				}
			}
			c.setResolution(c.width, uint32(resY))
		default:
			log.Warnf("invalid settings entry: %v=%v", key, value)
		}
	}
	return nil
}
