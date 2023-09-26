package async

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/reugn/async/internal/assert"
)

type syncList struct {
	list []string
	l    *OptimisticLock
}

func newSyncList() *syncList {
	return &syncList{
		l: NewOptimisticLock(),
	}
}

func (list *syncList) append(s string) {
	list.l.Lock()
	defer list.l.Unlock()
	list.list = append(list.list, s)
}

func (list *syncList) read() {
	ok := false
	for !ok {
		stamp := list.l.OptLock()
		_ = list.size()
		ok = list.l.OptUnlock(stamp)
	}
}

func (list *syncList) size() int {
	return len(list.list)
}

func TestList1(t *testing.T) {
	list := newSyncList()
	var wg sync.WaitGroup
	wg.Add(51)
	for i := 0; i < 50; i++ {
		go func() {
			for i := 0; i < 50; i++ {
				list.read()
			}
			wg.Done()
		}()
	}
	go func() {
		for i := 0; i < 50; i++ {
			list.append(strconv.Itoa(i))
			time.Sleep(time.Millisecond)
		}
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, 50, list.size())
}
