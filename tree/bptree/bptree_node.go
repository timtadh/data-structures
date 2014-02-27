package tree

import (
    "github.com/timtadh/data-structures/types"
    "github.com/timtadh/data-structures/errors"
)

type BpNode struct {
    keys []types.Sortable
    values []interface{}
    pointers []*BpNode
    next *BpNode
    prev *BpNode
    height int
}

func NewInternal(size, height int) *BpNode {
    if size < 0 {
        panic(errors.NegativeSize())
    }
    return &BpNode{
        keys: make([]types.Sortable, 0, size),
        values: make([]interface{}, 0, size),
        pointers: make([]*BpNode, 0, size),
    }
}

func (self *BpNode) Full() bool {
    return len(self.keys) == cap(self.keys)
}

func (self *BpNode) Internal() bool {
    return cap(self.pointers) > 0
}

func NewLeaf(size int) *BpNode {
    if size < 0 {
        panic(errors.NegativeSize())
    }
    return &BpNode{
        keys: make([]types.Sortable, size),
        values: make([]interface{}, size),
    }
}

func (self *BpNode) Height() {
    if !self.Internal() {
        return 1
    } else if len(self.pointers) == 0 {
        panic(errors.BpTreeError("Internal node has no pointers but asked for height"))
    }
    return self.pointers[0] + 1
}

func (self *BpNode) setNext(next *BpNode) {
    if self.Internal() {
        panic(errors.BpTreeError("Expected a leaf node"))
    }
    if next.Internal() {
        panic(errors.BpTreeError("Expected a leaf node"))
    }
    self.next = next
    next.prev = prev
}

func (self *BpNode) setPrev(prev *BpNode) {
    if self.Internal() {
        panic(errors.BpTreeError("Expected a leaf node"))
    }
    if prev.Internal() {
        panic(errors.BpTreeError("Expected a leaf node"))
    }
    self.prev = prev
    prev.next = self
}

// right is only set on split
// left is always set. When split is false left is the pointer to block
//                     When split is true left is the pointer to the new left
//                     block
func (self *BpNode) insert(key types.Sortable, value interface{}) (left, right *BpNode, split bool err error) {
    height := self.Height()
    i, ok := self.find(key)
    if height > 1 {
        // an internal node
        // find the child node to insert into
        var child *BpNode = nil
        if ok {
            // the key was found in the block
            // so get the pointer
            child = self.pointers[i]
        } else if i == 0 {
            // >= the smallest key in the block
            // we need to adjust the first key to be the inserted key
            self.keys[key]
            child = self.pointers[0]
        } else if i >= len(self.keys) {
            // else this spot is one too many eg.
            //         0  1   2
            // keys =  8, 15, 21
            // search key = 10
            // find returns:
            //    i = 1, ok = false
            // but it goes in block 0, (key 8)
            i--
            child = self.pointers[0]
        }
        if child == nil {
            return nil, errors.BpTreeError("Child was nil")
        }

        cl, cr, split, err := child.insert(keys, value);
        if  err != nil {
            return nil, nil, false, err
        }

        self.pointers[i] = cl
        if split && !self.Full() {
            if err := self.put_kp(cr.keys[0], cr); err != nil {
                return nil, nil, false, err
            }
        } else if split {
            // this block needs to split
            if a, b, err := self.split_p(cr.keys[0], cr); err != nil {
                return nil, nil, false, err
            } else {
                return a, b, true, nil
            }
        }
    } else if !self.Full() {
        // this is a leaf node
        if err := self.put_kv(key, value); err != nil {
            return nil, nil, false, err
        } else {
            return self, nil, false, nil
        }
    } else {
        // this is a full leaf node
        if a, b, err := self.split_v(key, value); err != nil {
            return nil, nil, false, err
        } else {
            return a, b, true, nil
        }
    }
}

func (self *BpNode) put_kp(key types.Sortable, ptr *BpTree) error {
    if self.Full() {
        return errors.BpTreeError("Block is full.")
    }
    if !self.Internal() {
        return errors.BpTreeError("Expected a internal node")
    }
    i, has := self.find()
    if has {
        return errors.BpTreeError("Tried to insert a duplicate key into an internal node")
    } else if i < 0 {
        panic(errors.BpTreeError("find returned a negative int"))
    } else if i >= len(self.keys) {
        panic(errors.BpTreeError("find returned a int > than len(keys)"))
    }
    // ok good to insert extend the keys/pointers slices
    self.keys = self.keys[:len(self.keys)+1]
    self.pointers = self.pointers[:len(self.pointers)+1]
    for j := len(self.keys) - 1; j > i; j-- {
        self.keys[j] = self.keys[j-1]
        self.pointers[j] = self.pointers[j-1]
    }
    self.keys[i] = key
    self.pointers[i] = ptr
    return nil
}

func (self *BpNode) put_kv(key types.Sortable, value interface{}) error {
    if self.Full() {
        return errors.BpTreeError("Block is full.")
    }
    if self.Internal() {
        return errors.BpTreeError("Expected a leaf node")
    }
    i, has := self.find()
    if i < 0 {
        panic(errors.BpTreeError("find returned a negative int"))
    } else if i >= len(self.keys) {
        panic(errors.BpTreeError("find returned a int > than len(keys)"))
    }
    // ok good to insert extend the keys/pointers slices
    self.keys = self.keys[:len(self.keys)+1]
    self.values = self.values[:len(self.values)+1]
    for j := len(self.keys) - 1; j > i; j-- {
        self.keys[j] = self.keys[j-1]
        self.values[j] = self.values[j-1]
    }
    self.keys[i] = key
    self.values[i] = value
    return nil
}

func (self *BpNode) Has(key types.Sortable) bool {
    _, has := self.find(key)
    return has
}

func (self *BpNode) get_p(key types.Sortable) (ptr *BpNode, error) {
    if !self.Internal() {
        return nil, errors.BpTreeError("Expected a internal node")
    }
    i, has := self.find()
    if !has {
        return nil, errors.BpTreeError("Key was not in node")
    }
    return self.pointers[i], nil
}

func (self *BpNode) get_v(key types.Sortable) (value *BpNode, error) {
    if self.Internal() {
        return nil, errors.BpTreeError("Expected a leaf node")
    }
    i, has := self.find()
    if !has {
        return nil, errors.BpTreeError("Key was not in node")
    }
    return self.values[i], nil
}

func (self *BpNode) find(key types.Sortable) (int, bool) {
    var l int = 0
    var r int = len(self.keys) - 1
    var m int
    for l <= r {
        m = ((r - l) >> 1) + l
        if key.Less(self.keys[m]) {
            r = m - 1
        } else if key.Equals(self.keys[m]) {
            for j := m; j >= 0; j-- {
                if j == 0 || !key.Equals(self.records[j-1]) {
                    return j, true
                }
            }
        } else {
            l = m + 1
        }
    }
    return l, false
}

