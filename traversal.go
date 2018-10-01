package rbtree

type ordering int

const (
	lessThan ordering = iota - 1
	equalTo
	greaterThan
)

// Returns the maximum-valued node in a given subtree
func min(n *node) *node {
	for n.HasLeftChild() {
		n = n.left
	}

	return n
}

// Returns the maximum-valued node in a given subtree
func max(n *node) *node {
	for n.HasRightChild() {
		n = n.right
	}

	return n
}

// Returns the in-order predecessor of a given node.
func predecessor(n *node) *node {
	if n.HasLeftChild() {
		return max(n.left)
	}

	for p := n.Parent(); p != nil; n, p = p, p.Parent() {
		if n.IsRightChildOf(p) {
			return p
		}
	}

	return nil
}

// Returns the in-order successor of a given node.
func successor(n *node) *node {
	if n.HasRightChild() {
		return min(n.right)
	}

	for p := n.Parent(); p != nil; n, p = p, p.Parent() {
		if n.IsLeftChildOf(p) {
			return p
		}
	}

	return nil
}

// get tries to find the highest node in the tree equal to target. If it
// succeeds, it returns the highest node and true. If it fails, it
// returns the insertion point of target (the node that would become the
// parent of target were target to be inserted into the tree).
// getCmp is equivalent to get but it returns an integer indicating
// whether node is less than, equal to, or greater than target.
func get(n *node, subject Item) (*node, ordering) {
	for {
		switch {
		case subject.Less(n.item):
			if !n.HasLeftChild() {
				return n, lessThan
			}

			n = n.left
		case n.item.Less(subject):
			if !n.HasRightChild() {
				return n, greaterThan
			}

			n = n.right
		default:
			return n, equalTo
		}
	}
}

// getUpper finds the rightmost position where an item could be inserted in the tree.
func getUpper(n *node, subject Item) (*node, ordering) {
	for {
		switch {
		case subject.Less(n.item):
			if !n.HasLeftChild() {
				return n, lessThan
			}

			n = n.left
		default:
			if !n.HasRightChild() {
				if n.item.Less(subject) {
					return n, greaterThan
				} else {
					return n, equalTo
				}
			}

			n = n.right
		}
	}
}

// getLower finds the leftmost position where an item could be inserted in the tree.
func getLower(n *node, target Item) (*node, ordering) {
	for {
		switch {
		case n.item.Less(target):
			if !n.HasRightChild() {
				return n, greaterThan
			}

			n = n.right
		default:
			if !n.HasLeftChild() {
				if target.Less(n.item) {
					return n, lessThan
				} else {
					return n, equalTo
				}
			}

			n = n.left
		}
	}
}
