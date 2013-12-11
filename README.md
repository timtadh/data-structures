# Go Data Structures

by Tim Henderson (tim.tadh@gmail.com)

Copyright 2013, Licensed under the GPL version 2.

## Purpose

To collect many important data structures for usage in go programs. Some of
these data structures may have implementations elsewhere but are collected here
for completeness and instructional usage.

# Current Collection

## Hash Tables

### Classic Separate Chaining Hash Table

See `hashtable/hashtable.go`. This implements the classic hash table with
separate chaining to handle collisions.

### Linear Hash Table with BST Buckets

See `hashtables/linhash.go`. This implement [Linear
Hashing](http://hackthology.com/linear-hashing.html), a technique usually used
for secondary storage hash tables. Often employed by databases and file systems
for hash indices. This version is mostly instructional see the
[accompanying blog post](
http://hackthology.com/an-in-memory-go-implementation-of-linear-hashing.html).
If you want the "real" disk backed version you want to check my
[file-structures](https://github.com/timtadh/file-structures) repository. See
the `linhash` directory.

