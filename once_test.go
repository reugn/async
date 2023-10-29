package async

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/reugn/async/internal/assert"
)

func TestOnce(t *testing.T) {
	var once Once[int32]
	var count int32

	for i := 0; i < 10; i++ {
		count, _ = once.Do(func() (int32, error) {
			count++
			return count, nil
		})
	}
	assert.Equal(t, count, 1)
}

func TestOnceConcurrent(t *testing.T) {
	var once Once[int32]
	var count atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, _ := once.Do(func() (int32, error) {
				newCount := count.Add(1)
				return newCount, nil
			})
			count.Store(result)
		}()
	}
	wg.Wait()
	assert.Equal(t, int(count.Load()), 1)
}

func TestOncePanic(t *testing.T) {
	var once Once[int32]
	var count int32
	var err error

	for i := 0; i < 10; i++ {
		count, err = once.Do(func() (int32, error) {
			count /= count
			return count, nil
		})
	}
	assert.Equal(t, err.Error(), "recovered runtime error: integer divide by zero")
	assert.Equal(t, count, 0)
}
