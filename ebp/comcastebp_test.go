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
package ebp

import (
	"bytes"
	"testing"
	"time"
)

func TestParseComcastEBP(t *testing.T) {
	// Read a valid EBP
	ebp, err := readComcastEbp(ComcastEBPBytes)
	if err != nil {
		t.Errorf("readComcastEbp() returned not null error on valid EBP %v", err)
	}
	if ebp == nil {
		t.Errorf("readComcastEbp() returned null EBP on valid EBP")
	}
	if !ebp.FragmentFlag() {
		t.Errorf("readComcastEbp() has incorrect fragment flag")
	}
	if ebp.SegmentFlag() {
		t.Errorf("readComcastEbp() has incorrect segment flag")
	}
	if !ebp.SapFlag() {
		t.Errorf("readComcastEbp() has incorrect SAP flag")
	}
	if !ebp.GroupingFlag() {
		t.Errorf("readComcastEbp() has incorrect grouping flag")
	}
	if !ebp.TimeFlag() {
		t.Errorf("readComcastEbp() has incorrect time flag")
	}
	if !ebp.DiscontinuityFlag() {
		t.Errorf("readComcastEbp() has incorrect disc flag")
	}
	if !ebp.ExtensionFlag() {
		t.Errorf("readComcastEbp() has incorrect extension flag")
	}
	if ebp.EBPTime() != time.Unix(0, 1396964696553818999).UTC() {
		t.Errorf("readComcastEbp() has incorrect time")
	}
	if ebp.ExtensionFlags != byte(0x01) {
		t.Errorf("readComcastEbp() has incorrect extension")
	}
	if ebp.Sap() != byte(0x02) {
		t.Errorf("readComcastEbp() has incorrect SAP")
	}
	if ebp.Grouping != byte(0x03) {
		t.Errorf("readComcastEbp() has incorrect grouping")
	}

	// Read an EBP with 0 length
	ebp, err = readComcastEbp([]byte{ComcastEbpTag, 0x00})
	if err != nil {
		t.Errorf("readComcastEbp() returned error %v on valid EBP", err)
	}
	if ebp == nil {
		t.Errorf("readComcastEbp() returned not null EBP on valid EBP")
	}

	// Read an EBP with an invalid length
	ebp, err = readComcastEbp([]byte{ComcastEbpTag, 0xff})
	if err == nil {
		t.Errorf("readComcastEbp() returned null error %v on invalid EBP", err)
	}
	if ebp != nil {
		t.Errorf("readComcastEbp() returned not null EBP on invalid EBP")
	}
}

func TestCreateComcastEBP(t *testing.T) {
	expected := ComcastEBPBytes
	ebp := CreateComcastEBP()
	ebp.SetFragmentFlag(true)
	ebp.SetDiscontinuityFlag(true)
	ebp.SetExtensionFlag(true)
	ebp.ExtensionFlags = 0x01
	ebp.SetSapFlag(true)
	ebp.SapType = 0x02
	ebp.SetGroupingFlag(true)
	ebp.Grouping = 0x03
	ebp.SetTimeFlag(true)
	ebp.SetEBPTime(time.Unix(0, 1396964696553818999).UTC())
	ebp.ReservedBytes = []byte{0x04, 0x05}

	generated := ebp.Data()
	if !bytes.Equal(generated, expected) {
		t.Errorf("Data() does not produce expected raw data\nExpected: %X\n     Got: %X",
			expected,
			generated,
		)
	}
}
