package data_structures

import "testing"
import (
    "github.com/timtadh/data-structures/types"
    "github.com/timtadh/data-structures/tree"
    "github.com/timtadh/data-structures/trie"
    "github.com/timtadh/data-structures/hashtable"
)

func TestAvlTreeCast(t *testing.T) {
    tree := tree.NewAvlTree()
    _ = types.TreeMap(tree)
}

func TestImmutableAvlTreeCast(t *testing.T) {
    tree := tree.NewImmutableAvlTree()
    _ = types.TreeMap(tree)
}

func TestHashtableCast(t *testing.T) {
    hash := hashtable.NewHashTable(16)
    _ = types.Sized(hash)
    _ = types.MapIterable(hash)
    _ = types.MapOperable(hash)
    _ = types.Map(hash)
}

func TestLinearHashtableCast(t *testing.T) {
    hash := hashtable.NewLinearHash()
    _ = types.Sized(hash)
    _ = types.MapIterable(hash)
    _ = types.MapOperable(hash)
    _ = types.Map(hash)
}

func TestTSTCast(t *testing.T) {
    hash := new(trie.TST)
    _ = types.MapIterable(hash)
}

