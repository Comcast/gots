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
package gots

import (
	"testing"
)

func TestExtractPCR(t *testing.T) {
	b := []byte{168, 149, 227, 4, 0, 126}
	pcr := ExtractPCR(b)
	if pcr != 1697037160926 {
		t.Errorf("ExtractPCR returned %v", pcr)
	}
}

func TestInsertPCR(t *testing.T) {
	b := make([]byte, 6)

	var pcr uint64 = 1697037160926
	InsertPCR(b, pcr)
	if ExtractPCR(b) != pcr {
		t.Errorf("Insert PCR test 1 failed: %v", b)
	}

	var pcr2 uint64 = 2576980377599
	InsertPCR(b, pcr2)
	if ExtractPCR(b) != pcr2 {
		t.Errorf("Insert PCR test 2 failed: %v (%v)", b, ExtractPCR(b))
	}
}
