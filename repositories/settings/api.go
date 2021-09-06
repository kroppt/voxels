package settings

import (
	"errors"
	"fmt"
	"io"
)

// ConstErr is a constant error type.
type ConstErr string

// Error returns the error message.
func (e ConstErr) Error() string {
	return string(e)
}

// ErrParseSyntax indicates that the settings failed to parse the syntax.
const ErrParseSyntax ConstErr = "syntax should be key=value format"

// ErrParseValue indicates that the settings failed to parse a value.
const ErrParseValue ConstErr = "value invalid"

// ErrParse indicates that the settings failed to parse.
type ErrParse struct {
	Line int
	Err  error
}

func (e ErrParse) Error() string {
	return fmt.Sprintf("%v at line %v", e.Err, e.Line)
}

// Is returns the value of performing errors.Is on the wrapped error.
func (e ErrParse) Is(err error) bool {
	return errors.Is(e.Err, err)
}

// SetFOV sets the vertical field of view.
func (r *Repository) SetFOV(degY float64) {
	r.c.setFOV(degY)
}

// GetFOV gets the vertical field of view.
func (r *Repository) GetFOV() float64 {
	return r.c.getFOV()
}

// SetResolution sets the width and height of the window in pixels.
func (r *Repository) SetResolution(width, height int32) {
	r.c.setResolution(width, height)
}

// GetResolution gets the width and height of the window in pixels.
func (r *Repository) GetResolution() (int32, int32) {
	return r.c.getResolution()
}

// SetFromReader sets repository value from a reader in key=value format.
func (r *Repository) SetFromReader(reader io.Reader) error {
	return r.c.setFromReader(reader)
}
