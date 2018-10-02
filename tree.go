// Package rbtree provides a conventional (not left-leaning)
// implementation of red-black trees. The internals are based on the
// Linux kernel red-black tree implementation
// (http://lxr.free-electrons.com/source/lib/rbtree.c).
//
// Red black trees are a type of self-balancing binary tree. As such,
// insertion, deletion, and search are O(log n) in both the worst and
// average cases. Compare this to hash tables (maps in go), which are
// amortized O(1) for these operations in the average case, but O(n) in
// the worst case.
//
// In practice, hash tables are much faster than binary trees for most use
// cases. Each step in a binary tree traversal requires at least one random
// memory access (this implementation uses interfaces to store values and
// provide comparison functions, so it requires even more), resulting in poor
// cache performance. Nevertheless, they can be useful if you are willing to
// sacrifice overall speed for improved worst case performance on very large
// data sets. Also, if you need to iterate over ranges of data, binary trees can
// do so efficiently.
package rbtree

type tree struct {
	root *node
	size int
}

// Returns true if the number of items in the tree is zero
func (t tree) Empty() bool {
	return t.root == nil
}

// Returns the minimum value in the tree or nil if the tree is empty.
// Runs in O(log n) time.
func (t tree) Min() Item {
	if t.Empty() {
		return nil
	}

	return min(t.root).item
}

// Returns the maximum value in the tree or nil if the tree is empty.
//
// Runs in O(log n) time.
func (t tree) Max() Item {
	if t.Empty() {
		return nil
	}

	return max(t.root).item
}

func (t tree) Size() int {
	return t.size
}

func (t tree) Find(item Item) (Iterator, bool) {
	if n, ord := get(t.root, item); ord == equalTo {
		return Iterator{n}, true
	} else {
		return t.End(), false
	}
}

func (t *tree) Insert(item Item) {
	n := newRedNode(item)
	t.size += 1

	if t.Empty() {
		n.SetBlack()
		t.root = n
		return
	}

	// The choice between rightmost and leftmost is arbitrary
	// TODO: benchmark?
	place, ord := getRightmostInsertionPoint(t.root, item)
	n.SetParent(place)

	// We know that place.item == item implies place.hasRightChild() == false
	// because otherwise getRightmostInsertionPoint would have continued to the
	// right.
	switch ord {
	case greaterThan, equalTo:
		place.right = n
	case lessThan:
		place.left = n
	}

	balanceAfterInsert(n, &t.root)
}

// Tries to insert a unique item into the tree. If the item already exists in the
// tree, does nothing and returns a pointer to the highest node in the
// hierarchy with the same item.
func (t *tree) insertUniqueOrReturnPlace(item Item) *node {
	if t.Empty() {
		n := newRedNode(item)
		n.SetBlack()
		t.size += 1
		t.root = n
		return nil
	}

	place, ord := get(t.root, item)
	if ord == equalTo {
		return place
	}

	n := newRedChildNode(item, place)
	t.size += 1
	switch ord {
	case greaterThan:
		place.right = n
	case lessThan:
		place.left = n
	}

	balanceAfterInsert(n, &t.root)
	return nil
}

// InsertUnique inserts an item into a tree and returns true if an
// equivalent item does not already exist. If an equivalent item does
// exist, InsertUnique returns false and does not modify the tree.
func (t *tree) InsertUnique(item Item) bool {
	return t.insertUniqueOrReturnPlace(item) == nil
}

func (t *tree) InsertOrReplace(item Item) Item {
	if place := t.insertUniqueOrReturnPlace(item); place != nil {
		// Swap the old item for the new
		item, place.item = place.item, item
		return item
	} else {
		return nil
	}
}

// Removes all items from the tree.
func (t *tree) Clear() {
	t.size = 0
	t.root = nil
}

// Delete looks for an item equivalent to target in the tree and deletes
// it, returning the value that was present in the tree. If no item was found,
// Delete returns nil and does not modify the tree.
func (t *tree) Delete(item Item) Item {
	n, ord := get(t.root, item)
	if ord != equalTo {
		return nil
	}

	item = deleteNode(n, &t.root)
	t.size -= 1

	// If we deleted the last element in the tree, we now have nilChild as the root pointer.
	if t.root == nilChild {
		t.root = nil
	}

	return item
}

// Returns an Iterator pointing to the first item in the tree,
//
// Runs in O(log n) time.
func (t tree) First() Iterator {
	if t.Empty() {
		return t.End()
	} else {
		return Iterator{min(t.root)}
	}
}

// Returns an Iterator pointing to the last item in the tree.
//
// Runs in O(log n) time.
func (t tree) Last() Iterator {
	if t.Empty() {
		return t.End()
	} else {
		return Iterator{max(t.root)}
	}
}

// Returns an invalid Iterator pointing one past the beginning/end of
// the tree. (it != tree.End()) implies it.IsValid().
func (t tree) End() Iterator {
	return Iterator{nil}
}

// Returns an Iterator pointing to the first item greater than or equal to target.
func (t tree) LowerBound(target Item) Iterator {
	n, ord := getLeftmostInsertionPoint(t.root, target)

	// If the target is greater than the insertion point, we actually want the
	// successor of the node.
	if ord == greaterThan {
		n = successor(n)
	}

	return Iterator{n}
}

// Returns an Iterator pointing to the first item greater than target.
func (t tree) UpperBound(target Item) Iterator {
	n, ord := getRightmostInsertionPoint(t.root, target)

	// If the target is greater than or equal to the insertion point, we
	// actually want the successor of the node.
	if ord != lessThan {
		n = successor(n)
	}

	return Iterator{n}
}
