package voxgl

func NewFrame(vertices []float32) (*Object, error) {
	prog, err := GetProgram(vertFrameShader, fragFrameShader, geoFrameShader)
	if err != nil {
		return nil, err
	}

	obj, err := NewObject(prog, vertices, []int32{4})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

const vertFrameShader = `
	#version 420 core

	layout (location = 0) in vec4 pos;

	out Vertex {
		float vbits;
	} OUT;

	void main()
	{
		gl_Position = vec4(pos.xyz, 1.0f);
		OUT.vbits = pos[3]; 
	}
`

const geoFrameShader = `
	#version 420 core

	layout(points) in;
	layout(line_strip, max_vertices = 30) out;

	layout (std140, binding = 0) uniform Matrices
	{
		mat4 view;
		mat4 projection;
	} cam;

	in Vertex {
		float vbits;
	} IN[];

	out Vertex {
		vec4 col;
	} OUT;

	void createVertex(vec4 p) {
		gl_Position = cam.projection * cam.view * p;
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

		float delta = 0.0001;
		vec4 deltavec = vec4(delta, delta, delta, 0.0);
		vec4 dx = vec4(1.0+delta*2, 0.0, 0.0, 0.0);
		vec4 dy = vec4(0.0, 1.0+delta*2, 0.0, 0.0);
		vec4 dz = vec4(0.0, 0.0, 1.0+delta*2, 0.0);
		vec4 p1 = origin - deltavec;
		vec4 p2 = p1 + dx + dy + dz + deltavec;

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
