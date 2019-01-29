package trie

import (
	"fmt"
	"strings"
)

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/tree"
	. "github.com/timtadh/data-structures/types"
)

type TST struct {
	heads [256]*TSTNode
}

func New() *TST {
	return &TST{}
}

func (self *TST) ValidateKey(key []byte) error {
	if key == nil {
		return errors.InvalidKey(key, "key is nil")
	}
	if len(key) == 0 {
		return errors.InvalidKey(key, "len(key) == 0")
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
				next = &entry{next.n.m, next.d + 1}
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

func (self *TST) Remove(key []byte) (value interface{}, err error) {
	if err := self.ValidateKey(key); err != nil {
		return nil, err
	}
	symbol := ByteSlice(append(key, END))
	check := func(n *TSTNode, err error) (*TSTNode, error) {
		if err != nil {
			return nil, err
		} else if n == nil {
			return nil, nil
		} else if !n.Internal() && n.key == nil {
			return nil, nil
		}
		return n, nil
	}
	var remove func(n *TSTNode, d int) (*TSTNode, error)
	remove = func(n *TSTNode, d int) (*TSTNode, error) {
		if n == nil {
			return nil, errors.NotFound(key)
		} else if n.Internal() {
			n = n.Copy()
			ch := symbol[d]
			if ch < n.ch {
				l, err := check(remove(n.l, d))
				if err != nil {
					return nil, err
				}
				n.l = l
			} else if ch == n.ch {
				m, err := check(remove(n.m, d+1))
				if err != nil {
					return nil, err
				}
				n.m = m
			} else if ch > n.ch {
				r, err := check(remove(n.r, d))
				if err != nil {
					return nil, err
				}
				n.r = r
			}
		} else if n.key.Equals(symbol) {
			// found it
			value = n.value
			return nil, nil
		} else {
			return nil, errors.NotFound(key)
		}
		return n, nil
	}
	n, err := remove(self.heads[symbol[0]], 1)
	if err != nil {
		return nil, err
	}
	self.heads[symbol[0]] = n
	return value, nil
}

func (self *TST) PrefixFind(prefix ByteSlice) KVIterator {
	if len(prefix) == 0 {
		return self.Iterate()
	}
	type entry struct {
		n *TSTNode
		d int
	}
	var root *TSTNode = nil
	next := &entry{self.heads[prefix[0]], 1}
	for {
		if next.n == nil {
			break
		} else if next.n.Internal() {
			if next.d == len(prefix) {
				root = next.n
				break
			}
			ch := prefix[next.d]
			if ch < next.n.ch {
				next = &entry{next.n.l, next.d}
			} else if ch == next.n.ch {
				next = &entry{next.n.m, next.d + 1}
			} else if ch > next.n.ch {
				next = &entry{next.n.r, next.d}
			}
		} else if next.n.accepting && next.n.key[:len(prefix)].Equals(prefix) {
			root = next.n
			break
		} else {
			break
		}
	}
	if root == nil {
		return func() (Hashable, interface{}, KVIterator) {
			return nil, nil, nil
		}
	}
	tni := tree.TraverseTreePreOrder(root)
	var kv_iterator KVIterator
	kv_iterator = func() (key Hashable, value interface{}, next KVIterator) {
		var tn TreeNode
		for {
			tn, tni = tni()
			if tni == nil {
				return nil, nil, nil
			}
			n := tn.(*TSTNode)
			if n.accepting {
				return n.key[:len(n.key)-1], n.value, kv_iterator
			}
		}
	}
	return kv_iterator
}

func (self *TST) Iterate() KVIterator {
	tnis := make([]TreeNodeIterator, 0, 256)
	for _, n := range self.heads {
		if n != nil {
			tnis = append(tnis, tree.TraverseTreePreOrder(n))
		}
	}
	tni := ChainTreeNodeIterators(tnis...)
	var kv_iterator KVIterator
	kv_iterator = func() (key Hashable, value interface{}, next KVIterator) {
		var tn TreeNode
		for {
			if tni == nil {
				return nil, nil, nil
			}
			tn, tni = tni()
			if tni == nil {
				return nil, nil, nil
			}
			n := tn.(*TSTNode)
			if n.accepting {
				return n.key[:len(n.key)-1], n.value, kv_iterator
			}
		}
	}
	return kv_iterator
}

func (self *TST) Items() (vi KIterator) {
	return MakeItemsIterator(self)
}

func (self *TST) Keys() KIterator {
	return MakeKeysIterator(self)
}

func (self *TST) Values() Iterator {
	return MakeValuesIterator(self)
}

func (self *TST) String() string {
	var nodes []string
	for i, n := range self.heads {
		if n == nil {
			continue
		}
		nodes = append(nodes, fmt.Sprintf("%x:(%v)", i, n))
	}
	return fmt.Sprintf("TST<%v>", strings.Join(nodes, ", "))
}

func (self *TST) Dotty() string {
	header := "digraph TST {\nrankdir=LR;\n"
	node_root := "%v[label=\"%v\", shape=\"rect\"];"
	node :=
		"%s[label=\"%v\", shape=\"circle\", fillcolor=\"#aaffff\", style=\"filled\"];"
	node_acc := "%v[label=\"%v\", fillcolor=\"#aaffaa\" style=\"filled\"];"
	edge := "%v -> %v [label=\"%v\"];"
	footer := "\n}\n"

	var nodes []string
	var edges []string

	name := func() string {
		return fmt.Sprintf("n%d", len(nodes))
	}

	var dotnode func(cur *TSTNode, parent, ch string)
	dotnode = func(cur *TSTNode, parent, ch string) {
		n := name()
		if cur.accepting {
			nodes = append(
				nodes,
				fmt.Sprintf(node_acc, n, string(cur.key[:len(cur.key)-1])),
			)
		} else if cur.ch == END {
			nodes = append(
				nodes,
				fmt.Sprintf(node, n, "\\\\0"),
			)
		} else {
			nodes = append(
				nodes,
				fmt.Sprintf(node, n, fmt.Sprintf("%c", cur.ch)),
			)
		}
		edges = append(
			edges,
			fmt.Sprintf(edge, parent, n, ch),
		)
		if cur.l != nil {
			dotnode(cur.l, n, "<")
		}
		if cur.m != nil {
			dotnode(cur.m, n, "=")
		}
		if cur.r != nil {
			dotnode(cur.r, n, ">")
		}
	}

	root := name()
	nodes = append(nodes, fmt.Sprintf(node_root, root, "heads"))

	for k, head := range self.heads {
		if head == nil {
			continue
		}
		dotnode(head, root, fmt.Sprintf("%c", byte(k)))
	}

	ret := header + strings.Join(nodes, "\n")
	ret += "\n" + strings.Join(edges, "\n")
	ret += footer
	return ret
}
