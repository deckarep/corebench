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
	"log"

	"github.com/deckarep/corebench/lib/providers"
	"github.com/spf13/cobra"
)

var ()

// Usage: ./corebench do bench -t=$TOKEN -k=$SSH_FINGERPRINT -git github.com/deckarep/golang-set
func init() {
	// digitalOceanBenchCmd.PersistentFlags().StringVarP(&keys,
	// 	"keys", "k", "", "keys allow you to embed ssh keys via their MD5 fingerprint id, comma delimited list")
	// digitalOceanBenchCmd.PersistentFlags().StringVarP(&cpu,
	// 	"cpu", "c", "", "cpu is a comma delimited list: -cpu=1,2,4,8 or -cpu=1-16")
	// digitalOceanBenchCmd.PersistentFlags().StringVarP(&git,
	// 	"git", "g", "", "git path to a git repo to clone from, this must be publicly accessable")
	// digitalOceanBenchCmd.PersistentFlags().BoolVarP(&term,
	// 	"term", "", true, "indicates whether corebench should auto-terminate instance(s) on complete")

	// TODO: -benchmem flag (like go tooling)
	// TODO: -regex flag (like go tooling)
	// TODO: -race flag (like go tooling)

	digitalOceanCmd.AddCommand(digitalOceanListCmd)
}

var digitalOceanListCmd = &cobra.Command{
	Use:   "list",
	Short: "lists corebench resources provisioned in digitalocean that are currently alive",
	Run: func(cmd *cobra.Command, args []string) {
		provider := providers.NewDigitalOceanProvider(token)
		ctx := context.Background()
		err := provider.List(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}
