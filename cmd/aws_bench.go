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
	"strings"
	"github.com/deckarep/corebench/pkg/providers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	keypair string = "corebench"
	instanceType string
)


// TODO: split out cores/instance types, map AMIs/regions/if HVM accordingly
//       add cmds for stackname, region, keypair (?)
func init() {
	awsBenchCmd.PersistentFlags().StringVarP(&instanceType,
		"instancetype", "", "t2.micro", "instance type to deploy (e.g. p2.xlarge)")
	awsBenchCmd.PersistentFlags().StringVarP(&regexString,
		"regex", "", "", "a regex to filter bench tests by")
	awsBenchCmd.PersistentFlags().BoolVarP(&leaveRunning,
		"leave-running", "", false, "indicates whether corebench should auto-terminate instance(s) on complete")
	awsBenchCmd.PersistentFlags().BoolVarP(&stat,
		"stat", "", false, "indicates whether corebench should generate benchstat summary")
	awsBenchCmd.PersistentFlags().BoolVarP(&benchMem,
		"benchmem", "", false, "indicates whether corebench include allocations just like the go tool")
	awsBenchCmd.PersistentFlags().StringVarP(&goVersion,
		"go", "", "1.10.1", "specifies the go version and must be a proper released version")
  awsCmd.AddCommand(awsBenchCmd)
}

var awsBenchCmd = &cobra.Command{
	Use:    "bench",
	Short:  "runs a remote benchmark on an aws instancetype",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if len(args) == 0 {
			log.WithField("example_repo", "github.com/foo/bar").Fatal("You must specify a git repo to bench")
		}

		settings := &providers.AwsSpinSettings{
			Git:              args[0],
			InstanceType:     instanceType,
			Benchmem:         benchMem,
			RegexFlag:        regexString,
			LeaveRunningFlag: leaveRunning,
			GoVersionFlag:    goVersion,
			CountFlag:        count,
			StatFlag:         stat,
		}

		provider := providers.NewAwsProvider()
		provider.SetKeys(strings.Split(keypair, ","))
		provider.Spinup(ctx, settings)
	},
}
