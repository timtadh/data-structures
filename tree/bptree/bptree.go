package bptree

import (
  "github.com/timtadh/data-structures/types"
)

type BpTree struct {
    root *BpNode
    size int
}

type loc_iterator func()(i int, leaf *BpNode, li loc_iterator)

func NewBpTree(node_size int) *BpTree {
    return &BpTree{
        root: NewLeaf(node_size),
        size: 0,
    }
}

func (self *BpTree) Size() int {
    return self.size
}

func (self *BpTree) Has(key types.Hashable) bool {
    j, l := self.root.get_start(key)
    return l.keys[j].Equals(key)
}

func (self *BpTree) Count(key types.Hashable) int {
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
    li := self.forward(key, key)
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
        li = self.forward(from, to)
    } else {
        li = self.backward(from, to)
    }
    kvi = func() (key types.Equatable, value interface{}, next types.KVIterator) {
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

func (self *BpTree) RemoveWhere(key types.Hashable, where types.WhereFunc) (value interface{}, err error) {
    panic("unimplemented")
}

func (self *BpTree) Keys() (ki types.KIterator) {
    return types.MakeKeysIterator(self)
}

func (self *BpTree) Values() (vi types.Iterator) {
    return types.MakeValuesIterator(self)
}

func (self *BpTree) Iterate() (kvi types.KVIterator) {
    li := self.all()
    kvi = func() (key types.Equatable, value interface{}, next types.KVIterator) {
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

func (self *BpTree) all() (li loc_iterator) {
    j := -1
    l := self.root.left_most_leaf()
    end := false
    j, l, end = next_location(j, l)
    li = func() (i int, leaf *BpNode, next loc_iterator) {
        if end {
            return -1, nil, nil
        }
        i = j
        leaf = l
        j, l, end = next_location(j, l)
        return i, leaf, li
    }
    return li
}

func (self *BpTree) forward(from, to types.Sortable) (li loc_iterator) {
    j, l := self.root.get_start(from)
    end := false
    li = func() (i int, leaf *BpNode, next loc_iterator) {
        if end || to.Less(l.keys[j]) {
            return -1, nil, nil
        }
        i = j
        leaf = l
        j, l, end = next_location(i, l)
        return i, leaf, li
    }
    return li
}

func (self *BpTree) backward(from, to types.Sortable) (li loc_iterator) {
    j, l := self.root.get_end(from)
    end := false
    li = func() (i int, leaf *BpNode, next loc_iterator) {
        if end || l.keys[j].Less(to) {
            return -1, nil, nil
        }
        i = j
        leaf = l
        j, l, end = prev_location(i, l)
        return i, leaf, li
    }
    return li
}

