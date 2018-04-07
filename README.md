# corebench
Benchmark utility that's intended to exercise benchmarks and how they scale with a large number of cores.

### TL;DR
How does your code scale and perform when running on high-core servers?

Let's find out:

```sh
./corebench bench -repo github.com/{user}/{your-code} --provider=DO --token=XXX --cpu=1,2,4,8,16,32,64,128
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

### Caution:
* This utility is unstable, API is in flux and is expected to change
* This project and its maintainers are NOT responsible for any monetary charges, overages, fees as a result of the auto-provision process during proper usage of the script, bugs in the script or because you decided to leave a cloud server running for months

