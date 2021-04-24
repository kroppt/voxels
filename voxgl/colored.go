package voxgl

// NewColoredObject returns a newly created Object with the given colors.
//
// Vertices should be vertices of format X, Y, Z, R, G, B.
// X, Y, and Z options should be in the range -1.0 to 1.0.
// R, G, and B should be in the range 0.0 to 1.0.
func NewColoredObject(vertices [][6]float32) (*Object, error) {
	verts := make([][]float32, 0, len(vertices))
	for i := 0; i < len(vertices); i++ {
		verts = append(verts, vertices[i][:])
	}

	obj, err := NewObject(vertColShader, fragColShader, verts, []int32{3, 3})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

const vertColShader = `
	#version 330 core

	layout (location = 0) in vec3 pos;
	layout (location = 1) in vec3 col;

	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;

	out vec3 color;

	void main()
	{
		gl_Position = projection * view * model * vec4(pos, 1.0f);
		color = col;
	}
`

const fragColShader = `
	#version 330

	in vec3 color;

	out vec4 frag_color;

	void main() {
		frag_color = vec4(color, 1.0);
	}
`
