// +build !test

package voxgl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kroppt/gfx"
	"github.com/kroppt/voxels/util"
)

// NewColoredObject returns a newly created Object with the given colors.
//
// Vertices should be vertices of format X, Y, Z, R, G, B, A.
// X, Y, and Z options should be in the range -1.0 to 1.0.
// R, G, B, and A should be in the range 0.0 to 1.0.
func NewColoredObject(vertices []float32) (*Object, error) {
	sw := util.Start()
	// prog, err := GetProgram(vertColShader, fragColShader, geoColShader)
	// if err != nil {
	// 	return nil, err
	// }
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
	sw.StopRecordAverage("program")

	obj, err := NewObject(prog, vertices, []int32{3, 4})
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
	layout(triangle_strip, max_vertices = 24) out;

	layout (std140, binding = 0) uniform Matrices
	{
		mat4 view;
		mat4 projection;
	} cam;
	uniform sampler3D adjaTex;

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

	void createQuad(vec4 p, vec4 d1, vec4 d2) {
		createVertex(p);
		createVertex(p+d1);
		createVertex(p+d2);
		createVertex(p+d1+d2);
		EndPrimitive();
	}

	ivec2 getChunkPos(vec4 p, int chunkSize) {
		int x = int(p.x);
		int z = int(p.z);
		if (p.x < 0) {
			x = x + 1;
		}
		if (p.z < 0) {
			z = z + 1;
		}
		x = x / chunkSize;
		z = z / chunkSize;
		if (p.x < 0) {
			x = x - 1;
		}
		if (p.z < 0) {
			z = z - 1;
		}
		return ivec2(x, z);
	}

	ivec3 getLocalChunkPos(vec4 p, int chunkSize) {
		ivec2 chunkPos = getChunkPos(p, chunkSize);
		return ivec3(int(p.x)-(chunkPos.x*chunkSize), int(p.y), int(p.z)-(chunkPos.y*chunkSize));
	}
	
	void main() {
		vec4 center = gl_in[0].gl_Position;
		int chunkSize = textureSize(adjaTex, 0).x;
		int height = textureSize(adjaTex, 0).y;
		
		vec3 indices = getLocalChunkPos(center, chunkSize);
		int i = int(indices.x);
		int j = int(indices.y);
		int k = int(indices.z);

		vec4 dx = vec4(1.0, 0.0, 0.0, 0.0);
		vec4 dy = vec4(0.0, 1.0, 0.0, 0.0);
		vec4 dz = vec4(0.0, 0.0, 1.0, 0.0);
		vec4 p1 = center;
		vec4 p2 = p1 + dx + dy + dz;

		if (k == chunkSize - 1 || texelFetch(adjaTex, ivec3(i, j, k+1), 0).r == 0) {
			createQuad(p2, -dx, -dy);
		}
		if (k == 0 || texelFetch(adjaTex, ivec3(i, j, k-1), 0).r == 0) {
			createQuad(p1, dy, dx);
		}
		if (j == height - 1 || texelFetch(adjaTex, ivec3(i, j+1, k), 0).r == 0) {
			createQuad(p2, -dz, -dx);
		}
		if (j == 0 || texelFetch(adjaTex, ivec3(i, j-1, k), 0).r == 0) {
			createQuad(p1, dx, dz);
		}
		if (i == chunkSize - 1 || texelFetch(adjaTex, ivec3(i+1, j, k), 0).r == 0) {
			createQuad(p2, -dy, -dz);
		}
		if (i == 0 || texelFetch(adjaTex, ivec3(i-1, j, k), 0).r == 0) {
			createQuad(p1, dz, dy);
		}
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
