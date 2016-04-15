package csbt

import (
	"container/list"
)

type _node_id struct {
	node   uint //internal_node_ptr of parent node
	level  int
	offset int
	path   *list.List
}

func (t *DCSBT) NodeAddLeafChild(n uint, index int, e uint) {
	t.NodeSetLeafChild(n, index, e)
}

func (t *DCSBT) NodeAddNodeChild(n uint, index int, e uint) {
	t.mb.Copy(t.mb.GetChild(n, index), e, 64)
}

// If the leafnode is inserted into the internal node
// we know that no more data will be allocated below thos level
func (t *DCSBT) NodeSetLeafChild(n uint, index int, e uint) {
	num_keys := t.mb.GetNodeKeyNum(n)
	t.mb.Copy(t.mb.GetChild(n, index), e, 64)
	var l, r uint //Leaf
	if index == 0 {
		l = t.mb.GetChild(n, 0)
		r = t.mb.GetChild(n, 1)
		t.mb.SetLeafRight(l, r)
		t.mb.SetLeafLeft(r, l)
	} else if index > 0 && index < num_keys {
		l = t.mb.GetChild(n, index-1)
		r = t.mb.GetChild(n, index)
		t.mb.SetLeafRight(l, r)
		t.mb.SetLeafLeft(r, l)
		l = t.mb.GetChild(n, index)
		r = t.mb.GetChild(n, index+1)
		t.mb.SetLeafRight(l, r)
		t.mb.SetLeafLeft(r, l)
	} else if index == num_keys {
		l = t.mb.GetChild(n, index-1)
		r = t.mb.GetChild(n, index)
		t.mb.SetLeafRight(l, r)
		t.mb.SetLeafLeft(r, l)
	}
}

func (t *DCSBT) NodePushLeafBack(n uint, e uint) {
	num_keys := t.mb.GetNodeKeyNum(n)
	t.NodeAddLeafChild(n, num_keys, e)
}

// When an internal node is set, we have to copy the
// children of the parameter node as well
func (t *DCSBT) NodeSetNodeChild(n uint, index int, e uint) {
	// Copy the value from the node but than explicitely copy
	// the structure
	t.mb.Copy(t.mb.GetChild(n, index), e, 64)
	t.mb.SetChildren(t.mb.GetChild(n, index), t.mb.GetChild(e, 0))

	t.mb.SetChildren(e, 0)
}

func (t *DCSBT) NodePushNodeBack(n uint, e uint) {
	num_keys := t.mb.GetNodeKeyNum(n)
	// Copys the memory for the node
	t.NodeSetNodeChild(n, num_keys, e)

	if t.mb.IsLeaf(t.mb.GetChild(t.mb.GetChild(n, num_keys), 0)) {
		var l, r uint //Leaf
		tmp := t.mb.GetChild(n, num_keys-1)
		l = t.mb.GetChild(tmp, t.mb.GetNodeKeyNum(tmp))
		tmp = t.mb.GetChild(n, num_keys)
		r = t.mb.GetChild(tmp, 0)
		t.mb.SetLeafRight(l, r)
		t.mb.SetLeafLeft(r, l)
	}
}

func (t *DCSBT) NodeInsertNodeForKey(n uint, k uint32, e uint) {

	// 1. Find the insertion point for the key
	// 2. move all elements one to the right
	// 3. insert the child
	tmp := t.mb.GetChild(n, 0) //Leaf

	insertion_pointer := int(0)
	num_keys := t.mb.GetNodeKeyNum(n)
	for t.mb.GetNodeKey(n, insertion_pointer) < k && insertion_pointer < num_keys {
		insertion_pointer++
	}
	for i := num_keys; i > insertion_pointer; i-- {
		t.mb.SetNodeKey(n, i, t.mb.GetNodeKey(n, i-1))
	}
	num_keys++
	t.mb.SetNodeKeyNum(n, num_keys)
	if insertion_pointer == 0 {
		if t.mb.GetLeafKey(tmp, 0) < k {
			t.mb.SetNodeKey(n, insertion_pointer, k)
		} else {
			t.mb.SetNodeKey(n, insertion_pointer, t.mb.GetLeafKey(tmp, 0))
		}
	} else {
		t.mb.SetNodeKey(n, insertion_pointer, k)
	}

	// Increment the number of used keys
	for i := num_keys; i > insertion_pointer; i-- {
		t.NodeSetLeafChild(n, i, t.mb.GetChild(n, i-1))
	}

	if t.mb.GetLeafKey(tmp, 0) < k {
		t.NodeSetLeafChild(n, insertion_pointer+1, e)
	} else {
		t.NodeSetLeafChild(n, insertion_pointer, e)
	}

}

