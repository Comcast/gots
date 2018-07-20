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

func CreateSegmentationDescriptor() SegmentationDescriptor {
	return &segmentationDescriptor{}
}

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
				UpidData = append(UpidData, byte(d.mid[i].UpidType))
				UpidData = append(UpidData, byte(d.mid[i].upidLen))
				UpidData = append(UpidData, d.mid[i].Upid...)
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

func (d *segmentationDescriptor) SetEventID(value uint32) {
	d.eventID = value
}

func (d *segmentationDescriptor) SetTypeID(value SegDescType) {
	d.typeID = value
	if d.typeID == 0x34 || d.typeID == 0x36 {
		d.hasSubSegments = true
	} else {
		d.hasSubSegments = false
	}
}

func (d *segmentationDescriptor) SetIsEventCanceled(value bool) {
	d.eventCancelIndicator = value
}

func (d *segmentationDescriptor) SetHasDuration(value bool) {
	d.hasDuration = value
}

func (d *segmentationDescriptor) SetDuration(value gots.PTS) {
	d.duration = value
}

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

func (d *segmentationDescriptor) SetUPID(value []byte) {
	// Check if this data can be set
	if d.UpidType == SegUPIDMID {
		return
	}
	d.Upid = value
}

func (d *segmentationDescriptor) SetSegmentNumber(value uint8) {
	d.segNum = value
}

func (d *segmentationDescriptor) SetSegmentsExpected(value uint8) {
	d.segsExpected = value
}

func (d *segmentationDescriptor) SetSubSegmentNumber(value uint8) {
	d.subSegNum = value
}

func (d *segmentationDescriptor) SetSubSegmentsExpected(value uint8) {
	d.subSegsExpected = value
}

func (d *segmentationDescriptor) SetHasProgramSegmentation(value bool) {
	d.programSegmentationFlag = value
}

func (d *segmentationDescriptor) SetIsDeliveryNotRestricted(value bool) {
	d.deliveryNotRestricted = value
}

func (d *segmentationDescriptor) SetIsWebDeliveryAllowed(value bool) {
	d.webDeliveryAllowedFlag = value
}

func (d *segmentationDescriptor) SetIsArchiveAllowed(value bool) {
	d.archiveAllowedFlag = value
}

func (d *segmentationDescriptor) SetHasNoRegionalBlackout(value bool) {
	d.noRegionalBlackoutFlag = value
}

func (d *segmentationDescriptor) SetDeviceRestrictions(value DeviceRestrictions) {
	d.deviceRestrictions = value
}

func (d *segmentationDescriptor) MID() []UPID {
	// Check if this data should exist.
	if d.UpidType != SegUPIDMID {
		return nil
	}
	mid := make([]UPID, len(d.mid))
	for i := range d.mid {
		mid[i] = &d.mid[i]
	}
	return mid
}

func (d *segmentationDescriptor) SetMID(value []UPID) {
	// Check if this data can be set
	if d.UpidType != SegUPIDMID {
		return
	}
	d.mid = make([]upidSt, len(value))
	for i := range value {
		d.mid[i].UpidType = value[i].UPIDType()
		d.mid[i].Upid = value[i].UPID()
		d.mid[i].upidLen = len(d.mid[i].Upid)
	}
}

func (d *segmentationDescriptor) Components() []ComponentOffset {
	components := make([]ComponentOffset, len(d.components))
	for i := range d.components {
		components[i] = &d.components[i]
	}
	return components
}

func (d *segmentationDescriptor) SetComponents(value []ComponentOffset) {
	d.components = make([]componentOffset, len(value))
	for i := range value {
		d.components[i].componentTag = value[i].ComponentTag()
		d.components[i].ptsOffset = value[i].PTSOffset()
	}
}
