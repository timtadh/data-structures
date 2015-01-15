package avl

import "testing"

import (
	"github.com/timtadh/data-structures/tree"
	"github.com/timtadh/data-structures/types"
)

func TestTraversals(t *testing.T) {
	var data []int = []int{
		1, 5, 7, 9, 12, 13, 17, 18, 19, 20,
	}
	var order []int = []int{
		6, 1, 8, 2, 4, 9, 5, 7, 0, 3,
	}
	var preorder []int = []int{
		17, 7, 5, 1, 12, 9, 13, 19, 18, 20,
	}
	var postorder []int = []int{
		1, 5, 9, 13, 12, 7, 18, 20, 19, 17,
	}

	test := func(T types.TreeMap) {
		t.Logf("%T", T)
		for j := range order {
			if err := T.Put(types.Int(data[order[j]]), order[j]); err != nil {
				t.Error(err)
			}
		}

		j := 0
		for tn, next := tree.TraverseBinaryTreeInOrder(T.Root().(types.BinaryTreeNode))(); next != nil; tn, next = next() {
			if int(tn.Key().(types.Int)) != data[j] {
				t.Error("key in wrong spot in-order")
			}
			j += 1
		}

		j = 0
		for tn, next := tree.TraverseTreePreOrder(T.Root())(); next != nil; tn, next = next() {
			if int(tn.Key().(types.Int)) != preorder[j] {
				t.Error("key in wrong spot pre-order")
			}
			j += 1
		}

		j = 0
		for tn, next := tree.TraverseTreePostOrder(T.Root())(); next != nil; tn, next = next() {
			if int(tn.Key().(types.Int)) != postorder[j] {
				t.Error("key in wrong spot post-order")
			}
			j += 1
		}
	}
	test(NewAvlTree())
	test(NewImmutableAvlTree())
}
