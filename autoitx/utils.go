package autoitx

import (
	"syscall"
)

// IsErrSuccess checks if an "error" returned is actually the
// success code 0x0 "The operation completed successfully."
//
// This is the optimal approach since the error messages are
// localized depending on the OS language.
func IsErrSuccess(err error) bool {
	if errno, ok := err.(syscall.Errno); ok {
		if errno == 0 {
			return true
		}
	}
	return false
}
