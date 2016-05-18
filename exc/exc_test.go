package exc

import "testing"

import "github.com/timtadh/data-structures/test"


func TestTry(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Throwf("test", "this is a test")
	}).Error()
	t.Assert(err != nil, "wanted a non nil error got %v", err)
}

func TestTryPropogate(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Try(func() {
			Throwf("test", "this is a test")
		}).Unwind()
	}).Error()
	t.Assert(err != nil, "wanted a non nil error got %v", err)
}

func TestTryCatch(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Throwf("test", "this is a test")
	}).Catch(&BaseException{}, func(e Exception) {
		t.Log("Caught", e.Name())
	}).Error()
	t.AssertNil(err)
}

func TestTryCatchReraise(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Throwf("test", "this is a test")
	}).Catch(&BaseException{}, func(e Exception) {
		t.Log("Caught", e.Name())
		Throwf("test2", "%v", e)
	}).Error()
	t.Assert(err != nil, "wanted a non nil error got %v", err)
	t.Assert(err.Name() == "test2", "wanted test2 got %v", err)
}

func TestTryCatchFinally(x *testing.T) {
	t := (*test.T)(x)
	finally := false
	err := Try(func() {
		Throwf("test", "this is a test")
	}).Catch(&BaseException{}, func(e Exception) {
		t.Log("Caught", e.Name())
	}).Finally(func() {
		t.Log("finally")
		finally = true
	}).Error()
	t.AssertNil(err)
	t.Assert(finally, "finally not run")
}

func TestTryCatchReraiseFinally(x *testing.T) {
	t := (*test.T)(x)
	finally := false
	err := Try(func() {
		Throwf("test", "this is a test")
	}).Catch(&BaseException{}, func(e Exception) {
		t.Log("Caught", e.Name())
		Throw(e)
	}).Finally(func() {
		t.Log("finally")
		finally = true
	}).Error()
	t.Assert(err != nil, "err != nil,  %v", err)
	t.Assert(finally, "finally not run")
}

func TestTryFinally(x *testing.T) {
	t := (*test.T)(x)
	finally := false
	err := Try(func() {
		Throwf("test", "this is a test")
	}).Finally(func() {
		t.Log("finally")
		finally = true
	}).Error()
	t.Assert(err != nil, "err != nil,  %v", err)
	t.Assert(finally, "finally not run")
}

