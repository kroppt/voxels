// +build !test

package voxgl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

// NewColoredObject returns a newly created Object with the given colors.
//
// Vertices should be vertices of format X, Y, Z, AdjacencyBits, R, G, B, A.
// X, Y, and Z options should be in the range -1.0 to 1.0.
// AdjacencyBits is a float where the least significant 6 bits
// represent whether each of the left right top bottom forward backward faces
// can be seen.
// R, G, B, and A should be in the range 0.0 to 1.0.
func NewColoredObject(vertices []float32) (*Object, error) {
	prog, err := GetProgram(vertColShader, fragColShader, geoColShader)
	if err != nil {
		return nil, err
	}

	obj, err := NewObject(prog, vertices, []int32{4, 4})
	if err != nil {
		return nil, err
	}

	return obj, nil
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

const vertColShader = `
	#version 420 core

	layout (location = 0) in vec4 pos;
	layout (location = 1) in vec4 col; // TODO delete me eventually

	out Vertex {
		vec4 color;
		float vbits;
	} OUT;

	void main()
	{
		gl_Position = vec4(pos.xyz, 1.0f);
		OUT.color = col;
		OUT.vbits = pos[3];
	}
`

const geoColShader = `
	#version 420 core

	layout(points) in;
	layout(triangle_strip, max_vertices = 24) out;

	layout (std140, binding = 0) uniform Matrices
	{
		mat4 view;
		mat4 projection;
	} cam;

	in Vertex {
		vec4 color;
		float vbits;
	} IN[];

	out Vertex {
		vec4 color;
		vec3 stdir;
		flat int blockType;
	} OUT;

	void createVertex(vec4 p) {
		vec3 center = (gl_in[0].gl_Position).xyz + 0.5;
		OUT.stdir = p.xyz - center;
		gl_Position = cam.projection * cam.view * p;
		OUT.color = IN[0].color;
		EmitVertex();
	}

	void createQuad(vec4 p, vec4 d1, vec4 d2) {
		createVertex(p);
		createVertex(p+d1);
		createVertex(p+d2);
		createVertex(p+d1+d2);
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
		OUT.blockType = (blockmask & bits) >> adjaBits;

		vec4 dx = vec4(1.0, 0.0, 0.0, 0.0);
		vec4 dy = vec4(0.0, 1.0, 0.0, 0.0);
		vec4 dz = vec4(0.0, 0.0, 1.0, 0.0);
		vec4 p1 = origin;
		vec4 p2 = p1 + dx + dy + dz;

		if ((bits & backwardmask) - backwardmask == 0) {
			createQuad(p2, -dx, -dy); // backward
		}
		if ((bits & forwardmask) - forwardmask == 0) {
			createQuad(p1, dy, dx); // forward
		}
		if ((bits & topmask) - topmask == 0) {
			createQuad(p2, -dz, -dx); // top
		}
		if ((bits & bottommask) - bottommask == 0) {
			createQuad(p1, dx, dz); // bottom
		}
		if ((bits & rightmask) - rightmask == 0) {
			createQuad(p2, -dy, -dz); // right
		}
		if ((bits & leftmask) - leftmask == 0) {
			createQuad(p1, dz, dy); // left
		}
	}
`

const fragColShader = `
	#version 400

	in Vertex {
		vec4 color;
		vec3 stdir;
		flat int blockType;
	} IN;
	uniform samplerCubeArray cubeMapArray;


	out vec4 frag_color;

	void main() {
		frag_color = texture(cubeMapArray, vec4(IN.stdir, IN.blockType));// * IN.color;
	}
`
