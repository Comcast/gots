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

	"github.com/comcast/gots"
)

type segmentationDescriptor struct {
	// common fields we care about for sorting/identifying, but is not necessarily needed for users of this lib
	typeID       SegDescType
	eventID      uint32
	hasDuration  bool
	duration     gots.PTS
	upidType     SegUPIDType
	upid         []byte
	segNum       uint8
	segsExpected uint8
	spliceInfo   SCTE35
}

type segCloseType uint8

// conditions for closing specific descriptor types
const (
	segCloseNormal segCloseType = iota
	segCloseNoBreakaway
	segCloseEventID
	segCloseBreakaway
	segCloseDiffPTS
	segCloseNotNested
	segCloseEventIDNotNested
	segCloseUnconditional
)

var segCloseRules map[SegDescType]map[SegDescType]segCloseType

// initialize the SegCloseRules map
func init() {
	segCloseRules = map[SegDescType]map[SegDescType]segCloseType{
		0x10: {0x10: segCloseNoBreakaway, 0x14: segCloseNormal, 0x17: segCloseNoBreakaway, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x11: {0x10: segCloseEventID, 0x14: segCloseEventID, 0x17: segCloseEventID, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x12: {0x10: segCloseEventID, 0x14: segCloseEventID, 0x17: segCloseEventID, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x13: {0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x14: {0x10: segCloseBreakaway, 0x17: segCloseBreakaway, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x20: {0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x21: {0x20: segCloseEventID, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x30: {0x30: segCloseNormal, 0x32: segCloseNormal},
		0x31: {0x30: segCloseEventID},
		0x32: {0x30: segCloseNormal, 0x32: segCloseNormal},
		0x33: {0x32: segCloseEventID},
		0x34: {0x30: segCloseDiffPTS, 0x32: segCloseDiffPTS, 0x34: segCloseNotNested, 0x36: segCloseNormal},
		0x35: {0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseEventIDNotNested, 0x36: segCloseNormal},
		0x36: {0x30: segCloseDiffPTS, 0x32: segCloseDiffPTS, 0x36: segCloseNotNested},
		0x37: {0x30: segCloseNormal, 0x32: segCloseNormal, 0x36: segCloseEventIDNotNested},
		0x40: {0x10: segCloseNormal, 0x14: segCloseNormal, 0x17: segCloseNormal, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal, 0x40: segCloseNormal},
		0x41: {0x10: segCloseNormal, 0x14: segCloseNormal, 0x17: segCloseNormal, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal, 0x40: segCloseEventID},
		0x50: {0x10: segCloseNormal, 0x14: segCloseNormal, 0x17: segCloseNormal, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal, 0x40: segCloseUnconditional, 0x50: segCloseNormal},
		0x51: {0x10: segCloseNormal, 0x14: segCloseNormal, 0x17: segCloseNormal, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal, 0x40: segCloseUnconditional, 0x50: segCloseEventID},
	}
	// add program breakaway rules to the close map.  Only ProgramResumption,
	// Unscheduled Event and Network signals can exit breakaway.
	// Program resumption should not close ProgramBreakaway as this logic will be
	// handled in state.  We want processing to stop so that it won't close
	// events after the breakaway
	segCloseRules[SegDescUnscheduledEventStart][0x13] = segCloseNormal
	segCloseRules[SegDescUnscheduledEventEnd][0x13] = segCloseNormal
	segCloseRules[SegDescNetworkStart][0x13] = segCloseNormal
	segCloseRules[SegDescNetworkEnd][0x13] = segCloseNormal
}

func (d *segmentationDescriptor) parseDescriptor(data []byte) error {
	buf := bytes.NewBuffer(data)
	// closure to ignore EOF error from buf.ReadByte().  We've already checked
	// the length, we don't need to continually check it after every ReadByte
	// call
	readByte := func() byte {
		b, _ := buf.ReadByte()
		return b
	}
	if binary.BigEndian.Uint32(buf.Next(4)) != segDescID {
		return gots.ErrSCTE35InvalidDescriptorID
	}
	d.eventID = binary.BigEndian.Uint32(buf.Next(4))
	if readByte()&0x80 == 0 { // Cancel indicator
		flags := readByte()
		if flags&0x80 == 0 {
			// skip over component info
			ct := readByte()
			if int(ct)*6 > buf.Len()-5 {
				return gots.ErrInvalidSCTE35Length
			}
			for ; ct > 0; ct-- {
				buf.Next(6)
			}
		}
		d.hasDuration = flags&0x40 != 0
		if d.hasDuration {
			if buf.Len() < 10 {
				return gots.ErrInvalidSCTE35Length
			}
			d.duration = uint40(buf.Next(5))
		}
		// upid unneeded now...
		d.upidType = SegUPIDType(readByte())
		upidLen := int(readByte())
		if buf.Len() < upidLen+3 {
			return gots.ErrInvalidSCTE35Length
		}
		d.upid = buf.Next(upidLen)
		d.typeID = SegDescType(readByte())
		d.segNum = readByte()
		d.segsExpected = readByte()
	}
	return nil
}

func (d *segmentationDescriptor) SCTE35() SCTE35 {
	return d.spliceInfo
}

func (d *segmentationDescriptor) EventID() uint32 {
	return d.eventID
}

func (d *segmentationDescriptor) TypeID() SegDescType {
	return d.typeID
}

func (d *segmentationDescriptor) IsOut() bool {
	switch d.TypeID() {
	case SegDescProgramStart, SegDescChapterStart,
		SegDescProviderAdvertisementStart, SegDescDistributorAdvertisementStart,
		SegDescProviderPOStart, SegDescDistributorPOStart,
		SegDescUnscheduledEventStart, SegDescNetworkStart,
		SegDescProgramOverlapStart, SegDescProgramResumption:
		return true
	default:
		return false
	}
}

func (d *segmentationDescriptor) IsIn() bool {
	switch d.TypeID() {
	case SegDescProgramEnd, SegDescChapterEnd,
		SegDescProviderAdvertisementEnd, SegDescDistributorAdvertisementEnd,
		SegDescProviderPOEnd, SegDescDistributorPOEnd,
		SegDescUnscheduledEventEnd, SegDescNetworkEnd,
		SegDescProgramRunoverPlanned, SegDescProgramRunoverUnplanned,
		SegDescProgramEarlyTermination, SegDescProgramBreakaway:
		return true
	default:
		return false
	}
}

func (d *segmentationDescriptor) HasDuration() bool {
	return d.hasDuration
}

func (d *segmentationDescriptor) Duration() gots.PTS {
	return d.duration
}

func (d *segmentationDescriptor) UPIDType() SegUPIDType {
	return d.upidType
}

func (d *segmentationDescriptor) UPID() []byte {
	return d.upid
}

func (d *segmentationDescriptor) CanClose(out SegmentationDescriptor) bool {
	inRules, ok := segCloseRules[d.TypeID()]
	// No rules associated with this signal means it can't close anything
	if !ok {
		return false
	}
	closeType, ok := inRules[out.TypeID()]
	// out signal type not found in rule set means d can't close out
	if !ok {
		return false
	}
	switch closeType {
	case segCloseNormal, segCloseUnconditional:
		return true
	// These close rules around what can/cannot go past a ProgramBreakaway.
	// Since we will handle program breakaway logic in State, we leave these as
	// automatic trues
	case segCloseBreakaway, segCloseNoBreakaway:
		return true
	case segCloseEventID:
		if d.EventID() == out.EventID() {
			return true
		}
	case segCloseDiffPTS:
		if d.SCTE35().PTS() != out.SCTE35().PTS() {
			return true
		}
	case segCloseEventIDNotNested:
		if d.EventID() != out.EventID() {
			return false
		}
		fallthrough
	case segCloseNotNested:
		if d.IsIn() {
			// if in, last desc is x/x
			if d.segNum == d.segsExpected {
				return true
			}
		} else if d.IsOut() {
			// if out, first descriptor in set closes existing open
			if d.segNum == 1 {
				return true
			}
		}
	}
	return false
}

// Determines equality for two segmentation descriptors
// Equality in this sense means that two "in" events are duplicates
// A lot of debate went in to what actually constitutes a "duplicate".  We get
// all sorts of things from different providers/transcoders, so in the end, we
// just settled on PTS, event id, and segmentation type
func (d *segmentationDescriptor) Equal(c SegmentationDescriptor) bool {
	if d == nil || c == nil {
		return false
	}
	if d.TypeID() != c.TypeID() {
		return false
	}
	if !d.SCTE35().HasPTS() || !c.SCTE35().HasPTS() {
		return false
	}
	if d.SCTE35().PTS() != c.SCTE35().PTS() {
		return false
	}
	if d.EventID() != c.EventID() {
		return false
	}
	return true
}
