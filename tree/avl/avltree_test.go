package avl

import (
	"testing"

	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"

	trand "github.com/timtadh/data-structures/rand"
	"github.com/timtadh/data-structures/test"
	"github.com/timtadh/data-structures/types"
)

var rand *mrand.Rand

func init() {
	seed := make([]byte, 8)
	if _, err := crand.Read(seed); err == nil {
		rand = trand.ThreadSafeRand(int64(binary.BigEndian.Uint64(seed)))
	} else {
		panic(err)
	}
}

func randstr(length int) types.String {
	return types.String(test.RandStr(length))
}

func TestAvlPutHasGetRemove(t *testing.T) {

	type record struct {
		key   types.String
		value types.String
	}

	records := make([]*record, 400)
	var tree *AvlNode
	var err error
	var val interface{}
	var updated bool

	ranrec := func() *record {
		return &record{randstr(20), randstr(20)}
	}

	for i := range records {
		r := ranrec()
		records[i] = r
		tree, updated = tree.Put(r.key, types.String(""))
		if updated {
			t.Error("should have not been updated")
		}
		tree, updated = tree.Put(r.key, r.value)
		if !updated {
			t.Error("should have been updated")
		}
		if tree.Size() != (i + 1) {
			t.Error("size was wrong", tree.Size(), i+1)
		}
	}

	for _, r := range records {
		if has := tree.Has(r.key); !has {
			t.Error("Missing key")
		}
		if has := tree.Has(randstr(12)); has {
			t.Error("Table has extra key")
		}
		if val, err := tree.Get(r.key); err != nil {
			t.Error(err, val.(types.String), r.value)
		} else if !(val.(types.String)).Equals(r.value) {
			t.Error("wrong value")
		}
	}

	for i, x := range records {
		if tree, val, err = tree.Remove(x.key); err != nil {
			t.Error(err)
		} else if !(val.(types.String)).Equals(x.value) {
			t.Error("wrong value")
		}
		for _, r := range records[i+1:] {
			if has := tree.Has(r.key); !has {
				t.Error("Missing key")
			}
			if has := tree.Has(randstr(12)); has {
				t.Error("Table has extra key")
			}
			if val, err := tree.Get(r.key); err != nil {
				t.Error(err)
			} else if !(val.(types.String)).Equals(r.value) {
				t.Error("wrong value")
			}
		}
		if tree.Size() != (len(records) - (i + 1)) {
			t.Error("size was wrong", tree.Size(), (len(records) - (i + 1)))
		}
	}
}

func TestImmutableAvlPutHasGetRemove(t *testing.T) {

	type record struct {
		key   types.String
		value types.String
	}

	records := make([]*record, 400)
	var tree *ImmutableAvlNode
	var err error
	var val interface{}
	var updated bool

	ranrec := func() *record {
		return &record{randstr(20), randstr(20)}
	}

	for i := range records {
		r := ranrec()
		records[i] = r
		tree, updated = tree.Put(r.key, types.String(""))
		if updated {
			t.Error("should have not been updated")
		}
		tree, updated = tree.Put(r.key, r.value)
		if !updated {
			t.Error("should have been updated")
		}
		if tree.Size() != (i + 1) {
			t.Error("size was wrong", tree.Size(), i+1)
		}
	}

	for _, r := range records {
		if has := tree.Has(r.key); !has {
			t.Error("Missing key")
		}
		if has := tree.Has(randstr(12)); has {
			t.Error("Table has extra key")
		}
		if val, err := tree.Get(r.key); err != nil {
			t.Error(err, val.(types.String), r.value)
		} else if !(val.(types.String)).Equals(r.value) {
			t.Error("wrong value")
		}
	}

	for i, x := range records {
		if tree, val, err = tree.Remove(x.key); err != nil {
			t.Error(err)
		} else if !(val.(types.String)).Equals(x.value) {
			t.Error("wrong value")
		}
		for _, r := range records[i+1:] {
			if has := tree.Has(r.key); !has {
				t.Error("Missing key")
			}
			if has := tree.Has(randstr(12)); has {
				t.Error("Table has extra key")
			}
			if val, err := tree.Get(r.key); err != nil {
				t.Error(err)
			} else if !(val.(types.String)).Equals(r.value) {
				t.Error("wrong value")
			}
		}
		if tree.Size() != (len(records) - (i + 1)) {
			t.Error("size was wrong", tree.Size(), (len(records) - (i + 1)))
		}
	}
}

func TestIterators(t *testing.T) {
	var data []int = []int{
		1, 5, 7, 9, 12, 13, 17, 18, 19, 20,
	}
	var order []int = []int{
		6, 1, 8, 2, 4, 9, 5, 7, 0, 3,
	}

	test := func(tree types.TreeMap) {
		t.Logf("%T", tree)
		for j := range order {
			if err := tree.Put(types.Int(data[order[j]]), order[j]); err != nil {
				t.Error(err)
			}
		}

		j := 0
		for k, v, next := tree.Iterate()(); next != nil; k, v, next = next() {
			if !k.Equals(types.Int(data[j])) {
				t.Error("Wrong key")
			}
			if v.(int) != j {
				t.Error("Wrong value")
			}
			j += 1
		}

		j = 0
		for k, next := tree.Keys()(); next != nil; k, next = next() {
			if !k.Equals(types.Int(data[j])) {
				t.Error("Wrong key")
			}
			j += 1
		}

		j = 0
		for v, next := tree.Values()(); next != nil; v, next = next() {
			if v.(int) != j {
				t.Error("Wrong value")
			}
			j += 1
		}
	}
	test(NewAvlTree())
	test(NewImmutableAvlTree())
}

func BenchmarkAvlTree(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   types.String
		value types.String
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{randstr(20), randstr(20)}
	}

	for i := range records {
		records[i] = ranrec()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t := NewAvlTree()
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}

func BenchmarkImmutableAvlTree(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   types.String
		value types.String
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{randstr(20), randstr(20)}
	}

	for i := range records {
		records[i] = ranrec()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t := NewImmutableAvlTree()
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}
