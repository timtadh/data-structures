package heap

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/types"
)

type entry struct {
	item     interface{}
	priority int
}

type PriorityQueue interface {
	types.Sized
	Push(priority int, item interface{})
	Peek() interface{}
	Pop() interface{}
}

// Notes:
// Parent of i : (i+1)/2 - 1
// Left Child of i : (i+1)*2 - 1
// Right Child of i : (i+1)*2

// A binary heap for Priority Queues. The priorities are modeled
// explicitly as integers. It can work either as a min heap or a max
// heap.
type Heap struct {
	min  bool
	list []entry
}

// Make a new binary heap.
// size : hint for the size of the heap
//        (should estimate the maximal size)
// min : false == max heap, true == min heap
func NewHeap(size int, min bool) *Heap {
	return &Heap{
		min:  min,
		list: make([]entry, 0, size),
	}
}

func NewMinHeap(size int) *Heap {
	return &Heap{
		min:  true,
		list: make([]entry, 0, size),
	}
}

func NewMaxHeap(size int) *Heap {
	return &Heap{
		min:  false,
		list: make([]entry, 0, size),
	}
}

// How many items in the heap?
func (h *Heap) Size() int {
	return len(h.list)
}

// Is this a min heap?
func (h *Heap) MinHeap() bool {
	return h.min
}

// Is this a max heap?
func (h *Heap) MaxHeap() bool {
	return !h.min
}

// Push an item with a priority
func (h *Heap) Push(priority int, item interface{}) {
	h.list = append(h.list, entry{item, priority})
	h.fixUp(len(h.list) - 1)
}

// Pop the highest (or lowest) priority item
func (h *Heap) Pop() interface{} {
	if len(h.list) == 0 {
		return nil
	}
	i := h.list[0].item
	h.list[0] = h.list[len(h.list)-1]
	h.list = h.list[:len(h.list)-1]
	h.fixDown(0)
	return i
}

// Peek at the highest (or lowest) priority item
func (h *Heap) Peek() interface{} {
	if len(h.list) == 0 {
		return nil
	}
	return h.list[0].item
}

func (h *Heap) Items() (it types.Iterator) {
	i := 0
	return func() (item interface{}, next types.Iterator) {
		var e entry
		if i < len(h.list) {
			e = h.list[i]
			i++
			return e.item, it
		}
		return nil, nil
	}
}

func (h *Heap) fixUp(k int) {
	parent := (k+1)/2 - 1
	for k > 0 {
		if h.gte(parent, k) {
			return
		}
		h.list[parent], h.list[k] = h.list[k], h.list[parent]
		k = parent
		parent = (k+1)/2 - 1
	}
}

func (h *Heap) fixDown(k int) {
	kid := (k+1)*2 - 1
	for kid < len(h.list) {
		if kid+1 < len(h.list) && h.lt(kid, kid+1) {
			kid++
		}
		if h.gte(k, kid) {
			break
		}
		h.list[kid], h.list[k] = h.list[k], h.list[kid]
		k = kid
		kid = (k+1)*2 - 1
	}
}

func (h *Heap) gte(i, j int) bool {
	if h.min {
		return h.list[i].priority <= h.list[j].priority
	} else {
		return h.list[i].priority >= h.list[j].priority
	}
}

func (h *Heap) lt(i, j int) bool {
	if h.min {
		return h.list[i].priority > h.list[j].priority
	} else {
		return h.list[i].priority < h.list[j].priority
	}
}

// Verify the heap is properly structured.
func (h *Heap) Verify() error {
	for i := 1; i < len(h.list); i++ {
		parent := (i+1)/2 - 1
		if h.lt(parent, i) {
			return errors.Errorf("parent %v '<' kid %v", h.list[parent], h.list[i])
		}
	}
	return nil
}