func (t *DCSBT) NewNode() uint {
	n := t.mb.NewNodes(1)
	t.mb.SetChildren(n, t.mb.NewNodes(14))
	return n
}

// This message is used to find an actual leaf node
// identified by the _node_id struct which contains a link
// to the parent, the level and the offset
func (t *DCSBT) FindLeaf(k uint32, withPath bool) _node_id {
	u := t.mb.GetRoot()
	parent := uint(0)
	i, level := 0, 0
	var path *list.List
	if withPath {
		path = list.New()
	}
	if t.mb.IsLeaf(t.mb.GetRoot()) {
		r := _node_id{
			node:   u,
			level:  0,
			offset: 0,
			path:   path,
		}
		return r
	}

	for !t.mb.IsLeaf(u) {
		i = 0

		for i < t.mb.GetNodeKeyNum(u) {
			if k >= t.mb.GetNodeKey(u, i) {
				i++
			} else {
				break
			}
		}

		parent = u

		// Add the path to the list
		if withPath {
			path.PushFront(parent)
		}
		level++
		u = t.mb.GetChild(u, i)
	}

	return _node_id{
		node:   parent,
		level:  level,
		offset: i,
		path:   path,
	}
}

func (t *DCSBT) Find(k uint32) uintptr {
	i := 0
	var ln uint //Leaf
	if t.mb.GetRoot() == 0 {
		return 0
	}
	if !t.mb.IsLeaf(t.mb.GetRoot()) {
		//We found the leaf node that should contain k
		id := t.FindLeaf(k, false)
		if id.node == 0 {
			return 0
		}
		if t.mb.IsLeaf(id.node) {
			ln = id.node
		} else {
			ln = t.mb.GetChild(id.node, id.offset)
		}
	} else {
		ln = t.mb.GetRoot()
	}

	for i = 0; i < t.mb.GetLeafKeyNum(ln); i++ {
		if t.mb.GetLeafKey(ln, i) == k {
			return t.mb.GetLeafValue(ln, i)
		}
	}

	return 0
}

func (t *DCSBT) Insert(k uint32, v uintptr) {
	//Check if the key already exists
	if t.mb.GetRoot() != 0 && t.Find(k) != 0 {
		return
	}

	// Check if the tree is blank
	// Now insert the leaf node, but make sure, we take
	// one from the right level
	if t.mb.GetRoot() == 0 {
		n := t.mb.NewLeaves(1) //Leaf
		t.mb.SetLeafKeyNum(n, 1)
		t.mb.SetLeafKey(n, 0, k)
		t.mb.SetLeafValue(n, 0, v)
		t.mb.SetRoot(n)
		t.mb.SetMin(n)
		return
	}

	// Find the first parent node of the leaf
	// we want to insert to find it with path
	id := t.FindLeaf(k, true)
	var leaf uint
	if id.node == 0 {
		leaf = t.mb.GetRoot()
	} else {
		if t.mb.IsLeaf(id.node) {
			leaf = id.node
		} else {
			leaf = t.mb.GetChild(id.node, id.offset)
		}
	}

	//No Split
	if t.mb.GetLeafKeyNum(leaf) < 3 {
		t.InsertIntoLeaf(leaf, k, v)
	} else {
		// Enter the splitting process
		t.InsertIntoLeafAfterSplitting(id.path, leaf, k, v)
	}

}

func CUT(x int) int {
	if x%2 == 0 {
		return x / 2
	}
	return x/2 + 1
}

