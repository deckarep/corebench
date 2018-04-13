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
	"os"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/deckarep/corebench/lib/ssh"
	"github.com/deckarep/corebench/lib/utility"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

const (
	doProviderInstanceNameFmt = "corebench-digitalocean-%s"
	cloudInitTemplate         = `
#cloud-config
runcmd:
  - echo "Setting up corebench for the first time..."
  - echo "Installing dependencies..."
  - apt-get -y install git
  - wget https://storage.googleapis.com/golang/${go-version}
  - tar -C /usr/local -xzf ${go-version}
  - export GOPATH=/root/go
  - mkdir -p $GOPATH
  - /usr/local/go/bin/go get ${git-repo}
  - touch $GOPATH/.core-init
  - echo "Finished corebench initialization"
`
	benchReadyScript     = "export GOPATH=/root/go && while [ ! -f $GOPATH/.core-init ]; do sleep 1; done"
	benchCommandTemplate = `cd $GOPATH/src/${git-repo} && /usr/local/go/bin/go version && /usr/local/go/bin/go test ${benchmem-setting}-cpu ${cpu-count} -bench=${bench-regex}`
	latestGoVersion      = "1.10.1"
	goVersion            = "go1.10.1.linux-amd64.tar.gz"
)

var (
	doDefaultPageOpts = &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	// doUSRegions are the US regions for digital ocean, let's start with this.
	doUSRegions = map[string]bool{
		"nyc1": true,
		"nyc2": true,
		"sfo1": true,
		"sfo2": true,
	}
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
		log.Info("No corebench droplets are provisioned on digitalocean")
		return nil
	}

	for _, d := range droplets {
		ip, _ := d.PublicIPv4()
		log.WithFields(
			log.Fields{
				"id":      d.ID,
				"name":    d.Name,
				"ip":      ip,
				"created": d.Created,
			},
		)
	}

	return nil
}

func filterUSRegions(regions []string) string {
	var results []string
	for _, reg := range regions {
		if _, ok := doUSRegions[reg]; ok {
			results = append(results, reg)
		}
	}
	return strings.Join(results, ", ")
}

func filterSizes(sizes []godo.Size, predicate func(slug string) bool) []godo.Size {
	var filteredSizes []godo.Size
	for _, s := range sizes {
		if predicate(s.Slug) {
			filteredSizes = append(filteredSizes, s)
		}
	}
	return filteredSizes
}

func displaySizes(category string, sizes []godo.Size) {
	const padding = 2
	const slugHdr = "Slug"
	const vcpuHdr = "VCpus"
	const mbHdr = "MB"
	const hourlyRateHdr = "$/HR"
	const availHdr = "Avail"
	const regHdr = "Regions"

	w := tabwriter.NewWriter(os.Stdout, 0, 8, padding, '\t', tabwriter.AlignRight)
	fmt.Println(category)
	fmt.Println()
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t", slugHdr, vcpuHdr, mbHdr, hourlyRateHdr, availHdr, regHdr))
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t",
		strings.Repeat("-", len(slugHdr)),
		strings.Repeat("-", len(vcpuHdr)),
		strings.Repeat("-", len(mbHdr)),
		strings.Repeat("-", len(hourlyRateHdr)),
		strings.Repeat("-", len(availHdr)),
		strings.Repeat("-", len(regHdr))))

	for _, sz := range sizes {
		avStatus := "yes"
		if !sz.Available {
			avStatus = "no"
		}
		fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%d\t%.2f\t%s\t%s\t", sz.Slug, sz.Vcpus, sz.Memory, sz.PriceHourly, avStatus, filterUSRegions(sz.Regions)))
	}
	w.Flush()
	fmt.Println()
}

type sizeTypes struct {
	optimized []godo.Size
	standard  []godo.Size
	flexible  []godo.Size
	count     int
}

func (p *DigitalOceanProvider) fetchSizes(ctx context.Context) (*sizeTypes, error) {
	sizes, _, err := p.client.Sizes.List(ctx, doDefaultPageOpts)
	if err != nil {
		return nil, err
	}

	optimizedSizes := filterSizes(sizes, func(slug string) bool {
		if strings.Contains(slug, "c-") {
			return true
		}
		return false
	})

	standardSizes := filterSizes(sizes, func(slug string) bool {
		if strings.Contains(slug, "s-") {
			return true
		}
		return false
	})

	flexSizes := filterSizes(sizes, func(slug string) bool {
		if !strings.Contains(slug, "s-") && !strings.Contains(slug, "c-") {
			return true
		}
		return false
	})

	st := sizeTypes{
		optimized: optimizedSizes,
		standard:  standardSizes,
		flexible:  flexSizes,
		count:     len(optimizedSizes) + len(standardSizes) + len(flexSizes),
	}
	return &st, nil
}

func (p *DigitalOceanProvider) Sizes(ctx context.Context) error {
	st, err := p.fetchSizes(ctx)
	if err != nil {
		log.Fatal("Error fetching sizes:", err)
	}

	fmt.Println()
	displaySizes("Standard Droplets:", st.standard)
	displaySizes("Flexible Droplets:", st.flexible)
	displaySizes("Optimized Droplets:", st.optimized)

	log.Infof("(%d) droplet sizes found\n", st.count)

	return nil
}

