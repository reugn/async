package benchmarks

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/reugn/async"
)

//go test -bench=. -benchmem -v

type syncList struct {
	list []string
	l    *async.OptimisticLock
	m    *sync.RWMutex
}

func newSyncList() *syncList {
	return &syncList{
		l: async.NewOptimisticLock(),
		m: &sync.RWMutex{},
	}
}

func (list *syncList) append(s string) {
	list.l.Lock()
	defer list.l.Unlock()
	list.list = append(list.list, s)
}

func (list *syncList) appendRW(s string) {
	list.m.Lock()
	defer list.m.Unlock()
	list.list = append(list.list, s)
}

func (list *syncList) readRW() {
	list.m.RLock()
	defer list.m.RUnlock()
	_ = list.size()
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

func benchOptimisticLock(b *testing.B, writes int) {
	list := newSyncList()
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(101)
		for i := 0; i < 100; i++ {
			go func() {
				for i := 0; i < 100; i++ {
					list.read()
				}
				wg.Done()
			}()
		}
		go func() {
			for i := 0; i < writes; i++ {
				list.append(strconv.Itoa(i))
				time.Sleep(time.Microsecond)
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

func benchRWLock(b *testing.B, writes int) {
	list := newSyncList()
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(101)
		for i := 0; i < 100; i++ {
			go func() {
				for i := 0; i < 100; i++ {
					list.readRW()
				}
				wg.Done()
			}()
		}
		go func() {
			for i := 0; i < writes; i++ {
				list.appendRW(strconv.Itoa(i))
				time.Sleep(time.Microsecond)
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

func BenchmarkRWLock1Write(b *testing.B) {
	benchRWLock(b, 1)
}

func BenchmarkOptimisticLock1Write(b *testing.B) {
	benchOptimisticLock(b, 1)
}

func BenchmarkRWLock100Writes(b *testing.B) {
	benchRWLock(b, 100)
}

func BenchmarkOptimisticLock100Writes(b *testing.B) {
	benchOptimisticLock(b, 100)
}
