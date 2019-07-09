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
	"testing"
	"time"

	"github.com/Comcast/gots"
)

var CableLabsEBPBytes = []byte{
	0xDF,                   // Tag
	0x14,                   // Length
	0x45, 0x42, 0x50, 0x30, // Format Identifier
	0xBD,       // Flags
	0x80,       // ExtensionFlags
	0x02,       // SapFlags
	0xFF, 0x7F, // GroupingFlag
	0xD6, 0xEE, 0x7B, 0xD8, // Time Second
	0x8D, 0xC7, 0x14, 0xFC, // Time Fraction
	0x03, // Partition
	0x04, 0x05}

var ComcastEBPBytes = []byte{
	0xA9,                   // Tag
	0x0E,                   // Length
	0xBD,                   // Flags
	0x01,                   // Extensions
	0x02,                   // SAP
	0x03,                   // Grouping
	0xD6, 0xEE, 0x7B, 0xD8, // Time Second
	0x8D, 0xC7, 0x14, 0xFC, // Time Fraction
	0x04, 0x05} // Reserved

func TestReadEncoderBoundaryPoint(t *testing.T) {
	ebp, err := ReadEncoderBoundaryPoint(CableLabsEBPBytes)
	if err != nil {
		t.Errorf("ReadEncoderBoundaryPoint() returned error on valid: %v", err)
	}
	if ebp.EBPType() != CableLabsEbpTag {
		t.Errorf("ReadEncoderBoundaryPoint() read wrong type of EBP")
	}

	ebp, err = ReadEncoderBoundaryPoint(ComcastEBPBytes)
	if err != nil {
		t.Errorf("ReadEncoderBoundaryPoint() returned error on valid: %v", err)
	}
	if ebp.EBPType() != ComcastEbpTag {
		t.Errorf("ReadEncoderBoundaryPoint() read wrong type of EBP")
	}

	ebp, err = ReadEncoderBoundaryPoint([]byte{0xAB})
	if err != gots.ErrUnrecognizedEbpType {
		t.Errorf("ReadEncoderBoundaryPoint() read wrong type of EBP")
	}

	ebp, err = ReadEncoderBoundaryPoint([]byte{})
	if err != gots.ErrNoEBPData {
		t.Errorf("ReadEncoderBoundaryPoint() should have returned gots.ErrNoEBPData")
	}
}

func TestExtractUtcTime(t *testing.T) {
	s := uint32(0xD6EE7BD8)
	f := uint32(0x8DC714FC)
	got := extractUtcTime(s, f)
	want := time.Unix(0, 1396964696553818999).UTC()
	if want != got {
		t.Errorf("TestUtcTime(), want=%v, got=%v, nanos=%d", want, got, got.UnixNano())
	}
}

func TestInsertUtcTime(t *testing.T) {
	s0 := uint32(0xD6EE7BD8)
	f0 := uint32(0x8DC714FC)
	want := time.Unix(0, 1396964696553818999).UTC()
	s1, f1 := insertUtcTime(want)
	if s0 != s1 {
		t.Errorf("secondsInserted=%d, secondsExpected=%d\n", s1, s0)
	}
	if f0 != f1 {
		t.Errorf("fractionInserted=%d, fractionExpected=%d\n", f1, f0)
	}

}

// TODO TestUtcTimeAfter2036
