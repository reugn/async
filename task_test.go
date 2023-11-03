package async

import (
	"errors"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/util"
)

func TestTaskSuccess(t *testing.T) {
	task := NewTask(func() (*string, error) {
		time.Sleep(10 * time.Millisecond)
		return util.Ptr("ok"), nil
	})
	res, err := task.Call().Join()

	assert.Equal(t, "ok", *res)
	assert.Equal(t, nil, err)
}

func TestTaskFailure(t *testing.T) {
	task := NewTask(func() (*string, error) {
		time.Sleep(10 * time.Millisecond)
		return nil, errors.New("error")
	})
	res, err := task.Call().Join()

	assert.Equal(t, nil, res)
	assert.ErrorContains(t, err, "error")
}
