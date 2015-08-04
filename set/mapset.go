package set

import (
	"github.com/timtadh/data-structures/types"
)


type MapSet struct {
	Set types.Set
}

func NewMapSet(set types.Set) *MapSet {
	return &MapSet{set}
}

/*
	Equals(b Equatable) bool

	Keys() KIterator
	Values() Iterator
	Iterate() KVIterator
	Items() KIterator

	Has(key Hashable) bool

	Add(item Hashable) (err error)
	Put(key Hashable, value interface{}) (err error)
	Get(key Hashable) (value interface{}, err error)
	Remove(key Hashable) (value interface{}, err error)
	Delete(item Hashable) (err error)
	Extend(items KIterator) (err error)


	Union(Set) (Set, error)
	Intersect(Set) (Set, error)
	Subtract(Set) (Set, error)
	Subset(Set) bool
	Superset(Set) bool
	ProperSubset(Set) bool
	ProperSuperset(Set) bool
*/
