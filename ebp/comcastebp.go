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

func (ebp comcastEbp) EBPType() byte {
	return ebp.DataFieldTag
}

func (ebp comcastEbp) FragmentFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x80 != 0
}

func (ebp comcastEbp) SegmentFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x40 != 0
}

func (ebp comcastEbp) SapFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x20 != 0
}

func (ebp comcastEbp) GroupingFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x10 != 0
}

func (ebp comcastEbp) TimeFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x08 != 0
}

func (ebp comcastEbp) DiscontinuityFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x04 != 0
}

func (ebp comcastEbp) ExtensionFlag() bool {
	return ebp.DataFieldLength != 0 && ebp.DataFlags&0x01 != 0
}

func (ebp comcastEbp) EBPTime() time.Time {
	return extractUtcTime(ebp.TimeSeconds, ebp.TimeFraction)
}

func (ebp comcastEbp) Sap() byte {
	return ebp.SapType
}

// Defines when the EBP was read successfully
func (ebp comcastEbp) EBPSuccessReadTime() time.Time {
	return ebp.SuccessReadTime
}

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
