package async

import (
	"fmt"
	"runtime/debug"
)

// GoroutineID returns current goroutine id
func GoroutineID() (uint, error) {
	var id uint
	var prefix string
	_, err := fmt.Sscanf(string(debug.Stack()), "%s %d", &prefix, &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
