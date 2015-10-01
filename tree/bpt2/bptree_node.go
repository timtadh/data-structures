package bpt2

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/types"
)

const (
	__internal__ uint16 = 1 << iota
	__leaf__
	__duplicates__
)

type Meta struct {
	flags uint16
	size  uint16
	count uint32
}

func (m *Meta) Size() int {
	return int(m.size)
}

func (m *Meta) Internal() bool {
	return m.flags & __internal__ != 0
}

func (m *Meta) Leaf() bool {
	return m.flags & __leaf__ != 0
}

func (m *Meta) Duplicates() bool {
	return m.flags & __duplicates__ != 0
}

type BpNode struct {
	meta Meta
	next *BpNode
	prev *BpNode
	data *list.Sorted
}

func NewInternal(meta Meta) *BpNode {
	meta.flags = (meta.flags & ^__leaf__) | __internal__
	meta.count = 0
	return &BpNode{
		meta: meta,
		data: list.NewFixedSorted(meta.size, false),
	}
}

func NewLeaf(meta Meta) *BpNode {
	meta.flags = (meta.flags & ^__internal__) | __leaf__
	meta.count = 0
	return &BpNode{
		meta: meta,
		data: list.NewFixedSorted(meta.size, false),
	}
}

func (n *BpNode) Full() bool {
	return n.data.Full()
}

func (n *BpNode) entry(i int) *types.Entry {
	e, err := n.data.Get(i)
	if err != nil {
		return errors.Panic(err)
	}
	return e.(*types.MapEntry)
}

func (n *BpNode) key(i int) (types.Hashable, error) {
	return n.entry(i).Key
}

func (n *BpNode) val(i int) (interface{}, error) {
	return n.entry(i).Value
}

func (n *BpNode) Pure() bool {
	if len(n.keys) == 0 {
		return true
	}
	k0 := n.key(0)
	for i, next := n.data.Items()(); next != nil; i, next = next() {
		k := i.(*types.MapEntry).Key
		if !k0.Equals(k) {
			return false
		}
	}
	return true
}

func (n *BpNode) Height() int {
	if !n.Internal() {
		return 1
	} else if n.data.Size() == 0 {
		panic(errors.BpTreeError("Internal node has no pointers but asked for height"))
	}
	return n.data.Get(0).(*types.MapEntry).Value.(*BpNode).Height() + 1
}

func (n *BpNode) Add(key, value types.Hashable) (delta int, root *BpNode, err error) {
	delta, a, b, err := n.insert(key, value)
	if err != nil {
		return 0, nil, err
	} else if b == nil {
		return delta, a, nil
	}
	root := NewInternal(n.meta)
	err = root.data.Add(&types.MapEntry{a.key(0), a})
	if err != nil {
		return 0, nil, err
	}
	err = root.data.Add(&types.MapEntry{b.key(0), b})
	if err != nil {
		return 0, nil, err
	}
	root.meta.size = a.meta.size + b.meta.size
	return delta, root, nil
}

func (n *BpNode) insert(key, value types.Hashable) (delta int, a, b *BpNode, err error) {
	if n.meta.Internal() {
		return n.internalInsert(key, value)
	} else if n.meta.Leaf() {
		return n.leafInsert(key, value)
	} else {
		return 0, nil, nil, errors.Errorf("Unlabeled *BpNode")
	}
}

func (n *BpNode) internalInsert(key, value types.Hashable) (delta int, a, b *BpNode, err error) {
	i, has, err := n.data.Find(key)
	if err != nil {
		return 0, nil, nil, err
	} else if !has && i > 0 {
		// if it doesn't have it and the index > 0 then
		//     i is the next block so subtract 1
		i--
	}
	entry := n.entry(i)
	kid := entry.(*BpNode)
	delta, p, q, err := kid.insert(key, value)
	if err != nil {
		return 0, nil, nil, err
	}
	// this mutation is safe as long as the invariants of the tree are not violated!
	entry.Key = p.key(0)
	entry.Value = p
	if q == nil {
		return delta, n, nil, nil
	}
	if n.Full() {
		return n.internalSplit(q.key(0), q)
	} else {
		err = n.data.Add(&types.MapEntry{q.key(0), q})
		if err != nil {
			return 0, nil, nil, err
		}
		return delta, n, nil, nil
	}
}

func (n *BpNode) leafInsert(key, value types.Hashable) (delta int, a, b *BpNode, err error) {
}

func (a *BpNode) internalSplit(key types.Hashable, kid *BpNode) (delta int, a, b *BpNode, err error) {
	b := NewInternal(a.meta)
	err := a.balance(b)
	if err != nil {
		return 0, nil, nil, err
	}
	if key.Less(b.key(0)) {
		err = a.data.Add(&types.MapEntry{key, kid})
	} else {
		err = b.data.Add(&types.MapEntry{key, kid})
	}
	if err != nil {
		return 0, nil, nil, err
	}
	return delta, a, b, nil
}
