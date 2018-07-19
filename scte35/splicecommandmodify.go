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

package scte35

import (
	"github.com/Comcast/gots"
)

func (c *spliceNull) Data() []byte {
	return []byte{}
}

func (c *spliceNull) SetHasPTS(value bool) {
	// do nothing
}

func (c *spliceNull) SetPTS(value gots.PTS) {
	// do nothing
}

func spliceTimeBytes(hasPTS bool, pts gots.PTS) []byte {
	if hasPTS {
		bytes := make([]byte, 5)
		bytes[0] = 0xFE                  // 1111 1110
		bytes[0] |= byte(pts>>32) & 0x01 // 0000 0001
		bytes[1] = byte(pts >> 24)       // 1111 1111
		bytes[2] = byte(pts >> 16)       // 1111 1111
		bytes[3] = byte(pts >> 8)        // 1111 1111
		bytes[4] = byte(pts)             // 1111 1111
		return bytes
	}
	// return 0111 1110
	return []byte{0x7E} // only reserved bits are set
}

func (c *timeSignal) Data() []byte {
	return spliceTimeBytes(c.HasPTS(), c.pts)
}

func (c *timeSignal) SetHasPTS(value bool) {
	c.hasPTS = value
}

func (c *timeSignal) SetPTS(value gots.PTS) {
	c.pts = value
}

func (c *spliceInsert) Data() []byte {
	bytes := make([]byte, 6)
	bytes[0] = byte(c.eventID >> 24)
	bytes[1] = byte(c.eventID >> 16)
	bytes[2] = byte(c.eventID >> 8)
	bytes[3] = byte(c.eventID)

	bytes[4] = 0x7F // reserved

	if c.eventCancelIndicator {
		bytes[4] |= 0x80
	}

	bytes[5] = 0x0F // reserved

	if c.outOfNetworkIndicator {
		bytes[5] |= 0x80
	}
	if c.isProgramSplice {
		bytes[5] |= 0x40
	}
	if c.isProgramSplice {
		bytes[5] |= 0x40
	}
	if c.hasDuration {
		bytes[5] |= 0x20
	}
	if c.spliceImmediate {
		bytes[5] |= 0x10
	}

	if c.isProgramSplice && !c.spliceImmediate {
		bytes = append(bytes, spliceTimeBytes(c.hasPTS, c.pts)...)
	}

	if !c.isProgramSplice {
		// TODO no support for components.
		// zero components
		bytes = append(bytes, []byte{0x00}...)
	}

	if c.hasDuration {
		durationBytes := make([]byte, 5)
		// break_duration() structure:

		durationBytes[0] |= 0x7E // reserved
		if c.autoReturn {
			durationBytes[0] |= 0x80
		}
		durationBytes[0] |= byte(c.duration>>32) & 0x01 // 0000 0001
		durationBytes[1] = byte(c.duration >> 24)       // 1111 1111
		durationBytes[2] = byte(c.duration >> 16)       // 1111 1111
		durationBytes[3] = byte(c.duration >> 8)        // 1111 1111
		durationBytes[4] = byte(c.duration)             // 1111 1111
		bytes = append(bytes, durationBytes...)
	}

	endingBytes := make([]byte, 4)
	endingBytes[0] = byte(c.uniqueProgramId >> 8)
	endingBytes[1] = byte(c.uniqueProgramId)
	endingBytes[2] = byte(c.availNum)
	endingBytes[3] = byte(c.availsExpected)

	bytes = append(bytes, endingBytes...)
	return bytes
}

func (c *spliceInsert) SetEventID(value uint32) {
	c.eventID = value
}

func (c *spliceInsert) SetIsOut(value bool) {
	c.outOfNetworkIndicator = value
}

func (c *spliceInsert) SetIsEventCanceled(value bool) {
	c.eventCancelIndicator = value
}

func (c *spliceInsert) SetHasPTS(value bool) {
	c.hasPTS = value
}

func (c *spliceInsert) SetPTS(value gots.PTS) {
	c.pts = value
}

func (c *spliceInsert) SetHasDuration(value bool) {
	c.hasDuration = value
}

func (c *spliceInsert) SetDuration(value gots.PTS) {
	c.duration = value
}

func (c *spliceInsert) SetIsAutoReturn(value bool) {
	c.autoReturn = value
}

func (c *spliceInsert) SetUniqueProgramId(value uint16) {
	c.uniqueProgramId = value
}

func (c *spliceInsert) SetAvailNum(value uint8) {
	c.availNum = value
}

func (c *spliceInsert) SetAvailsExpected(value uint8) {
	c.availsExpected = value
}

func (c *spliceInsert) SetIsProgramSplice(value bool) {
	c.isProgramSplice = value
}

func (c *spliceInsert) SetSpliceImmediate(value bool) {
	c.spliceImmediate = value
}
