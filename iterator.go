package rbtree

// Iterators are the most efficient way of enumerating the elements of a
// tree in sorted order. Iterators are bidirectional (can be incremented
// or decremented) but not random access (cannot be moved to an
// arbitrary element in constant time).
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

// Returns true if the iterator points to an element in the tree. Once
// the iterator is incremented/decremented past the last/first element
// in the tree, IsValid will return false.
func (it Iterator) IsValid() bool { return it.node != nil }

/*
// The ForEach family of functions call fn for each element in the range
// [lower, upper). Passing nil for either limit acts as though you
// passed an element less than (for lower) or greater than (for upper)
// all elements in the tree.
func (t Tree) ForEach(fn func(Item))                  { t.ForEachIn(nil, nil, fn) }
func (t Tree) ForEachFrom(lower Item, fn func(Item))  { t.ForEachIn(lower, nil, fn) }
func (t Tree) ForEachUntil(upper Item, fn func(Item)) { t.ForEachIn(nil, upper, fn) }
func (t Tree) ForEachIn(lower, upper Item, fn func(Item)) {
	// Ensure our lower bound is actually less than the upper bound.
	if lower != nil && upper != nil && !lower.Less(upper) {
		return
	}

	it, end := t.First(), t.End()
	if upper != nil {
		end = t.Iter(upper)
	}
	if lower != nil {
		it = t.Iter(lower)
	}

	for it != end {
		fn(it.Item())
		it.Next()
	}
}

// The range family of functions return a channel that is sent all the
// elements in a given tree in sorted order. This is to allow for-range
// loops to iterate over red-black trees. For a description of the
// limits, see ForEach.
func (t Tree) Range() <-chan Item                { return t.RangeIn(nil, nil) }
func (t Tree) RangeFrom(lower Item) <-chan Item  { return t.RangeIn(lower, nil) }
func (t Tree) RangeUntil(upper Item) <-chan Item { return t.RangeIn(nil, upper) }
func (t Tree) RangeIn(lower, upper Item) <-chan Item {
	ch := make(chan Item, 1024)
	cb := func(item Item) { ch <- item }
	go func() {
		t.ForEachIn(lower, upper, cb)
		close(ch)
	}()

	return ch
}
*/
