# Go Data Structures

by Tim Henderson (tim.tadh@gmail.com)

Copyright 2013, Licensed under the GPL version 2.

## Purpose

To collect many important data structures for usage in go programs. Some of
these data structures may have implementations elsewhere but are collected here
for completeness and instructional usage.

# Current Collection

## Trees

### Classic AVL Tree

An AVL tree is a height balanced binary search tree. It is commonly taught in
algorithms courses. The one I implemented here is from the exercise in CLRS\*:

> An AVL Tree is a binary search tree that is **height balanced**: for each node
> *x*, the heights of the left and right subtrees differ by at most 1. To
> implement an AVL tree, we maintain an extra field in each node: *h*[x] is the
> height of the node *x*.
>
> 2. To insert into an AVL tree, a node is first placed in the appropriated
>    place in the binary search tree order. After this insertion, the tree
>    may no longer be height blanced. Specifically, the heights of the left
>    and right children of some node may differ by 2. Describe a procedure
>    BALANCE(*x*), which takes a subtree rooted at *x* whole left and right
>    children are height balanced and have heights that differ by at most
>    2.

\* : Cormen, T., Leiserson, C., Rivest, R., Stein, C. *Introduction to
Algorithms*. McGraw-Hill Book Company. 2nd Edition. 2001.

## Hash Tables

### Classic Separate Chaining Hash Table

See `hashtable/hashtable.go`. An implementation of the classic hash table with
separate chaining to handle collisions.

### Linear Hash Table with AVL Tree Buckets

See `hashtables/linhash.go`. An implementation of [Linear
Hashing](http://hackthology.com/linear-hashing.html), a technique usually used
for secondary storage hash tables. Often employed by databases and file systems
for hash indices. This version is mostly instructional see the
[accompanying blog post](
http://hackthology.com/an-in-memory-go-implementation-of-linear-hashing.html).
If you want the "real" disk backed version you want to check my
[file-structures](https://github.com/timtadh/file-structures) repository. See
the `linhash` directory.

### Benchmarks

    $ go test -v -bench ".*" github.com/timtadh/data-structures/hashtable

    BenchmarkGoMap    100000             28973 ns/op
    BenchmarkHash      20000             76326 ns/op
    BenchmarkMLHash    50000             73516 ns/op
    BenchmarkLHash       500           3307005 ns/op

The performance of the in memory linear hash (MLHash) is slightly improved since
the [blog post](
http://hackthology.com/an-in-memory-go-implementation-of-linear-hashing.html) do
to the usage of an AVL Tree `tree/avltree.go` instead of an unbalanced binary
search tree.

