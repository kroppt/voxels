package world

import (
	"github.com/engoengine/glm"
	"github.com/engoengine/glm/geo"
)

type Octree struct {
	voxel    *Voxel
	aabb     *geo.AABB
	children *octreeLinkedList
}

type octreeLinkedList struct {
	node *Octree
	next *octreeLinkedList
}

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

func (tree *Octree) GetVoxel() *Voxel {
	return tree.voxel
}

func (tree *Octree) GetAABB() *geo.AABB {
	return tree.aabb
}

func (tree *Octree) GetChildren() *octreeLinkedList {
	return tree.children
}

func (tree *Octree) AddLeaf(voxel *Voxel) *Octree {
	if tree == nil {
		aabb := &geo.AABB{
			Center:     voxel.Pos,
			HalfExtend: glm.Vec3{0.5, 0.5, 0.5},
		}
		octree := &Octree{
			voxel: voxel,
			aabb:  aabb,
		}
		return octree
	}
	curr := tree
	for !WithinAABB(curr.aabb, voxel.Pos) {
		// create bigger bounding box to include the new voxel
		aabb := ExpandAABB(curr.aabb, voxel.Pos)
		ll := &octreeLinkedList{
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
		if WithinAABB(curr.node.aabb, voxel.Pos) {
			curr.node.addLeafRecurse(voxel)
			return
		}
		curr = curr.next
	}

	aabb := GetChildAABB(tree.aabb, voxel.Pos)
	newNode := &Octree{
		aabb: aabb,
	}
	newChild := &octreeLinkedList{
		node: newNode,
	}

	head := tree.children
	if head == nil {
		tree.children = head
	} else {
		newChild.next = head
		tree.children = newChild
	}

	if aabb.HalfExtend.X() == 0.5 {
		// node is a leaf
		newNode.voxel = voxel
	} else {
		newNode.addLeafRecurse(voxel)
	}
}
