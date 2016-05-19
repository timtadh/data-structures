package exc

import "testing"

import "github.com/timtadh/data-structures/test"

type myException struct {
	Exception
}

func TestTry(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Throwf("this is a test")
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
		Throwf("this is a test")
	}).Catch(&Exception{}, func(e Throwable) {
		t.Log("Caught", e)
	}).Error()
	t.Log(err)
	t.Log(err == nil)
	t.AssertNil(err)
}

func TestTryCatchMyExceptionFail(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Throwf("this is a test")
	}).Catch(&myException{}, func(e Throwable) {
		t.Log("Caught", e)
	}).Error()
	t.Assert(err != nil, "should not have caught")
}

func TestTryCatchMyExceptionSuccess(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Throw(&myException{*Errorf("this is a test of my Exception").Exception()})
	}).Catch(&myException{}, func(e Throwable) {
		t.Log("Caught", e)
	}).Error()
	t.Assert(err == nil, "should have caught %v", err)
}

func TestTryCatchReraise(x *testing.T) {
	t := (*test.T)(x)
	err := Try(func() {
		Throwf("this is a test")
	}).Catch(&Exception{}, func(e Throwable) {
		t.Log("Caught", e)
		Rethrow(e, Errorf("rethrow"))
	}).Exception()
	t.Log(err)
	t.Assert(err != nil, "wanted a non nil error got %v", err)
	t.Assert(len(err.Exc().Errors) == 2, "wanted 2 errors got %v", err)
}

func TestTryCatchFinally(x *testing.T) {
	t := (*test.T)(x)
	finally := false
	err := Try(func() {
		Throwf("this is a test")
	}).Catch(&Exception{}, func(e Throwable) {
		t.Log("Caught", e)
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
		Throwf("this is a test")
	}).Catch(&Exception{}, func(e Throwable) {
		t.Log("Caught", e)
		Rethrow(e, Errorf("rethrow"))
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
		Throwf("this is a test")
	}).Finally(func() {
		t.Log("finally")
		finally = true
	}).Error()
	t.Assert(err != nil, "err != nil,  %v", err)
	t.Assert(finally, "finally not run")
}

