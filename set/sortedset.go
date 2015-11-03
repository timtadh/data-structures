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

type MSortedSet struct {
	list.MSorted
}

func NewMSortedSet(s *SortedSet, marshal types.ItemMarshal, unmarshal types.ItemUnmarshal) *MSortedSet {
	return &MSortedSet{
		MSorted: *list.NewMSorted(&s.Sorted, marshal, unmarshal),
	}
}

func (m *MSortedSet) SortedSet() *SortedSet {
	return &SortedSet{*m.MSorted.Sorted()}
}

// SortedSet is a list.Sorted and therefore has all of the methods
// that list.Sorted has. So although they do not show up in the generated
// docs you can just do this:
//
//     s := NewSortedSet(10)
//     s.Add(types.Int(5))
//     s2 = s.Union(FromSlice([]types.Hashable{types.Int(7)}))
//     fmt.Println(s2.Has(types.Int(7)))
//     fmt.Println(s.Has(types.Int(7)))
//
type SortedSet struct {
	list.Sorted
}

func NewSortedSet(initialSize int) *SortedSet {
	return &SortedSet{*list.NewSorted(initialSize, false)}
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

func (s *SortedSet) Copy() *SortedSet {
	return &SortedSet{ *s.Sorted.Copy() }
}

func (s *SortedSet) Random() (item types.Hashable, err error) {
	if s.Size() <= 0 {
		return nil, errors.Errorf("Set is empty")
	} else if s.Size() <= 1 {
		return s.Get(0)
	}
	i := rand.Intn(s.Size())
	return s.Get(i)
}

// Unions s with o and returns a new Sorted Set
func (s *SortedSet) Union(other types.Set) (types.Set, error) {
	if o, ok := other.(*SortedSet); ok {
		return s.union(o)
	} else {
		return Union(s, other)
	}
}

func (s *SortedSet) union(o *SortedSet) (n *SortedSet, err error) {
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
			return nil, err
		}
	}
	return n, nil
}

// Unions s with o and returns a new Sorted Set
func (s *SortedSet) Intersect(other types.Set) (types.Set, error) {
	return Intersect(s, other)
}

// Unions s with o and returns a new Sorted Set
func (s *SortedSet) Subtract(other types.Set) (types.Set, error) {
	return Subtract(s, other)
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

// Is s a subset of o?
func (s *SortedSet) Subset(o types.Set) bool {
	return Subset(s, o)
}

// Is s a proper subset of o?
func (s *SortedSet) ProperSubset(o types.Set) bool {
	return ProperSubset(s, o)
}

// Is s a superset of o?
func (s *SortedSet) Superset(o types.Set) bool {
	return Superset(s, o)
}

// Is s a proper superset of o?
func (s *SortedSet) ProperSuperset(o types.Set) bool {
	return ProperSuperset(s, o)
}


