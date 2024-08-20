package async

import (
	"errors"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/util"
)

func TestTask_Success(t *testing.T) {
	task := NewTask(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "ok", nil
	})
	res, err := task.Call().Join()

	assert.Equal(t, "ok", res)
	assert.IsNil(t, err)
}

func TestTask_SuccessPtr(t *testing.T) {
	task := NewTask(func() (*string, error) {
		time.Sleep(10 * time.Millisecond)
		return util.Ptr("ok"), nil
	})
	res, err := task.Call().Join()

	assert.Equal(t, "ok", *res)
	assert.IsNil(t, err)
}

func TestTask_Failure(t *testing.T) {
	task := NewTask(func() (*string, error) {
		time.Sleep(10 * time.Millisecond)
		return nil, errors.New("error")
	})
	res, err := task.Call().Join()

	assert.IsNil(t, res)
	assert.ErrorContains(t, err, "error")
}
