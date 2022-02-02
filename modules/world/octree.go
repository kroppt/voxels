package world

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/kroppt/voxels/chunk"
)

type Octree struct {
	voxel    *chunk.VoxelCoordinate
	aabc     *AABC
	children *octreeLinkedList
	parent   *Octree
}

type octreeLinkedList struct {
	node *Octree
	next *octreeLinkedList
}

// FindClosestIntersect finds the closest voxel that the camera's
// look direction intersects, if any.
func (tree *Octree) FindClosestIntersect(eye, dirNorm mgl.Vec3) (vc chunk.VoxelCoordinate, dist float64, found bool) {
	if tree == nil {
		return chunk.VoxelCoordinate{}, -1, false
	}
	boxDist, hit := Intersect(*tree.aabc, eye, dirNorm)
	if !hit {
		return chunk.VoxelCoordinate{}, -1, false
	}
	if tree.children == nil {
		return *tree.voxel, boxDist, true
	}

	var bestDist float64
	var bestVox chunk.VoxelCoordinate
	bestVox = chunk.VoxelCoordinate{}
	assigned := false
	head := tree.children
	for head != nil {
		vox, vdist, hit := head.node.FindClosestIntersect(eye, dirNorm)
		if hit && (vdist < bestDist || !assigned) {
			assigned = true
			bestVox = vox
			bestDist = vdist
		}
		head = head.next
	}

	if bestDist < 0 {
		bestDist = 0
	}
	return bestVox, bestDist, assigned
}

// CountChildren returns how many children an Octree node has
func (tree *Octree) CountChildren() int {
	ll := tree.children
	if ll == nil {
		return 0
	}
	curr := ll
	len := 0
	for curr != nil {
		len += 1
		curr = curr.next
	}
	return len
}

func (tree *Octree) GetAABC() *AABC {
	return tree.aabc
}

func (tree *Octree) GetVoxel() *chunk.VoxelCoordinate {
	return tree.voxel
}

func (tree *Octree) String() string {
	str := ""
	str += fmt.Sprint("{")
	if tree != nil {
		if tree.voxel == nil {
			str += fmt.Sprintf("(%v)->", tree.aabc.Size)
		} else {
			v := *tree.voxel
			str += fmt.Sprintf("(%v,%v,%v)->", v.X, v.Y, v.Z)
		}
		head := tree.children
		if head == nil {
			str += "{}"
		}
		for head != nil {
			str += head.node.String()
			head = head.next
		}
	}
	return str + fmt.Sprintf("}")
}

// AddLeaf adds a leaf Octree node to an Octree an adds the
// corresponding bouncing boxes in between.
func (tree *Octree) AddLeaf(voxel *chunk.VoxelCoordinate) *Octree {
	if voxel == nil {
		panic("voxel in AddLeaf is nil")
	}
	if tree == nil {
		aabc := &AABC{
			Origin: mgl.Vec3{float64(voxel.X), float64(voxel.Y), float64(voxel.Z)},
			Size:   1,
		}
		octree := &Octree{
			voxel: voxel,
			aabc:  aabc,
		}
		return octree
	}
	curr := tree
	for !WithinAABC(*curr.aabc, mgl.Vec3{float64(voxel.X), float64(voxel.Y), float64(voxel.Z)}) {
		// create bigger bounding box to include the new voxel
		aabc := ExpandAABC(*curr.aabc, mgl.Vec3{float64(voxel.X), float64(voxel.Y), float64(voxel.Z)})
		ll := &octreeLinkedList{
			node: curr,
		}
		newParent := &Octree{
			aabc:     &aabc,
			children: ll,
		}
		curr.parent = newParent
		curr = newParent
	}
	// voxel must be in the current bounding box
	curr.addLeafRecurse(voxel)
	return curr
}

func (tree *Octree) addChild(child *Octree) {
	child.parent = tree
	listNode := &octreeLinkedList{
		node: child,
	}
	head := tree.children
	if head == nil {
		tree.children = listNode
	} else {
		listNode.next = head
		tree.children = listNode
	}
}

func (tree *Octree) removeChild(node *Octree) {
	head := tree.children
	var prev *octreeLinkedList
	for head != nil {
		if head.node == node {
			if prev == nil {
				// only child is head
				tree.children = head.next
			} else {
				// had siblings
				prev.next = head.next
			}
		}
		prev = head
		head = head.next
	}
}

// voxel is inside AABC of tree and the tree has at least one child
func (tree *Octree) addLeafRecurse(voxel *chunk.VoxelCoordinate) {
	// reassignment case
	if tree.voxel != nil && *tree.voxel == *voxel {
		tree.voxel = voxel
		return
	}
	curr := tree.children
	for curr != nil {
		if WithinAABC(*curr.node.aabc, mgl.Vec3{float64(voxel.X), float64(voxel.Y), float64(voxel.Z)}) {
			curr.node.addLeafRecurse(voxel)
			return
		}
		curr = curr.next
	}

	aabc := GetChildAABC(*tree.aabc, mgl.Vec3{float64(voxel.X), float64(voxel.Y), float64(voxel.Z)})
	newNode := &Octree{
		aabc: &aabc,
	}
	tree.addChild(newNode)

	if aabc.Size == 1.0 {
		// node is a leaf
		newNode.voxel = voxel
	} else {
		newNode.addLeafRecurse(voxel)
	}
}

// Remove removes a specified voxel from the tree, and returns whether
// the root was removed (entire tree cleared)
func (tree *Octree) Remove(pos chunk.VoxelCoordinate) (newRoot *Octree, removed bool) {
	if tree.voxel != nil && *tree.voxel == pos { // correct leaf
		return tree.prune(), true
	} else if tree.voxel == nil { // not a leaf
		head := tree.children
		for head != nil {
			root, removed := head.node.Remove(pos)
			if removed {
				return root, removed
			}
			head = head.next
		}
	}
	// wrong leaf, do nothing
	return nil, false
}

func (tree *Octree) getRoot() *Octree {
	if tree == nil {
		panic("getRoot of nil tree")
	} else if tree.parent == nil {
		return tree.shrinkRoot()
	} else {
		return tree.parent.getRoot()
	}
}

func (tree *Octree) shrinkRoot() *Octree {
	if tree.CountChildren() == 1 {
		return tree.children.node.shrinkRoot()
	} else {
		return tree
	}
}

func (tree *Octree) prune() *Octree {
	if tree.parent == nil {
		// this is the root node
		return nil
	}
	// remove self from parent's children
	tree.parent.removeChild(tree)
	numChildren := tree.parent.CountChildren()
	if numChildren == 0 {
		return tree.parent.prune()
	}
	return tree.getRoot()
}
