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

type cableLabsEbp struct {
	DataFieldTag     uint8
	DataFieldLength  uint8
	FormatIdentifier uint32
	DataFlags        uint8
	ExtensionFlags   uint8
	SapType          uint8
	Grouping         []uint8
	TimeSeconds      uint32
	TimeFraction     uint32
	PartitionFlags   uint8
	ReservedBytes    []byte
	SuccessReadTime  time.Time
}

// CreateCableLabsEbp returns a new cableLabsEbp with default values.
func CreateCableLabsEbp() cableLabsEbp {
	return cableLabsEbp{
		DataFieldTag:     CableLabsEbpTag,
		DataFieldLength:  1, // not empty by default
		FormatIdentifier: CableLabsFormatIdentifier,
	}
}

// EBPtype returns the type (what is the format) of the EBP.
func (ebp *cableLabsEbp) EBPType() byte {
	return ebp.DataFieldTag
}

// FragmentFlag returns true if the fragment flag is set.
func (ebp *cableLabsEbp) FragmentFlag() bool {
	return ebp.DataFlags&0x80 != 0
}

// SetFragmentFlag sets the fragment flag.
func (ebp *cableLabsEbp) SetFragmentFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x80
	}
}

// SegmentFlag returns true if the segment flag is set.
func (ebp *cableLabsEbp) SegmentFlag() bool {
	return ebp.DataFlags&0x40 != 0
}

// SetSegmentFlag sets the segment flag.
func (ebp *cableLabsEbp) SetSegmentFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x40
	}
}

// SapFlag returns true if the sap flag is set.
func (ebp *cableLabsEbp) SapFlag() bool {
	return ebp.DataFlags&0x20 != 0
}

// SetSapFlag sets the sap flag.
func (ebp *cableLabsEbp) SetSapFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x20
	}
}

// GroupingFlag returns true if the grouping flag is set.
func (ebp *cableLabsEbp) GroupingFlag() bool {
	return ebp.DataFlags&0x10 != 0
}

// SetGroupingFlag sets the grouping flag.
func (ebp *cableLabsEbp) SetGroupingFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x10
	}
}

// TimeFlag returns true if the time flag is set.
func (ebp *cableLabsEbp) TimeFlag() bool {
	return ebp.DataFlags&0x08 != 0
}

// SetTimeFlag sets the time flag
func (ebp *cableLabsEbp) SetTimeFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x08
	}
}

// ConcealmentFlag returns true if the concealment flag is set.
func (ebp *cableLabsEbp) ConcealmentFlag() bool {
	return ebp.DataFlags&0x04 != 0
}

// SetConcealmentFlag sets the concealment flag.
func (ebp *cableLabsEbp) SetConcealmentFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x04
	}
}

// ExtensionFlag returns true if the extension flag is set.
func (ebp *cableLabsEbp) ExtensionFlag() bool {
	return ebp.DataFlags&0x01 != 0
}

// SetExtensionFlag sets the extension flag.
func (ebp *cableLabsEbp) SetExtensionFlag(value bool) {
	if ebp.DataFieldLength != 0 && value {
		ebp.DataFlags |= 0x01
	}
}

// EBPTime returns the EBP time as a UTC time.
func (ebp *cableLabsEbp) EBPTime() time.Time {
	return extractUtcTime(ebp.TimeSeconds, ebp.TimeFraction)
}

// SetEBPTime sets the time of the EBP. Takes UTC time as an input.
func (ebp *cableLabsEbp) SetEBPTime(t time.Time) {
	ebp.TimeSeconds, ebp.TimeFraction = insertUtcTime(t)
}

// Sap returns the sap of the EBP.
func (ebp *cableLabsEbp) Sap() byte {
	return ebp.SapType
}

// SetSap sets the sap of the EBP.
func (ebp *cableLabsEbp) SetSap(sapType byte) {
	ebp.SapType = sapType
}

