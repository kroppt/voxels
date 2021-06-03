package world

import (
	"github.com/engoengine/glm"
	"github.com/engoengine/glm/geo"
)

type Octree struct {
	voxel    *Voxel
	aabb     *geo.AABB
	children *OctreeLinkedList
}

type OctreeLinkedList struct {
	node *Octree
	next *OctreeLinkedList
}

func (tree *Octree) addLeaf(voxel *Voxel) *Octree {
	if tree == nil {
		aabb := &geo.AABB{
			Center:     voxel.coordinates,
			HalfExtend: glm.Vec3{0.5, 0.5, 0.5},
		}
		octree := &Octree{
			voxel: voxel,
			aabb:  aabb,
		}
		return octree
	}
	curr := tree
	for !withinAABB(tree.aabb, voxel.coordinates) {
		// create bigger bounding box to include the new voxel
		aabb := getNextBiggestAABB(tree.aabb, voxel.coordinates)
		ll := &OctreeLinkedList{
			node: curr,
		}
		curr = &Octree{
			aabb:     aabb,
			children: ll,
		}
	}
	// voxel must be in the current bounding box
	curr.addLeafRecurse(voxel)
	return curr
}

// voxel is inside AABB of tree and the tree has at least one child
func (tree *Octree) addLeafRecurse(voxel *Voxel) {
	curr := tree.children
	for curr != nil {
		if withinAABB(curr.node.aabb, voxel.coordinates) {
			curr.node.addLeafRecurse(voxel)
			return
		}
		curr = curr.next
	}
	aabb := getChildAABB(tree.aabb, voxel.coordinates)
	next := tree.children
	var prev *OctreeLinkedList
	for next != nil {
		prev = next
		next = next.next
	}
	if prev == nil {
		panic("inconsistent state: tree has no children")
	}
	newNode := &Octree{
		aabb: aabb,
	}
	newChild := &OctreeLinkedList{
		node: newNode,
	}
	prev.next = newChild

	if aabb.HalfExtend.X() == 0.5 {
		// node is a leaf
		newNode.voxel = voxel
	} else {
		newNode.addLeafRecurse(voxel)
	}
}

func getNextBiggestAABB(aabb *geo.AABB, coords glm.Vec3) *geo.AABB {
	size := aabb.HalfExtend.X() * 2
	newAabb := &geo.AABB{
		Center:     glm.Vec3{aabb.Center.X(), aabb.Center.Y(), aabb.Center.Z()},
		HalfExtend: glm.Vec3{size, size, size},
	}
	if coords.X() < aabb.Center.X() {
		sub := &glm.Vec3{size, 0.0, 0.0}
		newAabb.Center = newAabb.Center.Sub(sub)
	}
	if coords.Y() < aabb.Center.Y() {
		sub := &glm.Vec3{0.0, size, 0.0}
		newAabb.Center = newAabb.Center.Sub(sub)
	}
	if coords.Z() < aabb.Center.Z() {
		sub := &glm.Vec3{0.0, 0.0, size}
		newAabb.Center = newAabb.Center.Sub(sub)
	}
	return newAabb
}

func getChildAABB(aabb *geo.AABB, coords glm.Vec3) *geo.AABB {
	size := aabb.HalfExtend.X()
	offset := glm.Vec3{0.0, 0.0, 0.0}
	if coords.X() >= aabb.Center.X()+aabb.HalfExtend.X() {
		add := &glm.Vec3{size, 0.0, 0.0}
		offset = offset.Add(add)
	}
	if coords.Y() >= aabb.Center.Y()+aabb.HalfExtend.Y() {
		add := &glm.Vec3{0.0, size, 0.0}
		offset = offset.Add(add)
	}
	if coords.Z() >= aabb.Center.Z()+aabb.HalfExtend.Z() {
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

func withinAABB(aabb *geo.AABB, pos glm.Vec3) bool {
	// the vertex associated with the bounding box is the bounding box's minimum coordinate vertex
	minx := aabb.Center.X()
	maxx := aabb.Center.X() + aabb.HalfExtend.X()*2
	if pos.X() >= maxx || pos.X() < minx {
		return false
	}

	miny := aabb.Center.Y()
	maxy := aabb.Center.Y() + aabb.HalfExtend.Y()*2
	if pos.Y() >= maxy || pos.Y() < miny {
		return false
	}

	minz := aabb.Center.Z()
	maxz := aabb.Center.Z() + aabb.HalfExtend.Z()*2
	if pos.Z() >= maxz || pos.Z() < minz {
		return false
	}
	return true
}
