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
	"strings"

	"github.com/Comcast/gots"
)

// This is the struct used for creating a Multiple UPID (MID)
type upidSt struct {
	upidType SegUPIDType
	upidLen  int
	upid     []byte
}

type Component struct {
	ComponentTag byte
	PtsOffset    uint64
}

func (c *Component) Bytes() []byte {
	bytes := make([]byte, 6)
	bytes[0] = c.ComponentTag
	bytes[1] |= byte(c.PtsOffset>>32) & 0x01 // 0000 0001
	bytes[2] = byte(c.PtsOffset >> 24)       // 1111 1111
	bytes[3] = byte(c.PtsOffset >> 16)       // 1111 1111
	bytes[4] = byte(c.PtsOffset >> 8)        // 1111 1111
	bytes[5] = byte(c.PtsOffset)             // 1111 1111
	return bytes
}

type segmentationDescriptor struct {
	// common fields we care about for sorting/identifying, but is not necessarily needed for users of this lib
	typeID                SegDescType
	eventID               uint32
	hasDuration           bool
	duration              gots.PTS
	upidType              SegUPIDType
	upid                  []byte
	mid                   []upidSt //A MID can contains `n` UPID uids in it.
	segNum                uint8
	segsExpected          uint8
	subSegNum             uint8
	subSegsExpected       uint8
	spliceInfo            SCTE35
	eventCancelIndicator  bool
	deliveryNotRestricted bool
	hasSubSegments        bool

	programSegmentationFlag  bool
	segmentationDurationFlag bool
	webDeliveryAllowedFlag   bool
	noRegionalBlackoutFlag   bool
	archiveAllowedFlag       bool
	deviceRestrictions       DeviceRestrictions

	components []Component
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
		0x10: {0x10: segCloseNoBreakaway, 0x14: segCloseNormal, 0x17: segCloseNoBreakaway, 0x19: segCloseNoBreakaway, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x11: {0x10: segCloseEventID, 0x14: segCloseEventID, 0x17: segCloseEventID, 0x19: segCloseEventID, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x12: {0x10: segCloseEventID, 0x14: segCloseEventID, 0x17: segCloseEventID, 0x19: segCloseEventID, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x13: {0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x14: {0x10: segCloseBreakaway, 0x17: segCloseBreakaway, 0x19: segCloseBreakaway, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x19: {0x10: segCloseNoBreakaway, 0x14: segCloseNormal, 0x17: segCloseNoBreakaway, 0x19: segCloseNoBreakaway, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x20: {0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x21: {0x20: segCloseEventID, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal},
		0x22: {0x22: segCloseNormal},
		0x23: {0x22: segCloseEventID},
		0x30: {0x30: segCloseNormal, 0x32: segCloseNormal},
		0x31: {0x30: segCloseEventID},
		0x32: {0x30: segCloseNormal, 0x32: segCloseNormal},
		0x33: {0x32: segCloseEventID},
		0x34: {0x30: segCloseDiffPTS, 0x32: segCloseDiffPTS, 0x34: segCloseNotNested, 0x36: segCloseNormal},
		0x35: {0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseEventIDNotNested, 0x36: segCloseNormal},
		0x36: {0x30: segCloseDiffPTS, 0x32: segCloseDiffPTS, 0x36: segCloseNotNested},
		0x37: {0x30: segCloseNormal, 0x32: segCloseNormal, 0x36: segCloseEventIDNotNested},
		0x40: {0x40: segCloseNormal},
		0x41: {0x40: segCloseEventID},
		0x50: {0x10: segCloseNormal, 0x14: segCloseNormal, 0x17: segCloseNormal, 0x19: segCloseNormal, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal, 0x40: segCloseUnconditional, 0x50: segCloseNormal},
		0x51: {0x10: segCloseNormal, 0x14: segCloseNormal, 0x17: segCloseNormal, 0x19: segCloseNormal, 0x20: segCloseNormal, 0x30: segCloseNormal, 0x32: segCloseNormal, 0x34: segCloseNormal, 0x36: segCloseNormal, 0x40: segCloseUnconditional, 0x50: segCloseEventID},
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
	d.eventCancelIndicator = readByte()&0x80 != 0
	if !d.eventCancelIndicator {
		flags := readByte()
		// 3rd bit in the byte
		// if delivery_not_restricted = 0 -> deliveryNotRestricted = false
		d.deliveryNotRestricted = flags&0x20 != 0
		d.hasDuration = flags&0x40 != 0
		d.programSegmentationFlag = flags&0x80 != 0
		if !d.deliveryNotRestricted {
			d.webDeliveryAllowedFlag = flags&0x10 != 0
			d.noRegionalBlackoutFlag = flags&0x08 != 0
			d.archiveAllowedFlag = flags&0x04 != 0
			d.deviceRestrictions = DeviceRestrictions(flags & 0x03)
		}
		if !d.programSegmentationFlag {
			// skip over component info
			ct := readByte()
			if int(ct)*6 > buf.Len()-5 {
				return gots.ErrInvalidSCTE35Length
			}
			for ; ct > 0; ct-- {
				buf.Next(6)
			}
		}
		if d.hasDuration {
			if buf.Len() < 10 {
				return gots.ErrInvalidSCTE35Length
			}
			d.duration = uint40(buf.Next(5))
		}
		// upid unneeded now...
		d.upidType = SegUPIDType(readByte())
		segUpidLen := int(readByte())
		d.mid = []upidSt{}
		// This is a Multiple PID, consisting of `n` UPID's
		if d.upidType == SegUPIDMID {
			// SCTE35 signal will either have an UPID or a MID
			// When we have a MID, UPID value in struct will be 0.
			d.upid = []byte{}
			// Iterate over the whole MID len(segUpidLen) to get all `n` UPIDs
			// segUpidLen is in bytes.
			for segUpidLen != 0 {
				upidElem := upidSt{}
				upidElem.upidType = SegUPIDType(readByte())
				segUpidLen -= 1
				upidElem.upidLen = int(readByte())
				segUpidLen -= 1
				upidElem.upid = buf.Next(upidElem.upidLen)
				segUpidLen -= upidElem.upidLen
				d.mid = append(d.mid, upidElem)
			}
		} else {
			// This is a UPID, not a MID
			// MID should be 0 as SCTE35 can either have a UPID or a MID
			if buf.Len() < segUpidLen+3 {
				return gots.ErrInvalidSCTE35Length
			}
			d.upid = buf.Next(segUpidLen)
		}
		d.typeID = SegDescType(readByte())
		d.segNum = readByte()
		d.segsExpected = readByte()

		// Backwards compatible support for the 2016 spec
		if buf.Len() > 0 && (d.typeID == 0x34 || d.typeID == 0x36) {
			d.subSegNum = readByte()
			d.subSegsExpected = readByte()
			d.hasSubSegments = true
		}
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

func (d *segmentationDescriptor) IsEventCanceled() bool {
	return d.eventCancelIndicator
}

func (d *segmentationDescriptor) IsOut() bool {
	switch d.TypeID() {
	case SegDescProgramStart,
		SegDescProgramResumption,
		SegDescProgramOverlapStart,
		SegDescProgramStartInProgress,
		SegDescChapterStart,
		SegDescBreakStart,
		SegDescProviderAdvertisementStart,
		SegDescDistributorAdvertisementStart,
		SegDescProviderPOStart,
		SegDescDistributorPOStart,
		SegDescUnscheduledEventStart,
		SegDescNetworkStart:
		return true
	default:
		return false
	}
}

func (d *segmentationDescriptor) IsIn() bool {
	switch d.TypeID() {
	case SegDescProgramEnd,
		SegDescProgramEarlyTermination,
		SegDescProgramBreakaway,
		SegDescProgramRunoverPlanned,
		SegDescProgramRunoverUnplanned,
		SegDescProgramBlackoutOverride,
		SegDescChapterEnd,
		SegDescBreakEnd,
		SegDescProviderAdvertisementEnd,
		SegDescDistributorAdvertisementEnd,
		SegDescProviderPOEnd,
		SegDescDistributorPOEnd,
		SegDescUnscheduledEventEnd,
		SegDescNetworkEnd:
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

func (d *segmentationDescriptor) StreamSwitchSignalId() (string, error) {
	var signalId string
	var err error
	// The VSS SignalId is present in the MID of len 2.
	// SignalId is the UPID value in the MID which has
	// delivery_not_restricted flag = 0 and
	// contains "BLACKOUT" at UpidType of 0x09 and
	// comcast:linear:licenserotation at 0x0E
	if len(d.mid) == 2 &&
		!d.deliveryNotRestricted &&
		(d.mid[0].upidType == SegUPIDADI) &&
		(strings.Contains(string(d.mid[0].upid), "BLACKOUT")) &&
		(d.mid[1].upidType == SegUPADSINFO) &&
		(strings.Contains(string(d.mid[1].upid), "comcast:linear:licenserotation")) {
		signalId = strings.TrimPrefix(string(d.mid[0].upid), "BLACKOUT:")
	} else {
		err = gots.ErrVSSSignalIdNotFound
	}
	return signalId, err
}

func (d *segmentationDescriptor) SegmentNum() uint8 {
	return d.segNum
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
		// this should also consider segnum == segexpected for IN signals closing an out signal.
		if d.IsIn() && d.EventID() == out.EventID() && d.SegmentNumber() == d.SegmentsExpected() {
			return true
		}
	case segCloseNotNested: // only applies to 0x34 and 0x36 with subsegments.
		if d.HasSubSegments() {
			if d.SubSegmentNumber() == d.SubSegmentsExpected() {
				return true
			} else {
				return false
			}
		}
		return true
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
	if d.SegmentNumber() != c.SegmentNumber() {
		return false
	}
	if d.SegmentsExpected() != c.SegmentsExpected() {
		return false
	}
	if d.HasSubSegments() != c.HasSubSegments() {
		return false
	}
	if d.HasSubSegments() && c.HasSubSegments() && (d.SubSegmentNumber() != c.SubSegmentNumber()) {
		return false
	}
	if d.HasSubSegments() && c.HasSubSegments() && (d.SubSegmentsExpected() != c.SubSegmentsExpected()) {
		return false
	}
	return true
}

func (d *segmentationDescriptor) SegmentNumber() uint8 {
	return d.segNum
}

func (d *segmentationDescriptor) SegmentsExpected() uint8 {
	return d.segsExpected
}

func (d *segmentationDescriptor) HasSubSegments() bool {
	return d.hasSubSegments
}

func (d *segmentationDescriptor) SubSegmentNumber() uint8 {
	return d.subSegNum
}

func (d *segmentationDescriptor) SubSegmentsExpected() uint8 {
	return d.subSegsExpected
}

func (d *segmentationDescriptor) Bytes() []byte {
	var data, eventData []byte
	data = make([]byte, 11)
	data[0] = segDescTag
	data[2] = byte(segDescID >> 24 & 0xFF)
	data[3] = byte(segDescID >> 16 & 0xFF)
	data[4] = byte(segDescID >> 8 & 0xFF)
	data[5] = byte(segDescID & 0xFF)
	data[6] = byte(d.eventID >> 24 & 0xFF)
	data[7] = byte(d.eventID >> 16 & 0xFF)
	data[8] = byte(d.eventID >> 8 & 0xFF)
	data[9] = byte(d.eventID & 0xFF)
	data[10] = 0x7F // reserved bits

	if d.eventCancelIndicator {
		data[10] |= 0x80
	} else {
		eventData = make([]byte, 1)

		if d.deliveryNotRestricted {
			eventData[0] |= 0x20 // 0010 0000
			eventData[0] |= 0x1F // 0001 1111 reserved fields
		} else {
			if d.webDeliveryAllowedFlag {
				eventData[0] |= 0x10 // 0001 0000
			}
			if d.noRegionalBlackoutFlag {
				eventData[0] |= 0x08 // 0000 1000
			}
			if d.archiveAllowedFlag {
				eventData[0] |= 0x04 // 0000 0100
			}
			eventData[0] |= byte(d.deviceRestrictions) & 0x03 // 0000 0011
		}

		if d.programSegmentationFlag {
			eventData[0] |= 0x80 // 1000 0000
		} else {
			componentsBytes := make([]byte, 1, 6*len(d.components)+1)
			componentsBytes[0] = byte(len(d.components)) // set component count
			for i := range d.components {
				componentsBytes = append(componentsBytes, d.components[i].Bytes()...)
			}
			eventData = append(eventData, componentsBytes...)
		}

		if d.segmentationDurationFlag {
			eventData[0] |= 0x40 // 0100 0000
			durationBytes := make([]byte, 5)
			durationBytes[0] = byte(d.duration >> 32)
			durationBytes[1] = byte(d.duration >> 24)
			durationBytes[2] = byte(d.duration >> 16)
			durationBytes[3] = byte(d.duration >> 8)
			durationBytes[4] = byte(d.duration)
			eventData = append(eventData, durationBytes...)
		}
		upidData := make([]byte, 2)
		upidData[0] = byte(d.upidType)
		upidData[1] = byte(len(d.upid))
		upidData = append(upidData, d.upid...)
		eventData = append(eventData, upidData...)

		segmentBytes := make([]byte, 3)
		segmentBytes[0] = byte(d.typeID)
		segmentBytes[1] = byte(d.segNum)
		segmentBytes[2] = byte(d.segsExpected)
		if d.typeID == 0x34 || d.typeID == 0x36 {
			segmentBytes = append(segmentBytes, d.subSegNum, d.subSegsExpected)
		}
		eventData = append(eventData, segmentBytes...)

		data = append(data, eventData...)

	}
	data[1] = byte(len(data))
	return data
}
