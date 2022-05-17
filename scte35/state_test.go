/*
MIT License

Copyright 2016 Comcast Cable Communications Management, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

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
package scte35

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/Comcast/gots"
)

// All signal data generated with scte_creator: https://github.comcast.com/mniebu200/scte_creator
var poOpen1 = []byte{
	0x00, 0xfc, 0x30, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x1d, 0x02, 0x1b, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xff, 0x00, 0x00, 0x0a, 0xff, 0x50, 0x09, 0x05, 0x54, 0x65, 0x73, 0x74, 0x31, 0x34, 0x01,
	0x01, 0x00, 0x00, 0xff, 0x31, 0x22, 0x36,
}
var poClose1 = []byte{
	0x00, 0xfc, 0x30, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x0d, 0xbf, 0x24, 0x00, 0x1d, 0x02, 0x1b, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xbf, 0x09, 0x0a, 0x54, 0x65, 0x73, 0x74, 0x31, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x35, 0x01,
	0x01, 0x00, 0x00, 0xfc, 0x53, 0xaf, 0x44,
}

var poClose12 = []byte{
	0x00, 0xfc, 0x30, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x0a, 0xff, 0x50, 0x00, 0x1d, 0x02, 0x1b, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xbf, 0x09, 0x0a, 0x54, 0x65, 0x73, 0x74, 0x31, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x35, 0x01,
	0x02, 0x00, 0x00, 0xfd, 0x6f, 0xe8, 0xc7,
}
var poClose22 = []byte{
	0x00, 0xfc, 0x30, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x0d, 0xbf, 0x24, 0x00, 0x1d, 0x02, 0x1b, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xbf, 0x09, 0x0a, 0x54, 0x65, 0x73, 0x74, 0x31, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x35, 0x02,
	0x02, 0x00, 0x00, 0xfb, 0x7c, 0x2e, 0xe1,
}
var progStart = []byte{
	0x00, 0xfc, 0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x1a, 0x02, 0x18, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xbf, 0x09, 0x09, 0x50, 0x72, 0x6f, 0x67, 0x53, 0x74, 0x61, 0x72, 0x74, 0x10, 0x01, 0x01,
	0xf9, 0x43, 0xc2, 0x2f,
}
var progEnd = []byte{
	0x00, 0xfc, 0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x1a, 0x02, 0x18, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xbf, 0x09, 0x09, 0x50, 0x72, 0x6f, 0x67, 0x53, 0x74, 0x61, 0x72, 0x74, 0x11, 0x01, 0x01,
	0xfa, 0x95, 0x2c, 0xcf,
}
var progBreakaway = []byte{
	0x00, 0xfc, 0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x05, 0x7f, 0xa8, 0x00, 0x1a, 0x02, 0x18, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xbf, 0x09, 0x09, 0x50, 0x72, 0x6f, 0x67, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x13, 0x01, 0x01,
	0xf8, 0xd9, 0x85, 0xa7,
}
var progResumption = []byte{
	0x00, 0xfc, 0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x0d, 0xbf, 0x24, 0x00, 0x1a, 0x02, 0x18, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xbf, 0x09, 0x09, 0x50, 0x72, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x14, 0x01, 0x01,
	0xfb, 0x4f, 0x7b, 0x70,
}

var ppoStartSubsegments = []byte{
	0x00, 0xfc, 0x30, 0x3d, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x21, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x27, 0x02, 0x25, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xff, 0x00, 0x00, 0xa4, 0xcb, 0x80, 0x09, 0x0f, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x50, 0x4f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x34, 0x01, 0x03, 0x01, 0x02, 0xfa, 0x06, 0x95,
	0x8f,
}

var dpoStartSubsegments = []byte{
	0x00, 0xfc, 0x30, 0x40, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x21, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x2a, 0x02, 0x28, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xff, 0x00, 0x00, 0xa4, 0xcb, 0x80, 0x09, 0x12, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x6f, 0x72, 0x50, 0x4f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x36, 0x01, 0x01, 0x02, 0x02,
	0xfb, 0x2f, 0xe6, 0x7c,
}

var dpoFirstEndSubsegments = []byte{
	0x00, 0xfc, 0x30, 0x3c, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x21, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x26, 0x02, 0x24, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xff, 0x00, 0x00, 0xa4, 0xcb, 0x80, 0x09, 0x10, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x6f, 0x72, 0x50, 0x4f, 0x45, 0x6e, 0x64, 0x37, 0x01, 0x02, 0xfa, 0x60, 0x45, 0xdd,
}

var dpoSecondEndSubsegments = []byte{
	0x00, 0xfc, 0x30, 0x3c, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x21, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x26, 0x02, 0x24, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xff, 0x00, 0x00, 0xa4, 0xcb, 0x80, 0x09, 0x10, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x6f, 0x72, 0x50, 0x4f, 0x45, 0x6e, 0x64, 0x37, 0x02, 0x02, 0xfe, 0x65, 0x89, 0x11,
}

var ppoEndSubsegments = []byte{
	0x00, 0xfc, 0x30, 0x3b, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x21, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x25, 0x02, 0x23, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0xff, 0x00, 0x00, 0xa4, 0xcb, 0x80, 0x09, 0x0f, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x50, 0x4f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x35, 0x00, 0x00, 0xfa, 0x0b, 0x30, 0xf0,
}

// 0x10 ProgramStart
var scteDesc10Signal = "/DBFAAAAABeOAP/wBQb+WxdHNQAvAi1DVUVJACSqm3/9ABNP17wMGURJU0MyMjA1NjVfMDAyXzAxXzU3MUEtMDEQAQEeP8NH"
// 0x11 ProgramEnd
var scteDesc11Signal = "/DBAAAAAABeOAP/wBQb+bmg4eQAqAihDVUVJACSqm3+9DBlESVNDMjIwNTY1XzAwMl8wMV81NzFBLTA2EQEB/G8f5w=="
// 0x20 ChapterStart
var scteDesc20Signal = "/DBAAAAAAtaWAAAABQb/d5cVcgAqAihDVUVJ/////3//AALEr24BFGJyYXZvX0VQMDExMzQ0MTkwMTIyIAEA4BXDFw=="
// 0x21 ChapterEnd
var scteDesc21Signal = "/DBAAAAAAtaWAAAABQb/elu5JQAqAihDVUVJ/////3//AALEr24BFGJyYXZvX0VQMDExMzQ0MTkwMTIyIQEA1Kv4MQ=="
// 0x22 BreakStart
var scteDesc22Signal = "/DBAAAAAAtaWAAAABQb/elu5JQAqAihDVUVJ/////3//AAEgZ5IBFGJyYXZvX0VQMDExMzQ0MTkwMTIyIgIERTkzgA=="
// 0x23 BreakStart
var scteDesc23Signal = "/DBAAAAAAtaWAAAABQb/e3rvuQAqAihDVUVJ/////3//AAEgZ5IBFGJyYXZvX0VQMDExMzQ0MTkwMTIyIwIEzcrAYA=="
// 0x30 ProviderAdStart
var scteDesc30Signal = "/DBIAAAAAtaWAAAABQb/elu5JQAyAjBDVUVJ/////3//AAApJfYPHHVybjpuYmN1bmkuY29tOmJyYzo1MjE1NDkxMTkwAgS9D8Yj"
// 0x31 ProviderAdEnd
var scteDesc31Signal = "/DBIAAAAAtaWAAAABQb/elu5JQAyAjBDVUVJ/////3//AAApJfYPHHVybjpuYmN1bmkuY29tOmJyYzo1MjE1NDkxMTkwAgS9D8Yj"
// 0x34 ProviderPlacementOpportunityStart
var scteDesc34Signal = "/DBfAAAAAAAA///wBQb/iRp43QBJAhxDVUVJ6tzJ0n//AAEhrJQICAAFH4Lq3MnSNAIDAilDVUVJAAAAAH+/DBpWTU5VAWCXNGVv9BHsmxsOQM8vwoUB+olIUQEAAKn6Lds="
// 0x35 ProviderPlacementOpportunityEnd
var scteDesc35Signal = "/DBaAAAAAAAA///wBQb/ijwo4wBEAhdDVUVJ6tzJ0n+/CAgABR+C6tzJ0jUAAAIpQ1VFSQAAAAB/vwwaVk1OVQFglzRlb/QR7JsbDkDPL8KFAPqJSFEBAABl+tWe"
// 0x36 DistributorPlacementOpportunityStart
var scteDesc36Signal ="/DBLAAEs/S0UAP/wBQb+AAAAAAA1AjNDVUVJT///9X//AACky4AJH1NJR05BTDoyWURWeCtSKzlWc0FBQUFBQUFBQkFRPT02AAD9DXQ/"
//0x40 UnscheduledEventStart
var scteDesc40Signal = "/DB7AAFfzZzVAP/wBQb+AAAAAABlAlJDVUVJAABeUX+XDUMJIUJMQUNLT1VUOjI1dU5ZeEl3UXVXclI5WUxDR2I0Y2c9PQ4eY29tY2FzdDpsaW5lYXI6bGljZW5zZXJvdGF0aW9uQAAAAg9DVUVJAABeUX+XAABBAAB1H+6A"
// 0x44 ProviderAdBlockStart
var scteDesc44Signal = "/DCJAAGqtdGtAP/wBQb/AAAAAQBzAnFDVUVJQAFfVH//AAApMuANXQ8edXJuOmNvbWNhc3Q6YWx0Y29uOmFkZHJlc3NhYmxlDyx1cm46bWVybGluOmxpbmVhcjpzdHJlYW06Mzg0OTU3MzQ4MjgzNzU2NTE2MwkNUE86MTA3MzgzMTc2NEQAAAAA5o/+Nw=="
var scteDesc44Signal1 ="/DCfAAGab0HiAP/wBQb/AAAAAQCJAodDVUVJQAFfVH//AAApMuANcw8edXJuOmNvbWNhc3Q6YWx0Y29uOmFkZHJlc3NhYmxlDyx1cm46bWVybGluOmxpbmVhcjpzdHJlYW06Mzg0OTU3MzQ4MjgzNzU2NTE2MwkNUE86MTA3MzgzMTc2NAkKQlJFQUs6MTIzNAgIAAAAADeUEfJEAQMAAAu+e9A="
// 0x45 ProviderAdBlockEnd
var scteDesc45Signal = "/DCEAAGq3wk7AP/wBQb/AAAAAQBuAmxDVUVJQAFfVH+/DV0PHnVybjpjb21jYXN0OmFsdGNvbjphZGRyZXNzYWJsZQ8sdXJuOm1lcmxpbjpsaW5lYXI6c3RyZWFtOjM4NDk1NzM0ODI4Mzc1NjUxNjMJDVBPOjEwNzM4MzE3NjRFAAAAAHbwDu4="
// 0x50 NetworkStart
var scteDesc50Signal = "/DBQAAAAAAAAAABwBQb/LrIZ4QA6AhtDVUVJQAAAAH+fCgwUd4vl4/YAAAAAAABQAAACG0NVRUlAAAABf48KDBR3vd4u/wAAAAAAAFEAAF4PZmg="
// 0x51 NetworkEnd
var scteDesc51Signal = "DBQAAAAAAAAAABwBQb/LrIZ4QA6AhtDVUVJQAAAAH+fCgwUd4vl4/YAAAAAAABQAAACG0NVRUlAAAABf48KDBR3vd4u/wAAAAAAAFEAAF4PZmg="

func TestOutIn(t *testing.T) {
	st := NewState()
	open, e := NewSCTE35(poOpen1)
	if e != nil {
		t.Fatal("NewSCTE35(poOpen1) returned err:", e)
	}
	c, e := st.ProcessDescriptor(open.Descriptors()[0])
	if e != nil {
		t.Error("ProcessDescriptor of out returned unexpected err:", e)
	}
	if len(c) != 0 {
		t.Error("ProcessDescriptor returned closed signals when none should exist")
	}
	if len(st.Open()) != 1 {
		t.Error("Open() returned unexpected number of descriptors")
	} else if st.Open()[0] != open.Descriptors()[0] {
		t.Error("Open returned unexpected descriptor")
	}
	close, e := NewSCTE35(poClose1)
	if e != nil {
		t.Fatal("NewSCTE35(poClose1) returned err:", e)
	}
	c, e = st.ProcessDescriptor(close.Descriptors()[0])
	if e != nil {
		t.Error("ProcessDescriptor of in returned unexpected err:", e)
	}
	if len(c) != 1 {
		t.Error("ProcessDescriptor returned unexpected number of closed descriptors")
	} else if c[0] != open.Descriptors()[0] {
		t.Error("ProcessDescriptor returned unexpected close descriptor")
	}
	if len(st.Open()) != 0 {
		t.Error("Unexpectedly signals are still open")
	}
}

func TestOutInIn(t *testing.T) {
	st := NewState()
	open, e := NewSCTE35(poOpen1)
	if e != nil {
		t.Fatal("NewSCTE35(poOpen1) returned err:", e)
	}
	c, e := st.ProcessDescriptor(open.Descriptors()[0])
	if e != nil {
		t.Error("ProcessDescriptor of out returned unexpected err:", e)
	}
	if len(c) != 0 {
		t.Error("ProcessDescriptor returned closed signals when none should exist")
	}
	if len(st.Open()) != 1 {
		t.Error("Open() returned unexpected number of descriptors")
	} else if st.Open()[0] != open.Descriptors()[0] {
		t.Error("Open returned unexpected descriptor")
	}
	close1, e := NewSCTE35(poClose12)
	if e != nil {
		t.Fatal("NewSCTE35(poClose12) returned unexpected err:", e)
	}
	c, e = st.ProcessDescriptor(close1.Descriptors()[0])
	if e != nil {
		t.Error("ProcessDescriptor of in 1 returned unexpected err:", e)
	}
	if len(c) != 0 {
		t.Fatal("Close 1/2 closed open, which is not correct behavior")
	}
	if len(st.Open()) != 1 {
		t.Error("Open() returned unexpected number of descriptors")
	} else if st.Open()[0] != open.Descriptors()[0] {
		t.Error("Open returned unexpected descriptor")
	}
	close2, e := NewSCTE35(poClose22)
	if e != nil {
		t.Fatal("NewSCTE35(poClose22) returned unexpected err:", e)
	}
	c, e = st.ProcessDescriptor(close2.Descriptors()[0])
	if e != nil {
		t.Error("ProcessDescriptor of in 2 returned unexpected err:", e)
	}
	if len(c) != 1 {
		t.Error("ProcessDescriptor returned unexpected number of closed descriptors")
	} else if c[0] != open.Descriptors()[0] {
		t.Error("ProcessDescriptor returned unexpected close descriptor")
	}
	if len(st.Open()) != 0 {
		t.Error("Unexpectedly signals are still open")
	}
}

func TestDuplicateOut(t *testing.T) {
	st := NewState()
	open, e := NewSCTE35(poOpen1)
	if e != nil {
		t.Fatal("NewSCTE35(poOpen1) returned err:", e)
	}
	_, e = st.ProcessDescriptor(open.Descriptors()[0])
	if e != nil {
		t.Error("ProcessDescriptor of out returned unexpected err:", e)
	}
	_, e = st.ProcessDescriptor(open.Descriptors()[0])
	if e != gots.ErrSCTE35DuplicateDescriptor {
		t.Error("ProcessDescriptor of out returned unexpected err:", e)
	}
	if len(st.Open()) != 1 {
		t.Error("There should have been 1 open signal, instead num(open):", len(st.Open()))
	}
}

func TestOutOutIn(t *testing.T) {
	st := NewState()

	// 0x30
	padOpenSignalBytes, _ := base64.StdEncoding.DecodeString("/DBsAAH/7m1eAAKQBQb+sX6o+wBWAlRDVUVJAAAAJ3//AAApMuANQAwOQU1DTiBMMDAxMjM0NTYJCFBPOjEyMzQ1DiRiY2IxZGQ4ZS1kMzIzLTQ1ODktOWQ3OC1hM2QxMzYyYWJiYjYwAQGtc8xr")
	padOpenSignal, err := NewSCTE35(append([]byte{0x0}, padOpenSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	_, err = st.ProcessDescriptor(padOpenSignal.Descriptors()[0])
	if err != nil {
		t.Error("ProcessDescriptor of out returned unexpected err:", err)
	}
	if len(st.Open()) != 1 {
		t.Error("There should have been 1 open signal, instead num(open):", len(st.Open()))
	}

	// 0x34
	ppoStartSignalBytes, _ := base64.StdEncoding.DecodeString("/DA0AAG3wyKO///wBQb/yRLo+gAeAhxDVUVJQAAAZH/PAAEOFNwICAAAAAArJkYdNAAAMcFkEw==")
	ppoStartSignal, err := NewSCTE35(append([]byte{0x0}, ppoStartSignalBytes...))
	if err != nil {
		t.Fatal("NewSCTE35(ppoStartSignal) return err:", err)
	}

	c, err := st.ProcessDescriptor(ppoStartSignal.Descriptors()[0])
	if err != nil {
		t.Error("ProcessDescriptor of out 2 returned unexpected err:", err)
	}
	if len(c) != 1 {
		t.Error("ppoStartSignal unexpectedly did not close the first signal")
	}
	if len(st.Open()) != 1 {
		t.Error("There should have been 1 open signal, instead num(open):", len(st.Open()))
	}
	// now pass through the close signal and check
	// 0x35 - close the 0x34
	closeSignalBytes, _ := base64.StdEncoding.DecodeString("/DAvAAG2uS3c///wBQb/yiD8XAAZAhdDVUVJQAAAZH+fCAgAAAAAKyZGHTUAAMzqBnE=")
	close, err := NewSCTE35(append([]byte{0x0},closeSignalBytes...))
	if err != nil {
		t.Fatal("NewSCTE35(closeSignalBytes) return unexpected err:", err)
	}

	c, err = st.ProcessDescriptor(close.Descriptors()[0])
	if err != nil {
		t.Error("Processing first desc of close returned unexpected err:", err)
	}

	if len(c) != 1 {
		t.Error("First desc unexpectedly did not close inner out")
	}
	if len(st.Open()) != 0 {
		t.Error("There should have been 0 open signals, instead num(open):", len(st.Open()))
	}
}

func TestOutOut(t *testing.T) {
	state := NewState()
	// 0x30 - Out signal 1
	padOpenSignalBytes, _ := base64.StdEncoding.DecodeString("/DBsAAH/7m1eAAKQBQb+sX6o+wBWAlRDVUVJAAAAJ3//AAApMuANQAwOQU1DTiBMMDAxMjM0NTYJCFBPOjEyMzQ1DiRiY2IxZGQ4ZS1kMzIzLTQ1ODktOWQ3OC1hM2QxMzYyYWJiYjYwAQGtc8xr")
	padOpenSignal, err := NewSCTE35(append([]byte{0x0}, padOpenSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	_, err = state.ProcessDescriptor(padOpenSignal.Descriptors()[0])
	if err != nil {
		t.Error("ProcessDescriptor of out returned unexpected err:", err)
	}
	if len(state.Open()) != 1 {
		t.Error("There should have been 1 open signal, instead num(open):", len(state.Open()))
	}

	// 0x36 - event_id: 1342177266 - seg_num: 0 - seg_expected: 0
	// Should close the earlier Out
	secondOutSignalBytes, _ := base64.StdEncoding.DecodeString("/DBLAAF0QXOWAP/wBQb+AAAAAAA1AjNDVUVJT///8n//AACky4AJH1NJR05BTDozR1NOanl3cE1sb0FBQUFBQUFBQkFRPT02AAA9gIK2")
	secondOutSignal, err := NewSCTE35(append([]byte{0x0}, secondOutSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(secondOutSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 1 {
		t.Errorf("One event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}
}

//Creates a PPO start, DPO start, DPO end, DPO end, PPO end.
func TestSubsegments(t *testing.T) {
	state := NewState()

	// 0x34
	ppoStart, err := NewSCTE35(ppoStartSubsegments)
	if err != nil {
		t.Fatal("NewSCTE35(ppoStartSubsegments) return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(ppoStart.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x36
	dpoStart, err := NewSCTE35(dpoStartSubsegments)
	if err != nil {
		t.Fatal("NewSCTE35(dpoStartSubsegments) return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(dpoStart.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 2 {
		t.Errorf("There should be two open signals (%d)", len(state.Open()))
	}

	// 0x37
	dpoFirstEnd, err := NewSCTE35(dpoFirstEndSubsegments)
	if err != nil {
		t.Fatal("NewSCTE35(dpoFirstEndSubsegments) return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(dpoFirstEnd.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 2 {
		t.Errorf("There should be two open signals (%d)", len(state.Open()))
	}

	// 0x37
	dpoSecondEnd, err := NewSCTE35(dpoSecondEndSubsegments)
	if err != nil {
		t.Fatal("NewSCTE35(dpoSecondEndSubsegments) return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(dpoSecondEnd.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 1 {
		t.Errorf("One event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x35
	ppoEnd, err := NewSCTE35(ppoEndSubsegments)
	if err != nil {
		t.Fatal("NewSCTE35(ppoEndSubsegments) return err:", err.Error())
	}
	closed, err = state.ProcessDescriptor(ppoEnd.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 1 {
		t.Errorf("One event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should be no open signals (%d)", len(state.Open()))
	}
}

// Test the logic for when a closing IN signal occurs after another OUT signal.
// 0x36 -> 0x37 (1/3) -> 0x37 (2/3) -> 0x30 -> 0x37 (3/3)
// End state should be just 0x30.
func TestOutInInOutIn(t *testing.T) {
	state := NewState()

	// 0x36 - event_id:0 - seg_num: 0 - seg_expected: 0
	outSignalBytes, _ := base64.StdEncoding.DecodeString("/DBLAAFztMbuAP/wBQb+AAAAAAA1AjNDVUVJAAAAAH//AACky4AJH1NJR05BTDozR1NOajNnb01sb0FBQUFBQUFBQkFRPT02AADO/OgI")
	outSignal, err := NewSCTE35(append([]byte{0x0}, outSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(outSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x37 - event_id: 0 - seg_num: 1 - seg_expected: 3
	firstInSignalBytes, _ := base64.StdEncoding.DecodeString("/DBGAAF0ByyuAP/wBQb+AAAAAAAwAi5DVUVJAAAAAH+/CR9TSUdOQUw6M0dTTmozZ29NbG9BQUFBQUFBQUJBZz09NwEDfTeSVQ==")
	firstInSignal, err := NewSCTE35(append([]byte{0x0}, firstInSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(firstInSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x37 - event_id: 0 - seg_num: 2 - seg_expected: 3
	secondInSignalBytes, _ := base64.StdEncoding.DecodeString("/DBGAAF0MF+OAP/wBQb+AAAAAAAwAi5DVUVJAAAAAH+/CR9TSUdOQUw6M0dTTmozZ29NbG9BQUFBQUFBQUJBdz09NwIDvefEqg==")
	secondInSignal, err := NewSCTE35(append([]byte{0x0}, secondInSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(secondInSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x30 Another Out
	padOpenSignalBytes, _ := base64.StdEncoding.DecodeString("/DBsAAH/7m1eAAKQBQb+sX6o+wBWAlRDVUVJAAAAJ3//AAApMuANQAwOQU1DTiBMMDAxMjM0NTYJCFBPOjEyMzQ1DiRiY2IxZGQ4ZS1kMzIzLTQ1ODktOWQ3OC1hM2QxMzYyYWJiYjYwAQGtc8xr")
	padOpenSignal, err := NewSCTE35(append([]byte{0x0}, padOpenSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(padOpenSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 2 {
		t.Errorf("There should be two open signals (%d)", len(state.Open()))
	}

	// 0x37 = event_id: 0 - seg_num: 3 - seg_expected: 3
	// This closed the 0x30 and the 0x36
	thirdInSignalBytes, _ := base64.StdEncoding.DecodeString("/DBGAAF0WZJuAP/wBQb+AAAAAAAwAi5DVUVJAAAAAH+/CR9TSUdOQUw6M0dTTmozZ29NbG9BQUFBQUFBQUJCQT09NwMDFkn/Gw==")
	thirdInSignal, err := NewSCTE35(append([]byte{0x0}, thirdInSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(thirdInSignal.Descriptors()[0])
	if err != nil {
		t.Error("ProcessDescriptor of out returned unexpected err:", err)
	}
	if len(closed) != 2 {
		t.Errorf("Two events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should be no open signal (%d)", len(state.Open()))
	}
}

// Test the logic for when we recieve a VSS signal
func TestVSS(t *testing.T) {
	state := NewState()

	outSignalBytes, _ := base64.StdEncoding.DecodeString("/DB7AAFe1ms7AP/wBQb+AAAAAABlAlJDVUVJAABeT3+XDUMJIUJMQUNLT1VUOlEza2dMYmx4UzlhTmh4S24wY1N0MlE9PQ4eY29tY2FzdDpsaW5lYXI6bGljZW5zZXJvdGF0aW9uQAAAAg9DVUVJAABeT3+XAABBAAC9uy+v")
	outSignal, err := NewSCTE35(append([]byte{0x0}, outSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	// 0x40
	closed, err := state.ProcessDescriptor(outSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}
	if state.Open()[0].TypeID() != SegDescUnscheduledEventStart {
		t.Errorf("Expected segmentation_type_id 0x40 but got %x", state.Open()[0].TypeID())
	}
	if state.Open()[0].EventID() != 24143 {
		t.Errorf("Expected event_id 24143 but got %d", state.Open()[0].EventID())
	}

	// 0x41
	closed, err = state.ProcessDescriptor(outSignal.Descriptors()[1])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 1 {
		t.Errorf("1 event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should not be any open signals (%d)", len(state.Open()))
	}
}

// Test the logic for when we recieve a VSS signal
// with same signalID and eventID as the earlier one
// and that we drop it.
func TestVSSSameSignalIdBackToBack(t *testing.T) {
	state := NewState()

	vss_1 := "/DB7AAH//y+UAP/wBQb+AAAAAABlAlJDVUVJAAALkn+XDUMJIUJMQUNLT1VUOnZUaDZqMUNDRFZ3QUFBQUFBQUFCQVE9PQ4eY29tY2FzdDpsaW5lYXI6bGljZW5zZXJvdGF0aW9uQAAAAg9DVUVJAAALkn+XAABBAABupj9l" //PTS: 8589881236 , EventID: 2962, SignalID:vTh6j1CCDVwAAAAAAAABAQ==
	vss_2 := "/DB7AAAAAOBaAP/wBQb+AAAAAABlAlJDVUVJAAALkn+XDUMJIUJMQUNLT1VUOnZUaDZqMUNDRFZ3QUFBQUFBQUFCQVE9PQ4eY29tY2FzdDpsaW5lYXI6bGljZW5zZXJvdGF0aW9uQAAAAg9DVUVJAAALkn+XAABBAAAuynMR" //PTS:57434, EventID:2962, SignalID:vTh6j1CCDVwAAAAAAAABAQ==

	outSignalBytes1, _ := base64.StdEncoding.DecodeString(vss_1)
	outSignal1, err := NewSCTE35(append([]byte{0x0}, outSignalBytes1...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	// 0x40
	closed1, err := state.ProcessDescriptor(outSignal1.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed1) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed1))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}
	if state.Open()[0].TypeID() != SegDescUnscheduledEventStart {
		t.Errorf("Expected segmentation_type_id 0x40 but got %x", state.Open()[0].TypeID())
	}
	if state.Open()[0].EventID() != 2962 {
		t.Errorf("Expected event_id 2962 but got %d", state.Open()[0].EventID())
	}

	// 0x41
	closed1, err = state.ProcessDescriptor(outSignal1.Descriptors()[1])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed1) != 1 {
		t.Errorf("1 event should have been closed (%d were)", len(closed1))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should not be any open signals (%d)", len(state.Open()))
	}

	// Send another VSS signal with the same signalID and same eventID
	// Check that we drop the signal.
	outSignalBytes2, _ := base64.StdEncoding.DecodeString(vss_2)
	outSignal2, err := NewSCTE35(append([]byte{0x0}, outSignalBytes2...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	// 0x40
	closed2, err := state.ProcessDescriptor(outSignal2.Descriptors()[0])
	if err != gots.ErrSCTE35DuplicateDescriptor {
		t.Errorf("ProcessDescriptor should have dropped this as it is a duplicate descriptor: %s", err.Error())
	}
	if len(state.Open()) != 0 {
		t.Errorf("This signal should have been dropped, there should be 0 open signals (%d)", len(state.Open()))
	}
	// 0x41
	closed2, err = state.ProcessDescriptor(outSignal2.Descriptors()[1])
	if err != gots.ErrSCTE35MissingOut {
		t.Errorf("0x40 was dropped so 0x41 should have been dropped as there was no matching out: %s", err.Error())
	}
	if len(closed2) != 0 {
		t.Errorf("0 events should have been closed (%d were)", len(closed2))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should not be any open signals (%d)", len(state.Open()))
	}
}

func printState(s State, header string) {
	fmt.Printf("\n%s\n", header)
	for _, open := range s.Open() {
		fmt.Printf("%X - %d - (%d/%d) - %s\n", open.TypeID(), open.EventID(), open.SegmentNumber(), open.SegmentsExpected(), base64.StdEncoding.EncodeToString(open.SCTE35().Data()))
	}
	println()
}

func TestOutInInOutIn36_37_10_37_37(t *testing.T) {
	state := NewState()

	// 0x36 - event_id:0 - seg_num: 0 - seg_expected: 0
	outSignalBytes, _ := base64.StdEncoding.DecodeString("/DBLAAEs/S0UAP/wBQb+AAAAAAA1AjNDVUVJT///9X//AACky4AJH1NJR05BTDoyWURWeCtSKzlWc0FBQUFBQUFBQkFRPT02AAD9DXQ/")
	outSignal, err := NewSCTE35(append([]byte{0x0}, outSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(outSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x37 - event_id: 0 - seg_num: 1 - seg_expected: 3
	firstInSignalBytes, _ := base64.StdEncoding.DecodeString("/DBGAAEtT5LUAP/wBQb+AAAAAAAwAi5DVUVJT///9X+/CR9TSUdOQUw6MllEVngrUis5VnNBQUFBQUFBQUJBZz09NwEDU/ktPg==")
	firstInSignal, err := NewSCTE35(append([]byte{0x0}, firstInSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(firstInSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x10 closes 0x36 as per the SCTE-35 2020 spec
	In35SignalBytes, _ := base64.StdEncoding.DecodeString("/DBFAAAAABeOAP/wBQb+WxdHNQAvAi1DVUVJACSqm3/9ABNP17wMGURJU0MyMjA1NjVfMDAyXzAxXzU3MUEtMDEQAQEeP8NH")
	In35Signal, err := NewSCTE35(append([]byte{0x0}, In35SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(In35Signal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 1 {
		t.Errorf("1 event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be 1 signal (%d)", len(state.Open()))
	}

	// 0x37 - event_id: 0 - seg_num: 2 - seg_expected: 3
	secondInSignalBytes, _ := base64.StdEncoding.DecodeString("/DBGAAEteMW0AP/wBQb+AAAAAAAwAi5DVUVJT///9X+/CR9TSUdOQUw6MllEVngrUis5VnNBQUFBQUFBQUJBdz09NwID1nPQRg==")
	secondInSignal, err := NewSCTE35(append([]byte{0x0}, secondInSignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(secondInSignal.Descriptors()[0])
	// ErrSCTE35MissingOut expected as 0x30 closed the 0x36, the 0x37 has nothing to close.
	if err != nil && err != gots.ErrSCTE35MissingOut {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should have been 1 open signal, but we saw (%d) open signals instead.", len(state.Open()))
	}
}

// Tests closing logic of 0x11 with 0x10, 0x20 ,0x30
func Test11ClosingLogic(t *testing.T) {
	state := NewState()

	// 0x10
	scteDesc10SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc10Signal)
	outSignal, err := NewSCTE35(append([]byte{0x0}, scteDesc10SignalBytes	...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(outSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be one open signal (%d)", len(state.Open()))
	}

	// 0x20
	scteDesc20SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc20Signal)
	outSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc20SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}
	closed, err = state.ProcessDescriptor(outSignal.Descriptors()[0])
	if len(closed) != 0 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 2 {
		t.Errorf("There should be two open signals (%d)", len(state.Open()))
	}

	// 0x30
	scteDesc30SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc30Signal)
	inSignal, err := NewSCTE35(append([]byte{0x0}, scteDesc30SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 3 {
		t.Errorf("There should be three open signals (%d)", len(state.Open()))
	}

	// 0x11, it should close the open 0x10, 0x20 ,0x30
	scteDesc11SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc11Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc11SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 3 {
		t.Errorf("Three events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should be no open signal (%d)", len(state.Open()))
	}
}

func Test11Closing22(t *testing.T) {
	state := NewState()

	// 0x22
	scteDesc22SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc22Signal)
	outSignal, err := NewSCTE35(append([]byte{0x0}, scteDesc22SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(outSignal.Descriptors()[0])
	if len(closed) != 0	 {
		t.Errorf("No events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be 1 open signal (%d)", len(state.Open()))
	}

	// 0x11, it should close the open 0x22
	scteDesc11SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc11Signal)
	inSignal, err := NewSCTE35(append([]byte{0x0}, scteDesc11SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 1 {
		t.Errorf("One event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should be no open signal (%d)", len(state.Open()))
	}
}

func Test11Closing34_36_40_44(t *testing.T) {
	state := NewState()

	//0x34
	scteDesc34SignalBytes , _ := base64.StdEncoding.DecodeString(scteDesc34Signal)
	inSignal, err := NewSCTE35(append([]byte{0x0}, scteDesc34SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be 1 open signal (%d)", len(state.Open()))
	}

	// 0x36
	scteDesc36SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc36Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc36SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 2 {
		t.Errorf("There should be 2 open signals (%d)", len(state.Open()))
	}

	// 0x40
	scteDesc40SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc40Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc40SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 3 {
		t.Errorf("There should be 3 open signals (%d)", len(state.Open()))
	}

	// 0x44
	scteDesc44SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc44Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc44SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 4 {
		t.Errorf("There should be 4 open signals (%d)", len(state.Open()))
	}

	// 0x11 - it should close the open 0x34, 0x36, 0x40, 0x44
	scteDesc11SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc11Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc11SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 4 {
		t.Errorf("4 events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should be no open signal (%d)", len(state.Open()))
	}
}

// Tests closing logic of 0x45
func Test45ClosingLogic(t *testing.T) {
	state := NewState()

	//0x44
	scteDesc44SignalBytes , _ := base64.StdEncoding.DecodeString(scteDesc44Signal)
	inSignal, err := NewSCTE35(append([]byte{0x0}, scteDesc44SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be 1 open signal (%d)", len(state.Open()))
	}

	// 0x30
	scteDesc30SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc30Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc30SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 2 {
		t.Errorf("There should be 2 open signals (%d)", len(state.Open()))
	}

	// 0x45 - closes 0x44 and 0x30
	scteDesc45SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc45Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc45SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 2 {
		t.Errorf("2 events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 0 {
		t.Errorf("There should be no open signals (%d)", len(state.Open()))
	}
}
// Tests closing logic of 0x44
func Test44ClosingLogic(t *testing.T) {
	state := NewState()

	// 0x44
	scteDesc44SignalBytes , _ := base64.StdEncoding.DecodeString(scteDesc44Signal)
	inSignal, err := NewSCTE35(append([]byte{0x0}, scteDesc44SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err := state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 {
		t.Errorf("There should be 1 open signal (%d)", len(state.Open()))
	}

	// 0x30
	scteDesc30SignalBytes, _ := base64.StdEncoding.DecodeString(scteDesc30Signal)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc30SignalBytes...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 0 {
		t.Errorf("No event should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 2 {
		t.Errorf("There should be 2 open signals (%d)", len(state.Open()))
	}

	// Different 0x44 - closes the earlier 0x44 and 0x30
	scteDesc44SignalBytes1, _ := base64.StdEncoding.DecodeString(scteDesc44Signal1)
	inSignal, err = NewSCTE35(append([]byte{0x0}, scteDesc44SignalBytes1...))
	if err != nil {
		t.Fatal("Error creating SCTE-35 signal, return err:", err.Error())
	}

	closed, err = state.ProcessDescriptor(inSignal.Descriptors()[0])
	if err != nil {
		t.Errorf("ProcessDescriptor returned an error: %s", err.Error())
	}
	if len(closed) != 2 {
		t.Errorf("2 events should have been closed (%d were)", len(closed))
	}
	if len(state.Open()) != 1 { // the last 0x44 sent should still be open
		t.Errorf("There should be 1 open signals (%d)", len(state.Open()))
	}
}
