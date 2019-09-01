package vw

// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Error errors returned by vowpal wabbig functions
type Error struct {
	message string
}

func (te *Error) Error() string {
	return te.message
}

func checkError(err vwError) *Error {
	if err.message != nil {
		defer C.free(unsafe.Pointer(err.message))
		return &Error{
			message: C.GoString(err.message),
		}
	}

	return nil
}
