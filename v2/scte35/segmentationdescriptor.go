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

	"github.com/Comcast/gots/"
)

// upidSt is the struct used for creating a Multiple UPID (MID)
type upidSt struct {
	upidType SegUPIDType
	upidLen  int
	upid     []byte
}

// CreateUPID will create a default UPID. (SegUPIDNotUsed)
func CreateUPID() UPID {
	return &upidSt{}
}

// UPIDType returns the type of UPID stored
func (u *upidSt) UPIDType() SegUPIDType {
	return u.upidType
}

// UPID returns the actual UPID
func (u *upidSt) UPID() []byte {
	return u.upid
}

// componentOffset is a structure in SegmentationDescriptor.
type componentOffset struct {
	componentTag byte
	ptsOffset    gots.PTS
}

// CreateComponentOffset will create a ComponentOffset structure that
// belongs in a SegmentationDescriptor.
func CreateComponentOffset() ComponentOffset {
	return &componentOffset{}
}

// ComponentTag returns the tag of the component.
func (c *componentOffset) ComponentTag() byte {
	return c.componentTag
}

// PTSOffset returns the PTS offset of the component.
func (c *componentOffset) PTSOffset() gots.PTS {
	return c.ptsOffset
}

// data returns the raw data bytes of the componentOffset
func (c *componentOffset) data() []byte {
	bytes := make([]byte, 6)
	bytes[0] = c.componentTag
	bytes[1] = 0xFE
	bytes[1] |= byte(c.ptsOffset>>32) & 0x01 // 0000 0001
	bytes[2] = byte(c.ptsOffset >> 24)       // 1111 1111
	bytes[3] = byte(c.ptsOffset >> 16)       // 1111 1111
	bytes[4] = byte(c.ptsOffset >> 8)        // 1111 1111
	bytes[5] = byte(c.ptsOffset)             // 1111 1111
	return bytes
}

// componentFromBytes will create a componentOffset from a byte slice.
// length of byte slice must be 5 or larger.
func componentFromBytes(bytes []byte) componentOffset {
	c := componentOffset{}
	c.componentTag = bytes[0]
	pts := (uint64(bytes[1]) << 32 & 0x01) | (uint64(bytes[2]) << 24) |
		(uint64(bytes[3]) << 16) | (uint64(bytes[4]) << 8) | uint64(bytes[5])
	c.ptsOffset = gots.PTS(pts)
	return c
}

// segmentationDescriptor is a strurture representing a segmentation descriptor in SCTE35
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

	programSegmentationFlag bool
	webDeliveryAllowedFlag  bool
	noRegionalBlackoutFlag  bool
	archiveAllowedFlag      bool
	deviceRestrictions      DeviceRestrictions

	components []componentOffset
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

// parseDescriptor will parse a slice of bytes and store the information in a
// parseDescriptor, if possible. If not it will return an error.
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
			// read component info
			ct := readByte()
			if int(ct)*6 > buf.Len()-5 {
				return gots.ErrInvalidSCTE35Length
			}
			for ; ct > 0; ct-- {
				d.components = append(d.components, componentFromBytes(buf.Next(6)))
			}
		}
		if d.hasDuration {
			if buf.Len() < 10 {
				return gots.ErrInvalidSCTE35Length
			}
			d.duration = uint40(buf.Next(5))
		}
		// Upid unneeded now...
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
				UpidElem := upidSt{}
				UpidElem.upidType = SegUPIDType(readByte())
				segUpidLen -= 1
				UpidElem.upidLen = int(readByte())
				segUpidLen -= 1
				UpidElem.upid = buf.Next(UpidElem.upidLen)
				segUpidLen -= UpidElem.upidLen
				d.mid = append(d.mid, UpidElem)
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

// SCTE35 returns the SCTE35 signal this segmentation descriptor was found in.
func (d *segmentationDescriptor) SCTE35() SCTE35 {
	return d.spliceInfo
}

// EventID returns the event id
func (d *segmentationDescriptor) EventID() uint32 {
	return d.eventID
}

