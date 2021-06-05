package chunk

type Octree struct {
	voxel    *Voxel
	aabc     *AABC
	children *octreeLinkedList
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

// Apply applies the function to every Octree node that is a leaf
// in a depth-first order.
func (tree *Octree) Apply(fn func(*Octree)) {
	if tree == nil {
		return
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
			Pos:  voxel.Pos,
			Size: 1,
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
		curr = &Octree{
			aabc:     aabc,
			children: ll,
		}
	}
	// voxel must be in the current bounding box
	curr.addLeafRecurse(voxel)
	return curr
}

// voxel is inside AABC of tree and the tree has at least one child
func (tree *Octree) addLeafRecurse(voxel *Voxel) {
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
	newChild := &octreeLinkedList{
		node: newNode,
	}

	head := tree.children
	if head == nil {
		tree.children = newChild
	} else {
		newChild.next = head
		tree.children = newChild
	}

	if aabc.Size == 1.0 {
		// node is a leaf
		newNode.voxel = voxel
	} else {
		newNode.addLeafRecurse(voxel)
	}
}
