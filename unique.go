package rbtree

type Tree struct {
	inner tree
}

// New returns an initialized red-black tree.
func New() Tree {
	return Tree{}
}

// Returns true if the number of items in the tree is zero
func (t Tree) Empty() bool {
	return t.inner.Empty()
}

// Returns the minimum value in the tree or nil if the tree is empty.
// Runs in O(log n) time.
func (t Tree) Min() Item {
	return t.inner.Min()
}

// Returns the maximum value in the tree or nil if the tree is empty.
//
// Runs in O(log n) time.
func (t Tree) Max() Item {
	return t.inner.Max()
}

// Returns the size of the tree. Runs in O(1) time.
func (t Tree) Size() int {
	return t.inner.Size()
}

// InsertUnique inserts an item into a tree and returns true if an
// equivalent item does not already exist. If an equivalent item does
// exist, InsertUnique returns false and does not modify the tree.
func (t *Tree) Insert(item Item) bool {
	return t.inner.InsertUnique(item)
}

func (t *Tree) InsertOrReplace(item Item) Item {
	return t.inner.InsertOrReplace(item)
}

// Removes all items from the tree.
func (t *Tree) Clear() {
	t.inner.Clear()
}

func (t Tree) Find(item Item) (Iterator, bool) {
	return t.inner.Find(item)
}

func (t Tree) FindItem(item Item) Item {
	if it, ok := t.inner.Find(item); ok {
		return it.Item()
	} else {
		return nil
	}
}

// Delete looks for an item equivalent to target in the tree and deletes
// it, returning the value that was present in the tree. If no item was found,
// Delete returns nil and does not modify the tree.
func (t *Tree) Delete(item Item) Item {
	return t.inner.Delete(item)
}

// Returns an invalid Iterator pointing one past the beginning/end of
// the tree. (it != tree.End()) implies it.IsValid().
func (t Tree) End() Iterator {
	return Iterator{nil}
}

// Returns an Iterator pointing to the first item in the tree.
//
// Runs in O(log n) time.
func (t Tree) First() Iterator {
	return t.inner.First()
}

// Returns an Iterator pointing to the last item in the tree.
//
// Runs in O(log n) time.
func (t Tree) Last() Iterator {
	return t.inner.Last()
}
