/*
Open Source Initiative OSI - The MIT License (MIT):Licensing
The MIT License (MIT)
Copyright (c) 2018 Ralph Caraveo (deckarep@gmail.com)
Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

/*
These basic benchmarks don't test anything within corebench but are rather used
so corebench can deploy it's own code and run reasonable benchmarks that don't
take an incredibly long time to complete.
*/

import (
	"sync"
	"sync/atomic"
	"testing"
)

// BenchmarkAtomicIn increments a `uint64`...atomically.
func BenchmarkAtomicAddInc(b *testing.B) {
	var myVar uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddUint64(&myVar, 1)
		}
	})
}

// BenchmarkAtomicLoadStoreInc increments a `uint64` using an atomic swap.
func BenchmarkAtomicLoadStoreInc(b *testing.B) {
	var myVar uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			current := atomic.LoadUint64(&myVar)
			current++
			atomic.StoreUint64(&myVar, current)
		}
	})
}

// BenchmarkAtomicSwapInc increments a `uint64` using an atomic compare and swap.
func BenchmarkAtomicSwapInc(b *testing.B) {
	var myVar uint64
	// var retries uint64
	// var swaps uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
		finished:
			for {
				oldValue := atomic.LoadUint64(&myVar)
				newValue := oldValue + 1
				if atomic.CompareAndSwapUint64(&myVar, oldValue, newValue) {
					//atomic.AddUint64(&swaps, 1)
					break finished
				}
				//atomic.AddUint64(&retries, 1)
			}
		}
	})
	//fmt.Println("Failure rate", retries/swaps)
	//fmt.Println("Swaps: ", swaps, "Retries:", retries)
}

// BenchmarkMutexInc increments a `uint64` with a sync.Mutex lock.
func BenchmarkMutexInc(b *testing.B) {
	var mu sync.Mutex
	var myVar uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			myVar++
			mu.Unlock()
		}
	})
}

// BenchmarkMutexDeferInc increments a `uint64`...with a sync.Mutex lock and
// unlocks with a defer.
func BenchmarkMutexDeferInc(b *testing.B) {
	var mu sync.Mutex
	var myVar uint64

	f := func() {
		mu.Lock()
		defer mu.Unlock()
		myVar++
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f()
		}
	})
}

// BenchmarkRWMutexInc increments a `uint64` with a sync.RWMutex lock.
func BenchmarkRWMutexInc(b *testing.B) {
	var mu sync.RWMutex
	var myVar uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			myVar++
			mu.Unlock()
		}
	})
}

// BenchmarkRWMutexDeferInc increments a `uint64`...with a sync.RWMutex lock and
// unlocks with a defer.
func BenchmarkRWMutexDeferInc(b *testing.B) {
	var mu sync.RWMutex
	var myVar uint64

	f := func() {
		mu.Lock()
		defer mu.Unlock()
		myVar++
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f()
		}
	})
}
