package hashtable

import "fmt"

type Hashable interface {
    Equals(b Hashable) bool
    Less(b Hashable) bool
    Hash() int
}

type HashTable interface {
    Get(key Hashable) (value interface{}, err error)
    Put(key Hashable, value interface{}) (err error)
    Has(key Hashable) (has bool)
    Remove(key Hashable) (value interface{}, err error)
    Size() int
}

var Errors map[string]error = map[string]error{
    "not-found":fmt.Errorf("Key was not in hash table"),
    "list-not-found":fmt.Errorf("Key was not in bucket when expected"),
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
    table []*bst
    n uint
    r uint
    i uint
}

func NewLinearHash() *linearhash {
    N := uint(32)
    I := uint(5)
    return &linearhash{
        table: make([]*bst, N),
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
    var bkt_a, bkt_b *bst
    self.n += 1
    if self.n > (1 << self.i) {
        self.i += 1
    }
    for key, value, next := old_bkt.Iterate()(); next != nil; key, value, next = next() {
        if self.bucket(key) == bkt_idx {
            bkt_a, _ = bkt_a.Put(key, value)
        } else {
            bkt_b, _ = bkt_b.Put(key, value)
        }
    }
    self.table[bkt_idx] = bkt_a
    self.table = append(self.table, bkt_b)
    return nil
}

type BSTIterator func()(key Hashable, value interface{}, next BSTIterator)
func (self *bst) Iterate() BSTIterator {
    pop := func(stack []*bst) ([]*bst, *bst) {
        if len(stack) <= 0 {
            return stack, nil
        } else {
            return stack[0:len(stack)-1], stack[len(stack)-1]
        }
    }
    procnode := func(stack []*bst, node *bst) []*bst {
        if node == nil {
            return stack
        }
        if node.right != nil {
            stack = append(stack, node.right)
        }
        if node.left != nil {
            stack = append(stack, node.left)
        }
        return stack
    }
    var make_iterator func(stack []*bst) BSTIterator
    make_iterator = func(stack []*bst) BSTIterator {
        return func()(Hashable, interface{}, BSTIterator){
            var node *bst
            stack, node = pop(stack)
            if node == nil {
                return nil, nil, nil
            }
            stack = procnode(stack, node)
            return node.key, node.value, make_iterator(stack)
        }
    }
    return make_iterator([]*bst{self})
}

func (self *bst) Has(key Hashable) (has bool) {
    if self == nil {
        return false
    }
    if self.key.Equals(key) {
        return true
    } else if key.Less(self.key) {
        return self.left.Has(key)
    } else {
        return self.right.Has(key)
    }
}

func (self *bst) Get(key Hashable) (value interface{}, err error) {
    if self == nil {
        return nil, Errors["not-found"]
    }
    if self.key.Equals(key) {
        return self.value, nil
    } else if key.Less(self.key) {
        return self.left.Get(key)
    } else {
        return self.right.Get(key)
    }
}

func (self *bst) Put(key Hashable, value interface{}) (_ *bst, updated bool) {
    if self == nil {
        return &bst{hash: key.Hash(), key: key, value: value}, false
    }
    if self.key.Equals(key) {
        self.value = value
        return self, true
    }

    if key.Less(self.key) {
        self.left, updated = self.left.Put(key, value)
    } else {
        self.right, updated = self.right.Put(key, value)
    }
    return self, updated
}

func (self *bst) Remove(key Hashable) (_ *bst, value interface{}, err error) {
    if self == nil {
        return nil, nil, Errors["not-found"]
    }

    if self.key.Equals(key) {
        if self.left != nil && self.right != nil {
            if self.left.Size() < self.right.Size() {
                lmd := self.right.lmd()
                lmd.left = self.left
                return self.right, self.value, nil
            } else {
                rmd := self.left.rmd()
                rmd.right = self.right
                return self.left, self.value, nil
            }
        } else if self.left == nil {
            return self.right, self.value, nil
        } else if self.right == nil {
            return self.left, self.value, nil
        } else {
            return nil, self.value, nil
        }
    }
    if key.Less(self.key) {
        self.left, value, err = self.left.Remove(key)
    } else {
        self.right, value, err = self.right.Remove(key)
    }
    return self, value, err
}

func (self *bst) Size() int {
    if self == nil {
        return 0
    }
    return 1 + self.left.Size() + self.right.Size()
}

func (self *bst) _md(side func(*bst)*bst) (*bst) {
    if self == nil {
        return nil
    } else if side(self) != nil {
        return side(self)._md(side)
    } else {
        return self
    }
}

func (self *bst) lmd() (*bst) {
    return self._md(func(node *bst)*bst { return node.left })
}

func (self *bst) rmd() (*bst) {
    return self._md(func(node *bst)*bst { return node.right })
}

