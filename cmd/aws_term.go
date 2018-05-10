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

	log "github.com/sirupsen/logrus"

	"github.com/deckarep/corebench/pkg/providers"
	"github.com/spf13/cobra"
)

// TODO: clean up AwsTermSettings-related stuff, nlr

func init() {
}

var (
	awsall       bool
	awsip        string
	instancename string
)

func init() {
	awsTermCmd.PersistentFlags().BoolVarP(&awsall,
		"all", "", false, "indicates if you would like to terminate all instances")
	awsCmd.AddCommand(awsTermCmd)
}

var awsTermCmd = &cobra.Command{
	Use:   "term",
	Short: "terminates corebench resources provisioned on aws that are currently alive",
	Run: func(cmd *cobra.Command, args []string) {

		settings := &providers.AwsTermSettings{
			AllFlag:  all,
			IPFlag:   ip,
			NameFlag: name,
		}

		provider := providers.NewAwsProvider()
		ctx := context.Background()
		err := provider.Term(ctx, settings)
		if err != nil {
			log.Fatal(err)
		}
	},
}
