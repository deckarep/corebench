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

package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/deckarep/corebench/lib/providers"
	"github.com/spf13/cobra"
)

var (
	keys         string
	cpu          string
	leaveRunning bool
	benchMem     bool
	regexString  string
)

// Usage: ./corebench do bench -t=$TOKEN -k=$SSH_FINGERPRINT -git github.com/deckarep/golang-set
func init() {
	digitalOceanBenchCmd.PersistentFlags().StringVarP(&keys,
		"ssh-fp", "", "", "ssh fingerprints allow you to embed ssh keys via their MD5 fingerprint id, comma delimited list")
	digitalOceanBenchCmd.PersistentFlags().StringVarP(&cpu,
		"cpu", "c", "`nproc`", "cpu is a comma delimited list: -cpu=1,2,4,8 or -cpu=1-16")
	digitalOceanBenchCmd.PersistentFlags().StringVarP(&regexString,
		"regex", "", "", "a regex to filter bench tests by")
	digitalOceanBenchCmd.PersistentFlags().BoolVarP(&leaveRunning,
		"leave-running", "", false, "indicates whether corebench should auto-terminate instance(s) on complete")
	digitalOceanBenchCmd.PersistentFlags().BoolVarP(&benchMem,
		"benchmem", "", false, "indicates whether corebench include allocations just like the go tool")

	// TODO: -race flag (like go tooling)
	digitalOceanCmd.AddCommand(digitalOceanBenchCmd)
}

// benchCmd executes a remote benchmark.
var digitalOceanBenchCmd = &cobra.Command{
	Use:   "bench",
	Short: "runs a remote benchmark on a multi-core cloud resource from digitalocean",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if len(args) == 0 {
			log.Fatal("You must specificy a git repo to bench: https://github.com/foo/bar")
		}

		settings := &providers.DoSpinSettings{
			Git:              args[0],
			Cpu:              cpu,
			Benchmem:         benchMem,
			RegexFlag:        regexString,
			LeaveRunningFlag: leaveRunning,
		}

		provider := providers.NewDigitalOceanProvider(token)
		// Maybe this SetKeys api method isn't ideal.
		if keys != "" {
			provider.SetKeys(strings.Split(keys, ","))
		}

		fmt.Println()
		provider.Spinup(ctx, settings)
	},
}
