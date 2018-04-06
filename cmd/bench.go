/*
Open Source Initiative OSI - The MIT License (MIT):Licensing
The MIT License (MIT)
Copyright (c) 2017 Ralph Caraveo (deckarep@gmail.com)
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

	"github.com/deckarep/corebench/lib/providers"
	"github.com/spf13/cobra"
)

var (
	provider string
	token    string
	keys     string
	cpu      string
	gitRepo  string
	file     string
)

// Usage: ./corebench bench -t=$TOKEN -k=$SSH_FINGERPRINT -repo github.com/deckarep/golang-set
func init() {
	benchCmd.PersistentFlags().StringVarP(&provider,
		"provider", "p", "", "cloud provider to launch the remote resource: {DO - DigitalOcean}")
	benchCmd.PersistentFlags().StringVarP(&token,
		"token", "t", "", "token is some cloud provider personal access token")
	benchCmd.PersistentFlags().StringVarP(&keys,
		"keys", "k", "", "keys allow you to embed ssh keys via their MD5 fingerprint id, comma delimited list")
	benchCmd.PersistentFlags().StringVarP(&cpu,
		"cpu", "c", "", "cpu is a comma delimited list: -cpu=1,2,4,8 or -cpu=1-16")
	benchCmd.PersistentFlags().StringVarP(&gitRepo,
		"git", "g", "", "gitRepo a path to a git repo to clone from, this must be publicly accessable")
	benchCmd.PersistentFlags().StringVarP(&file,
		"file", "f", "", "file is a path to save benchmark results")

	// TODO: -benchmem flag (like go tooling)
	// TODO: -regex flag (like go tooling)
	// TODO: -race flag (like go tooling)

	RootCmd.AddCommand(benchCmd)
}

// benchCmd executes a remote benchmark.
var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "bench runs a remote benchmark on a multi-core cloud resource",
	Run: func(cmd *cobra.Command, args []string) {
		println("we're about to bench!")
		println("spinning up a DO droplet")
		println("here's the token: ", token)

		sshKeys := strings.Split(keys, ",")

		provider := providers.NewDigitalOceanProvider(token, sshKeys)
		ctx := context.Background()
		provider.Spinup(ctx)
	},
}
