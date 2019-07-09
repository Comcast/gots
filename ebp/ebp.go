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
	"encoding/binary"
	"time"

	"github.com/Comcast/gots"
)

// EBP tags
const (
	ComcastEbpTag             = uint8(0xA9)
	CableLabsEbpTag           = uint8(0xDF)
	CableLabsFormatIdentifier = 0x45425030
)

var ebpEncoding = binary.BigEndian

// EncoderBoundaryPoint represents shared operations available on an all EBPs.
type EncoderBoundaryPoint interface {
	// SegmentFlag returns true if the segment flag is set.
	SegmentFlag() bool
	// SetSegmentFlag sets the segment flag.
	SetSegmentFlag(bool)
	// FragmentFlag returns true if the fragment flag is set.
	FragmentFlag() bool
	// SetFragmentFlag sets the fragment flag.
	SetFragmentFlag(bool)
	// TimeFlag returns true if the time flag is set.
	TimeFlag() bool
	// SetTimeFlag sets the time flag
	SetTimeFlag(bool)
	// GroupingFlag returns true if the grouping flag is set.
	GroupingFlag() bool
	// SetGroupingFlag sets the grouping flag.
	SetGroupingFlag(bool)
	// EBPTime returns the EBP time as a UTC time.
	EBPTime() time.Time
	// SetEBPTime sets the time of the EBP. Takes UTC time as an input.
	SetEBPTime(time.Time)
	// EBPSuccessReadTime defines when the EBP was read successfully.
	EBPSuccessReadTime() time.Time
	// SapFlag returns true if the sap flag is set.
	SapFlag() bool
	// SetSapFlag sets the sap flag.
	SetSapFlag(bool)
	// Sap returns the sap of the EBP.
	Sap() byte
	// SetSap sets the sap of the EBP.
	SetSap(byte)
	// ExtensionFlag returns true if the extension flag is set.
	ExtensionFlag() bool
	// SetExtensionFlag sets the extension flag.
	SetExtensionFlag(bool)
	// EBPtype returns the type (what is the format) of the EBP.
	EBPType() byte
	// IsEmpty returns if the EBP is empty (zero length)
	IsEmpty() bool
	// SetIsEmpty sets if the EBP is empty (zero length)
	SetIsEmpty(bool)
	// Data will return the raw bytes of the EBP
	Data() []byte
}

// baseEbp the base struct that is embedded in every ebp in gots
type baseEbp struct {
	DataFieldTag    uint8
	DataFieldLength uint8
	DataFlags       uint8
	ExtensionFlags  uint8
	SapType         uint8

	TimeSeconds  uint32
	TimeFraction uint32

	ReservedBytes   []uint8
	SuccessReadTime time.Time
}

// ReadEncoderBoundaryPoint parses and creates an EncoderBoundaryPoint from the given
// reader. If the bytes do not conform to a know EBP type, an error is returned.
func ReadEncoderBoundaryPoint(data []byte) (ebp EncoderBoundaryPoint, err error) {
	if len(data) == 0 {
		return nil, gots.ErrNoEBPData
	}

	switch data[0] {
	case ComcastEbpTag:
		ebp, err = readComcastEbp(data)
	case CableLabsEbpTag:
		ebp, err = readCableLabsEbp(data)
	default:
		ebp, err = nil, gots.ErrUnrecognizedEbpType
	}

	return ebp, err
}

func extractUtcTime(seconds uint32, fraction uint32) time.Time {
	nanos := uint64(seconds) * 1e9
	nanos += (uint64(fraction) * 1e9) >> 32 // truncateing off 3 or 2 bits.

	if seconds&0x80000000 != 0 { // if MSB set
		// If bit 0 is set, the UTC time is in the range 1968-2036 and UTC
		// time is reckoned from 0h 0m 0s UTC on 1 January 1900.
		return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(nanos))
	}
	// If bit 0 is not set, the time is in the range 2036-2104 and UTC
	// time is reckoned from 6h 28m 16s UTC on 7 February 2036.
	return time.Date(2036, 2, 7, 6, 28, 16, 0, time.UTC).Add(time.Duration(nanos))

}

func insertUtcTime(t time.Time) (seconds uint32, fraction uint32) {
	var startingTime time.Time
	if t.Before(time.Date(2036, 2, 7, 6, 28, 16, 0, time.UTC)) { // greater than or equal to
		startingTime = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		startingTime = time.Date(2036, 2, 7, 6, 28, 16, 0, time.UTC)
	}
	nanos := uint64(t.Sub(startingTime).Nanoseconds())
	seconds = uint32(nanos / 1e9)
	// running extractUtcTime and then insertUtcTime will produce
	// different results because of rounding that heppns twice.
	// 1 is added to avoid truncating the second time.
	fraction = uint32((((nanos % 1e9) + 1) << 32) / 1e9)
	return
}
