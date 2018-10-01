package rbtree

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

func ExampleTree_Insert() {
	tree := New()
	tree.Insert(String("world"))
	tree.Insert(String("hello"))

	// Strings are ordered lexicographically (e.g. "hello" < "world")
	fmt.Println(tree.Max().(String))
	// Output: world
}

func ExampleTree_Delete() {
	tree := New()
	tree.Insert(String("foo"))
	tree.Insert(String("bar"))
	fmt.Println(tree.Min().(String))

	tree.Delete(String("bar"))
	fmt.Println(tree.Min().(String))
	// Output:
	// bar
	// foo
}

type keyValue struct {
	key   int
	value string
}

func (o keyValue) Less(than Item) bool {
	return o.key < than.(keyValue).key
}

func ExampleItem_keyValue() {
	/*
		type keyValue struct {
			key   int
			value string
		}

		func (o keyValue) Less(than Item) {
			return o.key < than.(keyValue).key
		}
	*/

	tree := New()
	tree.Insert(keyValue{0, "zero"})
	tree.Insert(keyValue{1, "one"})
	item := tree.FindItem(keyValue{key: 1})
	fmt.Println(item.(keyValue).value)
	// Output: one
}

// Repeatedly inserts and removes items.
func TestTree(t *testing.T) {
	rand.Seed(42)

	nonmembers := make([]int, 1000)
	for i := range nonmembers {
		nonmembers[i] = i
	}

	members := make([]int, 0)

	tree := New()
	for i := 0; i < 1000000; i += 1 {
		will_insert := rand.Float64() < probabilityOfInsert(tree.Size())
		if will_insert {
			// Insert
			i := rand.Intn(len(nonmembers))
			item := swapBetween(i, &nonmembers, &members)
			// t.Log("Inserting", item)

			if !tree.Insert(Int(item)) {
				t.Fatal("Inserted unique element but InsertUnique failed")
			}
		} else {
			// Delete
			i := rand.Intn(len(members))
			item := swapBetween(i, &members, &nonmembers)
			// t.Log("Deleting", item)

			if tree.Delete(Int(item)) == nil {
				t.Fatal("Failed to find deleted item")
			}
		}

		checkTree(t, tree.inner, members)
	}
}

func TestMultiValuedTree(t *testing.T) {
	rand.Seed(43)

	members := make([]int, 0)
	tree := NewMultiValued()
	for i := 0; i < 1000000; i += 1 {
		will_insert := rand.Float64() < probabilityOfInsert(tree.Size())
		if will_insert {
			// Insert
			item := rand.Intn(100)
			members = append(members, item)
			// t.Log("Inserting", item)
			tree.Insert(Int(item))

		} else {
			// Delete
			i := rand.Intn(len(members))
			item := members[i]
			members[i] = members[len(members)-1]
			members = members[:len(members)-1]
			// t.Log("Deleting", item)

			if tree.Delete(Int(item)) == nil {
				t.Fatal("Failed to find deleted item")
			}
		}

		checkTree(t, tree.inner, members)
	}
}

// The chances of inserting vs deleting for a tree of a given size
func probabilityOfInsert(size int) float64 {
	switch {
	case size == 0:
		return 1.0
	case size == 1000:
		return 0.0
	case size < 4:
		return 0.7
	case size > 16:
		return 0.3
	default:
		return 0.5
	}
}

// Swaps the ith element of from into to, returning it
func swapBetween(i int, from, to *[]int) int {
	el := (*from)[i]
	*to = append(*to, el)
	(*from)[i] = (*from)[len(*from)-1]
	*from = (*from)[:len(*from)-1]
	return el
}

// Checks that a tree contains the correct members and does not violate the red/black invariants
func checkTree(t *testing.T, tree tree, members []int) {
	if tree.Size() != len(members) {
		t.Fatal("tree size was not updated properly")
	}

	checkTreeInvariants(t, tree.root)
	if t.Failed() {
		t.FailNow()
	}

	// Check iteration order for the whole tree.
	sort.Ints(members)
	// t.Log("Tree should contain", members)
	assertRangeEq(t, tree.First(), tree.End(), members)

	// Check iteration order for a subrange of the tree.

	if tree.Empty() {
		return
	}

	lo := rand.Intn(len(members))
	hi := rand.Intn(len(members)-lo) + lo

	// If hi was chosen to be the middle of a run of duplicates, advance it to the last duplicate.
	for hi+1 < len(members) && members[hi+1] == members[hi] {
		hi += 1
	}

	// Do the same for lo
	for lo-1 >= 0 && members[lo-1] == members[lo] {
		lo -= 1
	}

	begin := tree.LowerBound(Int(members[lo]))
	end := tree.UpperBound(Int(members[hi]))
	expect := members[lo : hi+1]
	assertRangeEq(t, begin, end, expect)
}

func checkTreeInvariants(t *testing.T, x *node) {
	if x == nil {
		return
	}

	expectedBlackAncestors := -1

	var check func(x *node, blackAncestors int)
	check = func(x *node, blackAncestors int) {
		if x.IsRoot() && x.IsRed() {
			t.Errorf("Root node must be black")
		}

		if x.IsRed() && (x.left != nil && x.left.IsRed() || x.right != nil && x.right.IsRed()) {
			t.Errorf("Both children of a red node must be black")
		}

		if x.IsBlack() {
			blackAncestors += 1
		}

		for _, child := range x.Children() {
			if child == nilChild {
				// Leaf node
				if expectedBlackAncestors == -1 {
					expectedBlackAncestors = blackAncestors
					continue
				} else if blackAncestors != expectedBlackAncestors {
					t.Errorf("Every path from a given node to any of its descendent nodes must contain the same number of black nodes")
					break
				}
			} else {
				if child.Parent() != x {
					t.Errorf("Invalid parent pointer")
					break
				}

				check(child, blackAncestors)
			}
		}

		return
	}

	check(x, 0)
}

func TestSuccessorPredecessor(t *testing.T) {
	tree := New()
	tree.Insert(Int(3))
	if nil != successor(tree.inner.root) {
		t.Fatal("Successor failed")
	}

	if nil != predecessor(tree.inner.root) {
		t.Fatal("Predecessor failed")
	}
}

func assertRangeEq(t *testing.T, begin, end Iterator, expected []int) {
	i := 0
	for it := begin; it != end && it.IsValid(); it.Next() {
		if i >= len(expected) {
			i += 1
			continue
		}

		item := int(it.Item().(Int))
		if item != expected[i] {
			t.Errorf("Expected item %d to be %d, got %d", i, expected[i], item)
		}

		i += 1
	}

	if i != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), i)
	}

	if t.Failed() {
		t.FailNow()
	}
}

// Precompute random numbers
func randRange(size, seed int) (slice []Int) {
	rng := rand.New(rand.NewSource(int64(seed)))
	slice = make([]Int, size)
	for i := 0; i < size; i++ {
		slice[i] = Int(rng.Int())
	}

	return
}

// Build a large tree of random integers.
func BenchmarkRBInsert(b *testing.B) {
	ints := randRange(1<<16, 43)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := New()
		for _, n := range ints {
			tree.Insert(n)
		}
	}
}

// Build a large tree of random integers, then delete every element one
// by one.
func BenchmarkRBDelete(b *testing.B) {
	ints := randRange(1<<16, 43)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Build the tree
		b.StopTimer()
		tree := New()
		for _, n := range ints {
			tree.Insert(n)
		}
		b.StartTimer()

		// Delete every item in the tree
		for _, n := range ints {
			tree.Delete(n)
		}
	}
}
