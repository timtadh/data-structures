# Go Data Structures

by Tim Henderson (tim.tadh@gmail.com)

Copyright 2013, Licensed under the GPL version 2. Please reach out to me
directly if you require another licensing option. I am willing to work with you.

## Purpose

To collect many important data structures for usage in go programs. Golang's
standard library lacks many useful and important structures. This library
attempts to fill the gap. I have implemented data-structure's as I have needed
them. If there is a missing structure or even just a missing (or incorrect)
method open an issue, send a pull request, or send an email patch.

The library also provides generic types to allow the user to swap out various
data structures transparently. The interfaces provide operation for adding,
removing, retrieving objects from collections as well as iterating over the
collection using functional iterators.

Finally, the tree sub-package provides a variety of generic tree traversals. The
tree traversals and other iterators in the package use a functional iteration
technique [detailed on my blog](
http://hackthology.com/functional-iteration-in-go.html).

I hope you find my library useful. If you are using it drop me a line I would
love to hear about it.

# Current Collection

[godoc](https://godoc.org/github.com/timtadh/data-structures)

## Lists

### Array List `list.List`

Similar to a Java ArrayList or a Python or Ruby "list". There is a version
(called Sortable) which integrates with the `"sort"` package from the standard
library.

### Sorted Array List `list.Sorted`

Keeps the ArrayList in sorted order for you.

### Sorted Set `set.SortedSet`

Built on top of `*list.Sorted`, it provides basic set operations. With
`set.SortedSet` you don't have to write code re-implementing sets with the
`map[type]` datatype. Supports: intersection, union, set difference and overlap
tests.

### Fixed Size Lists

Both `list.List` and `list.Sorted` have alternative constructors which make them
fixed size. This prevents them from growing beyond a certain size bound and is
useful for implementing other data structures on top of them.

### Serialization

`list.List`, `list.Sorted`, and `set.SortedSet` all can be serialized provided
their contents can be serialized. They are therefore suitable for being sent
over the wire. See this
[example](https://github.com/timtadh/data-structures/blob/master/set/example_serialize_test.go)
for how to use the serialization.

## Trees

### Classic AVL Tree `tree/avl.AvlTree`

An AVL tree is a height balanced binary search tree. It is commonly taught in
algorithms courses.

### Immutable AVL Tree `tree/avl.ImmutableAvlTree`

This version of the classic is immutable and should be thread safe due to
immutability. However, there is a performance hit:

    BenchmarkAvlTree           10000            166657 ns/op
    BenchmarkImmutableAvlTree   5000            333709 ns/op

### Ternary Search Trie `trie.TST`

A [ternary search trie](
http://hackthology.com/ternary-search-tries-for-fast-flexible-string-search-part-1.html)
is a symbol table specialized to byte strings. It can be used to build a
suffix tree for full text string indexing. However, even without a suffix tree
it is still a great structure for flexible prefix searches.

### B+Tree (with and without support for duplicate keys) `tree/bptree.BpTree`

A
[B+Tree](http://hackthology.com/lessons-learned-while-implementing-a-btree.html)
is a general symbol table usually used for database indices. This implementation
is not currently thread safe. However, unlike many B+Trees it fully supports
duplicate keys making it suitable for use as a Multi-Map. There is also a
variant which has unique keys, `bptree.BpMap`. B+Trees are storted and efficient
to iterate over making them ideal choices for storing a large amount of data
in sorted order. For storing a **very** large amount of data please utilize the
fs2 version, [fs2/bptree](https://github.com/timtadh/fs2#b-tree). fs2 utilizes
memory mapped files in order to allow you to store more data than your computer
has RAM.

## Hash Tables

### Classic Separate Chaining Hash Table `hashtable.Hash`

See `hashtable/hashtable.go`. An implementation of the classic hash table with
separate chaining to handle collisions.

### Linear Hash Table with AVL Tree Buckets `hashtable.LinearHash`

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

**Note**: these benchmarsk are fairly old and probably not easy to understand.
Look at the relative difference not the absolute numbers as they are misleading.
Each benchmark does many operations per "test" which makes it difficult to
compare these numbers to numbers found elsewhere.

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

# Related Projects

- [fs2](https://github.com/timtadh/fs2) Memory mapped datastructures. A B+Tree,
  a list, and a platform for implementing more.

- [file-structures](https://github.com/timtadh/file-structures) The previous
  version of fs2 of disk based file-structures. Also includes a linear virtual
  hashing implementation.