func (p *DigitalOceanProvider) Term(ctx context.Context, settings ProviderTermSettings) error {
	droplets, _, err := p.client.Droplets.ListByTag(ctx, "corebench", doDefaultPageOpts)
	if err != nil {
		return err
	}

	if len(droplets) == 0 {
		log.Info("No corebench droplets are alive to terminate on digitalocean")
		return nil
	}

	totalCount := len(droplets)
	termedCount := 0
	for _, droplet := range droplets {
		ip, _ := droplet.PublicIPv4()
		if settings.ShouldTerm(droplet.Name, ip) {
			log.Info("Terminating:", droplet.ID, droplet.Name, ip, "against match")
			_, err := p.client.Droplets.Delete(ctx, droplet.ID)
			if err != nil {
				log.WithField("id", droplet.ID).Warning("Failed to terminate droplet: need to retry or delete it manually or you will billed!!!")
				continue
			}
			termedCount++
		}
	}

	if termedCount == 0 {
		log.Warning("No instances were terminated that matched criteria")
	} else {
		log.Infof("Terminated (%d) droplets out of (%d) total droplets found\n", termedCount, totalCount)
	}

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
	benchCmd :=
		strings.Replace(benchCommandTemplate, "${git-repo}", settings.GitURL(), -1)
	benchCmd =
		strings.Replace(benchCmd, "${cpu-count}", settings.Cpus(), -1)
	benchCmd =
		strings.Replace(benchCmd, "${benchmem-setting}", settings.BenchMemString(), -1)
	benchCmd =
		strings.Replace(benchCmd, "${bench-regex}", settings.Regex(), -1)

	return fmt.Sprintf("%s && %s", benchReadyScript, benchCmd)
}

func (p *DigitalOceanProvider) Spinup(ctx context.Context, settings ProviderSpinSettings) error {

	//log.Fatal(p.processCloudInitTemplate(settings))

	//log.Fatal(p.processBenchCommandTemplate(settings))

	createRequest := &godo.DropletCreateRequest{
		Name: fmt.Sprintf(doProviderInstanceNameFmt, utility.NewInstanceID()),
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
		UserData: p.processCloudInitTemplate(settings),
	}

	if len(p.sshKeys) > 0 {
		var dropKeys []godo.DropletCreateSSHKey
		for _, k := range p.sshKeys {
			//println("adding ssh key: ", k)
			dropKeys = append(dropKeys, godo.DropletCreateSSHKey{
				Fingerprint: k,
			})
		}
		createRequest.SSHKeys = dropKeys
	}

	newDroplet, _, err := p.client.Droplets.Create(ctx, createRequest)
	if err != nil {
		log.Error("Failed to create droplet with err: %s\n", err)
		return err
	}

	log.Infof("Provisioning Droplet: %s ...\n", newDroplet.Name)
	log.Info("Slug:", createRequest.Size)
	log.Info("Region:", createRequest.Region)

	// fmt.Println(newDroplet.Name)
	// fmt.Println(newDroplet.ID)
	// fmt.Println(newDroplet.PublicIPv4())

	// Capture the droplets by id because we should delete at the end.
	// TODO: retry the deletes
	// TODO: recover panics and ensure the delete operation happens.
	var allDropletIds []int
	if !settings.LeaveRunning() {
		defer func() {
			p.cleanup(ctx, allDropletIds)
		}()
	}

	var chosenIP string
	const maxDialAttempts = 20

	// Spin wait - TODO: make this more graceful.
advance_to_ssh:
	for {
		droplets, _, err := p.client.Droplets.ListByTag(ctx, "corebench", doDefaultPageOpts)
		if err != nil {
			log.Fatal("Couldn't list all droplets with err: ", err)
		}

		// reset the ids on each outer iteration otherwise the list just grows.
		allDropletIds = nil
		for _, d := range droplets {
			allDropletIds = append(allDropletIds, d.ID)
			ip, _ := d.PublicIPv4()
			//fmt.Println(d.Name, d.ID, ip)
			// if we have an ip, start attempting...
			if ip != "" {
				chosenIP = ip
				err = ssh.PollSSH(chosenIP + ":22")
				if err == nil {
					//println("ssh dial success continuing!")
					break advance_to_ssh
				}
			}
		}
		time.Sleep(time.Second * 3)
	}

	log.Info("Droplet is provisioned and reachable at ip:", chosenIP)
	log.Info("Droplet benchmark starting momentarily...")
	fmt.Println()
	benchCmd := p.processBenchCommandTemplate(settings)
	err = ssh.ExecuteSSH(chosenIP, benchCmd)
	if err != nil {
		log.Fatalln("Failed to SSH: ", err)
	}

	return nil
}

func (p *DigitalOceanProvider) cleanup(ctx context.Context, ids []int) {
	for _, id := range ids {
		log.Info("Cleaning up droplet:", id)
		_, err := p.client.Droplets.Delete(ctx, id)
		if err != nil {
			log.Fatal("Failed to delete droplet: need to retry or delete it manually or you will billed!!!", id)
		}
	}
}
