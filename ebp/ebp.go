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

	"github.com/comcast/gots"
)

var ebpEncoding = binary.BigEndian

// ReadEncoderBoundaryPoint parses and creates an EncoderBoundaryPoint from the given
// reader. If the bytes do not conform to a know EBP type, an error is returned.
func ReadEncoderBoundaryPoint(data io.Reader) (ebp EncoderBoundaryPoint, err error) {
	dataFieldTag := []byte{byte(0)}
	if _, err := data.Read(dataFieldTag); err != nil {
		return nil, err
	}

	switch dataFieldTag[0] {
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
	nanos += (uint64(fraction) * 1e9) >> 32

	if seconds&0x80000000 > 0 { // if MSB set
		// If bit 0 is set, the UTC time is in the range 1968-2036 and UTC
		// time is reckoned from 0h 0m 0s UTC on 1 January 1900.
		return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(nanos))
	}
	// If bit 0 is not set, the time is in the range 2036-2104 and UTC
	// time is reckoned from 6h 28m 16s UTC on 7 February 2036.
	return time.Date(2036, 2, 7, 6, 28, 16, 0, time.UTC).Add(time.Duration(nanos))

}
