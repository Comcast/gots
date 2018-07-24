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

// SetUPIDType will set the type of the UPID
func (u *upidSt) SetUPIDType(value SegUPIDType) {
	u.upidType = value
}

// UPID set the actual UPID
func (u *upidSt) SetUPID(value []byte) {
	u.upid = value
}

// SetComponentTag sets the component tag, which is used for the identification of the component.
func (c *componentOffset) SetComponentTag(value byte) {
	c.componentTag = value
}

// SetPTS sets the PTS offset of the component.
func (c *componentOffset) SetPTSOffset(value gots.PTS) {
	c.ptsOffset = value
}

// CreateSegmentationDescriptor creates and returns a default
// SegmentationDescriptor.
func CreateSegmentationDescriptor() SegmentationDescriptor {
	return &segmentationDescriptor{}
}

// Data returns the raw data bytes of the SegmentationDescriptor
func (d *segmentationDescriptor) Data() []byte {
	var data, eventData []byte
	data = make([]byte, 11)
	data[0] = segDescTag
	data[2] = byte(segDescID >> 24 & 0xFF)
	data[3] = byte(segDescID >> 16 & 0xFF)
	data[4] = byte(segDescID >> 8 & 0xFF)
	data[5] = byte(segDescID & 0xFF)
	data[6] = byte(d.eventID >> 24)
	data[7] = byte(d.eventID >> 16)
	data[8] = byte(d.eventID >> 8)
	data[9] = byte(d.eventID)
	data[10] = 0x7F // 0111 1111 reserved bits set to 1

	if d.eventCancelIndicator {
		data[10] |= 0x80
	} else {
		eventData = make([]byte, 1)

		if d.deliveryNotRestricted {
			eventData[0] |= 0x20 // 0010 0000
			eventData[0] |= 0x1F // 0001 1111 reserved fields all set to 1
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
				componentsBytes = append(componentsBytes, d.components[i].data()...)
			}
			eventData = append(eventData, componentsBytes...)
		}

		if d.hasDuration {
			eventData[0] |= 0x40 // 0100 0000
			durationBytes := make([]byte, 5)
			durationBytes[0] = byte(d.duration >> 32)
			durationBytes[1] = byte(d.duration >> 24)
			durationBytes[2] = byte(d.duration >> 16)
			durationBytes[3] = byte(d.duration >> 8)
			durationBytes[4] = byte(d.duration)
			eventData = append(eventData, durationBytes...)
		}
		UpidData := make([]byte, 2)
		UpidData[0] = byte(d.UpidType)

		if len(d.Upid) > 0 {
			UpidData = append(UpidData, d.Upid...)
		} else {
			for i := range d.mid {
				UpidData = append(UpidData, byte(d.mid[i].upidType))
				UpidData = append(UpidData, byte(d.mid[i].upidLen))
				UpidData = append(UpidData, d.mid[i].upid...)
			}
		}
		UpidData[1] = byte(len(UpidData) - 2)

		eventData = append(eventData, UpidData...)

		segmentBytes := make([]byte, 3)
		segmentBytes[0] = byte(d.typeID)
		segmentBytes[1] = byte(d.segNum)
		segmentBytes[2] = byte(d.segsExpected)
		if (d.typeID == 0x34) || (d.typeID == 0x36) {
			segmentBytes = append(segmentBytes, d.subSegNum, d.subSegsExpected)
		}
		eventData = append(eventData, segmentBytes...)

		data = append(data, eventData...)

	}
	data[1] = byte(len(data) - 2)
	return data
}

// SetEventID sets the event id
func (d *segmentationDescriptor) SetEventID(value uint32) {
	d.eventID = value
}

// SetTypeID sets the type id
func (d *segmentationDescriptor) SetTypeID(value SegDescType) {
	d.typeID = value
	if d.typeID == 0x34 || d.typeID == 0x36 {
		d.hasSubSegments = true
	} else {
		d.hasSubSegments = false
	}
}

