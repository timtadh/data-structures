# Go Data Structures

by Tim Henderson (tim.tadh@gmail.com)

Copyright 2013, Licensed under the GPL version 2.

## Purpose

To collect many important data structures for usage in go programs. Some of
these data structures may have implementations elsewhere but are collected here
for completeness and instructional usage.

The library also provided generic types to allow the user to swap out various
data structures transparently. The interfaces provide operation for adding,
removing, retrieving objects from collections as well as iterating over the
collection using functional iterators.

Finally, the tree sub-package provides a variety of generic tree traversals. The
tree traversals and other iterators in the package use a functional iteration
technique [detailed on my blog](
http://hackthology.com/functional-iteration-in-go.html).



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

### Ternary Search Trie

A [ternary search trie](
http://hackthology.com/ternary-search-tries-for-fast-flexible-string-search-part-1.html)
is a symbol table specialized to byte strings. It can be used to build a
suffix tree for full text string indexing. However, even without a suffix tree
it is still a great structure for flexible prefix searches.

### B+Tree (with and without support for duplicate keys)

A
[B+Tree](http://hackthology.com/lessons-learned-while-implementing-a-btree.html)
is a general symbol table usually used for database indices. This implementation
is not currently thread safe. It uses the structure detailed in the link and was
ported from my file-structures repository.

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

Benchmarks Put + Remove

    $ go test -v -bench '.*' \
    >   github.com/timtadh/data-structures/hashtable
    >   github.com/timtadh/data-structures/tree/...
    >   github.com/timtadh/data-structures/trie

    BenchmarkGoMap             50000             30051 ns/op
    BenchmarkMLHash            20000             78840 ns/op
    BenchmarkHash              20000             81012 ns/op
    BenchmarkTST               10000            149985 ns/op
    BenchmarkBpTree            10000            185134 ns/op
    BenchmarkAvlTree           10000            193069 ns/op
    BenchmarkImmutableAvlTree   5000            367602 ns/op
    BenchmarkLHash              1000           2743693 ns/op

Benchmarks Put

    BenchmarkGoMap            100000             22036 ns/op
    BenchmarkMLHash            50000             52104 ns/op
    BenchmarkHash              50000             53426 ns/op
    BenchmarkTST               50000             69852 ns/op
    BenchmarkBpTree            20000             76124 ns/op
    BenchmarkAvlTree           10000            142104 ns/op
    BenchmarkImmutableAvlTree  10000            302196 ns/op
    BenchmarkLHash              1000           1739710 ns/op

The performance of the in memory linear hash (MLHash) is slightly improved since
the [blog post](
http://hackthology.com/an-in-memory-go-implementation-of-linear-hashing.html) do
to the usage of an AVL Tree `tree/avltree.go` instead of an unbalanced binary
search tree.

