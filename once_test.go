package async

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/reugn/async/internal/assert"
	"github.com/reugn/async/internal/util"
)

func TestOnce(t *testing.T) {
	var once Once[int]
	var count int

	for i := 0; i < 10; i++ {
		count, _ = once.Do(func() (int, error) {
			count++
			return count, nil
		})
	}
	assert.Equal(t, 1, count)
}

func TestOnce_Ptr(t *testing.T) {
	var once Once[*int]
	count := new(int)

	for i := 0; i < 10; i++ {
		count, _ = once.Do(func() (*int, error) {
			*count++
			return count, nil
		})
	}
	assert.Equal(t, util.Ptr(1), count)
}

func TestOnce_Concurrent(t *testing.T) {
	var once Once[*int32]
	var count atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, _ := once.Do(func() (*int32, error) {
				newCount := count.Add(1)
				return &newCount, nil
			})
			count.Store(*result)
		}()
	}
	wg.Wait()
	assert.Equal(t, 1, int(count.Load()))
}

func TestOnce_Panic(t *testing.T) {
	var once Once[*int]
	count := new(int)
	var err error

	for i := 0; i < 10; i++ {
		count, err = once.Do(func() (*int, error) {
			*count /= *count
			return count, nil
		})
	}
	assert.ErrorContains(t, err, "integer divide by zero")
}
