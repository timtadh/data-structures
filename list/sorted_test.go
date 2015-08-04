package list

import "testing"

import (
	"fmt"
	"math/rand"
)

import (
	"github.com/timtadh/data-structures/types"
)


func TestSortedAddHasRemove(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	set := NewSorted(10, false)
	items := make([]types.Int, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := types.Int(rand.Intn(1000))
		for set.Has(item) {
			item = types.Int(rand.Intn(1000))
		}
		items = append(items, item)
		t.assert_nil(set.Add(item))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), set.Has(item))
	}
	for _, item := range items {
		t.assert_nil(set.Remove(item))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), !set.Has(item))
	}
}

func TestSortedExtend(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	a := NewSorted(10, false)
	b := NewSorted(10, false)
	items := make([]types.ByteSlice, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := t.randslice(rand.Intn(10)+1)
		items = append(items, item)
		if i < SIZE/2 {
			t.assert_nil(a.Add(item))
		} else {
			t.assert_nil(b.Add(item))
		}
	}
	t.assert_nil(a.Extend(b.Items()))
	for i := SIZE - 1; i >= 0; i-- {
		err := a.Remove(items[i])
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, !a.Has(item)", i), !a.Has(items[i]))
	}
}



func TestSortedLess(x *testing.T) {
	t := (*T)(x)
	a := SortedFromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)}, false)
	b := SortedFromSlice([]types.Hashable{types.Int(3), types.Int(2), types.Int(3)}, false)
	c := SortedFromSlice([]types.Hashable{types.Int(1), types.Int(0), types.Int(3)}, false)
	d := SortedFromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)}, false)
	e := SortedFromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)}, false)
	small := SortedFromSlice([]types.Hashable{types.Int(2), types.Int(4)}, false)
	big := SortedFromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(4)}, false)
	t.assert("b < a", b.Less(a))
	t.assert("c < a", c.Less(a))
	t.assert("b < c", b.Less(c))
	t.assert("a !< d", a.Less(d) == false)
	t.assert("a !< e", e.Less(a) == false)
	t.assert("a < big", a.Less(big))
	t.assert("small < a", small.Less(a))
	t.assert("a !< small", a.Less(small) == false)
}

func TestSortedEqualsHash(x *testing.T) {
	t := (*T)(x)
	a := SortedFromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)}, false)
	b := SortedFromSlice([]types.Hashable{types.Int(3), types.Int(2), types.Int(3)}, false)
	c := SortedFromSlice([]types.Hashable{types.Int(1), types.Int(0), types.Int(3)}, false)
	d := SortedFromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)}, false)
	small := SortedFromSlice([]types.Hashable{types.Int(2), types.Int(4)}, false)
	empty := SortedFromSlice([]types.Hashable{}, false)
	t.assert("a != b", !a.Equals(b))
	t.assert("c != a", !c.Equals(a))
	t.assert("c != b", !c.Equals(b))
	t.assert("a == d", a.Equals(d))
	t.assert("c != small", !c.Equals(small))
	t.assert("a.Hash() != b.Hash()", a.Hash() != b.Hash())
	t.assert("c.Hash() != b.Hash()", c.Hash() != b.Hash())
	t.assert("a.Hash() != d.Hash()", a.Hash() == d.Hash())
	t.assert("d.Hash() != b.Hash()", d.Hash() != b.Hash())
	t.assert("d.Hash() != empty.Hash()", d.Hash() != empty.Hash())
}

