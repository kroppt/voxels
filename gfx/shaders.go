package gfx

const (
	SampleCubeVertex = `
	#version 330
	layout (location = 0) in vec3 position_in;
	layout (location = 1) in vec3 fragColor;
	out vec3 color;
	void main() {
		gl_Position = vec4(position_in, 1.0);
		color = fragColor;
	}`

	SampleCubeFragment = `
	#version 330
	out vec4 frag_color;
	in vec3 color;
	void main() {
		frag_color = vec4(color, 1.0);
	}`
)
