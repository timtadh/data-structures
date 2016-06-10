package heap

import "testing"

import (
	"github.com/timtadh/data-structures/test"
)

func heap1(min bool) *Heap {
	h := NewHeap(12, min)
	if min {
		h.list = []entry{
			{"x", 3},
			{"t", 7},
			{"o", 5},
			{"g", 18},
			{"s", 16},
			{"m", 19},
			{"n", 14},
			{"b", 28},
			{"e", 27},
			{"r", 29},
			{"a", 26},
			{"i", 23},
		}
	} else {
		h.list = []entry{
			{"x", 25},
			{"t", 22},
			{"o", 20},
			{"g", 18},
			{"s", 16},
			{"m", 19},
			{"n", 14},
			{"b", 8},
			{"e", 7},
			{"r", 9},
			{"a", 6},
			{"i", 3},
		}
	}
	return h
}

func TestMaxFixUp(x *testing.T) {
	t := (*test.T)(x)
	h := heap1(false)
	t.Log(h)
	h.list[10].priority = 30
	h.fixUp(10)
	t.Log(h)
	t.Assert(h.list[0].item.(string) == "a", "heap did not start with {a 30} %v", h)
}

func TestMinFixUp(x *testing.T) {
	t := (*test.T)(x)
	h := heap1(true)
	t.Log(h)
	h.list[10].priority = 1
	h.fixUp(10)
	t.Log(h)
	t.Assert(h.list[0].item.(string) == "a", "heap did not start with {a 1} %v", h)
}

func TestMaxFixDown(x *testing.T) {
	t := (*test.T)(x)
	h := heap1(false)
	t.Log(h)
	h.list[0].priority = 0
	h.fixDown(0)
	t.Log(h)
	t.Assert(h.list[0].item.(string) == "t", "heap did not start with {t 22} %v", h)
	t.Assert(h.list[7].item.(string) == "x", "heap[8] != {x 0} %v", h)
}

func TestMinFixDown(x *testing.T) {
	t := (*test.T)(x)
	h := heap1(true)
	t.Log(h)
	h.list[0].priority = 30
	h.fixDown(0)
	t.Log(h)
	t.Assert(h.list[0].item.(string) == "o", "heap did not start with {o 5} %v", h)
	t.Assert(h.list[6].item.(string) == "x", "heap[8] != {n 30} %v", h)
}

func TestPushMax(x *testing.T) {
	t := (*test.T)(x)
	h := NewHeap(12, false)
	h.Push(18, "g")
	h.Push(7, "e")
	h.Push(3, "i")
	h.Push(6, "a")
	h.Push(25, "x")
	h.Push(22, "t")
	h.Push(14, "n")
	h.Push(8, "b")
	h.Push(19, "m")
	h.Push(9, "r")
	h.Push(20, "o")
	h.Push(16, "s")
	t.Log(h)
	t.AssertNil(h.Verify())
}

func TestPushMin(x *testing.T) {
	t := (*test.T)(x)
	h := NewHeap(12, true)
	h.Push(18, "g")
	h.Push(7, "e")
	h.Push(3, "i")
	h.Push(6, "a")
	h.Push(25, "x")
	h.Push(22, "t")
	h.Push(14, "n")
	h.Push(8, "b")
	h.Push(19, "m")
	h.Push(9, "r")
	h.Push(20, "o")
	h.Push(16, "s")
	t.Log(h)
	t.AssertNil(h.Verify())
}
