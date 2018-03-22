package set

import "testing"

import (
	"fmt"
	"runtime/debug"
)

import (
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/types"
)

type T testing.T

func (t *T) assert(msg string, oks ...bool) {
	for _, ok := range oks {
		if !ok {
			t.Log("\n" + string(debug.Stack()))
			t.Error(msg)
			t.Fatal("assert failed")
		}
	}
}

func (t *T) assert_set(set types.Set, err error) types.Set {
	t.assert_nil(err)
	return set
}

func (t *T) assert_nil(errors ...error) {
	for _, err := range errors {
		if err != nil {
			t.Log("\n" + string(debug.Stack()))
			t.Fatal(err)
		}
	}
}

func TestAddMarshalUnmarshalHas(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	set := NewSortedSet(10)
	items := make([]types.Int, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := types.Int(rand.Intn(10) + 1)
		items = append(items, item)
		t.assert_nil(set.Add(item))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("!set.Has(%v) ", i), set.Has(item))
	}
	marshal, unmarshal := types.IntMarshals()
	mset1 := NewMSortedSet(set, marshal, unmarshal)
	bytes, err := mset1.MarshalBinary()
	t.assert_nil(err)
	mset2 := &MSortedSet{MSorted: list.MSorted{MList: list.MList{MarshalItem: marshal, UnmarshalItem: unmarshal}}}
	t.assert_nil(mset2.UnmarshalBinary(bytes))
	set2 := mset2.SortedSet()
	for i, item := range items {
		t.assert(fmt.Sprintf("!set.Has(%v)", i), set2.Has(item))
	}
}

func TestAddHasDeleteRandom(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	set := NewSortedSet(10)
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
	for i := 0; i < 10; i++ {
		item, err := set.Random()
		t.assert_nil(err)
		t.assert(fmt.Sprintf("!set.Has(%v)", i), set.Has(item))
	}
	for _, item := range items {
		t.assert_nil(set.Delete(item))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), !set.Has(item))
	}
	_, err := set.Random()
	t.assert(fmt.Sprintf("err == nil"), err != nil)
	t.assert_nil(set.Add(types.Int(1)))
	item, err := set.Random()
	t.assert_nil(err)
	t.assert(fmt.Sprintf("item == 1"), item.Equals(types.Int(1)))
}

func TestAddHasDelete(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	set := NewSortedSet(10)
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
		t.assert_nil(set.Delete(item))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), !set.Has(item))
	}
}

func TestLess(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(3), types.Int(2), types.Int(3)})
	c := FromSlice([]types.Hashable{types.Int(1), types.Int(0), types.Int(3)})
	d := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	e := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)})
	small := FromSlice([]types.Hashable{types.Int(2), types.Int(4)})
	big := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(4)})
	t.assert("b < a", b.Less(a))
	t.assert("c < a", c.Less(a))
	t.assert("b < c", b.Less(c))
	t.assert("a !< d", a.Less(d) == false)
	t.assert("a !< e", e.Less(a) == false)
	t.assert("a < big", a.Less(big))
	t.assert("small < a", small.Less(a))
	t.assert("a !< small", a.Less(small) == false)
}

func TestEqualsHash(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(3), types.Int(2), types.Int(3)})
	c := FromSlice([]types.Hashable{types.Int(1), types.Int(0), types.Int(3)})
	d := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(3)})
	small := FromSlice([]types.Hashable{types.Int(2), types.Int(4)})
	empty := FromSlice([]types.Hashable{})
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

func TestUnion(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)})
	c := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(3), types.Int(4)})
	t.assert("a | b == c", t.assert_set(a.Union(b)).Equals(c))
	t.assert("b | a == c", t.assert_set(b.Union(a)).Equals(c))
	t.assert("a | b == c", t.assert_set(Union(a, b)).Equals(c))
	t.assert("b | a == c", t.assert_set(Union(b, a)).Equals(c))
}

func TestIntersect(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)})
	c := FromSlice([]types.Hashable{types.Int(1), types.Int(2)})
	d := FromSlice([]types.Hashable{types.Int(50), types.Int(20), types.Int(30), types.Int(40), types.Int(10)})
	e := FromSlice([]types.Hashable{})
	t.assert("a & b == c", t.assert_set(a.Intersect(b)).Equals(c))
	t.assert("b & a == c", t.assert_set(b.Intersect(a)).Equals(c))
	t.assert("a & d == e", t.assert_set(a.Intersect(d)).Equals(e))
	t.assert("d & a == e", t.assert_set(d.Intersect(a)).Equals(e))
}

func TestSubtract(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)})
	c := FromSlice([]types.Hashable{types.Int(0), types.Int(3)})
	d := FromSlice([]types.Hashable{types.Int(4)})
	t.assert("a - b == c", t.assert_set(a.Subtract(b)).Equals(c))
	t.assert("b - a == d", t.assert_set(b.Subtract(a)).Equals(d))
}

func TestOverlap(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)})
	t.assert("a & b != 0", a.Overlap(b))
	c := FromSlice([]types.Hashable{types.Int(5), types.Int(4)})
	t.assert("a & c == 0", !a.Overlap(c))
	d := FromSlice([]types.Hashable{types.Int(-2), types.Int(4)})
	t.assert("a & d == 0", !a.Overlap(d))
}

func TestSubsetSuperset(x *testing.T) {
	t := (*T)(x)
	a := FromSlice([]types.Hashable{types.Int(0), types.Int(1), types.Int(2), types.Int(3)})
	b := FromSlice([]types.Hashable{types.Int(1), types.Int(2), types.Int(4)})
	c := FromSlice([]types.Hashable{types.Int(1), types.Int(2)})
	t.assert("a not subset b", !a.Subset(b))
	t.assert("b not subset a", !b.Subset(a))
	t.assert("c subset a", c.Subset(a))
	t.assert("c subset b", c.Subset(b))
	t.assert("c proper subset a", c.ProperSubset(a))
	t.assert("c proper subset b", c.ProperSubset(b))
	t.assert("a subset a", a.Subset(a))
	t.assert("a not proper subset a", !a.ProperSubset(a))
	t.assert("a superset a", a.Superset(a))
	t.assert("a not proper superset a", !a.ProperSuperset(a))
	t.assert("a superset c", a.Superset(c))
	t.assert("b superset c", b.Superset(c))
	t.assert("a superset c", a.ProperSuperset(c))
	t.assert("b superset c", b.ProperSuperset(c))
}
