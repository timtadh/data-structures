package heap

import (
	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
)

// This priority queue only allows unique entries. Internally this is
// implemented using a Hash set. All items added must be types.Hashable
type UniquePQ struct {
	pq  PriorityQueue
	set *set.SetMap
}

// Construct a new unique priority queue using the provided priority queue.
func NewUnique(pq PriorityQueue) *UniquePQ {
	return &UniquePQ{
		pq:  pq,
		set: set.NewSetMap(hashtable.NewLinearHash()),
	}
}

// How many items in the queue?
func (u *UniquePQ) Size() int {
	return u.pq.Size()
}

// Add an item to the priority queue. It must be hashable.
func (u *UniquePQ) Add(priority int, item types.Hashable) {
	if !u.set.Has(item) {
		u.set.Add(item)
		u.pq.Push(priority, item)
	}
}

// This method is provided so it implements the PriorityQueue interface. In
// reality item must be types.Hashable. The implementation simply does a type
// assertion on item and calls Add.
func (u *UniquePQ) Push(priority int, item interface{}) {
	u.Add(priority, item.(types.Hashable))
}

// Get the top element
func (u *UniquePQ) Peek() interface{} {
	return u.pq.Peek()
}

// Get and remove the top element
func (u *UniquePQ) Pop() interface{} {
	item := u.pq.Pop().(types.Hashable)
	u.set.Delete(item)
	return item
}
