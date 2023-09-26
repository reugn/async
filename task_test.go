package async

import (
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

func TestTask(t *testing.T) {
	task := NewTask(func() (string, error) {
		time.Sleep(1 * time.Second)
		return "ok", nil
	})
	res, err := task.Call().Join()

	assert.Equal(t, "ok", res)
	assert.Equal(t, err, nil)
}
