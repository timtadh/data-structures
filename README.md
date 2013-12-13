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
algorithms courses

### Immutable AVL Tree

This version of the classic is immutable and should be thread safe due to
immutability. However, there is a performance hit:

    BenchmarkAvlTree           10000            166657 ns/op
    BenchmarkImmutableAvlTree   5000            333709 ns/op

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
    $ go test -v -bench '.*' \
    >   github.com/timtadh/data-structures/tree
    >   github.com/timtadh/data-structures/hashtable

    BenchmarkGoMap    100000             28973 ns/op
    BenchmarkHash      20000             76326 ns/op
    BenchmarkMLHash    50000             73516 ns/op
    BenchmarkLHash       500           3307005 ns/op
    BenchmarkAvlTree   10000            144663 ns/op


The performance of the in memory linear hash (MLHash) is slightly improved since
the [blog post](
http://hackthology.com/an-in-memory-go-implementation-of-linear-hashing.html) do
to the usage of an AVL Tree `tree/avltree.go` instead of an unbalanced binary
search tree.

