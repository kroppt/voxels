package shapes

import (
	"github.com/kroppt/voxels/voxgl"
)

// NewColoredCube returns a colored object of an object.
//
// col should contain the colors (0.0-1.0) for each of the vertices of a cube:
// far bottom left, far bottom right, far top left, far top right, close bottom
// left, close bottom right, close top left, close top right.
func NewColoredCube(col [8][3]float32) (*voxgl.Object, error) {
	// close / far
	// top / bottom
	// left / right
	var fbl = [6]float32{-0.5, -0.5, 0.5, col[0][0], col[0][1], col[0][2]}
	var fbr = [6]float32{0.5, -0.5, 0.5, col[1][0], col[1][1], col[1][2]}
	var ftl = [6]float32{-0.5, 0.5, 0.5, col[2][0], col[2][1], col[2][2]}
	var ftr = [6]float32{0.5, 0.5, 0.5, col[3][0], col[3][1], col[3][2]}

	var cbl = [6]float32{-0.5, -0.5, -0.5, col[4][0], col[4][1], col[4][2]}
	var cbr = [6]float32{0.5, -0.5, -0.5, col[5][0], col[5][1], col[5][2]}
	var ctl = [6]float32{-0.5, 0.5, -0.5, col[6][0], col[6][1], col[6][2]}
	var ctr = [6]float32{0.5, 0.5, -0.5, col[7][0], col[7][1], col[7][2]}

	vertices := [][6]float32{
		// far face
		fbl, ftl, fbr,
		ftl, fbr, ftr,

		// left face
		ftl, ctl, fbl,
		ctl, fbl, cbl,

		// top face
		ftl, ctl, ftr,
		ctl, ftr, ctr,

		// right face
		fbr, ftr, cbr,
		ftr, cbr, ctr,

		// bottom face
		cbl, cbr, fbl,
		cbr, fbl, fbr,

		// close face
		ctl, ctr, cbr,
		cbl, ctl, cbr,
	}

	obj, err := voxgl.NewColoredObject(vertices)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
