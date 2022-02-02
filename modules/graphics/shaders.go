package graphics

const vertCrossShader = `
#version 420 core

layout (location = 0) in vec2 pos;
layout (location = 1) in vec4 col;

out Vertex {
	vec4 col;
} OUT;

void main()
{
	gl_Position = vec4(pos.xy, 0.0, 1.0);
	OUT.col = col;
}
`

const fragCrossShader = `
#version 400

in Vertex {
	vec4 col;
} IN;

out vec4 frag_color;

void main() {
	frag_color = IN.col;
}
`

const vertColShader = `
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

const geoColShader = `
	#version 420 core

	layout(points) in;
	layout(triangle_strip, max_vertices = 24) out;

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
		vec3 stdir;
		flat int blockType;
		flat uint faceLight;
	} OUT;

	void createVertex(vec4 p) {
		vec3 center = (gl_in[0].gl_Position).xyz + 0.5;
		OUT.stdir = p.xyz - center;
		gl_Position = vec4(cam.projection * cam.view * p);
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
		if (OUT.blockType == 0) {
			return; // render nothing if air block
		}

		uint nMaskBits = 4;
		uint lightFrontMask   = 15;         // The voxel's front face lighting bits.
		uint lightBackMask    = lightFrontMask << nMaskBits;  // The voxel's back face lighting bits.
		uint lightBottomMask  = lightFrontMask << (nMaskBits*2); // The voxel's bottom face lighting bits.
		uint lightTopMask     = lightFrontMask << (nMaskBits*3); // The voxel's top face lighting bits.
		uint lightLeftMask    = lightFrontMask << (nMaskBits*4); // The voxel's left face lighting bits.
		uint lightRightMask   = lightFrontMask << (nMaskBits*5); // The voxel's right face lighting bits.
		uint lbits = uint(IN[0].lighting);

		
		vec4 dx = vec4(1.0, 0.0, 0.0, 0.0);
		vec4 dy = vec4(0.0, 1.0, 0.0, 0.0);
		vec4 dz = vec4(0.0, 0.0, 1.0, 0.0);
		vec4 p1 = origin;
		vec4 p2 = p1 + dx + dy + dz;

		if ((bits & backwardmask) - backwardmask != 0) {
			OUT.faceLight = (lbits & lightBackMask) >> nMaskBits;
			createQuad(p2, -dx, -dy); // backward
		}
		if ((bits & forwardmask) - forwardmask != 0) {
			OUT.faceLight = lbits & lightFrontMask;
			createQuad(p1, dy, dx); // forward
		}
		if ((bits & topmask) - topmask != 0) {
			OUT.faceLight = (lbits & lightTopMask) >> (nMaskBits*3);
			createQuad(p2, -dz, -dx); // top
		}
		if ((bits & bottommask) - bottommask != 0) {
			OUT.faceLight = (lbits & lightBottomMask) >> (nMaskBits*2);
			createQuad(p1, dx, dz); // bottom
		}
		if ((bits & rightmask) - rightmask != 0) {
			OUT.faceLight = (lbits & lightRightMask) >> (nMaskBits*5);
			createQuad(p2, -dy, -dz); // right
		}
		if ((bits & leftmask) - leftmask != 0) {
			OUT.faceLight = (lbits & lightLeftMask) >> (nMaskBits*4);
			createQuad(p1, dz, dy); // left
		}
	}
`

const fragColShader = `
	#version 400

	in Vertex {
		vec3 stdir;
		flat int blockType;
		flat uint faceLight;
	} IN;
	uniform samplerCubeArray cubeMapArray;


	out vec4 frag_color;

	void main() {
		uint maxFaceLight = 8;
		uint correctedFaceLight = IN.faceLight;
		if (correctedFaceLight == 0) {
			correctedFaceLight = 1;
		}
		float lightFrac = 1;//float(correctedFaceLight) / float(maxFaceLight);
		vec4 fullBright = texture(cubeMapArray, vec4(IN.stdir, IN.blockType));
		frag_color = vec4(fullBright.xyz * lightFrac, fullBright.w);
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
