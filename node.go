package rbtree

// red-black tree properties:  http://en.wikipedia.org/wiki/Rbtree
//
//  1) A node is either red or black
//  2) The root is black
//  3) All leaves (NULL) are black
//  4) Both children of every red node are black
//  5) Every simple path from root to leaves contains the same number
//     of black nodes.

type node struct {
	black       bool
	parent      *node
	left, right *node

	item Item
}

// This sentinel represents the null leaf nodes of an rb tree. We could
// use nil as the child pointer, but having an actual node simplifies
// traversal and some other operations.
var nilChild = &node{black: true}

// Returns a new red node containing the given item with no parent or children.
func newRedNode(item Item) *node {
	return &node{
		item:  item,
		left:  nilChild,
		right: nilChild,
	}
}

// Returns a new red node with the given parent pointer
func newRedChildNode(item Item, parent *node) *node {
	return &node{
		item:   item,
		left:   nilChild,
		right:  nilChild,
		parent: parent,
	}
}

// Getters and setters for parent node and color.
//
// TODO: we could store the node's color in the low bit of the parent pointer, since
// nodes should be at least 2-byte aligned.
func (n *node) IsRoot() bool        { return n.parent == nil }
func (n *node) HasLeftChild() bool  { return n.left != nilChild }
func (n *node) HasRightChild() bool { return n.right != nilChild }
func (n *node) IsBlack() bool       { return n.black }
func (n *node) IsRed() bool         { return !n.black }
func (n *node) SetBlack()           { n.black = true }
func (n *node) SetRed()             { n.black = false }
func (n *node) CopyColorOf(o *node) { n.black = o.black }
func (n *node) Parent() *node       { return n.parent }
func (n *node) SetParent(p *node)   { n.parent = p }

func (n *node) IsLeftChildOf(p *node) bool  { return p.left == n }
func (n *node) IsRightChildOf(p *node) bool { return p.right == n }

func (n *node) Children() [2]*node {
	return [...]*node{n.left, n.right}
}

// Rotates the left child of root clockwise so that it becomes the new parent
// of root, without fixing the child pointer of root's previous parent.
//
//        r         p
//       / \       / \
//      p   b  -> a   r
//     / \           / \
//    a   o         o   b
//
// A rotation requires three steps:
//   1. Change ownership of orphan(o) from pivot(p) to root(r).
//   2. Make pivot the parent of root and root the child of pivot.
//   3. Update root's previous parent's child pointer to point to pivot.
//
// RotateRightNoFixup performs the first two steps, but leaves the third to the caller.
// The caller likely knows which child pointer (right or left) must be updated
// to point to pivot.
//
// If calling code has no extra information about which side of its parent root
// was on, fixupAfterRotate should be called to perform the update. It requires
// two conditional branches.
//
// TODO: These might not be inlined because, thanks to the node getters and
// setters, they're not leaf functions.
func rotateRightNoFixup(root *node) {
	pivot := root.left

	// Change ownership of orphan from pivot to root
	orphan := pivot.right
	root.left = orphan
	orphan.SetParent(root)

	// Make pivot the parent of root
	pivot.SetParent(root.Parent())
	pivot.right = root
	root.SetParent(pivot)
}

// Same as rotateRightNoFixup, but rotates the right child of root counterclockwise.
func rotateLeftNoFixup(root *node) {
	pivot := root.right

	// Change ownership of orphan from pivot to root
	orphan := pivot.left
	root.right = orphan
	orphan.SetParent(root)

	// Make pivot the parent of root
	pivot.SetParent(root.Parent())
	pivot.left = root
	root.SetParent(pivot)
}

