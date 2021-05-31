package voxgl

import (
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/world"
)

// TexturedObject is a renderable set of vertices.
type TexturedObject struct {
	*Object
	tex gfx.Texture
}

// NewTexturedObject returns a newly created TexturedObject with the given texture.
//
// Vertices should be vertices of format x, y, z, s, t.
func NewTexturedObject(vertices [][5]float32, texture gfx.Texture) (*TexturedObject, error) {
	verts := make([][]float32, 0, len(vertices))
	for _, vs := range vertices {
		verts = append(verts, vs[:])
	}

	obj, err := NewObject(vertTexShader, fragTexShader, verts, []int32{3, 2})
	if err != nil {
		return nil, err
	}

	return &TexturedObject{
		Object: obj,
		tex:    texture,
	}, nil
}

// Render generates an image of the object with OpenGL.
func (to *TexturedObject) Render(cam world.Camera) {
	to.tex.Bind()
	to.Object.Render(cam)
	to.tex.Unbind()
}

// Destroy frees external resources.
func (to *TexturedObject) Destroy() {
	to.Object.Destroy()
	to.tex.Destroy()
}

const vertTexShader = `
	#version 420 core

	layout (std140, binding = 0) uniform Matrices
	{
		mat4 view;
		mat4 projection;
	};
	uniform mat4 model;

	layout (location = 0) in vec3 pos;
	layout (location = 1) in vec2 tex_coord;

	out vec2 tex_coords;

	void main()
	{
		gl_Position = projection * view * model * vec4(pos, 1.0f);
		tex_coords = vec2(tex_coord.x, tex_coord.y);
	}
`

const fragTexShader = `
	#version 330

	uniform sampler2D frag_tex;

	in vec2 tex_coords;

	out vec4 frag_color;

	void main() {
		frag_color = texture(frag_tex, tex_coords);
	}
`
