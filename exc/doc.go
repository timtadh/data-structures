// Exceptions for Go as a Library
//
// Go (golang) lacks support for exceptions found in many other languages. There
// are good reasons for Go to not include exceptions. For instance, by making
// error handling explicit the programmer is forced to think concretely about
// the correct action to take. Fined grained control over the handling of errors
// using multiple return parameters is one of Go's strengths.
//
// However, there are cases where Go programs do not universally benefit from
// the explicit handling of errors. For instance, consider the following code:
//
//    func DoStuff(a, b, c interface{}) error {
//    	x, err := foo(a)
//    	if err != nil {
//    		return err
//    	}
//    	y, err := bar(b, x)
//    	if err != nil {
//    		return err
//    	}
//    	z, err := bax(c, x, y)
//    	if err != nil {
//    		return err
//    	}
//    	return baz(x, y, z)
//    }
//
// If Go had exceptions such code could be easily simplified:
//
//    func DoStuff(a, b, c interface{}) throws error {
//    	x := foo(a)
//    	y := bar(b, x)
//    	baz(x, y, bax(c, x, y)
//    }
//
// This library allow you to write go with exceptions and try-catch-finally
// blocks. It is not appropriate for all situations but can simplify some
// application code. Libraries and external APIs should continue to conform to
// the Go standard of returning error values.
//
// Here is an example of the DoStuff function where foo, bar and baz all throw
// exceptions instead of returning errors. (We will look at the case where they
// return errors that you want to turn into exceptions next). We want DoStuff to
// be an public API function and return an error:
//
//     func DoStuff(a, b, c interface{}) error {
//     	return Try(func() {
//     		x := foo(a)
//     		y := bar(b, x)
//     		baz(x, y, bax(c, x, y)
//     	}).Error()
//     }
//
// Now let's consider the case where we want to catch the exception log and
// reraise it:
//
//     func DoStuff(a, b, c interface{}) error {
//     	return exc.Try(func() {
//     		x := foo(a)
//     		y := bar(b, x)
//     		baz(x, y, bax(c, x, y)
//     	}).Catch(&exc.Exception{}, func(t exc.Throwable) {
//     		log.Log(t)
//     		exc.Rethrow(t, exc.Errorf("rethrow after logging"))
//     	}).Error()
//     }
//
// To be continued...
//
package exc
