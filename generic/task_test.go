// go:build go1.18

package generic

import (
	"testing"
	"time"

	"github.com/reugn/async/internal"
)

func TestAsyncTask(t *testing.T) {
	task := newAsyncTask[string](func() (string, error) {
		time.Sleep(1 * time.Second)
		return "ok", nil
	})
	res, err := task.call().Get()

	internal.AssertEqual(t, "ok", res.(string))
	internal.AssertEqual(t, err, nil)
}
