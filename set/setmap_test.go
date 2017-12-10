package set

import "testing"

import (
	"fmt"
)

import (
	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/types"
)

func TestNewSetMap(t *testing.T) {
	s := types.Set(NewSetMap(hashtable.NewLinearHash()))
	t.Log(s.Size())
}

func TestSetMapAddHasDelete(x *testing.T) {
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