func (t *DCSBT) SplitLeafNode(leaf uint, k uint32, v uintptr) uint {
	num_keys := t.mb.GetLeafKeyNum(leaf)
	tmpK := make([]uint32, 3+1)
	tmpV := make([]uintptr, 3+1)
	cut_point := CUT(num_keys + 1)
	var insertion_pointer, i, j int
	insertion_pointer = 0
	for t.mb.GetLeafKey(leaf, insertion_pointer) < k && insertion_pointer < 3 {
		insertion_pointer++
	}
	for i, j = 0, 0; i < num_keys; i++ {
		if j == insertion_pointer {
			j++
		}
		tmpK[j] = t.mb.GetLeafKey(leaf, i)
		tmpV[j] = t.mb.GetLeafValue(leaf, i)
		j++
	}
	tmpK[insertion_pointer] = k
	tmpV[insertion_pointer] = v

	// Create a new leaf node and copy the values
	result := t.mb.NewLeaves(1)
	t.mb.SetLeafKeyNum(leaf, 0)
	// The old leafs need to be modified
	for i = 0; i < cut_point; i++ {
		t.mb.SetLeafKey(leaf, i, tmpK[i])
		t.mb.SetLeafValue(leaf, i, tmpV[i])
	}
	// The new leaf must be modified
	for i = cut_point; i < num_keys+1; i++ {
		t.mb.SetLeafKey(result, i-cut_point, tmpK[i])
		t.mb.SetLeafValue(result, i-cut_point, tmpV[i])
		t.mb.SetLeafKeyNum(result, t.mb.GetLeafKeyNum(result)+1)
	}
	// Modify the number of keys for the input leaf
	t.mb.SetLeafKeyNum(leaf, cut_point)
	return result
}

func (t *DCSBT) NodeAddKey(n uint, k uint32) {
	num_keys := t.mb.GetNodeKeyNum(n)
	t.mb.SetNodeKey(n, num_keys, k)
	num_keys++
	t.mb.SetNodeKeyNum(n, num_keys)
}

/**
 * When inserting a node and splitting it at the same time, we have to
 * check if the nodegroup size is sufficient to hold another node on the same
 * level, if this is the case SOP follows, if not, we have to reallocate the
 * nodegroup and copy the complete region plus, updating the first_child pointer
 * on the parent node (or realloc). When inserting the parent from their the
 * parent has to decide for its own.
 */

func (t *DCSBT) InsertIntoLeafAfterSplitting(path *list.List, leaf uint, k uint32, v uintptr) {
	// Leaf is the leaf we need to split, while path keeps track of the parents
	new_leaf := t.SplitLeafNode(leaf, k, v)
	if path.Len() == 0 {
		tmp := t.NewNode()
		t.mb.SetRoot(tmp)
		//Add both childs
		t.NodeAddLeafChild(tmp, 0, leaf)

		t.NodeAddKey(tmp, t.mb.GetLeafKey(new_leaf, 0))
		t.NodeAddLeafChild(tmp, 1, new_leaf)
	} else {
		parent := path.Front().Value.(uint)
		if t.mb.GetNodeKeyNum(parent) < 13 {
			// The new key must be added into the parent at
			// the right position
			t.NodeInsertNodeForKey(parent, t.mb.GetLeafKey(new_leaf, 0), new_leaf)
		} else {
			t.InsertLeafIntoParentAfterSplitting(path, new_leaf, t.mb.GetLeafKey(new_leaf, 0))
		}
	}
}

