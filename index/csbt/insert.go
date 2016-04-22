package csbt

import (
	"container/list"
	"errors"
	"strconv"
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
	t.MB.Copy(t.MB.GetChild(n, index), e, 64)
}

// If the leafnode is inserted into the internal node
// we know that no more data will be allocated below thos level
func (t *DCSBT) NodeSetLeafChild(n uint, index int, e uint) {
	num_keys := t.MB.GetNodeKeyNum(n)
	t.MB.Copy(t.MB.GetChild(n, index), e, 64)
	var l, r uint //Leaf
	if index == 0 {
		l = t.MB.GetChild(n, 0)
		r = t.MB.GetChild(n, 1)
		t.MB.SetLeafRight(l, r)
		t.MB.SetLeafLeft(r, l)
	} else if index > 0 && index < num_keys {
		l = t.MB.GetChild(n, index-1)
		r = t.MB.GetChild(n, index)
		t.MB.SetLeafRight(l, r)
		t.MB.SetLeafLeft(r, l)
		l = t.MB.GetChild(n, index)
		r = t.MB.GetChild(n, index+1)
		t.MB.SetLeafRight(l, r)
		t.MB.SetLeafLeft(r, l)
	} else if index == num_keys {
		l = t.MB.GetChild(n, index-1)
		r = t.MB.GetChild(n, index)
		t.MB.SetLeafRight(l, r)
		t.MB.SetLeafLeft(r, l)
	}
}

func (t *DCSBT) NodePushLeafBack(n uint, e uint) {
	num_keys := t.MB.GetNodeKeyNum(n)
	t.NodeAddLeafChild(n, num_keys, e)
}

// When an internal node is set, we have to copy the
// children of the parameter node as well
func (t *DCSBT) NodeSetNodeChild(n uint, index int, e uint) {
	// Copy the value from the node but than explicitely copy
	// the structure
	t.MB.Copy(t.MB.GetChild(n, index), e, 64)
	t.MB.SetChildren(t.MB.GetChild(n, index), t.MB.GetChild(e, 0))

	t.MB.SetChildren(e, 0)
}

func (t *DCSBT) NodePushNodeBack(n uint, e uint) {
	num_keys := t.MB.GetNodeKeyNum(n)
	// Copys the memory for the node
	t.NodeSetNodeChild(n, num_keys, e)

	if t.MB.IsLeaf(t.MB.GetChild(t.MB.GetChild(n, num_keys), 0)) {
		var l, r uint //Leaf
		tmp := t.MB.GetChild(n, num_keys-1)
		l = t.MB.GetChild(tmp, t.MB.GetNodeKeyNum(tmp))
		tmp = t.MB.GetChild(n, num_keys)
		r = t.MB.GetChild(tmp, 0)
		t.MB.SetLeafRight(l, r)
		t.MB.SetLeafLeft(r, l)
	}
}

func (t *DCSBT) NodeInsertNodeForKey(n uint, k uint32, e uint) {

	// 1. Find the insertion point for the key
	// 2. move all elements one to the right
	// 3. insert the child
	tmp := t.MB.GetChild(n, 0) //Leaf

	insertion_pointer := int(0)
	num_keys := t.MB.GetNodeKeyNum(n)
	for t.MB.GetNodeKey(n, insertion_pointer) < k && insertion_pointer < num_keys {
		insertion_pointer++
	}
	for i := num_keys; i > insertion_pointer; i-- {
		t.MB.SetNodeKey(n, i, t.MB.GetNodeKey(n, i-1))
	}
	num_keys++
	t.MB.SetNodeKeyNum(n, num_keys)
	if insertion_pointer == 0 {
		if t.MB.GetLeafKey(tmp, 0) < k {
			t.MB.SetNodeKey(n, insertion_pointer, k)
		} else {
			t.MB.SetNodeKey(n, insertion_pointer, t.MB.GetLeafKey(tmp, 0))
		}
	} else {
		t.MB.SetNodeKey(n, insertion_pointer, k)
	}

	// Increment the number of used keys
	for i := num_keys; i > insertion_pointer; i-- {
		t.NodeSetLeafChild(n, i, t.MB.GetChild(n, i-1))
	}

	if t.MB.GetLeafKey(tmp, 0) < k {
		t.NodeSetLeafChild(n, insertion_pointer+1, e)
	} else {
		t.NodeSetLeafChild(n, insertion_pointer, e)
	}

}

