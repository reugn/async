# async.OptimisticLock vs sync.RWMutex
```
BenchmarkRWLock1Write-8                     2000            790698 ns/op             161 B/op          2 allocs/op
BenchmarkOptimisticLock1Write-8            10000            108738 ns/op             393 B/op          5 allocs/op
BenchmarkRWLock100Writes-8                   200           7808822 ns/op            8993 B/op          2 allocs/op
BenchmarkOptimisticLock100Writes-8           300           6405564 ns/op            7857 B/op          5 allocs/op
```
