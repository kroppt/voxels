package oldworld

// AABC represents an Axis Aligned Bounding Cube, where the
// position represents the corner minimum coordinates.
type AABC struct {
	Origin VoxelPos
	Size   int
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
		Origin: aabc.Origin,
		Size:   size * 2,
	}
	if target.X < aabc.Origin.X {
		sub := VoxelPos{size, 0, 0}
		newAabb.Origin = newAabb.Origin.Sub(sub)
	}
	if target.Y < aabc.Origin.Y {
		sub := VoxelPos{0, size, 0}
		newAabb.Origin = newAabb.Origin.Sub(sub)
	}
	if target.Z < aabc.Origin.Z {
		sub := VoxelPos{0, 0, size}
		newAabb.Origin = newAabb.Origin.Sub(sub)
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
	if target.X >= aabc.Origin.X+size {
		add := VoxelPos{size, 0, 0}
		offset = offset.Add(add)
	}
	if target.Y >= aabc.Origin.Y+size {
		add := VoxelPos{0, size, 0}
		offset = offset.Add(add)
	}
	if target.Z >= aabc.Origin.Z+size {
		add := VoxelPos{0, 0, size}
		offset = offset.Add(add)
	}
	newAabb := &AABC{
		Origin: aabc.Origin.Add(offset),
		Size:   size,
	}
	return newAabb
}

// WithinAABC returns whether the target point is within the bounds of the aabc,
// where minimum is inclusive and maximum is exclusive.
func WithinAABC(aabc *AABC, target VoxelPos) bool {
	// the vertex associated with the bounding box is the bounding box's minimum coordinate vertex
	minx := aabc.Origin.X
	maxx := aabc.Origin.X + aabc.Size
	if target.X >= maxx || target.X < minx {
		return false
	}

	miny := aabc.Origin.Y
	maxy := aabc.Origin.Y + aabc.Size
	if target.Y >= maxy || target.Y < miny {
		return false
	}

	minz := aabc.Origin.Z
	maxz := aabc.Origin.Z + aabc.Size
	if target.Z >= maxz || target.Z < minz {
		return false
	}
	return true
}
