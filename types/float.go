package types

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"unsafe"
)

type Float64 float64
type Float32 float32

func (self Float32) Equals(other Equatable) bool {
	if o, ok := other.(Float32); ok {
		return self == o
	} else {
		return false
	}
}

func (self Float32) Less(other Sortable) bool {
	if o, ok := other.(Float32); ok {
		return self < o
	} else {
		return false
	}
}

func (self Float32) Hash() int {
	f := float32(self)
	s := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&f)),
		Len:  4,
		Cap:  4,
	}
	b := (*[]byte)(unsafe.Pointer(s))
	h := fnv.New64a()
	_, err := h.Write(*b)
	if err != nil {
		// should never happen...
		panic(fmt.Errorf("could not write to hash %v", err))
	}
	return int(h.Sum64())
}

func (self Float64) Equals(other Equatable) bool {
	if o, ok := other.(Float64); ok {
		return self == o
	} else {
		return false
	}
}

func (self Float64) Less(other Sortable) bool {
	if o, ok := other.(Float64); ok {
		return self < o
	} else {
		return false
	}
}

func (self Float64) Hash() int {
	f := float64(self)
	s := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&f)),
		Len:  8,
		Cap:  8,
	}
	b := (*[]byte)(unsafe.Pointer(s))
	h := fnv.New64a()
	_, err := h.Write(*b)
	if err != nil {
		// should never happen...
		panic(fmt.Errorf("could not write to hash %v", err))
	}
	return int(h.Sum64())
}

