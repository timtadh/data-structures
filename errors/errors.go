package errors

import "fmt"

type ErrorFmter func(a ...interface{}) error

func NotFound(a ...interface{}) error {
    return fmt.Errorf("Key '%v' was not found.", a...)
}

func NotFoundInBucket(a ...interface{}) error {
    return fmt.Errorf("Key, '%v', was not in bucket when expected.", a...)
}

var Errors map[string]ErrorFmter = map[string]ErrorFmter {
    "not-found":NotFound,
    "not-found-in-bucket":NotFoundInBucket,
}