func (t *DCSBT) InsertLeafIntoParentAfterSplitting(path *list.List, first_child uint, k uint32) {
	i := 0
	node_to_split_e := path.Front()
	node_to_split := node_to_split_e.Value.(uint)
	path.Remove(node_to_split_e)

	//Find the point to split
	cut_point := CUT(t.mb.GetNodeKeyNum(node_to_split) + 1)

	//Create the new node
	new_node := t.NewNode()

	// As a next step the keys have to be copied from left ot right
	for i = cut_point; i < t.mb.GetNodeKeyNum(node_to_split); i++ {
		t.NodeAddLeafChild(new_node, i-cut_point, t.mb.GetChild(node_to_split, i))
		t.NodeAddKey(new_node, t.mb.GetLeafKey(node_to_split, i))

		//Unset the old node
		t.mb.Write(t.mb.GetChild(node_to_split, i), make([]byte, 64))
	}

	//Add the last child from the previous nodes
	t.NodeAddLeafChild(new_node, t.mb.GetNodeKeyNum(new_node), t.mb.GetChild(node_to_split, t.mb.GetNodeKeyNum(node_to_split)))

	//Set the Key amout for node to split
	t.mb.SetNodeKeyNum(node_to_split, cut_point-1)

	// Now we have to check where the new first_child has to be inserted
	var pivot uint32

	if k < t.mb.GetNodeKey(node_to_split, t.mb.GetNodeKeyNum(node_to_split)) {
		//Left
		t.NodeInsertNodeForKey(node_to_split, t.mb.GetLeafKey(first_child, 0), first_child)
		pivot = t.mb.GetLeafKey(t.mb.GetChild(new_node, 0), 0)
	} else {
		//Right
		t.NodeInsertNodeForKey(new_node, t.mb.GetLeafKey(first_child, 0), first_child)
		pivot = t.mb.GetLeafKey(t.mb.GetChild(new_node, 0), 0)
	}

	//Now push up
	if path.Len() == 0 {
		//New Root
		new_root := t.NewNode()
		t.NodeSetNodeChild(new_root, 0, node_to_split)
		t.mb.SetNodeKey(new_root, 0, pivot)
		t.mb.SetNodeKeyNum(new_root, 1)
		t.NodeSetNodeChild(new_root, 1, new_node)

		l := t.mb.GetChild(t.mb.GetChild(new_root, 0), t.mb.GetNodeKeyNum(t.mb.GetChild(new_root, 0)))
		r := t.mb.GetChild(t.mb.GetChild(new_root, 1), 0)

		t.mb.SetLeafRight(l, r)
		t.mb.SetLeafLeft(r, l)

		//Set the new root
		t.mb.SetRoot(new_root)
	} else {
		parent := path.Front().Value.(uint) //Node
		if t.mb.GetNodeKeyNum(parent) < 13 {
			t.NodeAddKey(parent, pivot)
			t.NodePushNodeBack(parent, new_node)
		} else {
			t.InsertNodeIntoParentAfterSplitting(path, new_node, pivot)
		}
	}
}

func (t *DCSBT) InsertNodeIntoParentAfterSplitting(path *list.List, first_child uint, k uint32) {
	var i int

	//Get the path and reduce by one
	node_to_split_e := path.Front()
	node_to_split := node_to_split_e.Value.(uint)
	path.Remove(node_to_split_e)

	cut_point := CUT(t.mb.GetNodeKeyNum(node_to_split) + 1)

	//Create the new node
	new_node := t.NewNode()

	for i = cut_point; i < t.mb.GetNodeKeyNum(node_to_split); i++ {
		t.NodeAddNodeChild(new_node, i-cut_point, t.mb.GetChild(node_to_split, i))
		t.NodeAddKey(new_node, t.mb.GetNodeKey(node_to_split, i))
	}

	t.NodeAddNodeChild(new_node, t.mb.GetNodeKeyNum(new_node), t.mb.GetChild(node_to_split, t.mb.GetNodeKeyNum(node_to_split)))
	t.NodeAddKey(new_node, k)
	t.NodePushNodeBack(new_node, first_child)
	t.mb.SetNodeKeyNum(node_to_split, cut_point)
	pivot := t.mb.GetNodeKey(new_node, 0)

	if path.Len() == 0 {
		new_root := t.NewNode()
		t.NodePushNodeBack(new_root, node_to_split)
		t.NodeAddKey(new_root, pivot)
		t.NodePushNodeBack(new_root, new_node)
		t.mb.SetRoot(new_root)
	} else {
		parent := path.Front().Value.(uint)
		if t.mb.GetNodeKeyNum(parent) < 13 {
			t.NodeAddKey(parent, pivot)
			t.NodePushNodeBack(parent, new_node)
		} else {
			t.InsertNodeIntoParentAfterSplitting(path, new_node, pivot)
		}
	}
}

func (t *DCSBT) InsertIntoLeaf(leaf uint, k uint32, v uintptr) {
	var insertion_pointer, i int
	insertion_pointer = 0
	for insertion_pointer < t.mb.GetLeafKeyNum(leaf) && t.mb.GetLeafKey(leaf, insertion_pointer) < k {
		insertion_pointer++
	}

	for i = t.mb.GetLeafKeyNum(leaf); i > insertion_pointer; i-- {
		t.mb.SetLeafKey(leaf, i, t.mb.GetLeafKey(leaf, i-1))
		t.mb.SetLeafValue(leaf, i, t.mb.GetLeafValue(leaf, i-1))
	}

	t.mb.SetLeafKeyNum(leaf, t.mb.GetLeafKeyNum(leaf)+1)
	t.mb.SetLeafKey(leaf, insertion_pointer, k)
	t.mb.SetLeafValue(leaf, insertion_pointer, v)
}
