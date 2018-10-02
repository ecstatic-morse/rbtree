package rbtree

// A red-black tree whose items are unique.
//
// See MultiValuedTree for a red-black tree which allows duplicate items.
type Tree struct {
	inner tree
}

// Returns a fully initialized red-black tree.
func New() Tree {
	return Tree{}
}

// Returns true if the number of items in the tree is zero
func (t Tree) Empty() bool {
	return t.inner.Empty()
}

// Returns the minimum value in the tree or nil if the tree is empty.
//
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

// Inserts an item into the tree if an equivalent one does not already exist.
// Returns true if the item was inserted, or false if a duplicate was found.
//
// Runs in O(log n) time.
func (t *Tree) Insert(item Item) bool {
	return t.inner.InsertUnique(item)
}

// Inserts an item into the tree, or replaces an equivalent item if one exists.
// Returns the item which was previously in the tree, or nil if none was found.
//
// Runs in O(log n) time.
func (t *Tree) InsertOrReplace(item Item) Item {
	return t.inner.InsertOrReplace(item)
}

// Removes all items from the tree.
func (t *Tree) Clear() {
	t.inner.Clear()
}

// Searches the tree, returning an Iterator to the item if an equivalent one was
// found, along with a boolean indicating whether the search was successful.
//
// Runs in O(log n) time.
func (t Tree) Find(item Item) (Iterator, bool) {
	return t.inner.Find(item)
}

// Searches the tree, returning the Item if the search was successful, or nil if
// none was found.
//
// Runs in O(log n) time.
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
//
// Runs in O(log n) time.
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

// Returns an Iterator pointing to the smallest item greater than or equal to
// target.
//
// Runs in O(log n) time.
func (t Tree) LowerBound(target Item) Iterator {
	return t.inner.LowerBound(target)
}

// Returns an Iterator pointing to the smallest item greater than target.
//
// Runs in O(log n) time.
func (t Tree) UpperBound(target Item) Iterator {
	return t.inner.UpperBound(target)
}
