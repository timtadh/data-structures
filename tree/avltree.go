package tree

import (
  "github.com/timtadh/data-structures/types"
  "github.com/timtadh/data-structures/errors"
)

func abs(i int) int {
    if i < 0 {
        return -i
    }
    return i
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

type AvlTree struct {
    key types.Sortable
    value interface{}
    height int
    left *AvlTree
    right *AvlTree
}

func NewAvlTree() *AvlTree {
    return nil
}

func (self *AvlTree) Has(key types.Sortable) (has bool) {
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

func (self *AvlTree) Get(key types.Sortable) (value interface{}, err error) {
    if self == nil {
        return nil, errors.NotFound(key)
    }
    if self.key.Equals(key) {
        return self.value, nil
    } else if key.Less(self.key) {
        return self.left.Get(key)
    } else {
        return self.right.Get(key)
    }
}

func (self *AvlTree) pop_node(node *AvlTree) *AvlTree {
    if node == nil {
        panic("node can't be nil")
    } else if node.left != nil && node.right != nil {
        panic("node must not have both left and right")
    }

    if self == nil {
        return nil
    } else if self == node {
        var n *AvlTree
        if node.left != nil {
            n = node.left
        } else if node.right != nil {
            n = node.right
        } else {
            n = nil
        }
        node.left = nil
        node.right = nil
        return n
    }

    if node.key.Less(self.key) {
        self.left = self.left.pop_node(node)
    } else {
        self.right = self.right.pop_node(node)
    }

    self.height = max(self.left.Height(), self.right.Height()) + 1
    return self
}

func (self *AvlTree) push_node(node *AvlTree) *AvlTree {
    if node == nil {
        panic("node can't be nil")
    } else if node.left != nil || node.right != nil {
        panic("node now be a leaf")
    }

    if self == nil {
        node.height = 1
        return node
    } else if node.key.Less(self.key) {
        self.left = self.left.push_node(node)
    } else {
        self.right = self.right.push_node(node)
    }
    self.height = max(self.left.Height(), self.right.Height()) + 1
    return self
}

func (self *AvlTree) rotate_right() *AvlTree {
    if self == nil {
        return self
    }
    if self.left == nil {
        return self
    }
    new_root := self.left.rmd()
    self = self.pop_node(new_root)
    new_root.left = self.left
    new_root.right = self.right
    self.left = nil
    self.right = nil
    return new_root.push_node(self)
}

func (self *AvlTree) rotate_left() *AvlTree {
    if self == nil {
        return self
    }
    if self.right == nil {
        return self
    }
    new_root := self.right.lmd()
    self = self.pop_node(new_root)
    new_root.left = self.left
    new_root.right = self.right
    self.left = nil
    self.right = nil
    return new_root.push_node(self)
}

func (self *AvlTree) balance() *AvlTree {
    if self == nil {
        return self
    }
    for abs(self.left.Height() - self.right.Height()) > 2 {
        if self.left.Height() > self.right.Height() {
            self = self.rotate_right()
        } else {
            self = self.rotate_left()
        }
    }
    return self
}

func (self *AvlTree) Put(key types.Sortable, value interface{}) (_ *AvlTree, updated bool) {
    if self == nil {
        return &AvlTree{key: key, value: value, height: 1}, false
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
    if !updated {
        self.height += 1
        return self.balance(), updated
    }
    return self, updated
}

func (self *AvlTree) Remove(key types.Sortable) (_ *AvlTree, value interface{}, err error) {
    if self == nil {
        return nil, nil, errors.NotFound(key)
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
    if err != nil {
        return self.balance(), value, err
    }
    return self, value, err
}

func (self *AvlTree) Height() int {
    if self == nil {
        return 0
    }
    return self.height
}

func (self *AvlTree) Size() int {
    if self == nil {
        return 0
    }
    return 1 + self.left.Size() + self.right.Size()
}


func pop(stack []*AvlTree) ([]*AvlTree, *AvlTree) {
    if len(stack) <= 0 {
        return stack, nil
    } else {
        return stack[0:len(stack)-1], stack[len(stack)-1]
    }
}

func procnode(stack []*AvlTree, node *AvlTree) []*AvlTree {
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

func (self *AvlTree) Iterate() types.KVIterator {
    stack := make([]*AvlTree, 0, 10)
    cur := self
    var kv_iterator types.KVIterator
    kv_iterator = func()(key types.Equatable, val interface{}, next types.KVIterator) {
        if len(stack) > 0 || cur != nil {
            for cur != nil {
                stack = append(stack, cur)
                cur = cur.left
            }
            stack, cur = pop(stack)
            key = cur.key
            val = cur.value
            cur = cur.right
            return key, val, kv_iterator
        } else {
            return nil, nil, nil
        }
    }
    return kv_iterator
}

func (self *AvlTree) Keys() types.KIterator {
    kv_iterator := self.Iterate()
    var k_iterator types.KIterator
    k_iterator = func() (key types.Equatable, next types.KIterator) {
        key, _, kv_iterator = kv_iterator()
        if kv_iterator == nil {
            return nil, nil
        }
        return key, k_iterator
    }
    return k_iterator
}

func (self *AvlTree) Values() types.Iterator {
    kv_iterator := self.Iterate()
    var v_iterator types.Iterator
    v_iterator = func() (value interface{}, next types.Iterator) {
        _, value, kv_iterator = kv_iterator()
        if kv_iterator == nil {
            return nil, nil
        }
        return value, v_iterator
    }
    return v_iterator
}



func (self *AvlTree) _md(side func(*AvlTree)*AvlTree) (*AvlTree) {
    if self == nil {
        return nil
    } else if side(self) != nil {
        return side(self)._md(side)
    } else {
        return self
    }
}

func (self *AvlTree) lmd() (*AvlTree) {
    return self._md(func(node *AvlTree)*AvlTree { return node.left })
}

func (self *AvlTree) rmd() (*AvlTree) {
    return self._md(func(node *AvlTree)*AvlTree { return node.right })
}

