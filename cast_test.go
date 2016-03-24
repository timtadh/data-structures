package data_structures

import "testing"
import (
	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/tree/avl"
	"github.com/timtadh/data-structures/tree/bptree"
	"github.com/timtadh/data-structures/trie"
	"github.com/timtadh/data-structures/types"
)

func TestAvlTreeCast(t *testing.T) {
	tree := avl.NewAvlTree()
	_ = types.TreeMap(tree)
}

func TestImmutableAvlTreeCast(t *testing.T) {
	tree := avl.NewImmutableAvlTree()
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
	tst := new(trie.TST)
	_ = types.MapIterable(tst)
}

func TestBpTreeCast(t *testing.T) {
	bpt := bptree.NewBpTree(17)
	_ = types.MapIterable(bpt)
	_ = types.MultiMapOperable(bpt)
	_ = types.MultiMap(bpt)
}

func TestListCast(t *testing.T) {
	s := list.New(17)
	_ = types.Hashable(s)
	_ = types.List(s)
}

func TestSortedCast(t *testing.T) {
	s := list.NewSorted(17, false)
	_ = types.Hashable(s)
	_ = types.OrderedList(s)
}

func TestSetCast(t *testing.T) {
	s := set.NewSortedSet(17)
	_ = types.Set(s)
	_ = types.Hashable(s)
}
