package shapes

import (
	"github.com/kroppt/voxels/voxgl"
)

// NewColoredCube returns a colored object of an object.
//
// col should contain the colors (0.0-1.0) for each of the vertices of a cube:
// far bottom left, far bottom right, far top left, far top right, close bottom
// left, close bottom right, close top left, close top right.
func NewColoredCube(x, y, z, r, g, b, a float32) (*voxgl.Object, error) {
	vertices := [7]float32{x, y, z, r, g, b, a}

	obj, err := voxgl.NewColoredObject(vertices)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