// SetPartitionFlag returns true if the partition flag.
func (ebp *cableLabsEbp) PartitionFlag() bool {
	return ebp.ExtensionFlag() && ebp.ExtensionFlags&0x80 != 0
}

// SetPartitionFlag sets the partition flag.
func (ebp *cableLabsEbp) SetPartitionFlag(value bool) {
	if ebp.ExtensionFlag() && value {
		ebp.ExtensionFlags |= 0x80
	}
}

// IsEmpty returns if the EBP is empty (zero length)
func (ebp *cableLabsEbp) IsEmpty() bool {
	return ebp.DataFieldLength == 0
}

// SetIsEmpty sets if the EBP is empty (zero length)
func (ebp *cableLabsEbp) SetIsEmpty(value bool) {
	if value {
		ebp.DataFieldLength = 0
	} else {
		ebp.DataFieldLength = 1
	}
}

// Defines when the EBP was read successfully
func (ebp *cableLabsEbp) EBPSuccessReadTime() time.Time {
	return ebp.SuccessReadTime
}

func readCableLabsEbp(data io.Reader) (ebp *cableLabsEbp, err error) {
	ebp = &cableLabsEbp{DataFieldTag: CableLabsEbpTag}

	if err = binary.Read(data, ebpEncoding, &ebp.DataFieldLength); err != nil {
		return nil, err
	}

	remaining := ebp.DataFieldLength

	if err = binary.Read(data, ebpEncoding, &ebp.FormatIdentifier); err != nil {
		return nil, err
	}

	if err = binary.Read(data, ebpEncoding, &ebp.DataFlags); err != nil {
		return nil, err
	}

	remaining -= uint8(5)

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
		var group byte
		if err = binary.Read(data, ebpEncoding, &group); err != nil {
			return nil, err
		}
		ebp.Grouping = append(ebp.Grouping, group)

		remaining -= uint8(1)
		for group&0x80 != 0 {
			if err = binary.Read(data, ebpEncoding, &group); err != nil {
				return nil, err
			}
			ebp.Grouping = append(ebp.Grouping, group)
			remaining -= uint8(1)
		}
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

	if ebp.PartitionFlag() {
		if err = binary.Read(data, ebpEncoding, &ebp.PartitionFlags); err != nil {
			return nil, err
		}
		remaining -= uint8(1)
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
func (ebp *cableLabsEbp) Data() []byte {
	requiredFields := new(bytes.Buffer)
	data := new(bytes.Buffer)
	binary.Write(requiredFields, ebpEncoding, ebp.DataFieldTag)

	if ebp.DataFieldLength == 0 {
		return data.Bytes()
	}

	binary.Write(data, ebpEncoding, ebp.FormatIdentifier)

	binary.Write(data, ebpEncoding, ebp.DataFlags)
	if ebp.ExtensionFlag() {
		binary.Write(data, ebpEncoding, ebp.ExtensionFlags)
	}

	if ebp.SapFlag() {
		binary.Write(data, ebpEncoding, ebp.SapType)
	}

	if ebp.GroupingFlag() {
		for i := range ebp.Grouping {
			ebp.Grouping[i] |= 0x80 // set flag because this is not the last ID
			if i == len(ebp.Grouping)-1 {
				ebp.Grouping[i] &= 0x7F // last index does not have this flag set
			}
			binary.Write(data, ebpEncoding, ebp.Grouping[i])
		}
	}

	if ebp.TimeFlag() {
		binary.Write(data, ebpEncoding, ebp.TimeSeconds)
		binary.Write(data, ebpEncoding, ebp.TimeFraction)
	}

	if ebp.PartitionFlag() {
		binary.Write(data, ebpEncoding, ebp.PartitionFlags)
	}

	binary.Write(data, ebpEncoding, ebp.ReservedBytes)

	ebp.DataFieldLength = uint8(data.Len())

	binary.Write(requiredFields, ebpEncoding, ebp.DataFieldLength)
	binary.Write(requiredFields, ebpEncoding, data.Bytes())

	return requiredFields.Bytes()
}
