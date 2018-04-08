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

package providers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/deckarep/corebench/lib/ssh"
	"github.com/deckarep/corebench/lib/utility"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

const (
	doProviderInstanceNameFmt = "corebench-digitalocean-%s"
)

var (
	doDefaultPageOpts = &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	goVersion         = "go1.9.1.linux-amd64.tar.gz"
	cloudInitTemplate = `
#cloud-config
runcmd:
  - echo "Setting up corebench for the first time..."
  - echo "Installing dependencies..."
  - apt-get -y install git
  - wget https://storage.googleapis.com/golang/${go-version}
  - tar -C /usr/local -xzf ${go-version}
  - git clone ${git-repo} /opt/corebench/${git-repo-last-path}
  - touch /opt/corebench/.core-init
  - echo "Finished corebench initialization"
`
	benchCommandTemplate = `while [ ! -f /opt/corebench/.core-init ]; do sleep 1; done && cd /opt/corebench/${git-repo-last-path} && /usr/local/go/bin/go test -cpu ${cpu-count} -bench=. -benchmem`
)

type DigitalOceanProvider struct {
	client       *godo.Client
	repoLastPath string
	// sshKeys can be optionally used to provision resources so you can log in and inspect the host.
	sshKeys []string
}

func NewDigitalOceanProvider(pat string) Provider {
	ts := NewDigitalOceanAuth(pat)
	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)
	return &DigitalOceanProvider{
		client: godo.NewClient(oauthClient),
	}
}

func (p *DigitalOceanProvider) SetKeys(keys []string) {
	p.sshKeys = keys
}

func (p *DigitalOceanProvider) List(ctx context.Context) error {
	droplets, _, err := p.client.Droplets.ListByTag(ctx, "corebench", doDefaultPageOpts)
	if err != nil {
		return err
	}

	if len(droplets) == 0 {
		fmt.Println("No corebench droplets are provisioned on digitalocean")
		return nil
	}

	for _, d := range droplets {
		ip, _ := d.PublicIPv4()
		fmt.Println(d.ID, d.Name, ip, d.Created)
	}

	return nil
}

func (p *DigitalOceanProvider) Term(ctx context.Context) error {
	droplets, _, err := p.client.Droplets.ListByTag(ctx, "corebench", doDefaultPageOpts)
	if err != nil {
		return err
	}

	if len(droplets) == 0 {
		fmt.Println("No corebench droplets to terminate on digitalocean")
		return nil
	}

	totalCount := len(droplets)
	termedCount := 0
	for _, droplet := range droplets {
		_, err := p.client.Droplets.Delete(ctx, droplet.ID)
		if err != nil {
			log.Println("Failed to terminate droplet: need to retry or delete it manually or you will billed!!!", droplet.ID)
			continue
		}
		termedCount++
	}

	fmt.Printf("Terminated (%d) droplets out of (%d) total droplets found\n", termedCount, totalCount)

	return nil
}

func (p *DigitalOceanProvider) processCloudInitTemplate(settings ProviderSpinSettings) string {

	p.repoLastPath = utility.GitPathLast(settings.GitURL())

	finalCloudTemplate :=
		strings.Replace(cloudInitTemplate, "${go-version}", goVersion, -1)
	finalCloudTemplate =
		strings.Replace(finalCloudTemplate, "${git-repo}", settings.GitURL(), -1)
	finalCloudTemplate =
		strings.Replace(finalCloudTemplate, "${git-repo-last-path}", p.repoLastPath, -1)

	return finalCloudTemplate
}

func (p *DigitalOceanProvider) processBenchCommandTemplate(settings ProviderSpinSettings) string {
	final :=
		strings.Replace(benchCommandTemplate, "${git-repo-last-path}", p.repoLastPath, -1)
	final =
		strings.Replace(final, "${cpu-count}", settings.Cpus(), -1)

	return final
}

func (p *DigitalOceanProvider) Spinup(ctx context.Context, settings ProviderSpinSettings) error {

	// Using this to show output before we run code.
	//log.Fatal(p.processBenchCommandTemplate(settings))

	createRequest := &godo.DropletCreateRequest{
		Name: "fake-droplet", //fmt.Sprintf(doProviderInstanceNameFmt, "7cfeebd"),
		// Costs: .01 penny to turn on (test with this)
		Region: "sfo2",
		Size:   "s-1vcpu-1gb",
		// Costs: .71 cents just to turn this beyatch on.
		//Region: "nyc1",
		//Size:   "48gb",
		Tags: []string{"corebench"},
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-14-04-x64",
		},
		// TODO: templatize UserData so things like go version can be swapped out.
		UserData: p.processCloudInitTemplate(settings),
	}

	if len(p.sshKeys) > 0 {
		var dropKeys []godo.DropletCreateSSHKey
		for _, k := range p.sshKeys {
			println("adding ssh key: ", k)
			dropKeys = append(dropKeys, godo.DropletCreateSSHKey{
				Fingerprint: k,
			})
		}
		createRequest.SSHKeys = dropKeys
	}

	println("About to create droplet")
	newDroplet, _, err := p.client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
		return err
	}
	println("Finished creating droplet")

	fmt.Println(newDroplet.Name)
	fmt.Println(newDroplet.ID)
	fmt.Println(newDroplet.PublicIPv4())

	// Capture the droplets because need to delete at the end.
	// TODO: retry the deletes
	// TODO: capture panics and ensure we delete even still
	// TODO: document the repo, user is responsible for charges
	var allDropletIds []int
	defer func() {
		for _, id := range allDropletIds {
			fmt.Println("Cleaning up droplet:", id)
			_, err := p.client.Droplets.Delete(ctx, id)
			if err != nil {
				log.Fatal("Failed to delete droplet: need to retry or delete it manually or you will billed!!!", id)
			}
		}
	}()

	var chosenIP string
	const maxDialAttempts = 20

	// Spin wait - TODO: make this more graceful.
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

advance_to_ssh:
	for {
		droplets, _, err := p.client.Droplets.ListByTag(ctx, "corebench", opt)
		if err != nil {
			log.Fatal("Couldn't list all droplets with err: ", err)
		}

		// reset the ids on each outer iteration otherwise the list just grows.
		allDropletIds = nil
		for _, d := range droplets {
			allDropletIds = append(allDropletIds, d.ID)
			ip, _ := d.PublicIPv4()
			fmt.Println(d.Name, d.ID, ip)
			// if we have an ip, start attempting...
			if ip != "" {
				chosenIP = ip
				err = ssh.PollSSH(chosenIP + ":22")
				//conn, err := net.DialTimeout("tcp", chosenIP+":22", time.Duration(time.Millisecond*500))
				if err == nil {
					println("ssh dial success continuing!")
					break advance_to_ssh
				}
			}
		}
		time.Sleep(time.Second * 3)
	}

	benchCmd := p.processBenchCommandTemplate(settings)
	err = ssh.ExecuteSSH(chosenIP, benchCmd)
	if err != nil {
		fmt.Println("Failed to SSH: ", err)
	}

	return nil
}
