[![Build Status](https://travis-ci.org/deckarep/corebench.svg?branch=master)](https://travis-ci.org/deckarep/corebench)

# corebench
Benchmark utility that's intended to exercise benchmarks and how they scale with a large number of cores.

### TL;DR
How does your code scale and perform when running on high-core servers?

### Demo
[![asciicast](https://asciinema.org/a/l7q1KEz6knaH07z6zq7x2Ki8L.png)](https://asciinema.org/a/l7q1KEz6knaH07z6zq7x2Ki8L)


### Features
* First Provider: DigitalOcean up to 48 cores currently.
* --cpu flag supported: specify cpu delimited list
* --benchmem flag supported: capture allocations
* --count flag supported: multiple iterations of each benchmark
* --stat flag supported: executes [benchstat](https://github.com/golang/perf/tree/master/cmd/benchstat) analysis
* --regex flag supported: limits which benchmarks are run
* --leave-running flag supported: leaves a box running so user can log on
* sizes command: lists DigitalOcean instance sizes
* term command: terminates instances created by corebench
* list command: lists active corebench provisioned instances

### Usage
```go
# Install corebench
go get github.com/deckarep/corebench
```

Next, as DigitalOcean is the ONLY provider supported you must do the following 3 steps
* Sign up for a DigitalOcean account if not already a member
* Create a DigitalOcean Personal Access Token to be used for: --DO_PAT={token-here}
* Add your SSH public key to DigitalOcean for SSH access: --ssh-fp={ssh-md5-signature}

Run corebench:
```go
// Fetch instance sizes
./corebench do sizes --DO_PAT=$DO_PAT

// Run a benchmark
./corebench do bench github.com/{user}/{repo} [OPTIONS] --DO_PAT=$DO_PAT --ssh-fp=$SF

// List active instances
./corebench do list --DO_PAT=$DO_PAT

// Terminate instances created by corebench
./corebench do term --DO_PAT=$DO_PAT --all
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
* If you work in Go chances are you care about concurrent and parallel performance. If you don't care why are using Go at all?
* Developers often benchmark their code on developer workstations, with a small number of cores
* Benchmarks on a small number of cores often times don't reflect the true nature of your application
* Sometimes an algorithm looks great on a few cores, but performance dramatically drops off when the core count gets higher
* A larger number of cores often times illustrates performances problems around:
* * Contention/locking bottlenecks
* * Cache coherence issues
* * Parallelization overhead or lack of parallelization at all
* * Multi-threading overhead: starvation, race conditions, live-locks and priority inversion
* * The list goes on...

### F.A.Q.
 - Q: Why is 48 the max amount of cores this utility supports?
 - A: Because DigitalOcean is the first provider and that is their beefiest box.

 - Q: What happens to the server and the code after the benchmark completes?
 - A: The default behavior is the server is destroyed along with the code and benchmark data. There is a setting that allows you to leave the server running if you'd like to log in and inspect the results using the --leave-running flag.

 - Q: Why is DigitalOcean the first provider?
 - A: Easy, because their droplets fire up *FAST* allowing a quick feedback loop during development of this project.

 - Q: When you will you add Google Cloud, AWS, {other-provider} next?
 - A: Google Cloud is next because they offer per minute billing which is great to save money. I'm hoping the community can help me build other providers along with refactoring as necessary to align the API.

 - Q: Why did you build this tool?
 - A: Because I wanted to quick way to execute remote benchmarks on cloud servers that are beefy (large number of cores).

 - Q: Will you eventually support other languages?
 - A: Meybe. :) Did I mention this code is open source?

 - Q: Why is your code sloppy?
 - A: Because I'm currently in rapid prototype mode...don't worry it will get a lot better. Also through the power of open-source...yada, yada, yada.

 - Q: Doesn't this cost money everytime you need to fire up a benchmark?
 - A: Yes, yes it does...you have been warned.
 - A: If you want to test-drive, you can use a weak single core instance which costs like a penny an hour.

### Caution:
* This utility is in active development, API is in flux and is expected to change
* This project and its maintainers are NOT responsible for any monetary charges, overages, fees as a result of the auto-provision process during proper usage of the script, bugs in the script or because you decided to leave a cloud server running for months

