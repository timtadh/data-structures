package errors

import "fmt"

type ErrorFmter func(a ...interface{}) error

func NotFound(a ...interface{}) error {
    // return fmt.Errorf("Key '%v' was not found.", a...)
    return fmt.Errorf("Key was not found.")
}

func NotFoundInBucket(a ...interface{}) error {
    return fmt.Errorf("Key, '%v', was not in bucket when expected.", a...)
}

func InvalidKey(a ...interface{}) error {
    return fmt.Errorf("Key, '%v', is invalid, %s", a...)
}

func TSTError(a ...interface{}) error {
    return fmt.Errorf("Internal TST error - " + a[0].(string), a[1:]...)
}

func NegativeSize(a ...interface{}) error {
    return fmt.Errorf("Negative size")
}

func BpTreeError(a ...interface{}) error {
    return fmt.Errorf("Internal B+ Tree error - " + a[0].(string), a[1:]...)
}

var Errors map[string]ErrorFmter = map[string]ErrorFmter {
    "not-found":NotFound,
    "not-found-in-bucket":NotFoundInBucket,
    "invalid-key":InvalidKey,
    "tst-error":TSTError,
    "negative-size":NegativeSize,
    "bptree-error":BpTreeError,
}

