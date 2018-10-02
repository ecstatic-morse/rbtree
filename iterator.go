package rbtree

// Iterators are an efficient way to enumerate the items contained within a
// tree. Iterators are bidirectional (they can be advanced forwards or backwards)
// but not random access (they cannot advanced by more than one step at a time).
type Iterator struct {
	node *node
}

// Advances an iterator to the previous element in the tree. Prev must
// not be called if the iterator is no longer valid.
func (it *Iterator) Prev() {
	it.node = predecessor(it.node)
}

// Advances an iterator to the next element in the tree. Next must
// not be called if the iterator is no longer valid.
func (it *Iterator) Next() {
	it.node = successor(it.node)
}

// Returns the item pointed to by the iterator. Item must not be called
// if the iterator is no longer valid.
func (it Iterator) Item() Item { return it.node.item }

// Returns true if the iterator points to an element in the tree. Once the
// iterator is advanced past the last (or first) element in the tree, IsValid
// will return false.
func (it Iterator) IsValid() bool { return it.node != nil }
