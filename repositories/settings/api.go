package settings

import (
	"errors"
	"fmt"
	"io"
)

type Interface interface {
	SetFOV(degY float64)
	GetFOV() float64
	SetResolution(width, height uint32)
	GetResolution() (uint32, uint32)
	GetRenderDistance() uint32
	SetFromReader(reader io.Reader) error
}

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
func (r *Repository) SetResolution(width, height uint32) {
	r.c.setResolution(width, height)
}

// GetResolution gets the width and height of the window in pixels.
func (r *Repository) GetResolution() (uint32, uint32) {
	return r.c.getResolution()
}

func (r *Repository) GetRenderDistance() uint32 {
	return r.c.getRenderDistance()
}

// SetFromReader sets repository value from a reader in key=value format.
func (r *Repository) SetFromReader(reader io.Reader) error {
	return r.c.setFromReader(reader)
}

type FnRepository struct {
	FnSetFOV            func(degY float64)
	FnGetFOV            func() float64
	FnSetResolution     func(width, height uint32)
	FnGetResolution     func() (uint32, uint32)
	FnSetFromReader     func(reader io.Reader) error
	FnGetRenderDistance func() uint32
}

func (fn FnRepository) SetFOV(degY float64) {
	if fn.FnSetFOV != nil {
		fn.FnSetFOV(degY)
	}
}

func (fn FnRepository) GetFOV() float64 {
	if fn.FnGetFOV != nil {
		return fn.FnGetFOV()
	}
	return 0
}

func (fn FnRepository) SetResolution(width, height uint32) {
	if fn.FnSetResolution != nil {
		fn.FnSetResolution(width, height)
	}
}

func (fn FnRepository) GetResolution() (uint32, uint32) {
	if fn.FnGetResolution != nil {
		return fn.FnGetResolution()
	}
	return 0, 0
}

func (fn FnRepository) GetRenderDistance() uint32 {
	if fn.FnGetRenderDistance != nil {
		return fn.FnGetRenderDistance()
	}
	return 0
}

func (fn FnRepository) SetFromReader(reader io.Reader) error {
	if fn.FnSetFromReader != nil {
		return fn.FnSetFromReader(reader)
	}
	return nil
}
