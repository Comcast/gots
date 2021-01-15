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
	"time"
)

// FragmentFlag returns true if the fragment flag is set.
func (ebp *baseEbp) FragmentFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x80 != 0
}

// SetFragmentFlag sets the fragment flag.
func (ebp *baseEbp) SetFragmentFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x80
	}
}

// SegmentFlag returns true if the segment flag is set.
func (ebp *baseEbp) SegmentFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x40 != 0
}

// SetSegmentFlag sets the segment flag.
func (ebp *baseEbp) SetSegmentFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x40
	}
}

// SapFlag returns true if the sap flag is set.
func (ebp *baseEbp) SapFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x20 != 0
}

// SetSapFlag sets the sap flag.
func (ebp *baseEbp) SetSapFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x20
	}
}

// GroupingFlag returns true if the grouping flag is set.
func (ebp *baseEbp) GroupingFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x10 != 0
}

// SetGroupingFlag sets the grouping flag.
func (ebp *baseEbp) SetGroupingFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x10
	}
}

// TimeFlag returns true if the time flag is set.
func (ebp *baseEbp) TimeFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x08 != 0
}

// SetTimeFlag sets the time flag
func (ebp *baseEbp) SetTimeFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x08
	}
}

// EBPTime returns the EBP time as a UTC time.
func (ebp *baseEbp) EBPTime() time.Time {
	return extractUtcTime(ebp.TimeSeconds, ebp.TimeFraction)
}

// SetEBPTime sets the time of the EBP. Takes UTC time as an input.
func (ebp *baseEbp) SetEBPTime(t time.Time) {
	ebp.TimeSeconds, ebp.TimeFraction = insertUtcTime(t)
}

// EBPSuccessReadTime defines when the EBP was read successfully.
func (ebp *baseEbp) EBPSuccessReadTime() time.Time {
	return ebp.SuccessReadTime
}

// Sap returns the sap of the EBP.
func (ebp *baseEbp) Sap() byte {
	return ebp.SapType
}

// SetSap sets the sap of the EBP.
func (ebp *baseEbp) SetSap(sapType byte) {
	ebp.SapType = sapType
}

// ExtensionFlag returns true if the extension flag is set.
func (ebp *baseEbp) ExtensionFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x01 != 0
}

// SetExtensionFlag sets the extension flag.
func (ebp *baseEbp) SetExtensionFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x01
	}
}

// IsEmpty returns if the EBP is empty (zero length)
func (ebp *baseEbp) IsEmpty() bool {
	return ebp.DataFieldLength == 0
}

// SetIsEmpty sets if the EBP is empty (zero length)
func (ebp *baseEbp) SetIsEmpty(value bool) {
	if value {
		ebp.DataFieldLength = 0
	} else {
		ebp.DataFieldLength = 1
	}
}

// StreamSyncSignal is used by some Tanscoder vendors to indicate the transcoders for a stream are in sync or not in sync. 
// The packagers can read this signal and then determine what to do. Ideally they will always be in sync if the transcoders are configured correctly and healthy.
// StreamSyncSignal checks the Grouping ID field and returns the Stream Sync Signal
// Note: The intention of the Grouping ID was very broad at first and we have had many discussions over the years as to how to use it,
// but Stream Sync and SCTE-35 were the only two ever to get implemented. MK was the only Transcoder to implement SCTE-35 within the
// Grouping ID section and we never leveraged that data on any downstream devices.
func (ebp *baseEbp) StreamSyncSignal() uint8 {
	// The first byte is Grouping ID, the second byte is the sync signal
	if len(ebp.Grouping) == 2 && ebp.Grouping[0] == 0x80 {
		return ebp.Grouping[1]
	}

	return InvalidStreamSyncSignal
}