// SetIsEventCanceled sets the the event cancel indicator
func (d *segmentationDescriptor) SetIsEventCanceled(value bool) {
	d.eventCancelIndicator = value
}

// SetHasDuration determines if a duration is present in the descriptor
func (d *segmentationDescriptor) SetHasDuration(value bool) {
	d.hasDuration = value
}

// SetDuration sets the duration of the descriptor
func (d *segmentationDescriptor) SetDuration(value gots.PTS) {
	d.duration = value
}

// SetUPIDType sets the type of upid, only works if UPIDType is not SegUPIDMID
func (d *segmentationDescriptor) SetUPIDType(value SegUPIDType) {
	d.UpidType = value
	// only one can be set at a time
	if d.UpidType == SegUPIDMID {
		d.Upid = []byte{}
	} else if d.UpidType == SegUPIDNotUsed {
		d.mid = []upidSt{}
		d.Upid = []byte{}
	} else {
		d.mid = []upidSt{}
	}
}

// SetUPID sets the upid of the descriptor
func (d *segmentationDescriptor) SetUPID(value []byte) {
	// Check if this data can be set
	if d.UpidType == SegUPIDMID {
		return
	}
	d.Upid = value
}

// SetSegmentNumber sets the segment number for this descriptor.
func (d *segmentationDescriptor) SetSegmentNumber(value uint8) {
	d.segNum = value
}

// SetSegmentsExpected sets the number of expected segments for this descriptor.
func (d *segmentationDescriptor) SetSegmentsExpected(value uint8) {
	d.segsExpected = value
}

// SetSubSegmentNumber sets the sub-segment number for this descriptor.
func (d *segmentationDescriptor) SetSubSegmentNumber(value uint8) {
	d.subSegNum = value
}

// SetSubSegmentsExpected sets the number of expected sub-segments for this descriptor.
func (d *segmentationDescriptor) SetSubSegmentsExpected(value uint8) {
	d.subSegsExpected = value
}

// SetHasProgramSegmentation if the descriptor has program segmentation
func (d *segmentationDescriptor) SetHasProgramSegmentation(value bool) {
	d.programSegmentationFlag = value
}

// SetIsDeliveryNotRestricted sets a flag that determines if the delivery is not restricted
func (d *segmentationDescriptor) SetIsDeliveryNotRestricted(value bool) {
	d.deliveryNotRestricted = value
}

// SetIsWebDeliveryAllowed sets a flag that determines if web delivery is allowed, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) SetIsWebDeliveryAllowed(value bool) {
	d.webDeliveryAllowedFlag = value
}

// SetIsArchiveAllowed sets a flag that determines if there are restrictions to storing/recording this segment, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) SetIsArchiveAllowed(value bool) {
	d.archiveAllowedFlag = value
}

// SetHasNoRegionalBlackout sets a flag that determines if there is no regional blackout, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) SetHasNoRegionalBlackout(value bool) {
	d.noRegionalBlackoutFlag = value
}

// SetDeviceRestrictions sets which device group the segment is restriced to, this field has no meaning if delivery is not restricted.
func (d *segmentationDescriptor) SetDeviceRestrictions(value DeviceRestrictions) {
	d.deviceRestrictions = value
}

// SetMID sets multiple UPIDs, only works if UPIDType is SegUPIDMID
func (d *segmentationDescriptor) SetMID(value []UPID) {
	// Check if this data can be set
	if d.UpidType != SegUPIDMID {
		return
	}
	d.mid = make([]upidSt, len(value))
	for i := range value {
		d.mid[i].upidType = value[i].UPIDType()
		d.mid[i].upid = value[i].UPID()
		d.mid[i].upidLen = len(d.mid[i].upid)
	}
}

// Components will set components' offsets
func (d *segmentationDescriptor) SetComponents(value []ComponentOffset) {
	d.components = make([]componentOffset, len(value))
	for i := range value {
		d.components[i].componentTag = value[i].ComponentTag()
		d.components[i].ptsOffset = value[i].PTSOffset()
	}
}
