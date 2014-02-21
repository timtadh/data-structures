package types

import (
    "reflect"
)

func IsNil(object interface{}) bool {
    return object == nil || reflect.ValueOf(object).IsNil()
}

func MakeKVIteratorFromTreeNodeIterator(tni TreeNodeIterator) KVIterator {
    var kv_iterator KVIterator
    kv_iterator = func() (key Equatable, value interface{}, next KVIterator) {
        var tn TreeNode
        tn, tni = tni()
        if tni == nil {
            return nil, nil, nil
        }
        return tn.Key(), tn.Value(), kv_iterator
    }
    return kv_iterator
}

func ChainTreeNodeIterators(tnis ...TreeNodeIterator) TreeNodeIterator {
    var make_tni func(int)TreeNodeIterator
    make_tni = func(i int) TreeNodeIterator {
        var tni TreeNodeIterator = nil
        if i < len(tnis) {
            tni = tnis[i]
        }
        return func() (tn TreeNode, next TreeNodeIterator) {
            if tni == nil {
                return nil, nil
            }
            tn, next = tni()
            for next == nil && i < len(tnis) {
                tni = make_tni(i+1)
                tn, next = tni()
            }
            if next == nil {
                return nil, nil
            }
            return tn, next
        }
    }
    return make_tni(0)
}

func MakeKeysIterator(obj KVIterable) KIterator {
    kv_iterator := obj.Iterate()
    var k_iterator KIterator
    k_iterator = func() (key Equatable, next KIterator) {
        key, _, kv_iterator = kv_iterator()
        if kv_iterator == nil {
            return nil, nil
        }
        return key, k_iterator
    }
    return k_iterator
}

func MakeValuesIterator(obj KVIterable) Iterator {
    kv_iterator := obj.Iterate()
    var v_iterator Iterator
    v_iterator = func() (value interface{}, next Iterator) {
        _, value, kv_iterator = kv_iterator()
        if kv_iterator == nil {
            return nil, nil
        }
        return value, v_iterator
    }
    return v_iterator
}

func make_child_slice(node BinaryTreeNode) []BinaryTreeNode {
    nodes := make([]BinaryTreeNode, 0, 2)
    if !IsNil(node) {
        if !IsNil(node.Left()) {
            nodes = append(nodes, node.Left())
        }
        if !IsNil(node.Right()) {
            nodes = append(nodes, node.Right())
        }
    }
    return nodes
}

func DoGetChild(node BinaryTreeNode, i int) TreeNode {
    return make_child_slice(node)[i]
}

func DoChildCount(node BinaryTreeNode) int {
    return len(make_child_slice(node))
}

func MakeChildrenIterator(node BinaryTreeNode) TreeNodeIterator {
    nodes := make_child_slice(node)
    var make_tn_iterator func(int) TreeNodeIterator
    make_tn_iterator = func(i int) TreeNodeIterator {
        return func() (kid TreeNode, next TreeNodeIterator) {
            if i < len(nodes) {
                return nodes[i], make_tn_iterator(i+1)
            }
            return nil, nil
        }
    }
    return make_tn_iterator(0)
}

