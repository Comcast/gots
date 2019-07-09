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
	"time"

	"github.com/Comcast/gots"
)

// cableLabsEbp is an encoder boundary point
type comcastEbp struct {
	baseEbp
	Grouping uint8
}

// CreateComcastEBP returns a new comcastEbp with default values.
func CreateComcastEBP() comcastEbp {
	return comcastEbp{
		baseEbp: baseEbp{
			DataFieldTag:    ComcastEbpTag,
			DataFieldLength: 1, // not empty by default
		},
	}
}

// EBPtype returns the type (what is the format) of the EBP.
func (ebp *comcastEbp) EBPType() byte {
	return ebp.DataFieldTag
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

// readComcastEbp will parse raw bytes without the tag into a Comcast EBP.
func readComcastEbp(data []byte) (ebp *comcastEbp, err error) {
	ebp = &comcastEbp{
		baseEbp: baseEbp{DataFieldTag: ComcastEbpTag},
	}

	// We read 2 bytes not based on flags
	if len(data) < 2 {
		return nil, gots.ErrNoPayload
	}

	index := uint8(0)

	ebp.DataFieldTag = data[index]
	index += uint8(1)

	ebp.DataFieldLength = data[index]
	index += uint8(1)

	// Check if the data is as advertised
	if ebp.DataFieldLength > 0 {
		if len(data) >= 3 {
			ebp.DataFlags = data[index]
			index += uint8(1)
		} else {
			return nil, gots.ErrInvalidEBPLength
		}
	}

	if ebp.ExtensionFlag() {
		ebp.ExtensionFlags = data[index]
		index += uint8(1)
	}

	if ebp.SapFlag() {
		ebp.SapType = data[index]
		index += uint8(1)
	}

	if ebp.GroupingFlag() {
		ebp.Grouping = data[index]
		index += uint8(1)
	}

	if ebp.TimeFlag() {
		ebp.TimeSeconds = binary.BigEndian.Uint32(data[index : index+4])
		index += uint8(4)

		ebp.TimeFraction = binary.BigEndian.Uint32(data[index : index+4])
		index += uint8(4)
	}

	if index < ebp.DataFieldLength+2 {
		if int(ebp.DataFieldLength+2) > len(data) {
			return nil, gots.ErrInvalidEBPLength
		}
		ebp.ReservedBytes = data[index : ebp.DataFieldLength+2]
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
