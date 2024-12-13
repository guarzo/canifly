package testutil

import (
	"fmt"
	"runtime"
)

// CaptureCallStack returns a formatted string of the call stack, skipping the specified number of frames.
// skip determines how many stack frames to skip (useful if you call this from within a helper function).
func CaptureCallStack(skip int) string {
	const size = 32
	var pcs [size]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var callStack string
	for {
		frame, more := frames.Next()
		callStack += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return callStack
}
