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

// CreateSpliceInsertCommand will create a default SpliceInsertCommand.
func CreateSpliceInsertCommand() SpliceInsertCommand {
	return &spliceInsert{
		// enable this to have less variable sized fields by
		// default, do not use component splice mode by default
		isProgramSplice: true,
	}
}

// CreateTimeSignalCommand will create a default TimeSignalCommand
func CreateTimeSignalCommand() TimeSignalCommand {
	return &timeSignal{}
}

// CreateSpliceNull will create a Null SpliceCommand
func CreateSpliceNull() SpliceCommand {
	return &spliceNull{}
}

// SetComponentTag sets the component tag, which is used for the identification of the component.
func (c *component) SetComponentTag(value byte) {
	c.componentTag = value
}

// SetHasPTS sets a flag that determines if the component has a PTS.
func (c *component) SetHasPTS(value bool) {
	c.hasPts = value
}

// SetPTS sets the PTS of the component.
func (c *component) SetPTS(value gots.PTS) {
	c.pts = value & 0x01ffffffff
}

// Data returns the bytes of this splice command.
func (c *spliceNull) Data() []byte {
	return []byte{} // return empty slice
}

// SetHasPTS sets the flag that indicates if there is a pts time on the command.
func (c *spliceNull) SetHasPTS(value bool) {
	// do nothing
}

// SetPTS sets the PTS.
func (c *spliceNull) SetPTS(value gots.PTS) {
	// do nothing
}

// spliceTimeBytes returns the raw data bytes of a splice time.
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

// Data returns the bytes of this splice command.
func (c *timeSignal) Data() []byte {
	return spliceTimeBytes(c.hasPTS, c.pts)
}

// SetHasPTS sets the flag that indicates if there is a pts time on the command.
func (c *timeSignal) SetHasPTS(value bool) {
	c.hasPTS = value
}

// SetPTS sets the PTS.
func (c *timeSignal) SetPTS(value gots.PTS) {
	c.pts = value & 0x01ffffffff
}

// Data returns the bytes of this splice command.
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
		componentCount := byte(len(c.components))
		componentsBytes := []byte{componentCount}
		for _, component := range c.components {
			componentBytes := make([]byte, 1)
			componentBytes[0] = component.ComponentTag()
			if c.spliceImmediate {
				componentBytes = append(componentBytes, spliceTimeBytes(component.HasPTS(), component.PTS())...)
			}
			componentsBytes = append(componentsBytes, componentBytes...)
		}
		bytes = append(bytes, componentsBytes...)
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

// SetEventID sets the event id.
func (c *spliceInsert) SetEventID(value uint32) {
	c.eventID = value
}

// SetIsOut sets the value of the out of network indicator
func (c *spliceInsert) SetIsOut(value bool) {
	c.outOfNetworkIndicator = value
}

// SetIsEventCanceled sets the the event cancel indicator
func (c *spliceInsert) SetIsEventCanceled(value bool) {
	c.eventCancelIndicator = value
}

// SetHasPTS sets the flag that indicates if there is a pts time on the command.
func (c *spliceInsert) SetHasPTS(value bool) {
	c.hasPTS = value
}

// SetPTS sets the PTS.
func (c *spliceInsert) SetPTS(value gots.PTS) {
	c.pts = value & 0x01ffffffff
}

// SetHasDuration sets the duration flag
func (c *spliceInsert) SetHasDuration(value bool) {
	c.hasDuration = value
}

// SetDuration sets the PTS duration of the command
func (c *spliceInsert) SetDuration(value gots.PTS) {
	c.duration = value
}

// SetIsAutoReturn sets the boolean value of the auto return field
func (c *spliceInsert) SetIsAutoReturn(value bool) {
	c.autoReturn = value
}

// SetUniqueProgramId sets the unique program Id
func (c *spliceInsert) SetUniqueProgramId(value uint16) {
	c.uniqueProgramId = value
}

// SetAvailNum sets the avail_num field, zero if unused. otherwise index of the avail
func (c *spliceInsert) SetAvailNum(value uint8) {
	c.availNum = value
}

// SetAvailsExpected sets the avails_expected field, number of avails for program
func (c *spliceInsert) SetAvailsExpected(value uint8) {
	c.availsExpected = value
}

// SetIsProgramSplice sets the program splice flag
func (c *spliceInsert) SetIsProgramSplice(value bool) {
	c.isProgramSplice = value
}

// SetSpliceImmediate sets the splice immediate flag
func (c *spliceInsert) SetSpliceImmediate(value bool) {
	c.spliceImmediate = value
}
