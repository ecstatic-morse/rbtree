package rbtree

// A red-black tree which allows multiple items with the same value to be
// inserted.
//
// See Tree for a red-black tree whose items are unique.
type MultiValuedTree struct {
	inner tree
}

// Returns a fully initialized red-black tree which allows for duplicate items.
func NewMultiValued() MultiValuedTree {
	return MultiValuedTree{}
}

// Returns true if the number of items in the tree is zero
func (t MultiValuedTree) Empty() bool {
	return t.inner.Empty()
}

// Returns the minimum value in the tree or nil if the tree is empty.
// Runs in O(log n) time.
func (t MultiValuedTree) Min() Item {
	return t.inner.Min()
}

// Returns the maximum value in the tree or nil if the tree is empty.
//
// Runs in O(log n) time.
func (t MultiValuedTree) Max() Item {
	return t.inner.Max()
}

// Returns the size of the tree. Runs in O(1) time.
func (t MultiValuedTree) Size() int {
	return t.inner.Size()
}

// Inserts an item into the tree.
//
// Runs in O(log n) time.
func (t *MultiValuedTree) Insert(item Item) {
	t.inner.Insert(item)
}

// Removes all items from the tree.
func (t *MultiValuedTree) Clear() {
	t.inner.Clear()
}

// Searches the tree, returning the Item if the search was successful, or nil if
// none was found.
//
// Runs in O(log n) time.
func (t MultiValuedTree) FindItem(item Item) Item {
	if it, ok := t.inner.Find(item); ok {
		return it.Item()
	} else {
		return nil
	}
}

// Delete looks for an item equivalent to target in the tree and deletes
// it, returning the value that was present in the tree. If no item was found,
// Delete returns nil and does not modify the tree.
func (t *MultiValuedTree) Delete(item Item) Item {
	return t.inner.Delete(item)
}

// Returns an Iterator pointing to the first item in the tree.
//
// Runs in O(log n) time.
func (t MultiValuedTree) First() Iterator {
	return t.inner.First()
}

// Returns an Iterator pointing to the last item in the tree.
//
// Runs in O(log n) time.
func (t MultiValuedTree) Last() Iterator {
	return t.inner.Last()
}

// Returns an invalid Iterator pointing one past the beginning/end of
// the tree. (it != tree.End()) implies it.IsValid().
func (t MultiValuedTree) End() Iterator {
	return t.inner.End()
}

// Returns an Iterator pointing to the smallest item greater than or equal to
// target.
//
// Runs in O(log n) time.
func (t MultiValuedTree) LowerBound(target Item) Iterator {
	return t.inner.LowerBound(target)
}

// Returns an Iterator pointing to the smallest item greater than target.
//
// Runs in O(log n) time.
func (t MultiValuedTree) UpperBound(target Item) Iterator {
	return t.inner.UpperBound(target)
}
