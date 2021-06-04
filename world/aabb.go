package world

import (
	"github.com/engoengine/glm"
	"github.com/engoengine/glm/geo"
)

// ExpandAABB doubles the dimensions of the aabb, moving the center
// in the corresponding coordinates if the target is less than any of them.
// It is assumed that the target is not within the aabb, and will otherwise panic.
func ExpandAABB(aabb *geo.AABB, target glm.Vec3) *geo.AABB {
	if WithinAABB(aabb, target) {
		panic("unintended use: target is within aabb")
	}
	size := aabb.HalfExtend.X() * 2
	newAabb := &geo.AABB{
		Center:     glm.Vec3{aabb.Center.X(), aabb.Center.Y(), aabb.Center.Z()},
		HalfExtend: glm.Vec3{size, size, size},
	}
	if target.X() < aabb.Center.X() {
		sub := &glm.Vec3{size, 0.0, 0.0}
		newAabb.Center = newAabb.Center.Sub(sub)
	}
	if target.Y() < aabb.Center.Y() {
		sub := &glm.Vec3{0.0, size, 0.0}
		newAabb.Center = newAabb.Center.Sub(sub)
	}
	if target.Z() < aabb.Center.Z() {
		sub := &glm.Vec3{0.0, 0.0, size}
		newAabb.Center = newAabb.Center.Sub(sub)
	}
	return newAabb
}

// GetChildAABB returns the smaller aabb from the appropriate octant from within
// the larger aabb. The smaller aabb is gauranteed to contain the target.
// It is assumed that the target is within the aabb, and will otherwise panic.
func GetChildAABB(aabb *geo.AABB, target glm.Vec3) *geo.AABB {
	if !WithinAABB(aabb, target) {
		panic("unintended use: target is not within aabb")
	}
	size := aabb.HalfExtend.X()
	offset := glm.Vec3{0.0, 0.0, 0.0}
	if target.X() >= aabb.Center.X()+aabb.HalfExtend.X() {
		add := &glm.Vec3{size, 0.0, 0.0}
		offset = offset.Add(add)
	}
	if target.Y() >= aabb.Center.Y()+aabb.HalfExtend.Y() {
		add := &glm.Vec3{0.0, size, 0.0}
		offset = offset.Add(add)
	}
	if target.Z() >= aabb.Center.Z()+aabb.HalfExtend.Z() {
		add := &glm.Vec3{0.0, 0.0, size}
		offset = offset.Add(add)
	}
	newSize := size / 2
	newAabb := &geo.AABB{
		Center:     aabb.Center.Add(&offset),
		HalfExtend: glm.Vec3{newSize, newSize, newSize},
	}
	return newAabb
}

// WithinAABB returns whether the target point is within the bounds of the aabb,
// where minimum is inclusive and maximum is exclusive.
func WithinAABB(aabb *geo.AABB, target glm.Vec3) bool {
	// the vertex associated with the bounding box is the bounding box's minimum coordinate vertex
	minx := aabb.Center.X()
	maxx := aabb.Center.X() + aabb.HalfExtend.X()*2
	if target.X() >= maxx || target.X() < minx {
		return false
	}

	miny := aabb.Center.Y()
	maxy := aabb.Center.Y() + aabb.HalfExtend.Y()*2
	if target.Y() >= maxy || target.Y() < miny {
		return false
	}

	minz := aabb.Center.Z()
	maxz := aabb.Center.Z() + aabb.HalfExtend.Z()*2
	if target.Z() >= maxz || target.Z() < minz {
		return false
	}
	return true
}
