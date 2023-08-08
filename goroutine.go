package async

import (
	"bytes"
	"runtime"
	"strconv"
)

// GoroutineID returns the current goroutine id.
func GoroutineID() (uint, error) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	return uint(n), err
}
