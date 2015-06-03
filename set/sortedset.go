package set

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

import (
	"github.com/timtadh/data-structures/types"
)

type SortedSet struct {
	set []types.Hashable
}

func NewSortedSet(initialSize int) *SortedSet {
	return &SortedSet{
		set: make([]types.Hashable, 0, initialSize),
	}
}

func (s *SortedSet) Clear() {
	s.set = s.set[:0]
}

func (s *SortedSet) Size() int {
	return len(s.set)
}

func (s *SortedSet) Has(item types.Hashable) (has bool) {
	_, has = s.find(item)
	return has
}

func (s *SortedSet) Extend(other *SortedSet) (err error) {
	for _, item := range other.set {
		if err := s.Add(item); err != nil {
			return err
		}
	}
	return nil
}

func (s *SortedSet) Add(item types.Hashable) (err error) {
	i, has := s.find(item)
	if !has {
		s.insert(i, item)
	}
	return nil
}

func (s *SortedSet) Remove(item types.Hashable) (err error) {
	i, has := s.find(item)
	if !has {
		return fmt.Errorf("item %v not in the table", item)
	}
	s.remove(i)
	return nil
}

func (s *SortedSet) Equals(b types.Equatable) bool {
	if o, ok := b.(types.SetOperable); ok {
		return s.equals(o)
	} else {
		return false
	}
}

func (s *SortedSet) equals(o types.SetOperable) bool {
	if s.Size() != o.Size() {
		return false
	}
	for v, next := s.Items()(); next != nil; v, next = next() {
		item := v.(types.Hashable)
		if !o.Has(item) {
			return false
		}
	}
	return true
}

func (s *SortedSet) Less(b types.Sortable) bool {
	if o, ok := b.(types.Set); ok {
		return s.less(o)
	} else {
		return false
	}
}

func (s *SortedSet) less(o types.Set) bool {
	if s.Size() < o.Size() {
		return true
	} else if s.Size() > o.Size() {
		return false
	}
	cs, si := s.Items()()
	co, oi := o.Items()()
	for si != nil || oi != nil {
		if cs.(types.Hashable).Less(co.(types.Hashable)) {
			return true
		} else if !cs.Equals(co) {
			return false
		}
		cs, si = si()
		co, oi = oi()
	}
	return true
}

func (s *SortedSet) Hash() int {
	h := fnv.New32a()
	if len(s.set) == 0 {
		return 0
	}
	bs := make([]byte, 4)
	for _, item := range s.set {
		binary.LittleEndian.PutUint32(bs, uint32(item.Hash()))
		h.Write(bs)
	}
	return int(h.Sum32())
}

// Unions s with o and returns a new Sorted Set
func (s *SortedSet) Union(o *SortedSet) (n *SortedSet) {
	n = NewSortedSet(s.Size() + o.Size() + 10)
	cs, si := s.Items()()
	co, oi := o.Items()()
	for si != nil || oi != nil {
		if si == nil {
			n.set = append(n.set, co.(types.Hashable))
			co, oi = oi()
		} else if oi == nil {
			n.set = append(n.set, cs.(types.Hashable))
			cs, si = si()
		} else if cs.(types.Hashable).Less(co.(types.Hashable)) {
			n.set = append(n.set, cs.(types.Hashable))
			cs, si = si()
		} else {
			n.set = append(n.set, co.(types.Hashable))
			co, oi = oi()
		}
	}
	return n
}

// intersect s with o and returns a new Sorted Set
func (s *SortedSet) Intersect(o *SortedSet) (n *SortedSet) {
	n = NewSortedSet(cap(s.set))
	for v, next := s.Items()(); next != nil; v, next = next() {
		item := v.(types.Hashable)
		if o.Has(item) {
			n.set = append(n.set, item)
		}
	}
	return n
}

// Are there any overlapping elements?
func (s *SortedSet) Overlap(o *SortedSet) bool {
	cs, si := s.Items()()
	co, oi := o.Items()()
	for si != nil && oi != nil {
		s := cs.(types.Hashable)
		o := co.(types.Hashable)
		if s.Equals(o) {
			return true
		} else if s.Less(o) {
			cs, si = si()
		} else {
			co, oi = oi()
		}
	}
	return false
}

// subtract o from s and return new Sorted Set
func (s *SortedSet) Subtract(o *SortedSet) (n *SortedSet) {
	n = NewSortedSet(cap(s.set))
	for v, next := s.Items()(); next != nil; v, next = next() {
		item := v.(types.Hashable)
		if !o.Has(item) {
			n.set = append(n.set, item)
		}
	}
	return n
}

func (s *SortedSet) Items() (it types.KIterator) {
	i := 0
	return func() (item types.Equatable, next types.KIterator) {
		if i < len(s.set) {
			item = s.set[i]
			i++
			return item, it
		}
		return nil, nil
	}
}

func (s *SortedSet) insert(i int, item types.Hashable) {
	if len(s.set) == cap(s.set) {
		s.expand()
	}
	s.set = s.set[:len(s.set)+1]
	for j := len(s.set) - 1; j > 0; j-- {
		if j == i {
			s.set[i] = item
			break
		}
		s.set[j] = s.set[j-1]
	}
	if i == 0 {
		s.set[i] = item
	}
}

func (s *SortedSet) remove(i int) {
	for j := i; j+1 < len(s.set); j++ {
		s.set[j] = s.set[j+1]
	}
	s.set = s.set[:len(s.set)-1]
	s.shrink()
}

func (s *SortedSet) expand() {
	set := s.set
	if cap(set) < 100 {
		s.set = make([]types.Hashable, len(set), cap(set)*2)
	} else {
		s.set = make([]types.Hashable, len(set), cap(set)+100)
	}
	copy(s.set, set)
}

func (s *SortedSet) shrink() {
	if (len(s.set)-1)*2 >= cap(s.set) {
		return
	}
	set := s.set
	s.set = make([]types.Hashable, len(set), cap(set)/2)
	copy(s.set, set)
}

func (s *SortedSet) find(item types.Hashable) (int, bool) {
	var l int = 0
	var r int = len(s.set) - 1
	var m int
	for l <= r {
		m = ((r - l) >> 1) + l
		if item.Less(s.set[m]) {
			r = m - 1
		} else if item.Equals(s.set[m]) {
			for j := m; j >= 0; j-- {
				if j == 0 || !item.Equals(s.set[j-1]) {
					return j, true
				}
			}
		} else {
			l = m + 1
		}
	}
	return l, false
}
