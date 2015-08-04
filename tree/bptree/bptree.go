package bptree

import (
	"github.com/timtadh/data-structures/types"
)

/* A BpTree is a B+Tree with support for duplicate keys. This makes it behave as
 * a MultiMap. Additionally you can use the Range operator to select k/v in a
 * range. If from > to it will iterate backwards.
 */
type BpTree struct {
	root *BpNode
	size int
}

type loc_iterator func() (i int, leaf *BpNode, li loc_iterator)

func NewBpTree(node_size int) *BpTree {
	return &BpTree{
		root: NewLeaf(node_size, false),
		size: 0,
	}
}

func (self *BpTree) Size() int {
	return self.size
}

func (self *BpTree) Has(key types.Hashable) bool {
	if len(self.root.keys) == 0 {
		return false
	}
	j, l := self.root.get_start(key)
	return l.keys[j].Equals(key)
}

func (self *BpTree) Count(key types.Hashable) int {
	if len(self.root.keys) == 0 {
		return 0
	}
	j, l := self.root.get_start(key)
	count := 0
	end := false
	for !end && l.keys[j].Equals(key) {
		count++
		j, l, end = next_location(j, l)
	}
	return count
}

func (self *BpTree) Add(key types.Hashable, value interface{}) (err error) {
	new_root, err := self.root.put(key, value)
	if err != nil {
		return err
	}
	self.root = new_root
	self.size += 1
	return nil
}

func (self *BpTree) Replace(key types.Hashable, where types.WhereFunc, value interface{}) (err error) {
	li := self.root.forward(key, key)
	for i, leaf, next := li(); next != nil; i, leaf, next = next() {
		if where(leaf.values[i]) {
			leaf.values[i] = value
		}
	}
	return nil
}

func (self *BpTree) Find(key types.Hashable) (kvi types.KVIterator) {
	return self.Range(key, key)
}

func (self *BpTree) Range(from, to types.Hashable) (kvi types.KVIterator) {
	var li loc_iterator
	if !to.Less(from) {
		li = self.root.forward(from, to)
	} else {
		li = self.root.backward(from, to)
	}
	kvi = func() (key types.Hashable, value interface{}, next types.KVIterator) {
		var i int
		var leaf *BpNode
		i, leaf, li = li()
		if li == nil {
			return nil, nil, nil
		}
		return leaf.keys[i], leaf.values[i], kvi
	}
	return kvi
}

func (self *BpTree) RemoveWhere(key types.Hashable, where types.WhereFunc) (err error) {
	ns := self.root.NodeSize()
	new_root, err := self.root.remove(key, where)
	if err != nil {
		return err
	}
	if new_root == nil {
		self.root = NewLeaf(ns, false)
	} else {
		self.root = new_root
	}
	self.size -= 1
	return nil
}

func (self *BpTree) Keys() (ki types.KIterator) {
	li := self.root.all()
	var prev types.Equatable
	ki = func() (key types.Hashable, next types.KIterator) {
		var i int
		var leaf *BpNode
		i, leaf, li = li()
		if li == nil {
			return nil, nil
		}
		if leaf.keys[i].Equals(prev) {
			return ki()
		}
		prev = leaf.keys[i]
		return leaf.keys[i], ki
	}
	return ki
}

func (self *BpTree) Values() (vi types.Iterator) {
	return types.MakeValuesIterator(self)
}

func (self *BpTree) Items() (vi types.KIterator) {
	return types.MakeItemsIterator(self)
}

func (self *BpTree) Iterate() (kvi types.KVIterator) {
	li := self.root.all()
	kvi = func() (key types.Hashable, value interface{}, next types.KVIterator) {
		var i int
		var leaf *BpNode
		i, leaf, li = li()
		if li == nil {
			return nil, nil, nil
		}
		return leaf.keys[i], leaf.values[i], kvi
	}
	return kvi
}

func (self *BpTree) Backward() (kvi types.KVIterator) {
	li := self.root.all_backward()
	kvi = func() (key types.Hashable, value interface{}, next types.KVIterator) {
		var i int
		var leaf *BpNode
		i, leaf, li = li()
		if li == nil {
			return nil, nil, nil
		}
		return leaf.keys[i], leaf.values[i], kvi
	}
	return kvi
}
