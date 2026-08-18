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

import "testing"

// TestReadEncoderBoundaryPointTruncated verifies that a truncated EBP whose
// flag byte advertises fields that are not actually present is rejected with
// an error rather than panicking with an index/slice out of range. These
// inputs were found by fuzzing ReadEncoderBoundaryPoint.
func TestReadEncoderBoundaryPointTruncated(t *testing.T) {
	cases := []struct {
		name string
		data []byte
	}{
		// CableLabs (tag 0xDF): TimeFlag (0x08) set but the 8 time bytes are
		// missing. Previously panicked in readCableLabsEbp.
		{"cablelabs_time", []byte{0xDF, 0x08, 0x00, 0x00, 0x00, 0x00, 0x08}},
		// CableLabs: SapFlag (0x20) set but sap byte missing.
		{"cablelabs_sap", []byte{0xDF, 0x08, 0x00, 0x00, 0x00, 0x00, 0x20}},
		// CableLabs: GroupingFlag (0x10) with the extension bit (0x80) set,
		// runs off the end of the buffer in the grouping loop.
		{"cablelabs_grouping", []byte{0xDF, 0x08, 0x00, 0x00, 0x00, 0x00, 0x10, 0x80}},
		// Comcast (tag 0xA9): TimeFlag set but time bytes missing.
		{"comcast_time", []byte{0xA9, 0x08, 0x08}},
		// Comcast: SapFlag set but sap byte missing.
		{"comcast_sap", []byte{0xA9, 0x02, 0x20}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Must not panic; a truncated buffer is expected to error.
			if _, err := ReadEncoderBoundaryPoint(tc.data); err == nil {
				t.Fatalf("expected error for truncated EBP, got nil")
			}
		})
	}
}
