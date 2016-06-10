package heap

import "testing"

import (
	"github.com/timtadh/data-structures/test"
	"github.com/timtadh/data-structures/types"
)

func TestUniquePushMax(x *testing.T) {
	t := (*test.T)(x)
	h := NewUnique(NewHeap(12, false))
	h.Push(18, types.String("g"))
	h.Push(3, types.String("g"))
	h.Push(5, types.String("g"))
	h.Push(7, types.String("e"))
	h.Push(3, types.String("i"))
	h.Push(6, types.String("a"))
	h.Push(25, types.String("x"))
	h.Push(22, types.String("t"))
	h.Push(14, types.String("n"))
	h.Push(8, types.String("b"))
	h.Push(19, types.String("m"))
	h.Push(9, types.String("r"))
	h.Push(20, types.String("o"))
	h.Push(16, types.String("s"))
	t.Log(h.pq)
	t.AssertNil(h.pq.(*Heap).Verify())
	t.Assert(h.pq.(*Heap).list[3].item.(types.String) == "g", "heap[3] != {g 18} %v", h.pq.(*Heap).list[3])
	t.Assert(h.pq.(*Heap).list[3].priority == 18, "heap[3] != {g 18} %v", h.pq.(*Heap).list[3])
}
