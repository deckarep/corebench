# corebench
Benchmark utility

How does your code scale and perform when running on high-core servers?

Let's find out:

```sh
./corebench bench -repo github.com/deckarep/golang-set --provider=DO --token=XXX --cpu=1,2,4,8,16,32,64,128
```

Here's what happens:
1. A command like above will provision an on-demand high-performance computing server
2. It will install Go, and clone your repository
3. It will run Go's benchmark tooling against your repo and generate a comprehensive report demonstrating just how well your code scales across a large number of cores
4. It will immediately decomission the computing resource so you only pay for a fraction of the cost

Here's what you need:
1. API/Credential access to at least one provider - Digital Ocean is the first provider to exist
2. The ability to pay for your own computing resources for which ever providers you choose
3. Comprehensive benchmarks to run against this suite
