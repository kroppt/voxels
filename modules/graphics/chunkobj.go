package graphics

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

type ChunkObject struct {
	program gfx.Program
	vao     gfx.VAO
}

// NewChunkObject returns a newly created ChunkObject with the given vertices.
func NewChunkObject() (*ChunkObject, error) {
	prog, err := GetProgram(vertShader, fragFrameShader, geoFrameShader)
	if err != nil {
		return nil, err
	}
	vao := gfx.NewVAO(gl.POINTS, []int32{4, 1})

	return &ChunkObject{
		program: prog,
		vao:     *vao,
	}, nil
}

// SetData uploads data to OpenGL.
func (co *ChunkObject) SetData(data []float32) {
	err := co.vao.Load(data, gl.STATIC_DRAW)
	if err != nil {
		panic("failed to set data")
	}
}

// Render generates an image of the object with OpenGL.
func (co *ChunkObject) Render() {
	co.program.Bind()
	co.vao.Draw()
	co.program.Unbind()
}

// Destroy frees external resources.
func (co *ChunkObject) Destroy() {
	// o.program.Destroy() // TODO store and delete in world
	co.vao.Destroy()
}

var progMap map[string]gfx.Program

func GetProgram(vshadstr, fshadstr, gshadstr string) (gfx.Program, error) {
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

const vertShader = `
	#version 420 core

	layout (location = 0) in vec4 pos;
	layout (location = 1) in float lighting;

	out Vertex {
		float vbits;
		float lighting;
	} OUT;

	void main()
	{
		gl_Position = vec4(pos.xyz, 1.0f);
		OUT.vbits = pos[3];
		OUT.lighting = lighting;
	}
`

const geoFrameShader = `
	#version 420 core

	layout(points) in;
	layout(line_strip, max_vertices = 30) out;

	layout (std140, binding = 0) uniform Matrices
	{
		dmat4 view;
		dmat4 projection;
	} cam;

	in Vertex {
		float vbits;
		float lighting;
	} IN[];

	out Vertex {
		vec4 col;
	} OUT;

	void createVertex(vec4 p) {
		gl_Position = vec4(cam.projection * cam.view * p);
		OUT.col = vec4(0.8, 0.8, 0.8, 1.0);
		EmitVertex();
	}

	void createQuad(vec4 p, vec4 d1, vec4 d2) {
		createVertex(p);
		createVertex(p+d1);
		createVertex(p+d1+d2);
		createVertex(p+d2);
		createVertex(p);
		EndPrimitive();
	}

	void main() {
		vec4 origin = gl_in[0].gl_Position;

		// bottom 6 bits are for adjacency
		int adjaBits = 6;
		// bit order = right left top bottom backward forward
		int rightmask = 0x20;
		int leftmask = 0x10;
		int topmask = 0x08;
		int bottommask = 0x04;
		int backwardmask = 0x02;
		int forwardmask = 0x01;
		int bits = int(IN[0].vbits);

		// top 26 bits are for block types (for now)
		int blockmask = 0xFFFFFFC0;
		int blockType = (blockmask & bits) >> adjaBits;
		if (blockType == 0) {
			return; // render nothing if air block
		}

		vec4 dx = vec4(1.0, 0.0, 0.0, 0.0);
		vec4 dy = vec4(0.0, 1.0, 0.0, 0.0);
		vec4 dz = vec4(0.0, 0.0, 1.0, 0.0);
		vec4 p1 = origin;
		vec4 p2 = p1 + dx + dy + dz;

		if ((bits & backwardmask) - backwardmask != 0) {
			createQuad(p2, -dx, -dy); // backward
		}
		if ((bits & forwardmask) - forwardmask != 0) {
			createQuad(p1, dy, dx); // forward
		}
		if ((bits & topmask) - topmask != 0) {
			createQuad(p2, -dz, -dx); // top
		}
		if ((bits & bottommask) - bottommask != 0) {
			createQuad(p1, dx, dz); // bottom
		}
		if ((bits & rightmask) - rightmask != 0) {
			createQuad(p2, -dy, -dz); // right
		}
		if ((bits & leftmask) - leftmask != 0) {
			createQuad(p1, dz, dy); // left
		}
	}
`

const fragFrameShader = `
	#version 400

	in Vertex {
		vec4 col;
	} IN;

	out vec4 frag_color;

	void main() {
		frag_color = IN.col;
	}
`
