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

The library also provides generic
[types](https://godoc.org/github.com/timtadh/data-structures/types) to allow the
user to swap out various data structures transparently. The interfaces provide
operation for adding, removing, retrieving objects from collections as well as
iterating over the collection using functional iterators.

The tree sub-package provides a variety of generic tree traversals. The tree
traversals and other iterators in the package use a functional iteration
technique [detailed on my blog](
http://hackthology.com/functional-iteration-in-go.html).

I hope you find my library useful. If you are using it drop me a line I would
love to hear about it.

# Current Collection

[![GoDoc](https://godoc.org/github.com/timtadh/data-structures?status.svg)](https://godoc.org/github.com/timtadh/data-structures)

## Lists

### Doubly Linked List [`linked.LinkedList`](https://godoc.org/github.com/timtadh/data-structures/linked#LinkedList)

A simple an extensible doubly linked list. It is
[Equatable](https://godoc.org/github.com/timtadh/data-structures/types#Equatable)
[Sortable](https://godoc.org/github.com/timtadh/data-structures/types#Sortable),
and [Hashable](https://godoc.org/github.com/timtadh/data-structures/types#Hashable)
as are the [Node](https://godoc.org/github.com/timtadh/data-structures/linked#Node)s.

### Array List [`list.List`](https://godoc.org/github.com/timtadh/data-structures/list#List)

Similar to a Java ArrayList or a Python or Ruby "list". There is a version
(called Sortable) which integrates with the `"sort"` package from the standard
library.

### Sorted Array List [`list.Sorted`](https://godoc.org/github.com/timtadh/data-structures/list#Sorted)

Keeps the ArrayList in sorted order for you.

### Sorted Set [`set.SortedSet`](https://godoc.org/github.com/timtadh/data-structures/set#SortedSet)

Built on top of `*list.Sorted`, it provides basic set operations. With
`set.SortedSet` you don't have to write code re-implementing sets with the
`map[type]` datatype. Supports: intersection, union, set difference and overlap
tests.

### Map Set [`set.MapSet`](https://godoc.org/github.com/timtadh/data-structures/set#MapSet)

Construct a
[`types.Map`](https://godoc.org/github.com/timtadh/data-structures/types#Map)
from any [`types.Set`](https://godoc.org/github.com/timtadh/data-structures/types#Set).

### Set Map [`set.SetMap`](https://godoc.org/github.com/timtadh/data-structures/set#SetMap)

Construct a set from any
[`types.Map`](https://godoc.org/github.com/timtadh/data-structures/types#Map).

### Unique Deque [`linked.UniqueDeque`](https://godoc.org/github.com/timtadh/data-structures/linked#UniqueDeque)

A double ended queue that only allows unique items inside. Constructed from a
doubly linked list and a linear hash table.

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


## Heaps and Priority Queues

### Binary Heap [`heap/Heap`](https://godoc.org/github.com/timtadh/data-structures/heap#Heap)

This is a binary heap for usage as a priority queue. The priorities are given to
items in the queue on insertion and cannot be changed after insertion. It can be
used as both a min heap and a max heap.

### Unique Priority Queue [`heap/UniquePQ`](https://godoc.org/github.com/timtadh/data-structures/heap#UniquePQ)

A priority queue which only allows unique entries.

## Trees

### AVL Tree [`tree/avl.AvlTree`](https://godoc.org/github.com/timtadh/data-structures/tree/avl#AvlTree)

An [AVL Tree](https://en.wikipedia.org/wiki/AVL_tree) is a height balanced
binary search tree. Insertion and retrieval are both O(log(n)) where n is the
number items in the tree.

### Immutable AVL Tree [`tree/avl.ImmutableAvlTree`](https://godoc.org/github.com/timtadh/data-structures/tree/avl#ImmutableAvlTree)

This version of the classic is immutable and should be thread safe due to
immutability. However, there is a performance hit:

    BenchmarkAvlTree           10000            166657 ns/op
    BenchmarkImmutableAvlTree   5000            333709 ns/op

### Ternary Search Trie [`trie.TST`](https://godoc.org/github.com/timtadh/data-structures/trie#TST)

A [ternary search trie](
http://hackthology.com/ternary-search-tries-for-fast-flexible-string-search-part-1.html)
is a symbol table specialized to byte strings.  Ternary Search Tries (TSTs)
are a particularly fast version of the more common R-Way Trie. They utilize less
memory allowing them to store more data while still retaining all of the
flexibility of the R-Way Trie. TSTs can be used to build a suffix tree for full
text string indexing by storing every suffix of each string in addition to the
string. However, even without storing all of the suffixes it is still a great
structure for flexible prefix searches. For instance, TSTs can be used to
implement extremely fast auto-complete functionality.

### B+Tree [`tree/bptree.BpTree`](https://godoc.org/github.com/timtadh/data-structures/tree/bptree)

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

### Separate Chaining Hash Table [`hashtable.Hash`](https://godoc.org/github.com/timtadh/data-structures/hashtable#Hash)

See `hashtable/hashtable.go`. An implementation of the classic hash table with
separate chaining to handle collisions.

### Linear Hash Table with AVL Tree Buckets [`hashtable.LinearHash`](https://godoc.org/github.com/timtadh/data-structures/hashtable#LinearHash)

See `hashtables/linhash.go`. An implementation of [Linear
Hashing](http://hackthology.com/linear-hashing.html), a technique usually used
for secondary storage hash tables. Often employed by databases and file systems
for hash indices. This version is mostly instructional see the [accompanying
blog post](
http://hackthology.com/an-in-memory-go-implementation-of-linear-hashing.html).
If you want a disk backed version check out my
[file-structures](https://github.com/timtadh/file-structures) repository. See
the `linhash` directory.

## Exceptions, Errors, and Testing

### Errors [`errors`](https://godoc.org/github.com/timtadh/data-structures/errors)

By default, most errors in Go programs to not track where they were created.
Many programmers quickly discover they need to have stack traces associated with
their errors. This is a light weight package which adds stack traces to errors.
It also provides a very very simple logging function that reports where in your
code your printed out the log. This is not a full featured logging solution but
rather a replacement for using fmt.Printf when debugging.

### Test Support [`test`](https://godoc.org/github.com/timtadh/data-structures/test)

The test package provides two minimal assertions and a way to get random strings
and data. It also seeds the math/rand number generator. I consider this to the
bare minimum of what is often needed when testing go code particularly
data-structures. Since this package seeks to be entirely self contained with no
dependencies no external testing package is used. This package is slowly being
improved to encompass more common functionality between the different tests.

### Exceptions as a Library [`exc`](https://github.com/timtadh/data-structures/tree/master/exc)

- [![GoDoc](https://godoc.org/github.com/timtadh/data-structures/exc?status.svg)](https://godoc.org/github.com/timtadh/data-structures/exc)
- [Explanation of Implementation](http://hackthology.com/exceptions-for-go-as-a-library.html)
- [Explanation of Inheritance](http://hackthology.com/object-oriented-inheritance-in-go.html)

The [`exc`](https://github.com/timtadh/data-structures/tree/master/exc) package
provides support for exceptions. They work very similarly to the way unchecked
exceptions work in Java. They are built on top of the built-in `panic` and
`recover` functions. See the README in the package for more information or
checkout the documentation. They should play nice with the usual way of handling
errors in Go and provide an easy way to create public APIs which return errors
rather than throwing these non-standard exceptions.


## Benchmarks

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

