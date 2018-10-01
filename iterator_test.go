package rbtree

import (
	"fmt"
	"testing"
)

func TestRange(t *testing.T) {
	find := func(t Tree, i int) Iterator {
		it, _ := t.Find(Int(i))
		return it
	}

	tree := New()
	tree.Insert(Int(2))
	tree.Insert(Int(4))
	tree.Insert(Int(1))
	tree.Insert(Int(5))
	tree.Insert(Int(3))

	assertRangeEq(t, tree.First(), tree.End(), []int{1, 2, 3, 4, 5})
	assertRangeEq(t, find(tree, 2), find(tree, 5), []int{2, 3, 4})
	assertRangeEq(t, tree.First(), find(tree, 4), []int{1, 2, 3})
	assertRangeEq(t, find(tree, 3), tree.End(), []int{3, 4, 5})
	assertRangeEq(t, find(tree, 5), tree.End(), []int{5})
}

func ExampleIterator() {
	tree := New()
	tree.Insert(Int(2))
	tree.Insert(Int(4))
	tree.Insert(Int(1))
	tree.Insert(Int(5))
	tree.Insert(Int(3))

	end, _ := tree.Find(Int(4))
	for it := tree.First(); it != end; it.Next() {
		fmt.Printf("%d ", it.Item().(Int))
	}
	// Output: 1 2 3
}

func ExampleIterator_reverse() {
	tree := New()
	tree.Insert(Int(2))
	tree.Insert(Int(4))
	tree.Insert(Int(1))
	tree.Insert(Int(5))
	tree.Insert(Int(3))

	for it, _ := tree.Find(Int(3)); it.IsValid(); it.Prev() {
		fmt.Printf("%d ", it.Item().(Int))
	}
	// Output: 3 2 1
}

func ExampleIterator_UpperBound() {
	tree := NewMultiValued()
	tree.Insert(Int(2))
	tree.Insert(Int(1))
	tree.Insert(Int(2))
	tree.Insert(Int(3))
	tree.Insert(Int(2))

	for it, end := tree.LowerBound(Int(2)), tree.UpperBound(Int(2)); it != end; it.Next() {
		fmt.Printf("%d ", it.Item().(Int))
	}
	// Output: 2 2 2
}
