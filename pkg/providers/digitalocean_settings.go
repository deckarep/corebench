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
	"strconv"
	"strings"
)

type DoSpinSettings struct {
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

func (do *DoSpinSettings) GoVersion() string {
	return do.GoVersionFlag
}

func (do *DoSpinSettings) InstanceTypeString() string {
	return do.InstanceType
}

func (do *DoSpinSettings) BenchMemString() string {
	if do.Benchmem {
		return "-benchmem "
	}
	return ""
}

func (do *DoSpinSettings) Count() int {
	return do.CountFlag
}

func (do *DoSpinSettings) GitURL() string {
	return do.Git
}

func (do *DoSpinSettings) Cpus() string {

	return do.Cpu
}

func (do *DoSpinSettings) MaxCpu() int {
	cpus := strings.Split(do.Cpu, ",")
	var maxCpu int
	for _, c := range cpus {
		cpu, _ := strconv.Atoi(strings.TrimSpace(c))
		if cpu > maxCpu {
			maxCpu = cpu
		}
	}
	return maxCpu
}

func (do *DoSpinSettings) Regex() string {
	if do.RegexFlag == "" {
		return "."
	}
	return do.RegexFlag
}

func (do *DoSpinSettings) LeaveRunning() bool {
	return do.LeaveRunningFlag
}

func (do *DoSpinSettings) Stat() bool {
	return do.StatFlag
}

type DoTermSettings struct {
	AllFlag  bool
	IPFlag   string
	NameFlag string
}

func (do *DoTermSettings) ShouldTerm(name, ip string) bool {
	if do.AllFlag || do.NameFlag == name || do.IPFlag == ip {
		return true
	}
	return false
}
