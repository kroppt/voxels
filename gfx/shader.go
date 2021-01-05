package gfx

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/voxels/log"
)

type Shader struct {
	id uint32
}

// ErrCompileShader indicates that a shader failed to compile
const ErrCompileShader log.ConstErr = "failed to compile shader"

// ErrCreateShader indicates that a shader couldn't be created
const ErrCreateShader log.ConstErr = "failed to create shader"

// NewShader attempts to compile the given shader source code as a shader
// of type shaderType (ex: gl.FRAGMENT_SHADER)
func NewShader(source string, shaderType uint32) (Shader, error) {
	shader := gl.CreateShader(shaderType)
	if shader == 0 {
		return Shader{}, ErrCreateShader
	}

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := string(make([]byte, logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return Shader{}, fmt.Errorf("%w: %v", ErrCompileShader, log)
	}

	return Shader{shader}, nil
}
