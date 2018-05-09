/*
Open Source Initiative OSI - The MIT License (MIT):Licensing
The MIT License (MIT)
Copyright (c) 2018 Ralph Caraveo (deckarep@gmail.com)
Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated awscumentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to aws
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
	"strconv"
	"strings"
)

// TODO: clean up AwsTermSettings-related stuff, nlr

type AwsSpinSettings struct {
	Benchmem         bool
	CountFlag        int
  InstanceType     string
	Cpu              string
	Git              string
	GoVersionFlag    string
	LeaveRunningFlag bool
	RegexFlag        string
	StatFlag         bool
}

func (aws *AwsSpinSettings) GoVersion() string {
	return aws.GoVersionFlag
}

func (aws *AwsSpinSettings) InstanceTypeString() string {
	return aws.InstanceType
}

func (aws *AwsSpinSettings) BenchMemString() string {
	if aws.Benchmem {
		return "-benchmem "
	}
	return ""
}

func (aws *AwsSpinSettings) Count() int {
	return aws.CountFlag
}

func (aws *AwsSpinSettings) GitURL() string {
	return aws.Git
}

func (aws *AwsSpinSettings) Cpus() string {
	aws.Cpu = "2"

	return aws.Cpu
}

func (aws *AwsSpinSettings) MaxCpu() int {
	cpus := strings.Split(aws.Cpu, ",")
	var maxCpu int
	for _, c := range cpus {
		cpu, _ := strconv.Atoi(strings.TrimSpace(c))
		if cpu > maxCpu {
			maxCpu = cpu
		}
	}
	return maxCpu
}

func (aws *AwsSpinSettings) Regex() string {
	if aws.RegexFlag == "" {
		return "."
	}
	return aws.RegexFlag
}

func (aws *AwsSpinSettings) LeaveRunning() bool {
	return aws.LeaveRunningFlag
}

func (aws *AwsSpinSettings) Stat() bool {
	return aws.StatFlag
}

type AwsTermSettings struct {
	AllFlag  bool
	IPFlag   string
	NameFlag string
}

func (aws *AwsTermSettings) ShouldTerm(name, ip string) bool {
	if aws.AllFlag || aws.NameFlag == name || aws.IPFlag == ip {
		return true
	}
	return false
}
