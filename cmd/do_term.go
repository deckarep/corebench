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

var (
	all  bool
	ip   string
	name string
)

func init() {
	digitalOceanTermCmd.PersistentFlags().BoolVarP(&all,
		"all", "", false, "indicates if you would like to terminate all instances")
	digitalOceanTermCmd.PersistentFlags().StringVarP(&name,
		"name", "n", "", "terminate instance by droplet name")
	digitalOceanTermCmd.PersistentFlags().StringVarP(&ip,
		"ip", "i", "", "terminate instance by ip address")

	digitalOceanCmd.AddCommand(digitalOceanTermCmd)
}

var digitalOceanTermCmd = &cobra.Command{
	Use:   "term",
	Short: "terminates corebench resources provisioned on digitalocean that are currently alive",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if ip == "" && name == "" && !all {
			log.Fatal("You must choose an option to terminate instances: either --all, --ip, or --name")
		}

		if all && (ip != "" || name != "") {
			log.Fatal("You cannot choose --all and specify an --ip or --name at the same time.")
		}

		if ip != "" && name != "" {
			log.Fatal("You can only terminate instances by their --ip or --name but not both.")
		}

		settings := &providers.DoTermSettings{
			AllFlag:  all,
			IPFlag:   ip,
			NameFlag: name,
		}

		provider := providers.NewDigitalOceanProvider(token)

		err := provider.Term(ctx, settings)
		if err != nil {
			log.Fatal(err)
		}
	},
}
