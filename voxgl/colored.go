package voxgl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
)

// NewColoredObject returns a newly created Object with the given colors.
//
// Vertices should be vertices of format X, Y, Z, R, G, B, A.
// X, Y, and Z options should be in the range -1.0 to 1.0.
// R, G, B, and A should be in the range 0.0 to 1.0.
func NewColoredObject(vertices [7]float32) (*Object, error) {
	vshad, err := gfx.NewShader(vertColShader, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fshad, err := gfx.NewShader(fragColShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	gshad, err := gfx.NewShader(geoColShader, gl.GEOMETRY_SHADER_ARB)
	if err != nil {
		return nil, err
	}

	prog, err := gfx.NewProgram(vshad, fshad, gshad)
	if err != nil {
		return nil, err
	}

	obj, err := NewObject(prog, vertices[:], []int32{3, 4})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

const vertColShader = `
	#version 420 core

	layout (location = 0) in vec3 pos;
	layout (location = 1) in vec4 col;

	out Vertex {
		vec4 color;
	} OUT;

	void main()
	{
		gl_Position = vec4(pos, 1.0f);
		OUT.color = col;
	}
`

const geoColShader = `
	#version 420 core

	layout(points) in;
	layout(triangle_strip, max_vertices = 14) out;

	layout (std140, binding = 0) uniform Matrices
	{
		mat4 view;
		mat4 projection;
	} cam;

	in Vertex {
		vec4 color;
	} IN[];

	out Vertex {
		vec4 color;
	} OUT;

	void createVertex(vec4 p) {
		gl_Position = cam.projection * cam.view * p;
		OUT.color = IN[0].color;
		EmitVertex();
	}
	
	void main() {
		vec4 center = gl_in[0].gl_Position;

		vec4 dx = vec4(1.0, 0.0, 0.0, 0.0);
		vec4 dy = vec4(0.0, 1.0, 0.0, 0.0);
		vec4 dz = vec4(0.0, 0.0, 1.0, 0.0);

		vec4 p1 = center;
		vec4 p2 = center + dx;
		vec4 p3 = center + dy;
		vec4 p4 = p2 + dy;
		vec4 p5 = p1 + dz;
		vec4 p6 = p2 + dz;
		vec4 p7 = p3 + dz;
		vec4 p8 = p4 + dz;

		createVertex(p7);
		createVertex(p8);
		createVertex(p5);
		createVertex(p6);
		createVertex(p2);
		createVertex(p8);
		createVertex(p4);
		createVertex(p7);
		createVertex(p3);
		createVertex(p5);
		createVertex(p1);
		createVertex(p2);
		createVertex(p3);
		createVertex(p4);
	}
`

const fragColShader = `
	#version 330

	in Vertex {
		vec4 color;
	} IN;

	out vec4 frag_color;

	void main() {
		frag_color = IN.color;
	}
`
