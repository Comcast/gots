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
	"bytes"
	"encoding/binary"

	"github.com/Comcast/gots"
)

type timeSignal struct {
	hasPTS bool
	pts    gots.PTS
}

func (c *timeSignal) CommandType() SpliceCommandType {
	return TimeSignal
}

// parseTimeSignal extracts a time_signal() command from a bytes buffer. It
// returns a timeSignal describing the command.
func parseTimeSignal(buf *bytes.Buffer) (cmd *timeSignal, err error) {
	cmd = &timeSignal{}
	hasPTS, pts, err := parseSpliceTime(buf)
	if !hasPTS {
		return nil, gots.ErrSCTE35UnsupportedSpliceCommand
	}
	cmd.hasPTS = hasPTS
	cmd.pts = pts
	return cmd, nil
}

func (c *timeSignal) HasPTS() bool {
	return c.hasPTS
}

func (c *timeSignal) PTS() gots.PTS {
	return c.pts
}

type spliceNull struct {
}

func (c *spliceNull) CommandType() SpliceCommandType {
	return SpliceNull
}

func (c *spliceNull) HasPTS() bool {
	return false
}

func (c *spliceNull) PTS() gots.PTS {
	return 0
}

type spliceInsert struct {
	eventID               uint32
	eventCancelIndicator  bool
	outOfNetworkIndicator bool
	hasPTS                bool
	pts                   gots.PTS
	hasDuration           bool
	duration              gots.PTS
	autoReturn            bool
	uniqueProgramId       uint16
	availNum              uint8
	availsExpected        uint8
}

func (c *spliceInsert) CommandType() SpliceCommandType {
	return SpliceInsert
}

// parseSpliceInsert extracts a splice_insert() command from a bytes buffer.
// It returns a spliceInsert describing the command.
func parseSpliceInsert(buf *bytes.Buffer) (*spliceInsert, error) {
	cmd := &spliceInsert{}
	if err := cmd.parse(buf); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *spliceInsert) parse(buf *bytes.Buffer) error {
	baseFields := buf.Next(5)
	if len(baseFields) < 5 { // length of required fields
		return gots.ErrInvalidSCTE35Length
	}
	c.eventID = binary.BigEndian.Uint32(baseFields[:4])
	// splice_event_cancel_indicator 1
	// reserved 7
	c.eventCancelIndicator = baseFields[4]&0x80 == 0x80
	if c.eventCancelIndicator {
		return nil
	}
	// out_of_network_indicator 1
	// program_splice_flag 1
	// duration_flag 1
	// splice_immediate_flag 1
	// reserved 4
	flags, err := buf.ReadByte()
	if err != nil {
		return gots.ErrInvalidSCTE35Length
	}
	c.outOfNetworkIndicator = flags&0x80 == 0x80
	isProgramSplice := flags&0x40 == 0x40
	c.hasDuration = flags&0x20 == 0x20
	spliceImmediate := flags&0x10 == 0x10

	if isProgramSplice && !spliceImmediate {
		hasPTS, pts, err := parseSpliceTime(buf)
		if err != nil {
			return err
		}
		if !hasPTS {
			return gots.ErrSCTE35UnsupportedSpliceCommand
		}
		c.hasPTS = hasPTS
		c.pts = pts
	}
	if !isProgramSplice {
		cc, err := buf.ReadByte()
		if err != nil {
			return gots.ErrInvalidSCTE35Length
		}
		// skip components for now
		for ; cc > 0; cc-- {
			// component_tag
			if _, err := buf.ReadByte(); err != nil {
				return gots.ErrInvalidSCTE35Length
			}
			if !spliceImmediate {
				if _, _, err := parseSpliceTime(buf); err != nil {
					return err
				}
			}
		}
	}
	if c.hasDuration {
		data := buf.Next(5)
		if len(data) < 5 {
			return gots.ErrInvalidSCTE35Length
		}
		// break_duration() structure:
		c.autoReturn = data[0]&0x80 == 0x80
		c.duration = uint40(data) & 0x01ffffffff
	}
	progInfo := buf.Next(4)
	if len(progInfo) < 4 {
		return gots.ErrInvalidSCTE35Length
	}
	c.uniqueProgramId = binary.BigEndian.Uint16(progInfo[:2])
	c.availNum = progInfo[2]
	c.availsExpected = progInfo[3]
	return nil
}

func (c *spliceInsert) EventID() uint32 {
	return c.eventID
}

func (c *spliceInsert) IsOut() bool {
	return c.outOfNetworkIndicator
}

func (c *spliceInsert) IsEventCanceled() bool {
	return c.eventCancelIndicator
}

func (c *spliceInsert) HasPTS() bool {
	return c.hasPTS
}

func (c *spliceInsert) PTS() gots.PTS {
	return c.pts
}

func (c *spliceInsert) HasDuration() bool {
	return c.hasDuration
}

func (c *spliceInsert) Duration() gots.PTS {
	return c.duration
}

func (c *spliceInsert) IsAutoReturn() bool {
	return c.autoReturn
}

func (c *spliceInsert) UniqueProgramId() uint16 {
	return c.uniqueProgramId
}

func (c *spliceInsert) AvailNum() uint8 {
	return c.availNum
}

func (c *spliceInsert) AvailsExpected() uint8 {
	return c.availsExpected
}

// parseSpliceTime parses a splice_time() structure and returns the values of
// time_specified_flag and pts_time.
// If the time_specified_flag is 0, pts will have a value of gots.PTS(0).
func parseSpliceTime(buf *bytes.Buffer) (timeSpecified bool, pts gots.PTS, err error) {
	flags, err := buf.ReadByte()
	if err != nil {
		err = gots.ErrInvalidSCTE35Length
		return false, gots.PTS(0), err
	}
	timeSpecified = flags&0x80 == 0x80
	if !timeSpecified {
		// Time isn't specified, assume PTS of 0.
		return false, gots.PTS(0), nil
	}
	// unread prev byte because it contains the top bit of the pts offset
	if err = buf.UnreadByte(); err != nil {
		return true, gots.PTS(0), err
	}
	if buf.Len() < 5 {
		err = gots.ErrInvalidSCTE35Length
		return
	}
	pts = uint40(buf.Next(5)) & 0x01ffffffff
	return true, pts, nil
}