// Performs step 3 of a rotation.
//
// Calling rotate{Left,Right}NoFixup followed by fixupAfterRotate performs a full rotation.
func fixupAfterRotate(oldRoot *node, treeRoot **node) {
	// At this point, the old root and the new root (pivot) have had their
	// parent pointers updated, but the new root's parent's child pointer has
	// not yet been updated.
	newRoot := oldRoot.Parent()
	parent := newRoot.Parent()
	switch {
	case parent == nil:
		*treeRoot = newRoot
	case parent.left == oldRoot:
		parent.left = newRoot
	case parent.right == oldRoot:
		parent.right = newRoot
	}
}

// Balances a tree after inserting a node n, returning a pointer to the new
// root node of the tree, or nil if the tree root remains unchanged.
func balanceAfterInsert(x *node, treeRoot **node) {
	for {
		// Loop invariant: node x is red

		// Case 1: If x is the root node, set its color to black and return.
		if x.IsRoot() {
			x.SetBlack()
			*treeRoot = x
			return
		}

		parent := x.Parent()

		// Case 2: If the node's parent is black, the tree is valid.
		if parent.IsBlack() {
			return
		}

		// The root node is black, and parent is not black, therefore
		// parent is not the root node, and grandparent always exists.
		gparent := parent.Parent()

		if parent.IsLeftChildOf(gparent) {
			//       G            g
			//      / \          / \
			//     p   u  -->   P   U
			//    /            /
			//   x            x
			//
			// Case 3:
			// If uncle(u) is red along with parent(p), we flip the color of the
			// grandparent(g) to black and recurse with g as the new node.
			uncle := gparent.right
			if uncle.IsRed() { // Leaf nodes are always black
				parent.SetBlack()
				uncle.SetBlack()
				gparent.SetRed()
				x = gparent
				continue
			}

			//      G             G
			//     / \           / \
			//    p   U  -->    x   U
			//     \           /
			//      x         p
			//
			// Case 4:
			// If x is the right child of p, rotate left at p. x and p have now
			// swapped positions in the hierarchy.
			if x.IsRightChildOf(parent) {
				rotateLeftNoFixup(parent)
				gparent.left = x
				parent = x
			}

			//        G           P
			//       / \         / \
			//      p   U  -->  x   g
			//     /                 \
			//    x                   U
			//
			// Case 5:
			// Right rotate at grandparent.
			parent.SetBlack()
			gparent.SetRed()
			rotateRightNoFixup(gparent)
			fixupAfterRotate(gparent, treeRoot)
			return
		} else { // parent.IsRightChildOf(gparent)
			uncle := gparent.left

			// Case 3
			if uncle.IsRed() {
				parent.SetBlack()
				uncle.SetBlack()
				gparent.SetRed()
				x = gparent
				continue
			}

			// Case 4
			if parent.left == x {
				rotateRightNoFixup(parent)
				gparent.right = x
				parent = x
			}

			// Case 5
			parent.SetBlack()
			gparent.SetRed()
			rotateLeftNoFixup(gparent)
			fixupAfterRotate(gparent, treeRoot)
			return
		}
	}
}

