package world

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

type Voxel struct {
	Position Position
	Color    Color
}
