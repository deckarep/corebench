[![Build Status](https://travis-ci.org/deckarep/corebench.svg?branch=master)](https://travis-ci.org/deckarep/corebench)

# corebench
Benchmark utility that's intended to exercise benchmarks and how they scale with a large number of cores.

### TL;DR
How does your code scale and perform when running on high-core servers?

Let's find out:

### Example

Create some parallel benchmarks for your codebase
```go
// BenchmarkSomething utilizes the `b.RunParallel` feature of Go's benchmarking suite.
func BenchmarkSomething(b *testing.B) {
    templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
    b.RunParallel(func(pb *testing.PB) {
        var buf bytes.Buffer
        for pb.Next() {
            buf.Reset()
            templ.Execute(&buf, "World")
        }
    })
}
```

Run this command
```sh
./corebench do bench --git github.com/{user}/{your-code} --cpu=1,2,4,8,16,32,48
```

See this output
```sh
Provision 48-core droplet..
Droplet created: x.x.x.x -- running benchmark...

goos: linux
goarch: amd64
BenchmarkStoreRegular-2               	 2000000	       529 ns/op	      88 B/op	       0 allocs/op
BenchmarkStoreRegular-4               	 3000000	       452 ns/op	      61 B/op	       0 allocs/op
BenchmarkStoreRegular-8               	 5000000	       419 ns/op	      70 B/op	       0 allocs/op
BenchmarkStoreRegular-16               	 5000000	       419 ns/op	      70 B/op	       0 allocs/op
BenchmarkStoreRegular-32               	 5000000	       419 ns/op	      70 B/op	       0 allocs/op
BenchmarkStoreRegular-48               	 5000000	       419 ns/op	      70 B/op	       0 allocs/op
BenchmarkStoreSync-2                  	 1000000	      2250 ns/op	     179 B/op	       5 allocs/op
BenchmarkStoreSync-4                  	 1000000	      1807 ns/op	     179 B/op	       5 allocs/op
BenchmarkStoreSync-8                  	 1000000	      1637 ns/op	     179 B/op	       5 allocs/op
BenchmarkStoreSync-16                  	 1000000	      1637 ns/op	     179 B/op	       5 allocs/op
BenchmarkStoreSync-32                  	 1000000	      1637 ns/op	     179 B/op	       5 allocs/op
BenchmarkStoreSync-48                  	 1000000	      1637 ns/op	     179 B/op	       5 allocs/op
BenchmarkDeleteRegular-2              	10000000	       222 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeleteRegular-4              	10000000	       224 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeleteRegular-8              	10000000	       228 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeleteRegular-16              	10000000	       228 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeleteRegular-32              	10000000	       228 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeleteRegular-48              	10000000	       228 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeleteSync-2                   []
...
...
...

Benchmark completed...tearing down droplet.
```

### Here's what happens:
* A command like above will provision an on-demand high-performance computing server
* Installs Go, and clones your repository
* It will run Go's benchmark tooling against your repo and generate a comprehensive report demonstrating just how well your code scales across a large number of cores
* It will immediately decomission the computing resource so you only pay for a fraction of the cost

### Here's what you need:
* API/Credential access to at least one provider - Digital Ocean is the first provider to exist
* The ability to pay for your own computing resources for which ever providers you choose
* A source repo with comprehensive benchmarks to run against this suite

### Why benchmark on a large set of cores?
* Developers often benchmark their code on developer workstations, with a small number of cores
* Benchmarks on a small number of cores often times don't reflect the true nature of your application
* Sometimes an algorithm looks great on a few cores, but performance dramatically drops off when the core count gets higher
* A larger number of cores often times illustrates performances problems around:
* * Contention/locking bottlenecks
* * Cache incoherence
* * Parallelization overhead or lack of parallelization at all
* * Multi-threading overhead: starvation, race conditions, live-locks and priority inversion
* * The list goes on...

### Caution:
* This utility is unstable, API is in flux and is expected to change
* This project and its maintainers are NOT responsible for any monetary charges, overages, fees as a result of the auto-provision process during proper usage of the script, bugs in the script or because you decided to leave a cloud server running for months

