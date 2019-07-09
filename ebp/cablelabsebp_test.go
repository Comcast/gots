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

func TestParseCableLabsEBP(t *testing.T) {
	ebp, err := readCableLabsEbp(CableLabsEBPBytes)
	if err != nil {
		t.Errorf("readCableLabsEbp() returned not null error on valid EBP %v", err)
	}
	if ebp == nil {
		t.Errorf("readCableLabsEbp() returned null EBP on valid EBP")
	}
	if !ebp.FragmentFlag() {
		t.Errorf("readCableLabsEbp() has incorrect fragment flag")
	}
	if ebp.SegmentFlag() {
		t.Errorf("readCableLabsEbp() has incorrect segment flag")
	}
	if !ebp.SapFlag() {
		t.Errorf("readCableLabsEbp() has incorrect SAP flag")
	}
	if !ebp.GroupingFlag() {
		t.Errorf("readCableLabsEbp() has incorrect grouping flag")
	}
	if !ebp.TimeFlag() {
		t.Errorf("readCableLabsEbp() has incorrect time flag")
	}
	if !ebp.ConcealmentFlag() {
		t.Errorf("readCableLabsEbp() has incorrect concealment flag")
	}
	if !ebp.ExtensionFlag() {
		t.Errorf("readCableLabsEbp() has incorrect extension flag")
	}
	if !ebp.PartitionFlag() {
		t.Errorf("readCableLabsEbp() has incorrect partition flag")
	}
	if ebp.EBPTime() != time.Unix(0, 1396964696553818999).UTC() {
		t.Errorf("readCableLabsEbp() has incorrect time")
	}
	if ebp.Sap() != byte(0x02) {
		t.Errorf("readCableLabsEbp() has incorrect SAP")
	}
	if len(ebp.Grouping) != 2 {
		t.Errorf("readCableLabsEbp() has incorrect grouping, %v", ebp.Grouping)
	}
	if len(ebp.ReservedBytes) != 2 {
		t.Errorf("readCableLabsEbp() has incorrect reserved, %v", ebp.ReservedBytes)
	}

	ebp, err = readCableLabsEbp([]byte{0xff, 0xf3})
	if err == nil {
		t.Errorf("readCableLabsEbp() returned null error %v on invalid EBP", err)
	}
	if ebp != nil {
		t.Errorf("readCableLabsEbp() returned not null EBP on invalid EBP")
	}
}

func TestCableLabsEBP(t *testing.T) {
	expected := CableLabsEBPBytes
	ebp := CreateCableLabsEbp()
	ebp.SetFragmentFlag(true)
	ebp.SetConcealmentFlag(true)
	ebp.SetExtensionFlag(true)
	ebp.SetPartitionFlag(true)
	ebp.SetSapFlag(true)
	ebp.SapType = 0x02
	ebp.SetGroupingFlag(true)
	ebp.Grouping = []byte{0x7F, 0xFF} // wrong on purpose, will be corrected to 0xFF, 0x7F
	ebp.SetTimeFlag(true)
	ebp.SetEBPTime(time.Unix(0, 1396964696553818999).UTC())
	ebp.PartitionFlags = 0x03
	ebp.ReservedBytes = []byte{0x04, 0x05}

	generated := ebp.Data()
	if !bytes.Equal(generated, expected) {
		t.Errorf("Data() does not produce expected raw data\nExpected: %X\n     Got: %X",
			expected,
			generated,
		)
	}
}
