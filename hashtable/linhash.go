package hashtable

import (
    . "github.com/timtadh/data-structures/types"
    "github.com/timtadh/data-structures/tree"
)

type HashTable interface {
    Get(key Hashable) (value interface{}, err error)
    Put(key Hashable, value interface{}) (err error)
    Has(key Hashable) (has bool)
    Remove(key Hashable) (value interface{}, err error)
    Size() int
}

const (
    UTILIZATION = .75
    RECORDS_PER_BLOCK = 8
)

type bst struct {
    hash int
    key Hashable
    value interface{}
    left *bst
    right *bst
}

type linearhash struct {
    table []*tree.AvlTree
    n uint
    r uint
    i uint
}

func NewLinearHash() *linearhash {
    N := uint(32)
    I := uint(5)
    return &linearhash{
        table: make([]*tree.AvlTree, N),
        n: N,
        r: 0,
        i: I,
    }
}

func (self *linearhash) bucket(key Hashable) uint {
    m := uint(key.Hash() & ((1<<self.i)-1))
    if m < self.n {
        return m
    } else {
        return m ^ (1<<(self.i-1))
    }
}


func (self *linearhash) Size() int {
    return int(self.r)
}

func (self *linearhash) Put(key Hashable, value interface{}) (err error) {
    var updated bool
    bkt_idx := self.bucket(key)
    self.table[bkt_idx], updated = self.table[bkt_idx].Put(key, value)
    if !updated {
        self.r += 1
    }
    if float64(self.r) > UTILIZATION * float64(self.n) * float64(RECORDS_PER_BLOCK) {
        return self.split()
    }
    return nil
}

func (self *linearhash) Get(key Hashable) (value interface{}, err error) {
    bkt_idx := self.bucket(key)
    return self.table[bkt_idx].Get(key)
}

func (self *linearhash) Has(key Hashable) (bool) {
    bkt_idx := self.bucket(key)
    return self.table[bkt_idx].Has(key)
}

func (self *linearhash) Remove(key Hashable) (value interface{}, err error) {
    bkt_idx := self.bucket(key)
    self.table[bkt_idx], value, err = self.table[bkt_idx].Remove(key)
    if err == nil {
        self.r -= 1
    }
    return
}

func (self *linearhash) split() (err error) {
    bkt_idx := self.n % (1 << (self.i - 1))
    old_bkt := self.table[bkt_idx]
    var bkt_a, bkt_b *tree.AvlTree
    self.n += 1
    if self.n > (1 << self.i) {
        self.i += 1
    }
    for key, value, next := old_bkt.Iterate()(); next != nil; key, value, next = next() {
        if self.bucket(key.(Hashable)) == bkt_idx {
            bkt_a, _ = bkt_a.Put(key.(Sortable), value)
        } else {
            bkt_b, _ = bkt_b.Put(key.(Sortable), value)
        }
    }
    self.table[bkt_idx] = bkt_a
    self.table = append(self.table, bkt_b)
    return nil
}