func (t *DCSBT) NewNode() uint {
	n := t.MB.NewNodes(1)
	t.MB.SetChildren(n, t.MB.NewNodes(14))
	return n
}

// This message is used to find an actual leaf node
// identified by the _node_id struct which contains a link
// to the parent, the level and the offset
func (t *DCSBT) FindLeaf(k uint32, withPath bool) _node_id {
	u := t.MB.GetRoot()
	parent := uint(0)
	i, level := 0, 0
	var path *list.List
	if withPath {
		path = list.New()
	}
	if t.MB.IsLeaf(t.MB.GetRoot()) {
		r := _node_id{
			node:   u,
			level:  0,
			offset: 0,
			path:   path,
		}
		return r
	}

	for !t.MB.IsLeaf(u) {
		i = 0

		for i < t.MB.GetNodeKeyNum(u) {
			if k >= t.MB.GetNodeKey(u, i) {
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
		u = t.MB.GetChild(u, i)
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
	if t.MB.GetRoot() == 0 {
		return 0
	}
	if !t.MB.IsLeaf(t.MB.GetRoot()) {
		//We found the leaf node that should contain k
		id := t.FindLeaf(k, false)
		if id.node == 0 {
			return 0
		}
		if t.MB.IsLeaf(id.node) {
			ln = id.node
		} else {
			ln = t.MB.GetChild(id.node, id.offset)
		}
	} else {
		ln = t.MB.GetRoot()
	}

	for i = 0; i < t.MB.GetLeafKeyNum(ln); i++ {
		if t.MB.GetLeafKey(ln, i) == k {
			return t.MB.GetLeafValue(ln, i)
		}
	}

	return 0
}

func (t *DCSBT) Insert(k uint32, v uintptr) error {
	//Check if the key already exists
	if c := t.Select(k); c != nil {
		if a, b := c.Read(); a == 0 && b == 0 { //Deleted
			p := uintptr(uint(v) >> 24)
			offset := uint(v) & 0x0000000000111111
			c.Write(p, offset)
			return nil
		}
		return errors.New("Key " + strconv.Itoa(int(k)) + " already exists.")
	}

	// Check if the tree is blank
	// Now insert the leaf node, but make sure, we take
	// one from the right level
	if t.MB.GetRoot() == 0 {
		n := t.MB.NewLeaves(1) //Leaf
		t.MB.SetLeafKeyNum(n, 1)
		t.MB.SetLeafKey(n, 0, k)
		t.MB.SetLeafValue(n, 0, v)
		t.MB.SetRoot(n)
		t.MB.SetMin(n)
		return nil
	}

	// Find the first parent node of the leaf
	// we want to insert to find it with path
	id := t.FindLeaf(k, true)
	var leaf uint
	if id.node == 0 {
		leaf = t.MB.GetRoot()
	} else {
		if t.MB.IsLeaf(id.node) {
			leaf = id.node
		} else {
			leaf = t.MB.GetChild(id.node, id.offset)
		}
	}

	//No Split
	if t.MB.GetLeafKeyNum(leaf) < 3 {
		t.InsertIntoLeaf(leaf, k, v)
	} else {
		// Enter the splitting process
		t.InsertIntoLeafAfterSplitting(id.path, leaf, k, v)
	}

	return nil

}

func CUT(x int) int {
	if x%2 == 0 {
		return x / 2
	}
	return x/2 + 1
}

func (t *DCSBT) SplitLeafNode(leaf uint, k uint32, v uintptr) uint {
	num_keys := t.MB.GetLeafKeyNum(leaf)
	tmpK := make([]uint32, 3+1)
	tmpV := make([]uintptr, 3+1)
	cut_point := CUT(num_keys + 1)
	var insertion_pointer, i, j int
	insertion_pointer = 0
	for t.MB.GetLeafKey(leaf, insertion_pointer) < k && insertion_pointer < 3 {
		insertion_pointer++
	}
	for i, j = 0, 0; i < num_keys; i++ {
		if j == insertion_pointer {
			j++
		}
		tmpK[j] = t.MB.GetLeafKey(leaf, i)
		tmpV[j] = t.MB.GetLeafValue(leaf, i)
		j++
	}
	tmpK[insertion_pointer] = k
	tmpV[insertion_pointer] = v

	// Create a new leaf node and copy the values
	result := t.MB.NewLeaves(1)
	t.MB.SetLeafKeyNum(leaf, 0)
	// The old leafs need to be modified
	for i = 0; i < cut_point; i++ {
		t.MB.SetLeafKey(leaf, i, tmpK[i])
		t.MB.SetLeafValue(leaf, i, tmpV[i])
	}
	// The new leaf must be modified
	for i = cut_point; i < num_keys+1; i++ {
		t.MB.SetLeafKey(result, i-cut_point, tmpK[i])
		t.MB.SetLeafValue(result, i-cut_point, tmpV[i])
		t.MB.SetLeafKeyNum(result, t.MB.GetLeafKeyNum(result)+1)
	}
	// Modify the number of keys for the input leaf
	t.MB.SetLeafKeyNum(leaf, cut_point)
	return result
}

func (t *DCSBT) NodeAddKey(n uint, k uint32) {
	num_keys := t.MB.GetNodeKeyNum(n)
	t.MB.SetNodeKey(n, num_keys, k)
	num_keys++
	t.MB.SetNodeKeyNum(n, num_keys)
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
		t.MB.SetRoot(tmp)
		//Add both childs
		t.NodeAddLeafChild(tmp, 0, leaf)
		t.MB.SetMin(t.MB.GetChild(tmp, 0))

		t.NodeAddKey(tmp, t.MB.GetLeafKey(new_leaf, 0))
		t.NodeAddLeafChild(tmp, 1, new_leaf)
	} else {
		parent := path.Front().Value.(uint)
		if t.MB.GetNodeKeyNum(parent) < 13 {
			// The new key must be added into the parent at
			// the right position
			t.NodeInsertNodeForKey(parent, t.MB.GetLeafKey(new_leaf, 0), new_leaf)
		} else {
			t.InsertLeafIntoParentAfterSplitting(path, new_leaf, t.MB.GetLeafKey(new_leaf, 0))
		}
	}
}

func (t *DCSBT) InsertLeafIntoParentAfterSplitting(path *list.List, first_child uint, k uint32) {
	i := 0
	node_to_split_e := path.Front()
	node_to_split := node_to_split_e.Value.(uint)
	path.Remove(node_to_split_e)

	//Find the point to split
	cut_point := CUT(t.MB.GetNodeKeyNum(node_to_split) + 1)

	//Create the new node
	new_node := t.NewNode()

	// As a next step the keys have to be copied from left ot right
	for i = cut_point; i < t.MB.GetNodeKeyNum(node_to_split); i++ {
		t.NodeAddLeafChild(new_node, i-cut_point, t.MB.GetChild(node_to_split, i))
		t.NodeAddKey(new_node, t.MB.GetLeafKey(node_to_split, i))

		//Unset the old node
		t.MB.Write(t.MB.GetChild(node_to_split, i), make([]byte, 64))
	}

	//Add the last child from the previous nodes
	t.NodeAddLeafChild(new_node, t.MB.GetNodeKeyNum(new_node), t.MB.GetChild(node_to_split, t.MB.GetNodeKeyNum(node_to_split)))

	//Set the Key amout for node to split
	t.MB.SetNodeKeyNum(node_to_split, cut_point-1)

	// Now we have to check where the new first_child has to be inserted
	var pivot uint32

	if k < t.MB.GetNodeKey(node_to_split, t.MB.GetNodeKeyNum(node_to_split)) {
		//Left
		t.NodeInsertNodeForKey(node_to_split, t.MB.GetLeafKey(first_child, 0), first_child)
		pivot = t.MB.GetLeafKey(t.MB.GetChild(new_node, 0), 0)
	} else {
		//Right
		t.NodeInsertNodeForKey(new_node, t.MB.GetLeafKey(first_child, 0), first_child)
		pivot = t.MB.GetLeafKey(t.MB.GetChild(new_node, 0), 0)
	}

	//Now push up
	if path.Len() == 0 {
		//New Root
		new_root := t.NewNode()
		t.NodeSetNodeChild(new_root, 0, node_to_split)
		t.MB.SetNodeKey(new_root, 0, pivot)
		t.MB.SetNodeKeyNum(new_root, 1)
		t.NodeSetNodeChild(new_root, 1, new_node)

		l := t.MB.GetChild(t.MB.GetChild(new_root, 0), t.MB.GetNodeKeyNum(t.MB.GetChild(new_root, 0)))
		r := t.MB.GetChild(t.MB.GetChild(new_root, 1), 0)

		t.MB.SetLeafRight(l, r)
		t.MB.SetLeafLeft(r, l)

		t.MB.SetMin(t.MB.GetChild(t.MB.GetChild(new_root, 0), 0))
		//Set the new root
		t.MB.SetRoot(new_root)
	} else {
		parent := path.Front().Value.(uint) //Node
		if t.MB.GetNodeKeyNum(parent) < 13 {
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

	cut_point := CUT(t.MB.GetNodeKeyNum(node_to_split) + 1)

	//Create the new node
	new_node := t.NewNode()

	for i = cut_point; i < t.MB.GetNodeKeyNum(node_to_split); i++ {
		t.NodeAddNodeChild(new_node, i-cut_point, t.MB.GetChild(node_to_split, i))
		t.NodeAddKey(new_node, t.MB.GetNodeKey(node_to_split, i))
	}

	t.NodeAddNodeChild(new_node, t.MB.GetNodeKeyNum(new_node), t.MB.GetChild(node_to_split, t.MB.GetNodeKeyNum(node_to_split)))
	t.NodeAddKey(new_node, k)
	t.NodePushNodeBack(new_node, first_child)
	t.MB.SetNodeKeyNum(node_to_split, cut_point)
	pivot := t.MB.GetNodeKey(new_node, 0)

	if path.Len() == 0 {
		new_root := t.NewNode()
		t.NodePushNodeBack(new_root, node_to_split)
		t.NodeAddKey(new_root, pivot)
		t.NodePushNodeBack(new_root, new_node)
		t.MB.SetRoot(new_root)
	} else {
		parent := path.Front().Value.(uint)
		if t.MB.GetNodeKeyNum(parent) < 13 {
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
	for insertion_pointer < t.MB.GetLeafKeyNum(leaf) && t.MB.GetLeafKey(leaf, insertion_pointer) < k {
		insertion_pointer++
	}

	for i = t.MB.GetLeafKeyNum(leaf); i > insertion_pointer; i-- {
		t.MB.SetLeafKey(leaf, i, t.MB.GetLeafKey(leaf, i-1))
		t.MB.SetLeafValue(leaf, i, t.MB.GetLeafValue(leaf, i-1))
	}

	t.MB.SetLeafKeyNum(leaf, t.MB.GetLeafKeyNum(leaf)+1)
	t.MB.SetLeafKey(leaf, insertion_pointer, k)
	t.MB.SetLeafValue(leaf, insertion_pointer, v)
}
