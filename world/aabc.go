package world

// AABC represents an Axis Aligned Bounding Cube, where the
// position represents the corner minimum coordinates.
type AABC struct {
	Pos  VoxelPos
	Size int32
}

// ExpandAABC doubles the dimensions of the aabc, moving the center
// in the corresponding coordinates if the target is less than any of them.
// It is assumed that the target is not within the aabc, and will otherwise panic.
func ExpandAABC(aabc *AABC, target VoxelPos) *AABC {
	if WithinAABC(aabc, target) {
		panic("unintended use: target is within aabc")
	}
	size := aabc.Size
	newAabb := &AABC{
		Pos:  aabc.Pos,
		Size: size * 2,
	}
	if target.X < aabc.Pos.X {
		sub := VoxelPos{size, 0, 0}
		newAabb.Pos = newAabb.Pos.Sub(sub)
	}
	if target.Y < aabc.Pos.Y {
		sub := VoxelPos{0, size, 0}
		newAabb.Pos = newAabb.Pos.Sub(sub)
	}
	if target.Z < aabc.Pos.Z {
		sub := VoxelPos{0, 0, size}
		newAabb.Pos = newAabb.Pos.Sub(sub)
	}
	return newAabb
}

// GetChildAABC returns the smaller aabc from the appropriate octant from within
// the larger aabc. The smaller aabc is gauranteed to contain the target.
// It is assumed that the target is within the aabc, and will otherwise panic.
func GetChildAABC(aabc *AABC, target VoxelPos) *AABC {
	if !WithinAABC(aabc, target) {
		panic("unintended use: target is not within aabc (probably adding duplicate voxel)")
	}
	size := aabc.Size / 2.0
	offset := VoxelPos{0, 0, 0}
	if target.X >= aabc.Pos.X+size {
		add := VoxelPos{size, 0, 0}
		offset = offset.Add(add)
	}
	if target.Y >= aabc.Pos.Y+size {
		add := VoxelPos{0, size, 0}
		offset = offset.Add(add)
	}
	if target.Z >= aabc.Pos.Z+size {
		add := VoxelPos{0, 0, size}
		offset = offset.Add(add)
	}
	newAabb := &AABC{
		Pos:  aabc.Pos.Add(offset),
		Size: size,
	}
	return newAabb
}

// WithinAABC returns whether the target point is within the bounds of the aabc,
// where minimum is inclusive and maximum is exclusive.
func WithinAABC(aabc *AABC, target VoxelPos) bool {
	// the vertex associated with the bounding box is the bounding box's minimum coordinate vertex
	minx := aabc.Pos.X
	maxx := aabc.Pos.X + aabc.Size
	if target.X >= maxx || target.X < minx {
		return false
	}

	miny := aabc.Pos.Y
	maxy := aabc.Pos.Y + aabc.Size
	if target.Y >= maxy || target.Y < miny {
		return false
	}

	minz := aabc.Pos.Z
	maxz := aabc.Pos.Z + aabc.Size
	if target.Z >= maxz || target.Z < minz {
		return false
	}
	return true
}