// Balances a tree after deleting a node which used to occupy the same place in
// the tree as x.
func balanceAfterDelete(x *node, treeRoot **node) {
	for {
		// Case 1: If x is the root node, the tree is balanced.
		if x.IsRoot() {
			*treeRoot = x
			return
		}

		parent := x.Parent()

		if x.IsLeftChildOf(parent) {
			sibling := parent.right

			//     P               S
			//    / \             / \
			//   X   s    -->    p   Rn
			//      / \         / \
			//     Ln  Rn      X   Ln
			//
			// Case 2 - left rotate at parent
			if sibling.IsRed() {
				parent.SetRed()
				sibling.SetBlack()
				rotateLeftNoFixup(parent)
				fixupAfterRotate(parent, treeRoot)
				sibling = parent.right
			}

			//    (p)           (p)
			//    / \           / \
			//   X   S    -->  X   s
			//      / \           / \
			//     Ln  Rn        Ln  Rn
			//
			// Case 3 - sibling color flip
			// (p could be either color here)
			//
			// This leaves us violating 5) which
			// can be fixed by flipping p to black
			// if it was red, or by recursing at p.
			// p is red when coming from Case 2.
			leftNiece, rightNiece := sibling.left, sibling.right
			if sibling.IsBlack() && leftNiece.IsBlack() && rightNiece.IsBlack() {
				sibling.SetRed()
				if parent.IsRed() {
					parent.SetBlack()
					return
				} else {
					x = parent
					continue
				}
			}

			//   (p)           (p)
			//   / \           / \
			//  X   S    -->  X   Ln
			//     / \             \
			//    ln  Rn            s
			//                       \
			//                        Rn
			//
			// Case 4 - right rotate at sibling
			// (p could be either color here)
			if leftNiece.IsRed() && rightNiece.IsBlack() {
				leftNiece.SetBlack()
				sibling.SetRed()
				rotateRightNoFixup(sibling)
				parent.right = leftNiece
				sibling, leftNiece, rightNiece = leftNiece, leftNiece.left, sibling
			}

			//      (p)             (s)
			//      / \             / \
			//     X   S     -->   P   Rn
			//        / \         / \
			//      (ln) rn      X  (ln)
			//
			// Case 5 - left rotate at parent + color flips
			// (p and ln could be either color here.
			// After rotation, p becomes black, s acquires
			// p's color, and ln keeps its color)
			sibling.CopyColorOf(parent)
			parent.SetBlack()
			rightNiece.SetBlack()
			rotateLeftNoFixup(parent)
			fixupAfterRotate(parent, treeRoot)
			return
		} else { // x == parent.right
			sibling := parent.left

			// Case 2
			if sibling.IsRed() {
				parent.SetRed()
				sibling.SetBlack()
				rotateRightNoFixup(parent)
				fixupAfterRotate(parent, treeRoot)
				sibling = parent.left
			}

			// Case 3
			leftNiece, rightNiece := sibling.left, sibling.right
			if sibling.IsBlack() && leftNiece.IsBlack() && rightNiece.IsBlack() {
				sibling.SetRed()
				if parent.IsRed() {
					parent.SetBlack()
					return
				} else {
					x = parent
					continue
				}
			}

			// Case 4
			if leftNiece.IsBlack() && rightNiece.IsRed() {
				rightNiece.SetBlack()
				sibling.SetRed()
				rotateLeftNoFixup(sibling)
				parent.left = rightNiece
				sibling, rightNiece, leftNiece = rightNiece, rightNiece.right, sibling
			}

			// Case 5
			sibling.CopyColorOf(parent)
			parent.SetBlack()
			leftNiece.SetBlack()
			rotateRightNoFixup(parent)
			fixupAfterRotate(parent, treeRoot)
			return
		}
	}
}

func deleteNode(x *node, treeRoot **node) (deleted Item) {
	deleted = x.item

	// If node to be deleted has two non-leaf children, replace its item with
	// that of its in-order successor (or predecessor) and delete the
	// successor.
	if x.HasLeftChild() && x.HasRightChild() {
		succ := min(x.right)
		x.item = succ.item
		x = succ
	}

	// x now has at most one non-leaf child
	child := x.left
	if !x.HasLeftChild() {
		child = x.right
	}

	// Replace x with its non-leaf child (or a leaf if both children are leaves)
	parent := x.Parent()
	child.SetParent(parent)

	// If x was the root node, there's no child pointer to update, and we can make its child the new root.
	if x.IsRoot() {
		child.SetBlack()
		*treeRoot = child
		return
	}

	if x.IsLeftChildOf(parent) {
		parent.left = child
	} else {
		parent.right = child
	}

	// If x was a red node, we can replace it with its child without altering the number of
	// black nodes in a path.
	if x.IsRed() {
		return
	}

	// If x was black but its child is red, simply recolor the child.
	if child.IsRed() {
		child.SetBlack()
		return
	}

	// Otherwise we need to do a recursive reblance.
	balanceAfterDelete(child, treeRoot)
	return
}
