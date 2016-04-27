package linked

import "testing"

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"os"
)

import (
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/types"
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

func randstr(length int) types.String {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		panic(err)
	} else {
		slice := make([]byte, length)
		if _, err := urandom.Read(slice); err != nil {
			panic(err)
		}
		urandom.Close()
		return types.String(slice)
	}
	panic("unreachable")
}

func randhex(length int) types.String {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		panic(err)
	} else {
		slice := make([]byte, length/2)
		if _, err := urandom.Read(slice); err != nil {
			panic(err)
		}
		urandom.Close()
		return types.String(hex.EncodeToString(slice))
	}
	panic("unreachable")
}

func TestEquals(t *testing.T) {
	al := list.New(10)
	ll := New()
	for i := 0; i < 10; i++ {
		ll.Push(randhex(4))
		al.Push(ll.Last())
	}
	if !ll.Equals(al) {
		t.Fatalf("ll != al, %v != %v", ll, al)
	}
	t.Logf("ll %v", ll)
	t.Logf("al %v", ll)
}

func TestPushPopSize(t *testing.T) {
	al := list.New(10)
	ll := New()
	for i := 0; i < 10; i++ {
		ll.Push(randhex(4))
		al.Push(ll.Last())
	}
	if !ll.Equals(al) {
		t.Fatalf("ll != al, %v != %v", ll, al)
	}
	for i := 0; i < 10; i++ {
		llItem, err := ll.Pop()
		if err != nil {
			t.Fatal(err)
		}
		alItem, err := al.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if !alItem.Equals(llItem) {
			t.Fatalf("llItem != alItem, %v != %v", llItem, alItem)
		}
	}
	if !ll.Equals(al) {
		t.Fatalf("ll != al, %v != %v", ll, al)
	}
}

func TestFrontPushPopSize(t *testing.T) {
	al := list.New(10)
	ll := New()
	for i := 0; i < 10; i++ {
		ll.EnqueFront(randhex(4))
		al.Insert(0, ll.First())
	}
	if !ll.Equals(al) {
		t.Fatalf("ll != al, %v != %v", ll, al)
	}
	for i := 0; i < 10; i++ {
		llItem, err := ll.DequeFront()
		if err != nil {
			t.Fatal(err)
		}
		alItem, err := al.Get(0)
		if err != nil {
			t.Fatal(err)
		}
		err = al.Remove(0)
		if err != nil {
			t.Fatal(err)
		}
		if !alItem.Equals(llItem) {
			t.Fatalf("llItem != alItem, %v != %v", llItem, alItem)
		}
	}
	if !ll.Equals(al) {
		t.Fatalf("ll != al, %v != %v", ll, al)
	}
}
