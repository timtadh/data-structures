package test

import (
	"runtime/debug"
	"testing"

	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mrand "math/rand"

	trand "github.com/timtadh/data-structures/rand"
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
	slice := make([]byte, length)
	if _, err := crand.Read(slice); err != nil {
		panic(err)
	}
	return slice
}

func RandHex(length int) string {
	return hex.EncodeToString(RandSlice(length / 2))
}

func RandStr(length int) string {
	return string(RandSlice(length))
}
