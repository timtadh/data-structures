package set

import (
	"encoding/binary"
	"math/rand"
	"log"
	"os"
)

import (
	"github.com/timtadh/data-structures/types"
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/errors"
)

func init() {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		panic(err)
	} else {
		seed := make([]byte, 8)
		if _, err := urandom.Read(seed); err == nil {
			rand.Seed(int64(binary.BigEndian.Uint64(seed)))
		}
		urandom.Close()
	}
}

type SortedSet list.List

func NewSortedSet(initialSize int) *SortedSet {
	return (*SortedSet)(list.New(initialSize))
}

func FromSlice(items []types.Hashable) *SortedSet {
	s := NewSortedSet(len(items))
	for _, item := range items {
		err := s.Add(item)
		if err != nil {
			log.Panic(err)
		}
	}
	return s
}

func (s *SortedSet) Clear() {
	(*list.List)(s).Clear()
}

func (s *SortedSet) Size() int {
	return (*list.List)(s).Size()
}

func (s *SortedSet) Has(item types.Hashable) (has bool) {
	_, has = s.find(item)
	return has
}

func (s *SortedSet) Extend(other types.KIterator) (err error) {
	for item, next := other(); next != nil; item, next = next() {
		err := s.Add(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SortedSet) Add(item types.Hashable) (err error) {
	i, has := s.find(item)
	if !has {
		return (*list.List)(s).Insert(i, item)
	}
	return nil
}

func (s *SortedSet) Remove(item types.Hashable) (err error) {
	i, has := s.find(item)
	if !has {
		return errors.Errorf("item %v not in the table", item)
	}
	return (*list.List)(s).Remove(i)
	return nil
}

func (s *SortedSet) Random() (item types.Hashable, err error) {
	l := (*list.List)(s)
	if l.Size() <= 0 {
		return nil, errors.Errorf("Set is empty")
	} else if l.Size() <= 1 {
		return l.Get(0)
	}
	i := rand.Intn(l.Size())
	return l.Get(i)
}

func (s *SortedSet) Equals(b types.Equatable) bool {
	return (*list.List)(s).Equals(b)
}

func (s *SortedSet) Less(b types.Sortable) bool {
	return (*list.List)(s).Less(b)
}

func (s *SortedSet) Hash() int {
	return (*list.List)(s).Hash()
}

// Unions s with o and returns a new Sorted Set
func (s *SortedSet) Union(o *SortedSet) (n *SortedSet) {
	n = NewSortedSet(s.Size() + o.Size() + 10)
	cs, si := s.Items()()
	co, oi := o.Items()()
	for si != nil || oi != nil {
		var err error
		if si == nil {
			err = n.Add(co)
			co, oi = oi()
		} else if oi == nil {
			err = n.Add(cs)
			cs, si = si()
		} else if cs.Less(co) {
			err = n.Add(cs)
			cs, si = si()
		} else {
			err = n.Add(co)
			co, oi = oi()
		}
		if err != nil {
			log.Panic(err)
		}
	}
	return n
}

// intersect s with o and returns a new Sorted Set
func (s *SortedSet) Intersect(o *SortedSet) (n *SortedSet) {
	n = NewSortedSet(s.Size() + o.Size())
	l := (*list.List)(n)
	for v, next := s.Items()(); next != nil; v, next = next() {
		item := v.(types.Hashable)
		if o.Has(item) {
			err := l.Append(item)
			if err != nil {
				log.Panic(err)
			}
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
	n = NewSortedSet(s.Size() + o.Size())
	l := (*list.List)(n)
	for v, next := s.Items()(); next != nil; v, next = next() {
		item := v.(types.Hashable)
		if !o.Has(item) {
			err := l.Append(item)
			if err != nil {
				log.Panic(err)
			}
		}
	}
	return n
}

func (s *SortedSet) Items() (it types.KIterator) {
	return (*list.List)(s).Items()
}

func (s *SortedSet) find(item types.Hashable) (int, bool) {
	x := (*list.List)(s)
	var l int = 0
	var r int = x.Size() - 1
	var m int
	for l <= r {
		m = ((r - l) >> 1) + l
		im, err := x.Get(m)
		if err != nil {
			log.Panic(err)
		}
		if item.Less(im) {
			r = m - 1
		} else if item.Equals(im) {
			for j := m; j > 0; j-- {
				ij_1, err := x.Get(j-1)
				if err != nil {
					log.Panic(err)
				}
				if !item.Equals(ij_1) {
					return j, true
				}
			}
			return 0, true
		} else {
			l = m + 1
		}
	}
	return l, false
}

