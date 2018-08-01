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
	"encoding/binary"
	"io"
	"time"
)

type comcastEbp struct {
	DataFieldTag    uint8
	DataFieldLength uint8
	DataFlags       uint8
	ExtensionFlags  uint8
	SapType         uint8
	Grouping        uint8
	TimeSeconds     uint32
	TimeFraction    uint32
	ReservedBytes   []uint8
	SuccessReadTime time.Time
}

// CreateComcastEBP returns a new comcastEbp with default values.
func CreateComcastEBP() comcastEbp {
	return comcastEbp{
		DataFieldTag:    ComcastEbpTag,
		DataFieldLength: 1, // not empty by default
	}
}

// EBPtype returns the type (what is the format) of the EBP.
func (ebp *comcastEbp) EBPType() byte {
	return ebp.DataFieldTag
}

// FragmentFlag returns true if the fragment flag is set.
func (ebp *comcastEbp) FragmentFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x80 != 0
}

// SetFragmentFlag sets the fragment flag.
func (ebp *comcastEbp) SetFragmentFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x80
	}
}

// SegmentFlag returns true if the segment flag is set.
func (ebp *comcastEbp) SegmentFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x40 != 0
}

// SetSegmentFlag sets the segment flag.
func (ebp *comcastEbp) SetSegmentFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x40
	}
}

// SapFlag returns true if the sap flag is set.
func (ebp *comcastEbp) SapFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x20 != 0
}

// SetSapFlag sets the sap flag.
func (ebp *comcastEbp) SetSapFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x20
	}
}

// GroupingFlag returns true if the grouping flag is set.
func (ebp *comcastEbp) GroupingFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x10 != 0
}

// SetGroupingFlag sets the grouping flag.
func (ebp *comcastEbp) SetGroupingFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x10
	}
}

// TimeFlag returns true if the time flag is set.
func (ebp *comcastEbp) TimeFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x08 != 0
}

// SetTimeFlag sets the time flag
func (ebp *comcastEbp) SetTimeFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x08
	}
}

// DiscontinuityFlag returns true if the discontinuity flag is set.
func (ebp *comcastEbp) DiscontinuityFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x04 != 0
}

// SetDiscontinuityFlag sets the discontinuity flag.
func (ebp *comcastEbp) SetDiscontinuityFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x04
	}
}

// ExtensionFlag returns true if the extension flag is set.
func (ebp *comcastEbp) ExtensionFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x01 != 0
}

// SetExtensionFlag sets the extension flag.
func (ebp *comcastEbp) SetExtensionFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x01
	}
}

// EBPTime returns the EBP time as a UTC time.
func (ebp *comcastEbp) EBPTime() time.Time {
	return extractUtcTime(ebp.TimeSeconds, ebp.TimeFraction)
}

// SetEBPTime sets the time of the EBP. Takes UTC time as an input.
func (ebp *comcastEbp) SetEBPTime(t time.Time) {
	ebp.TimeSeconds, ebp.TimeFraction = insertUtcTime(t)
}

// Sap returns the sap of the EBP.
func (ebp *comcastEbp) Sap() byte {
	return ebp.SapType
}

// SetSap sets the sap of the EBP.
func (ebp *comcastEbp) SetSap(sapType byte) {
	ebp.SapType = sapType
}

// IsEmpty returns if the EBP is empty (zero length)
func (ebp *comcastEbp) IsEmpty() bool {
	return ebp.DataFieldLength == 0
}

// SetIsEmpty sets if the EBP is empty (zero length)
func (ebp *comcastEbp) SetIsEmpty(value bool) {
	if value {
		ebp.DataFieldLength = 0
	} else {
		ebp.DataFieldLength = 1
	}
}

// EBPSuccessReadTime defines when the EBP was read successfully.
func (ebp *comcastEbp) EBPSuccessReadTime() time.Time {
	return ebp.SuccessReadTime
}

// readComcastEbp will parse raw bytes without the tag into a Comcast EBP.
func readComcastEbp(data io.Reader) (ebp *comcastEbp, err error) {
	ebp = &comcastEbp{DataFieldTag: ComcastEbpTag}

	if err = binary.Read(data, ebpEncoding, &ebp.DataFieldLength); err != nil {
		return nil, err
	}

	remaining := ebp.DataFieldLength

	if remaining == 0 {
		return ebp, nil
	}

	if err = binary.Read(data, ebpEncoding, &ebp.DataFlags); err != nil {
		return nil, err
	}

	remaining -= uint8(1)

	if ebp.ExtensionFlag() {
		if err = binary.Read(data, ebpEncoding, &ebp.ExtensionFlags); err != nil {
			return nil, err
		}
		remaining -= uint8(1)
	}

	if ebp.SapFlag() {
		if err = binary.Read(data, ebpEncoding, &ebp.SapType); err != nil {
			return nil, err
		}
		remaining -= uint8(1)
	}

	if ebp.GroupingFlag() {
		if err = binary.Read(data, ebpEncoding, &ebp.Grouping); err != nil {
			return nil, err
		}
		remaining -= uint8(1)
	}

	if ebp.TimeFlag() {
		if err = binary.Read(data, ebpEncoding, &ebp.TimeSeconds); err != nil {
			return nil, err
		}
		if err = binary.Read(data, ebpEncoding, &ebp.TimeFraction); err != nil {
			return nil, err
		}
		remaining -= uint8(8)
	}

	if remaining > 0 {
		ebp.ReservedBytes = make([]byte, remaining)
		if err = binary.Read(data, ebpEncoding, &ebp.ReservedBytes); err != nil {
			return nil, err
		}

	}

	// update the successful read time
	ebp.SuccessReadTime = time.Now()

	return ebp, nil
}

// Data will return the raw bytes of the EBP
func (ebp *comcastEbp) Data() []byte {
	requiredFields := new(bytes.Buffer)
	data := new(bytes.Buffer)
	binary.Write(requiredFields, ebpEncoding, ebp.DataFieldTag)

	if ebp.DataFieldLength == 0 {
		return data.Bytes()
	}

	binary.Write(data, ebpEncoding, ebp.DataFlags)

	if ebp.ExtensionFlag() {
		binary.Write(data, ebpEncoding, ebp.ExtensionFlags)
	}

	if ebp.SapFlag() {
		binary.Write(data, ebpEncoding, ebp.SapType)
	}

	if ebp.GroupingFlag() {
		binary.Write(data, ebpEncoding, ebp.Grouping)
	}

	if ebp.TimeFlag() {
		binary.Write(data, ebpEncoding, ebp.TimeSeconds)
		binary.Write(data, ebpEncoding, ebp.TimeFraction)
	}

	binary.Write(data, ebpEncoding, ebp.ReservedBytes)

	ebp.DataFieldLength = uint8(data.Len())

	binary.Write(requiredFields, ebpEncoding, ebp.DataFieldLength)
	binary.Write(requiredFields, ebpEncoding, data.Bytes())

	return requiredFields.Bytes()
}