// TypeID returns the segmentation type for descriptor
func (d *segmentationDescriptor) TypeID() SegDescType {
	return d.typeID
}

// SetIsEventCanceled sets the the event cancel indicator
func (d *segmentationDescriptor) IsEventCanceled() bool {
	return d.eventCancelIndicator
}

// IsOut returns true if a signal is an out
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

// IsIn returns true if a signal is an in
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

// HasDuration returns true if there is a duration associated with the descriptor
func (d *segmentationDescriptor) HasDuration() bool {
	return d.hasDuration
}

// Duration returns the duration of the descriptor, 40 bit unsigned integer.
func (d *segmentationDescriptor) Duration() gots.PTS {
	return d.duration
}

// UPIDType returns the type of the upid
func (d *segmentationDescriptor) UPIDType() SegUPIDType {
	return d.upidType
}

// UPID returns the upid of the descriptor, if the UPIDType is not SegUPIDMID
func (d *segmentationDescriptor) UPID() []byte {
	// Check if this data should exist
	if d.upidType == SegUPIDMID {
		return []byte{}
	}
	return d.upid
}

// StreamSwitchSignalId returns the signalID of streamswitch signal if
// present in the descriptor
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

// SegmentNum is deprecated, use SegmentNumber instead.
// SegmentNum returns the segment_num descriptor field.
func (d *segmentationDescriptor) SegmentNum() uint8 {
	return d.segNum
}

// CanClose returns true if this descriptor can close the passed in descriptor
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

// Equal determines equality for two segmentation descriptors
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

// HasProgramSegmentation returns if the descriptor has program segmentation
func (d *segmentationDescriptor) HasProgramSegmentation() bool {
	return d.programSegmentationFlag
}

// SegmentNumber returns the segment number for this descriptor.
func (d *segmentationDescriptor) SegmentNumber() uint8 {
	return d.segNum
}

// SegmentsExpected returns the number of expected segments for this descriptor.
func (d *segmentationDescriptor) SegmentsExpected() uint8 {
	return d.segsExpected
}

// HasSubSegments returns true if this segmentation descriptor has subsegment fields.
func (d *segmentationDescriptor) HasSubSegments() bool {
	return d.hasSubSegments
}

// SubSegmentNumber returns the sub-segment number for this descriptor.
func (d *segmentationDescriptor) SubSegmentNumber() uint8 {
	return d.subSegNum
}

// SubSegmentsExpected returns the number of expected sub-segments for this descriptor.
func (d *segmentationDescriptor) SubSegmentsExpected() uint8 {
	return d.subSegsExpected
}

// SetIsDeliveryNotRestricted sets a flag that determines if the delivery is not restricted
func (d *segmentationDescriptor) IsDeliveryNotRestricted() bool {
	return d.deliveryNotRestricted
}

// IsWebDeliveryAllowed returns if web delivery is allowed, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) IsWebDeliveryAllowed() bool {
	return d.webDeliveryAllowedFlag
}

// IsArchiveAllowed returns true if there are restrictions to storing/recording this segment, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) IsArchiveAllowed() bool {
	return d.archiveAllowedFlag
}

// HasNoRegionalBlackout returns true if there is no regional blackout, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) HasNoRegionalBlackout() bool {
	return d.noRegionalBlackoutFlag
}

// DeviceRestrictions returns which device group the segment is restriced to, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) DeviceRestrictions() DeviceRestrictions {
	return d.deviceRestrictions
}

// MID returns multiple UPIDs, if UPIDType is SegUPIDMID
func (d *segmentationDescriptor) MID() []UPID {
	// Check if this data should exist.
	if d.upidType != SegUPIDMID {
		return nil
	}
	mid := make([]UPID, len(d.mid))
	for i := range d.mid {
		mid[i] = &d.mid[i]
	}
	return mid
}

// Components will return components' offsets
func (d *segmentationDescriptor) Components() []ComponentOffset {
	components := make([]ComponentOffset, len(d.components))
	for i := range d.components {
		components[i] = &d.components[i]
	}
	return components
}
