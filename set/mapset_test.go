package set

import "testing"

import (
	"fmt"
)

import (
	"github.com/timtadh/data-structures/types"
)

func TestMapSetAddHasDeleteRandom(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	a := NewMapSet(NewSortedSet(10))
	items := make([]*types.MapEntry, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := &types.MapEntry{types.Int(rand.Intn(1000)), types.Int(rand.Intn(1000))}
		for a.Has(item) {
			item = &types.MapEntry{types.Int(rand.Intn(1000)), types.Int(rand.Intn(1000))}
		}
		items = append(items, item)
		t.assert_nil(a.Add(item))
	}
	set := NewMapSet(NewSortedSet(10))
	t.assert_nil(set.Extend(a.Items()))
	t.assert("set == a", set.Equals(a))
	for i, item := range items {
		thing, err := set.Item(item)
		t.assert_nil(err)
		t.assert(fmt.Sprintf("i %v, thing == item", i), thing.Equals(item))
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), set.Has(item))
	}
	for _, item := range items {
		t.assert_nil(set.Delete(item))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item)", i), !set.Has(item))
	}
}

func TestMapSetMapFunc(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	set := types.Map(NewMapSet(NewSortedSet(10)))
	items := make([]*types.MapEntry, 0, SIZE)
	for i := 0; i < SIZE; i++ {
		item := &types.MapEntry{types.Int(rand.Intn(1000)), types.Int(rand.Intn(1000))}
		for set.Has(item.Key) {
			item = &types.MapEntry{types.Int(rand.Intn(1000)), types.Int(rand.Intn(1000))}
		}
		items = append(items, item)
		t.assert_nil(set.Put(item.Key, item.Value))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, !set.Has(item, %v)", i, item), set.Has(item.Key))
	}
	t.assert("size == 100", set.Size() == 100)
	for _, item := range items {
		rm, err := set.Remove(item.Key)
		t.assert_nil(err)
		t.assert("item == rm", item.Value.(types.Int).Equals(rm.(types.Int)))
	}
	for i, item := range items {
		t.assert(fmt.Sprintf("i %v, set.Has(item)", i), !set.Has(item.Key))
	}
}

func TestMapSetIterators(x *testing.T) {
	t := (*T)(x)
	SIZE := 100
	set := NewMapSet(NewSortedSet(10))
	for i := 0; i < SIZE; i++ {
		item := &types.MapEntry{types.Int(rand.Intn(1000)), types.Int(rand.Intn(1000))}
		t.assert_nil(set.Put(item.Key, item.Value))
	}
	item, items := set.Items()()
	key, keys := set.Keys()()
	for items != nil && keys != nil {
		if !item.(*types.MapEntry).Key.Equals(key) {
			t.Fatal("item != key")
		}
		item, items = items()
		key, keys = keys()
	}
	t.assert("uneven iteration", items == nil && keys == nil)
	item, items = set.Items()()
	val, vals := set.Values()()
	for items != nil && vals != nil {
		if !item.(*types.MapEntry).Value.(types.Int).Equals(val.(types.Int)) {
			t.Fatal("item != key")
		}
		item, items = items()
		val, vals = vals()
	}
	t.assert("uneven iteration", items == nil && keys == nil)
}
