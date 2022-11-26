package async

import (
	"testing"
	"time"

	"github.com/reugn/async/internal"
)

func TestTask(t *testing.T) {
	task := NewTask(func() (string, error) {
		time.Sleep(1 * time.Second)
		return "ok", nil
	})
	res, err := task.Call().Join()

	internal.AssertEqual(t, "ok", res)
	internal.AssertEqual(t, err, nil)
}
