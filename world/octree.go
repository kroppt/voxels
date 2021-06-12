package world

type Octree struct {
	voxel    *Voxel
	aabc     *AABC
	children *octreeLinkedList
	parent   *Octree
}

type octreeLinkedList struct {
	node *Octree
	next *octreeLinkedList
}

// Find finds based on a predicate function, returning a list of
// all candidate voxels and a boolean indicating if any were found,
// in a depth-first search.
func (tree *Octree) Find(fn func(*Octree) bool) ([]*Voxel, bool) {
	if tree == nil {
		return nil, false
	}
	if !fn(tree) {
		return nil, false
	}
	if tree.children == nil {
		return []*Voxel{tree.voxel}, true
	}
	head := tree.children
	var voxels []*Voxel
	for head != nil {
		if vox, ok := head.node.Find(fn); ok {
			voxels = append(voxels, vox...)
		}
		head = head.next
	}
	return voxels, voxels != nil
}

// FindClosestIntersect finds the closest voxel that the camera's
// look direction intersects, if any.
func (tree *Octree) FindClosestIntersect(cam *Camera) (block *Voxel, dist float32, found bool) {
	if tree == nil {
		return nil, -1, false
	}
	boxDist, hit := Intersect(*tree.aabc, cam.GetPosition(), cam.GetLookForward())
	if !hit {
		return nil, -1, false
	}
	if tree.children == nil {
		return tree.voxel, boxDist, true
	}

	var bestDist float32
	var bestVox *Voxel
	head := tree.children
	for head != nil {
		vox, vdist, hit := head.node.FindClosestIntersect(cam)
		if hit && (vdist < bestDist || bestVox == nil) {
			bestVox = vox
			bestDist = vdist
		}
		head = head.next
	}
	return bestVox, bestDist, bestVox != nil
}

// Apply applies the function to every Octree node that is a leaf
// in a depth-first order.
func (tree *Octree) Apply(fn func(*Octree)) {
	if tree == nil {
		return
	}
	if tree.children == nil && tree.voxel == nil {
		panic("broken tree")
	}
	if tree.children == nil {
		fn(tree)
		return
	}
	head := tree.children
	for head != nil {
		head.node.Apply(fn)
		head = head.next
	}
}

// ApplyAll applies the function to all Octree nodes in a
// depth-first order.
func (tree *Octree) ApplyAll(fn func(*Octree)) {
	if tree == nil {
		return
	}
	fn(tree)
	head := tree.children
	for head != nil {
		head.node.ApplyAll(fn)
		head = head.next
	}
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

func (tree *Octree) GetVoxel() *Voxel {
	return tree.voxel
}

func (tree *Octree) GetAABC() *AABC {
	return tree.aabc
}

func (tree *Octree) GetChildren() *octreeLinkedList {
	return tree.children
}

// AddLeaf adds a leaf Octree node to an Octree an adds the
// corresponding bouncing boxes in between.
func (tree *Octree) AddLeaf(voxel *Voxel) *Octree {
	if voxel == nil {
		panic("voxel in AddLeaf is nil")
	}
	if tree == nil {
		aabc := &AABC{
			Origin: voxel.Pos,
			Size:   1,
		}
		octree := &Octree{
			voxel: voxel,
			aabc:  aabc,
		}
		return octree
	}
	curr := tree
	for !WithinAABC(curr.aabc, voxel.Pos) {
		// create bigger bounding box to include the new voxel
		aabc := ExpandAABC(curr.aabc, voxel.Pos)
		ll := &octreeLinkedList{
			node: curr,
		}
		newParent := &Octree{
			aabc:     aabc,
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
func (tree *Octree) addLeafRecurse(voxel *Voxel) {
	// reassignment case
	if tree.voxel != nil && tree.voxel.Pos == voxel.Pos {
		tree.voxel = voxel
		return
	}
	curr := tree.children
	for curr != nil {
		if WithinAABC(curr.node.aabc, voxel.Pos) {
			curr.node.addLeafRecurse(voxel)
			return
		}
		curr = curr.next
	}

	aabc := GetChildAABC(tree.aabc, voxel.Pos)
	newNode := &Octree{
		aabc: aabc,
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
func (tree *Octree) Remove(pos VoxelPos) (newRoot *Octree, removed bool) {
	if tree.voxel != nil && tree.voxel.Pos == pos { // correct leaf
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
