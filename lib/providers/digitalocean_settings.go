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

type DoSpinSettings struct {
	Git      string
	Cpu      string
	Benchmem bool
}

func (do *DoSpinSettings) BenchMemString() string {
	if do.Benchmem {
		return "-benchmem "
	}
	return ""
}

func (do *DoSpinSettings) GitURL() string {
	return do.Git
}

func (do *DoSpinSettings) Cpus() string {
	return do.Cpu
}

type DoTermSettings struct {
	AllFlag  bool
	IPFlag   string
	NameFlag string
}

func (do *DoTermSettings) All() bool {
	return do.AllFlag
}

func (do *DoTermSettings) ByIP() string {
	return do.IPFlag
}

func (do *DoTermSettings) ByName() string {
	return do.NameFlag
}
