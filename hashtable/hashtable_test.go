package hashtable

import (
	"testing"

	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"

	trand "github.com/timtadh/data-structures/rand"
	"github.com/timtadh/data-structures/test"
	"github.com/timtadh/data-structures/tree/avl"
	. "github.com/timtadh/data-structures/types"
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

func randstr(length int) String {
	return String(test.RandStr(length))
}

func TestMake(t *testing.T) {
	NewHashTable(12)
}

func TestHashable(t *testing.T) {
	a := String("asdf")
	b := String("asdf")
	c := String("csfd")
	if !a.Equals(b) {
		t.Error("a != b")
	}
	if a.Hash() != b.Hash() {
		t.Error("hash(a) != hash(b)")
	}
	if a.Equals(c) {
		t.Error("a == c")
	}
	if a.Hash() == c.Hash() {
		t.Error("hash(a) == hash(c)")
	}
}

func TestPutHasGetRemove(t *testing.T) {

	type record struct {
		key   String
		value String
	}

	ranrec := func() *record {
		return &record{
			String(randstr(20)),
			String(randstr(20)),
		}
	}

	test := func(table MapOperable) {
		records := make([]*record, 400)
		for i := range records {
			r := ranrec()
			records[i] = r
			err := table.Put(r.key, String(""))
			if err != nil {
				t.Error(err)
			}
			err = table.Put(r.key, r.value)
			if err != nil {
				t.Error(err)
			}
			if table.Size() != (i + 1) {
				t.Error("size was wrong", table.Size(), i+1)
			}
		}

		for _, r := range records {
			if has := table.Has(r.key); !has {
				t.Error(table, "Missing key")
			}
			if has := table.Has(randstr(12)); has {
				t.Error("Table has extra key")
			}
			if val, err := table.Get(r.key); err != nil {
				t.Error(err)
			} else if !(val.(String)).Equals(r.value) {
				t.Error("wrong value")
			}
		}

		for i, x := range records {
			if val, err := table.Remove(x.key); err != nil {
				t.Error(err)
			} else if !(val.(String)).Equals(x.value) {
				t.Error("wrong value")
			}
			for _, r := range records[i+1:] {
				if has := table.Has(r.key); !has {
					t.Error("Missing key")
				}
				if has := table.Has(randstr(12)); has {
					t.Error("Table has extra key")
				}
				if val, err := table.Get(r.key); err != nil {
					t.Error(err)
				} else if !(val.(String)).Equals(r.value) {
					t.Error("wrong value")
				}
			}
			if table.Size() != (len(records) - (i + 1)) {
				t.Error("size was wrong", table.Size(), (len(records) - (i + 1)))
			}
		}
	}

	test(NewHashTable(64))
	test(NewLinearHash())
}

func TestIterate(t *testing.T) {

	test := func(table Map) {
		t.Logf("%T", table)
		for k, v, next := table.Iterate()(); next != nil; k, v, next = next() {
			t.Errorf("Should never reach here %v %v %v", k, v, next)
		}
		records := make(map[String]String)
		for i := 0; i < 100; i++ {
			k := randstr(8)
			v := randstr(8)
			records[k] = v
			err := table.Put(k, String(""))
			if err != nil {
				t.Error(err)
			}
			err = table.Put(k, v)
			if err != nil {
				t.Error(err)
			}
			if table.Size() != (i + 1) {
				t.Error("size was wrong", table.Size(), i+1)
			}
		}
		newrecs := make(map[String]String)
		for k, v, next := table.Iterate()(); next != nil; k, v, next = next() {
			if v2, has := records[k.(String)]; !has {
				t.Error("bad key in table")
			} else if !v2.Equals(v.(Equatable)) {
				t.Error("values don't agree")
			}
			newrecs[k.(String)] = v.(String)
		}
		if len(records) != len(newrecs) {
			t.Error("iterate missed records")
		}
		for k, v := range records {
			if v2, has := newrecs[k]; !has {
				t.Error("key went missing")
			} else if !v2.Equals(v) {
				t.Error("values don't agree")
			}
		}
	}
	test(NewHashTable(64))
	test(NewLinearHash())
	test(avl.NewAvlTree())
	test(avl.NewImmutableAvlTree())
}

func BenchmarkGoMap(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   string
		value string
	}

	records := make([]*record, 100)

	ranrec := func() *record {
		return &record{string(randstr(20)), string(randstr(20))}
	}

	for i := range records {
		records[i] = ranrec()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]string)
		for _, r := range records {
			m[r.key] = r.value
		}
		for _, r := range records {
			delete(m, r.key)
		}
	}
}

func BenchmarkHash(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   String
		value String
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
		t := NewHashTable(128)
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}

func BenchmarkMLHash(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   String
		value String
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
		t := NewLinearHash()
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}

func BenchmarkMLHashBetter(b *testing.B) {
	b.StopTimer()

	type record struct {
		key   String
		value String
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
		t := NewLinearHash()
		for _, r := range records {
			t.Put(r.key, r.value)
		}
		for _, _, next := t.Iterate()(); next != nil; _, _, next = next() {
		}
		for _, next := t.Keys()(); next != nil; _, next = next() {
		}
		for _, next := t.Values()(); next != nil; _, next = next() {
		}
		for _, next := t.Values()(); next != nil; _, next = next() {
		}
		for _, next := t.Values()(); next != nil; _, next = next() {
		}
		for _, r := range records {
			t.Remove(r.key)
		}
	}
}
