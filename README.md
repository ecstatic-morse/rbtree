rbtree
======

Package rbtree provides a conventional (not left-leaning) implementation
of red-black trees written in go. The internals are based on the [Linux kernel
red-black tree
implementation](http://lxr.free-electrons.com/source/lib/rbtree.c),
however a few features of the go language make this package much
easier to use than the kernel implementation.

Overview
--------

Red black trees are a type of self-balancing binary tree. As such,
insertion, deletion, and search are O(log n) in both the worst and
average cases. Compare this to hash tables (maps in go), which are
amortized O(1) for these operations in the average case, but O(n) in
the worst case.

In practice, hash tables are much faster than binary trees for most
use cases. Each step in a binary search requires at least one random
memory access (this implementation uses interfaces to store values
and provide comparison functions, so it requires even more),
resulting in poor cache performance. Nevertheless, they can be useful
if you are willing to sacrifice overall speed for improved worst case
performance on very large data sets. Also, if you need to iterate
over ranges of data, binary trees can do so efficiently.

License
-------

This project and all the files contained within are licensed under
GPLv2. See LICENSE.txt for details.
