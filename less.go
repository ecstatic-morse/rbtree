package rbtree

// This package also provides wrappers around a few common types to make
// them suitable for use in a tree, much like the convenience functions
// provided by 'sort'.

// All types to be stored in a red-black tree must implement a Less
// method which defines a strict weak ordering on the set of possible
// instances of that type.
//
// Specifically, for all x
// 	x.Less(x) == false
// and for all x and y,
// 	if x.Less(y) {
// 		y.Less(x) == false
// 	}
//
// Two items are equal if and only if neither is less than the other.
type Item interface {
	Less(than Item) bool
}

// Int wraps integers to provide a Less method.
type Int int

func (item Int) Less(than Item) bool {
	return item < than.(Int)
}

// Float64 wraps floating point numbers to provide a Less method.
type Float64 float64

func (item Float64) Less(than Item) bool {
	return item < than.(Float64)
}

// String wraps strings to provide a Less method.
type String string

func (item String) Less(than Item) bool {
	return item < than.(String)
}
