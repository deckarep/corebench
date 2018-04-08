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
	"github.com/spf13/cobra"
)

var (
	token string
	file  string
)

// Usage: ./corebench bench -t=$TOKEN -k=$SSH_FINGERPRINT -repo github.com/deckarep/golang-set
func init() {
	digitalOceanCmd.PersistentFlags().StringVarP(&token,
		"DO_PAT", "", "", "token is some cloud provider personal access token")
	digitalOceanCmd.PersistentFlags().StringVarP(&file,
		"file", "f", "", "file is a path to save benchmark results")
	RootCmd.AddCommand(digitalOceanCmd)
}

// benchCmd executes a remote benchmark.
var digitalOceanCmd = &cobra.Command{
	Use:   "do",
	Short: "is the digitalocean corebench toolkit",
}
