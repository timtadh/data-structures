package trie

import (
  "fmt"
  "strings"
)

import "github.com/timtadh/data-structures/errors"

type TST struct {
    heads [256]*TSTNode
}

func (self *TST) ValidateKey(key []byte) error {
    if key == nil {
        return errors.InvalidKey(key, "key is nil")
    }
    for _, ch := range key {
        if ch == 0 {
            return errors.InvalidKey(key, "key contains a null byte")
        }
    }
    return nil
}

func (self *TST) Put(key []byte, value interface{}) (err error) {
    if err := self.ValidateKey(key); err != nil {
        return err
    }
    symbol := append(key, END)
    node, err := self.heads[symbol[0]].insert(symbol, value, 1)
    if err != nil {
        return err
    }
    self.heads[symbol[0]] = node
    return nil
}

func (self *TST) Has(key []byte) bool {
    if _, err := self.Get(key); err != nil {
        return false
    }
    return true
}

func (self *TST) Get(key []byte) (value interface{}, err error) {
    type entry struct {
        n *TSTNode
        d int
    }
    if err := self.ValidateKey(key); err != nil {
        return nil, err
    }
    symbol := append(key, END)
    next := &entry{self.heads[symbol[0]], 1}
    for next != nil {
        if next.n == nil {
            return nil, errors.NotFound(key)
        } else if next.n.Internal() {
            ch := symbol[next.d]
            if ch < next.n.ch {
                next = &entry{next.n.l, next.d}
            } else if ch == next.n.ch {
                next = &entry{next.n.m, next.d+1}
            } else if ch > next.n.ch {
                next = &entry{next.n.r, next.d}
            }
        } else if next.n.KeyEq(symbol) {
            return next.n.value, nil
        } else {
            return nil, errors.NotFound(key)
        }
    }
    // should never reach ...
    return nil, errors.NotFound(key)
}

func (self *TST) String() string {
    var nodes []string
    for i,n := range self.heads {
        if n == nil { continue }
        nodes = append(nodes, fmt.Sprintf("%x:(%v)", i, n))
    }
    return fmt.Sprintf("TST<%v>", strings.Join(nodes, ", "))
}

