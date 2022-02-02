package graphics

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

type glObject struct {
	program gfx.Program
	vao     gfx.VAO
}

// newChunkObject returns a renderable chunk.
func newChunkObject() (*glObject, error) {
	// swap these lines for frame outlines or solid blocks
	prog, err := getProgram(vertColShader, fragColShader, geoColShader)
	// prog, err := GetProgram(vertColShader, fragFrameShader, geoFrameShader)

	if err != nil {
		return nil, err
	}
	vao := gfx.NewVAO(gl.POINTS, []int32{4, 1})

	return &glObject{
		program: prog,
		vao:     *vao,
	}, nil
}

// newFrameObject returns a renderable frame.
func newFrameObject() (*glObject, error) {
	prog, err := getProgram(vertColShader, fragFrameShader, geoFrameShader)

	if err != nil {
		return nil, err
	}
	vao := gfx.NewVAO(gl.POINTS, []int32{4, 1})

	return &glObject{
		program: prog,
		vao:     *vao,
	}, nil
}

func newCrosshairObject(size, aspect float32) (*glObject, error) {
	vao := gfx.NewVAO(gl.LINES, []int32{2, 4})
	vertices := []float32{
		-size / aspect, 0, 0.0, 1.0, 1.0, 1.0,
		size / aspect, 0, 1.0, 1.0, 0.0, 1.0,
		0, -size, 1.0, 0.0, 1.0, 1.0,
		0, size, 0.0, 1.0, 0.0, 1.0,
	}
	vao.Load(vertices, gl.STATIC_DRAW)

	vshad, err := gfx.NewShader(vertCrossShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fshad, err := gfx.NewShader(fragCrossShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog, err := gfx.NewProgram(vshad, fshad)
	if err != nil {
		return nil, err
	}

	prog.UploadUniform("aspect", aspect)

	return &glObject{
		program: prog,
		vao:     *vao,
	}, nil
}

// setData uploads data to OpenGL.
func (co *glObject) setData(data []float32) {
	err := co.vao.Load(data, gl.STATIC_DRAW)
	if err != nil {
		panic("failed to set data")
	}
}

// render generates an image of the object with OpenGL.
func (co *glObject) render() {
	co.program.Bind()
	co.vao.Draw()
	co.program.Unbind()
}

// destroy frees external resources.
func (co *glObject) destroy() {
	// o.program.Destroy() // TODO store and delete in world
	co.vao.Destroy()
}

var progMap map[string]gfx.Program

func getProgram(vshadstr, fshadstr, gshadstr string) (gfx.Program, error) {
	if progMap == nil {
		progMap = make(map[string]gfx.Program)
	}
	key := vshadstr + fshadstr + gshadstr
	if prog, ok := progMap[key]; ok {
		return prog, nil
	}
	vshad, err := gfx.NewShader(vshadstr, gl.VERTEX_SHADER)
	if err != nil {
		return gfx.Program{}, err
	}

	fshad, err := gfx.NewShader(fshadstr, gl.FRAGMENT_SHADER)
	if err != nil {
		return gfx.Program{}, err
	}

	gshad, err := gfx.NewShader(gshadstr, gl.GEOMETRY_SHADER_ARB)
	if err != nil {
		return gfx.Program{}, err
	}

	prog, err := gfx.NewProgram(vshad, fshad, gshad)
	if err != nil {
		return gfx.Program{}, err
	}
	progMap[key] = prog
	return prog, nil
}
