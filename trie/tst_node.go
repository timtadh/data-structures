package trie

import (
	"fmt"
)

import (
	"github.com/timtadh/data-structures/errors"
	. "github.com/timtadh/data-structures/types"
)

const END = 0

type KV struct {
	key   ByteSlice
	value interface{}
}

func (self *KV) Key() Hashable {
	return self.key
}

func (self *KV) Value() interface{} {
	return self.value
}

func (self *KV) KeyEq(key ByteSlice) bool {
	if len(self.key) != len(key) {
		return false
	}
	for i := range self.key {
		if self.key[i] != key[i] {
			return false
		}
	}
	return true
}

type TSTNode struct {
	KV
	ch        byte     "byte to check at this node"
	l         *TSTNode "left, < side"
	m         *TSTNode "middle, == side"
	r         *TSTNode "right, > side"
	accepting bool     "is this an accepting node"
}

func NewTSTNode(ch byte) *TSTNode {
	return &TSTNode{
		ch: ch,
	}
}

func NewAcceptingTSTNode(ch byte, key ByteSlice, value interface{}) *TSTNode {
	return &TSTNode{
		KV: KV{
			key:   key,
			value: value,
		},
		ch:        ch,
		accepting: true,
	}
}

func (self *TSTNode) Copy() *TSTNode {
	node := &TSTNode{
		KV: KV{
			key:   self.key,
			value: self.value,
		},
		ch:        self.ch,
		l:         self.l,
		m:         self.m,
		r:         self.r,
		accepting: self.accepting,
	}
	return node
}

func (self *TSTNode) Internal() bool {
	return self.l != nil || self.m != nil || self.r != nil
}

func (self *TSTNode) make_child_slice() []*TSTNode {
	nodes := make([]*TSTNode, 0, 3)
	if self != nil {
		if self.l != nil {
			nodes = append(nodes, self.l)
		}
		if self.m != nil {
			nodes = append(nodes, self.m)
		}
		if self.r != nil {
			nodes = append(nodes, self.r)
		}
	}
	return nodes
}

func (self *TSTNode) Children() TreeNodeIterator {
	nodes := self.make_child_slice()
	var make_tn_iterator func(int) TreeNodeIterator
	make_tn_iterator = func(i int) TreeNodeIterator {
		return func() (kid TreeNode, next TreeNodeIterator) {
			if i < len(nodes) {
				return nodes[i], make_tn_iterator(i + 1)
			}
			return nil, nil
		}
	}
	return make_tn_iterator(0)
}

func (self *TSTNode) GetChild(i int) TreeNode {
	return self.make_child_slice()[i]
}

func (self *TSTNode) ChildCount() int {
	return len(self.make_child_slice())
}

func (self *TSTNode) String() string {
	if self == nil {
		return "-"
	}
	ch := fmt.Sprintf("%x", self.ch)
	key := ""
	if self.ch == END {
		ch = "00"
	}
	if self.key != nil {
		key = string(self.key[:len(self.key)-1])
	}
	if self.accepting {
		return fmt.Sprintf("[%v %x]", ch, key)
	}
	return fmt.Sprintf("%v(%v, %v, %v)", ch, self.l, self.m, self.r)
}

func (n *TSTNode) insert(key ByteSlice, val interface{}, d int) (*TSTNode, error) {
	if d >= len(key) {
		return nil, errors.TSTError("depth exceeds key length")
	}
	if key[len(key)-1] != END {
		return nil, errors.TSTError("key must end in 0")
	}
	if n == nil {
		// if the node is nil we found teh spot, make a new node and return it
		return NewAcceptingTSTNode(key[d], key, val), nil
	} else if !n.Internal() {
		// if it is a leaf node we either have found the symbol or we need to
		// split the node
		if n.accepting && n.KeyEq(key) {
			n = n.Copy()
			n.value = val
			return n, nil
		} else {
			return n.split(NewAcceptingTSTNode(key[d], key, val), d)
		}
	} else {
		// it is an internal node
		ch := key[d]
		n = n.Copy()
		if ch < n.ch {
			l, err := n.l.insert(key, val, d)
			if err != nil {
				return nil, err
			}
			n.l = l
		} else if ch == n.ch {
			if d+1 == len(key) && ch == END {
				n.m = n.m.Copy()
				n.m.value = val
			} else {
				m, err := n.m.insert(key, val, d+1)
				if err != nil {
					return nil, err
				}
				n.m = m
			}
		} else if ch > n.ch {
			r, err := n.r.insert(key, val, d)
			if err != nil {
				return nil, err
			}
			n.r = r
		}
		return n, nil
	}
}

/* a is the new (conflicting node)
 * b is the node that needs to be split
 * d is the depth
 *
 * both a and b must be accepting nodes.
 */
func (b *TSTNode) split(a *TSTNode, d int) (t *TSTNode, err error) {
	if !a.accepting {
		return nil, errors.TSTError("`a` must be an accepting node")
	} else if !b.accepting {
		return nil, errors.TSTError("`b` must be an accepting node")
	}
	if d >= len(b.key) {
		return nil,
			errors.TSTError("depth of split exceeds key length of b")
	}
	t = NewTSTNode(b.ch)
	b = b.Copy()
	a = a.Copy()
	if d+1 < len(b.key) {
		b.ch = b.key[d+1]
	}
	a.ch = a.key[d]
	if a.ch < t.ch {
		t.m = b
		t.l = a
	} else if a.ch == t.ch {
		m, err := b.split(a, d+1)
		if err != nil {
			return nil, err
		}
		t.m = m
	} else if a.ch > t.ch {
		t.m = b
		t.r = a
	}
	if t.m == nil {
		panic("m is nil")
	}
	return t, nil
}
