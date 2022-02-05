package view

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl64"
)

// AABC represents an Axis Aligned Bounding Cube, where the
// position represents the corner minimum coordinates.
type AABC struct {
	Origin mgl.Vec3
	Size   int
}

// ExpandAABC doubles the dimensions of the aabc, moving the center
// in the corresponding coordinates if the target is less than any of them.
// It is assumed that the target is not within the aabc, and will otherwise panic.
func ExpandAABC(aabc AABC, target mgl.Vec3) AABC {
	if WithinAABC(aabc, target) {
		panic("unintended use: target is within aabc")
	}
	size := aabc.Size
	newAabb := AABC{
		Origin: aabc.Origin,
		Size:   size * 2,
	}
	if target.X() < aabc.Origin.X() {
		sub := mgl.Vec3{float64(size), 0, 0}
		newAabb.Origin = newAabb.Origin.Sub(sub)
	}
	if target.Y() < aabc.Origin.Y() {
		sub := mgl.Vec3{0, float64(size), 0}
		newAabb.Origin = newAabb.Origin.Sub(sub)
	}
	if target.Z() < aabc.Origin.Z() {
		sub := mgl.Vec3{0, 0, float64(size)}
		newAabb.Origin = newAabb.Origin.Sub(sub)
	}
	return newAabb
}

// GetChildAABC returns the smaller aabc from the appropriate octant from within
// the larger aabc. The smaller aabc is gauranteed to contain the target.
// It is assumed that the target is within the aabc, and will otherwise panic.
func GetChildAABC(aabc AABC, target mgl.Vec3) AABC {
	if !WithinAABC(aabc, target) {
		panic("unintended use: target is not within aabc (probably adding duplicate voxel)")
	}
	size := aabc.Size / 2.0
	offset := mgl.Vec3{0, 0, 0}
	if target.X() >= aabc.Origin.X()+float64(size) {
		add := mgl.Vec3{float64(size), 0, 0}
		offset = offset.Add(add)
	}
	if target.Y() >= aabc.Origin.Y()+float64(size) {
		add := mgl.Vec3{0, float64(size), 0}
		offset = offset.Add(add)
	}
	if target.Z() >= aabc.Origin.Z()+float64(size) {
		add := mgl.Vec3{0, 0, float64(size)}
		offset = offset.Add(add)
	}
	newAabb := AABC{
		Origin: aabc.Origin.Add(offset),
		Size:   size,
	}
	return newAabb
}

// WithinAABC returns whether the target point is within the bounds of the aabc,
// where minimum is inclusive and maximum is exclusive.
func WithinAABC(aabc AABC, target mgl.Vec3) bool {
	// the vertex associated with the bounding box is the bounding box's minimum coordinate vertex
	minx := aabc.Origin.X()
	maxx := aabc.Origin.X() + float64(aabc.Size)
	if target.X() >= maxx || target.X() < minx {
		return false
	}

	miny := aabc.Origin.Y()
	maxy := aabc.Origin.Y() + float64(aabc.Size)
	if target.Y() >= maxy || target.Y() < miny {
		return false
	}

	minz := aabc.Origin.Z()
	maxz := aabc.Origin.Z() + float64(aabc.Size)
	if target.Z() >= maxz || target.Z() < minz {
		return false
	}
	return true
}

// Intersect returns whether the given ray intersects the given box and the
// distance if it does.
func Intersect(box AABC, pos, dir mgl.Vec3) (dist float64, hit bool) {
	boxPos := box.Origin
	boxSize := float64(box.Size)
	boxmin := func(d int) float64 {
		return boxPos[d]
	}
	boxmax := func(d int) float64 {
		return boxPos[d] + boxSize
	}

	invx := 1.0 / dir[0]
	tx1 := (boxmin(0) - pos[0]) * invx
	tx2 := (boxmax(0) - pos[0]) * invx
	txmin := math.Min(tx1, tx2)
	txmax := math.Max(tx1, tx2)
	min := txmin
	max := txmax

	invy := 1.0 / dir[1]
	ty1 := (boxmin(1) - pos[1]) * invy
	ty2 := (boxmax(1) - pos[1]) * invy
	tymin := math.Min(ty1, ty2)
	tymax := math.Max(ty1, ty2)
	min = math.Max(min, tymin)
	max = math.Min(max, tymax)

	invz := 1.0 / dir[2]
	tz1 := (boxmin(2) - pos[2]) * invz
	tz2 := (boxmax(2) - pos[2]) * invz
	tzmin := math.Min(tz1, tz2)
	tzmax := math.Max(tz1, tz2)
	min = math.Max(min, tzmin)
	max = math.Min(max, tzmax)

	hit = (max >= min) && max >= 0.0
	dist = min

	return
}
