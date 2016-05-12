package test

import "testing"

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"os"
	"runtime/debug"
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

type T testing.T

func (t *T) Assert(ok bool, msg string, vars ...interface{}) {
	if !ok {
		t.Log("\n" + string(debug.Stack()))
		t.Fatalf(msg, vars...)
	}
}

func (t *T) AssertNil(errors ...error) {
	any := false
	for _, err := range errors {
		if err != nil {
			any = true
			t.Log("\n" + string(debug.Stack()))
			t.Error(err)
		}
	}
	if any {
		t.Fatal("assert failed")
	}
}

func RandSlice(length int) []byte {
	if urandom, err := os.Open("/dev/urandom"); err != nil {
		panic(err)
	} else {
		slice := make([]byte, length)
		if _, err := urandom.Read(slice); err != nil {
			panic(err)
		}
		urandom.Close()
		return slice
	}
	panic("unreachable")
}

func RandHex(length int) string {
	return hex.EncodeToString(RandSlice(length/2))
}

func RandStr(length int) string {
	return string(RandSlice(length))
}

