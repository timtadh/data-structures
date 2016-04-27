package linked

import "testing"

import (
	"github.com/timtadh/data-structures/list"
)

func TestUniquePushPopSize(t *testing.T) {
	al := list.New(10)
	ll := NewUniqueDeque()
	for i := 0; i < 10; i++ {
		ll.Push(randhex(4))
		ll.Push(ll.Last())
		al.Push(ll.Last())
	}
	if !ll.Equals(al) {
		t.Fatalf("ll != al, %v != %v", ll, al)
	}
	t.Logf("ll %v", ll)
	t.Logf("al %v", al)
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

func TestUniqueFrontPushPopSize(t *testing.T) {
	al := list.New(10)
	ll := NewUniqueDeque()
	for i := 0; i < 10; i++ {
		ll.EnqueFront(randhex(4))
		ll.EnqueFront(ll.First())
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
