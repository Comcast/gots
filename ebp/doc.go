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

// Package ebp is used to detect and extract EBP information. EBP information is stored in a packet adaptation field in the private data section.
package ebp

import "time"

// EBP tags
const (
	ComcastEbpTag             = uint8(0xA9)
	CableLabsEbpTag           = uint8(0xDF)
	CableLabsFormatIdentifier = 0x45425030
)

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
