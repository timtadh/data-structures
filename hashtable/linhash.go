package hashtable

import (
	"fmt"
	"strings"
)

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/types"
)

const (
	UTILIZATION       = .75
	RECORDS_PER_BLOCK = 16
)

type bst struct {
	hash  int
	key   types.Hashable
	value interface{}
	left  *bst
	right *bst
}

type LinearHash struct {
	table []*list.List
	n     uint
	r     uint
	i     uint
}

func NewLinearHash() *LinearHash {
	N := uint(32)
	I := uint(5)
	h := &LinearHash{
		table: make([]*list.List, N),
		n:     N,
		r:     0,
		i:     I,
	}
	for i := range h.table {
		h.table[i] = list.New(RECORDS_PER_BLOCK)
	}
	return h
}

func (self *LinearHash) bucket(key types.Hashable) uint {
	m := uint(key.Hash() & ((1 << self.i) - 1))
	if m < self.n {
		return m
	} else {
		return m ^ (1 << (self.i - 1))
	}
}

func (self *LinearHash) Size() int {
	return int(self.r)
}

func (self *LinearHash) Put(key types.Hashable, value interface{}) (err error) {
	e := &types.MapEntry{key, value}
	bkt_idx := self.bucket(key)
	bkt := self.table[bkt_idx]
	i, has, err := list.Find(bkt, e)
	if err != nil {
		return err
	} else if !has {
		bkt.Insert(i, e)
		self.r++
	} else {
		if item, err := bkt.Get(i); err != nil {
			return err
		} else {
			item.(*types.MapEntry).Value = value
		}
	}
	if float64(self.r) > UTILIZATION*float64(self.n)*float64(RECORDS_PER_BLOCK) {
		return self.split()
	}
	return nil
}

func (self *LinearHash) Get(key types.Hashable) (value interface{}, err error) {
	bkt_idx := self.bucket(key)
	bkt := self.table[bkt_idx]
	i, has, err := list.Find(bkt, key)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.Errorf("Key not found! '%v'", key)
	}
	item, err := bkt.Get(i)
	if err != nil {
		return nil, err
	}
	return item.(*types.MapEntry).Value, nil
}

func (self *LinearHash) Has(key types.Hashable) bool {
	bkt_idx := self.bucket(key)
	bkt := self.table[bkt_idx]
	_, has, _ := list.Find(bkt, key)
	return has
}

func (self *LinearHash) Remove(key types.Hashable) (value interface{}, err error) {
	bkt_idx := self.bucket(key)
	bkt := self.table[bkt_idx]
	i, has, err := list.Find(bkt, key)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.Errorf("Key not found! '%v'", key)
	}
	item, err := bkt.Get(i)
	if err != nil {
		return nil, err
	}
	err = bkt.Remove(i)
	if err != nil {
		return nil, err
	}
	if !item.(*types.MapEntry).Key.Equals(key) {
		return nil, errors.Errorf("assert fail: %v != %v", key, item)
	}
	self.r--
	return item.(*types.MapEntry).Value, nil
}

func (self *LinearHash) split() (err error) {
	bkt_idx := self.n % (1 << (self.i - 1))
	old_bkt := self.table[bkt_idx]
	bkt_a := list.New(RECORDS_PER_BLOCK)
	bkt_b := list.New(RECORDS_PER_BLOCK)
	self.n += 1
	if self.n > (1 << self.i) {
		self.i += 1
	}
	for item, next := old_bkt.Items()(); next != nil; item, next = next() {
		if self.bucket(item) == bkt_idx {
			bkt_a.Append(item)
		} else {
			bkt_b.Append(item)
		}
	}
	self.table[bkt_idx] = bkt_a
	self.table = append(self.table, bkt_b)
	return nil
}

func (self *LinearHash) Iterate() (kvi types.KVIterator) {
	table := self.table
	i := 0
	iter := table[i].Items()
	kvi = func() (key types.Hashable, val interface{}, _ types.KVIterator) {
		var item types.Hashable
		item, iter = iter()
		for iter == nil {
			i++
			if i >= len(table) {
				return nil, nil, nil
			}
			item, iter = table[i].Items()()
		}
		e := item.(*types.MapEntry)
		return e.Key, e.Value, kvi
	}
	return kvi
}

func (self *LinearHash) Items() (vi types.KIterator) {
	return types.MakeItemsIterator(self)
}

func (self *LinearHash) Keys() types.KIterator {
	return types.MakeKeysIterator(self)
}

func (self *LinearHash) Values() types.Iterator {
	return types.MakeValuesIterator(self)
}

func (self *LinearHash) String() string {
	if self.Size() <= 0 {
		return "{}"
	}
	items := make([]string, 0, self.Size())
	for item, next := self.Items()(); next != nil; item, next = next() {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return "{" + strings.Join(items, ", ") + "}"
}

