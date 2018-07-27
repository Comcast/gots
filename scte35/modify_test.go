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
	"github.com/Comcast/gots"
	"github.com/Comcast/gots/psi"
	"testing"
)

var testScteCreate = []byte{
	0x00, 0xFC, 0x30, 0x27, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xF0, 0x05, 0x06, 0xFE,
	0x86, 0xDF, 0x75, 0x50, 0x00, 0x11, 0x02, 0x0F, 0x43, 0x55, 0x45, 0x49, 0x41, 0x42, 0x43, 0x44,
	0x7F, 0x8F, 0x00, 0x00, 0x10, 0x01, 0x01, 0x0B, 0xFD, 0xD1, 0x40,
}

var testScteCreate2 = []byte{
	0x00, 0xfc, 0x30, 0x30, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x20, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x1a, 0x02, 0x18, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x02,
	0x7f, 0xff, 0x00, 0x00, 0x0a, 0xff, 0x50, 0x09, 0x04, 0x54, 0x45, 0x53, 0x54, 0x40, 0x00, 0x00,
	0x25, 0x12, 0xF4, 0x01,
}

var testScteCreate3 = []byte{
	0x00, 0xfc, 0x30, 0x30, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x20, 0x00, 0x00, 0x00, 0x05, 0x06, 0xfe,
	0x00, 0x02, 0xbf, 0xd4, 0x00, 0x1a, 0x02, 0x18, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x02,
	0x7f, 0xff, 0x00, 0x00, 0x0a, 0xff, 0x50, 0x09, 0x04, 0x54, 0x45, 0x53, 0x54, 0x51, 0x00, 0x00,
	0x39, 0x40, 0x90, 0xF6,
}

var testScteCreate4 = []byte{
	0x00, 0xFC, 0x30, 0x53, 0x00, 0x00, 0x00, 0x02, 0xDD, 0x20, 0x00, 0xFF, 0xF0, 0x05, 0x06, 0xFE,
	0x00, 0x08, 0x95, 0x44, 0x00, 0x3D, 0x02, 0x3B, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x02,
	0x7F, 0x1F, 0x02, 0x01, 0xFE, 0x00, 0x2D, 0xD2, 0x00, 0x02, 0xFE, 0x00, 0x00, 0x01, 0xE8, 0x09,
	0x1F, 0x53, 0x49, 0x47, 0x4E, 0x41, 0x4C, 0x3A, 0x59, 0x38, 0x6F, 0x30, 0x44, 0x33, 0x7A, 0x70,
	0x54, 0x78, 0x53, 0x30, 0x4C, 0x54, 0x31, 0x65, 0x77, 0x2B, 0x77, 0x75, 0x69, 0x77, 0x3D, 0x3D,
	0x36, 0x00, 0x00, 0x56, 0x50, 0xE1, 0xED,
}

var testRollOverScteAdjustment = []byte{
	0xFC, 0x30, 0x16, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, 0x00, 0xFF, 0xF0, 0x05, 0x06, 0xFF, 0xFF,
	0xFF, 0xF1, 0xF1, 0x00, 0x00, 0x1B, 0xD0, 0x87, 0x0B,
}

func TestPTSAdjustmentWithRollover(t *testing.T) {
	target := testRollOverScteAdjustment
	scte := CreateSCTE35()
	initialPTS := gots.PTS(0x1FFFFF1F1)
	adjustmentPTS := gots.PTS(0xF0F)          // should be derived by scte35
	finalPTS := initialPTS.Add(adjustmentPTS) // should be 256 (0x100)
	scte.SetAdjustPTS(finalPTS)
	cmd := CreateTimeSignalCommand()
	cmd.SetHasPTS(true)
	cmd.SetPTS(initialPTS)
	scte.SetCommandInfo(cmd)
	generated := scte.UpdateData()
	if !bytes.Equal(target, generated) {
		t.Errorf("Generated packet data does not match expected data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestDistributorPoStartCreateEncode(t *testing.T) {
	target := testScteCreate4
	scte := CreateSCTE35()
	scte.SetAdjustPTS(0xB7264)
	scte.SetTier(0xFFF)
	cmd := CreateTimeSignalCommand()
	cmd.SetHasPTS(true)
	cmd.SetPTS(0x89544)
	scte.SetCommandInfo(cmd)

	descriptor := CreateSegmentationDescriptor()
	descriptor.SetEventID(0x2)
	descriptor.SetIsEventCanceled(false)
	descriptor.SetHasProgramSegmentation(false)
	descriptor.SetHasDuration(false)
	descriptor.SetIsDeliveryNotRestricted(false)
	descriptor.SetIsWebDeliveryAllowed(true)
	descriptor.SetHasNoRegionalBlackout(true)
	descriptor.SetIsArchiveAllowed(true)
	descriptor.SetDeviceRestrictions(RestrictNone)
	descriptor.SetUPIDType(SegUPIDADI)
	descriptor.SetUPID([]byte("SIGNAL:Y8o0D3zpTxS0LT1ew+wuiw=="))
	descriptor.SetTypeID(SegDescDistributorPOStart)

	components := make([]ComponentOffset, 0)
	component := CreateComponentOffset()
	component.SetComponentTag(0x1)
	component.SetPTSOffset(0x2DD200)
	components = append(components, component)
	component = CreateComponentOffset()
	component.SetComponentTag(0x2)
	component.SetPTSOffset(0x1E8)
	components = append(components, component)
	descriptor.SetComponents(components)

	descriptor.SetSegmentNumber(0)
	descriptor.SetSegmentsExpected(0)
	descriptor.SetSubSegmentNumber(0)
	descriptor.SetSubSegmentsExpected(0)

	scte.SetDescriptors([]SegmentationDescriptor{descriptor})

	generated := append(psi.NewPointerField(0), scte.UpdateData()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated packet data does not match expected data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestDistributorPoStartDecodeEncode(t *testing.T) {
	target := testScteCreate4
	scte, err := NewSCTE35(target)
	if err != nil {
		t.Error(err.Error())
	}
	scte.UpdateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Original packet data does not match Generated data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestProgramStartDecodeEncode(t *testing.T) {
	target := testScteCreate4
	scte, err := NewSCTE35(target)
	if err != nil {
		t.Error(err.Error())
	}
	scte.UpdateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Original packet data does not match Generated data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestNetworkEndCreateEncode(t *testing.T) {
	target := testScteCreate3
	scte := CreateSCTE35()
	scte.SetAdjustPTS(0x59CF4)
	scte.SetTier(0x0)
	cmd := CreateTimeSignalCommand()
	cmd.SetHasPTS(true)
	cmd.SetPTS(0x2BFD4)
	scte.SetCommandInfo(cmd)

	descriptor := CreateSegmentationDescriptor()
	descriptor.SetEventID(0x2)
	descriptor.SetIsEventCanceled(false)
	descriptor.SetHasProgramSegmentation(true)
	descriptor.SetHasDuration(true)
	descriptor.SetDuration(0xAFF50)
	descriptor.SetIsDeliveryNotRestricted(true)
	descriptor.SetUPIDType(SegUPIDADI)
	descriptor.SetUPID([]byte("TEST"))
	descriptor.SetTypeID(SegDescNetworkEnd)
	descriptor.SetSegmentNumber(0)
	descriptor.SetSegmentsExpected(0)

	scte.SetDescriptors([]SegmentationDescriptor{descriptor})

	generated := append(psi.NewPointerField(0), scte.UpdateData()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated packet data does not match expected data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestNetworkEndDecodeEncode(t *testing.T) {
	target := testScteCreate3
	scte, err := NewSCTE35(target)
	if err != nil {
		t.Error(err.Error())
	}
	scte.UpdateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Original packet data does not match Generated data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestUnscheduledEventStartCreateEncode(t *testing.T) {
	target := testScteCreate2
	scte := CreateSCTE35()
	scte.SetTier(0x0)
	cmd := CreateTimeSignalCommand()
	cmd.SetHasPTS(true)
	cmd.SetPTS(0x2BFD4)
	scte.SetCommandInfo(cmd)
	scte.SetAdjustPTS(0x59CF4)

	descriptor := CreateSegmentationDescriptor()
	descriptor.SetEventID(0x2)
	descriptor.SetIsEventCanceled(false)
	descriptor.SetHasProgramSegmentation(true)
	descriptor.SetHasDuration(true)
	descriptor.SetDuration(0xAFF50)
	descriptor.SetIsDeliveryNotRestricted(true)
	descriptor.SetUPIDType(SegUPIDADI)
	descriptor.SetUPID([]byte("TEST"))
	descriptor.SetTypeID(SegDescUnscheduledEventStart)
	descriptor.SetSegmentNumber(0)
	descriptor.SetSegmentsExpected(0)

	scte.SetDescriptors([]SegmentationDescriptor{descriptor})

	generated := append(psi.NewPointerField(0), scte.UpdateData()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated packet data does not match expected data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestUnscheduledEventStartDecodeEncode(t *testing.T) {
	target := testScteCreate2
	scte, err := NewSCTE35(target)
	if err != nil {
		t.Error(err.Error())
	}
	scte.UpdateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Original packet data does not match Generated data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestSCPCreateEncode(t *testing.T) {
	target := csp
	scte := CreateSCTE35()
	scte.SetTier(0xFFF)
	scte.SetCommandInfo(CreateSpliceNull())

	descriptor := CreateSegmentationDescriptor()
	descriptor.SetEventID(0xC0000000)
	descriptor.SetIsEventCanceled(false)
	descriptor.SetHasProgramSegmentation(true)
	descriptor.SetHasDuration(false)
	descriptor.SetIsDeliveryNotRestricted(true)
	descriptor.SetUPIDType(SegUPIDURN)
	descriptor.SetUPID([]byte("urn:merlin:linear:stream:8987205474424984163"))
	descriptor.SetTypeID(SegDescContentIdentification)
	descriptor.SetSegmentNumber(0)
	descriptor.SetSegmentsExpected(0)

	scte.SetDescriptors([]SegmentationDescriptor{descriptor})

	generated := append(psi.NewPointerField(0), scte.UpdateData()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated packet data does not match expected data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestSCPDecodeEncode(t *testing.T) {
	target := csp
	scte, err := NewSCTE35(target)
	if err != nil {
		t.Error(err.Error())
	}
	scte.UpdateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Original packet data does not match Generated data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestScteCreateEncode(t *testing.T) {
	target := testScteCreate
	scte := CreateSCTE35()
	scte.SetTier(0xFFF)
	cmd := CreateTimeSignalCommand()
	scte.SetCommandInfo(cmd)
	cmd.SetHasPTS(true)
	cmd.SetPTS(0x86DF7550)
	scte.SetAdjustPTS(0x86DF7550)

	descriptor := CreateSegmentationDescriptor()
	descriptor.SetEventID(0x41424344)
	descriptor.SetIsEventCanceled(false)
	descriptor.SetHasProgramSegmentation(true)
	descriptor.SetHasDuration(false)
	descriptor.SetIsDeliveryNotRestricted(false)
	descriptor.SetIsWebDeliveryAllowed(false)
	descriptor.SetHasNoRegionalBlackout(true)
	descriptor.SetIsArchiveAllowed(true)
	descriptor.SetDeviceRestrictions(RestrictNone)
	descriptor.SetUPIDType(SegUPIDNotUsed)

	descriptor.SetTypeID(0x10)
	descriptor.SetSegmentNumber(1)
	descriptor.SetSegmentsExpected(1)

	scte.SetDescriptors([]SegmentationDescriptor{descriptor})

	generated := append(psi.NewPointerField(0), scte.UpdateData()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated packet data does not match expected data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestScteDecodeEncode(t *testing.T) {
	target := testScteCreate
	scte, err := NewSCTE35(target)
	if err != nil {
		t.Error(err.Error())
	}
	scte.UpdateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Original packet data does not match Generated data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestVssCreateEncode(t *testing.T) {
	target := vss
	scte := CreateSCTE35()
	scte.SetTier(0xFFF)
	cmd := CreateTimeSignalCommand()
	scte.SetCommandInfo(cmd)
	cmd.SetHasPTS(true)
	cmd.SetPTS(0x00000000)
	scte.SetAdjustPTS(0x6D71C7EF)

	descriptors := make([]SegmentationDescriptor, 0, 2)

	descriptor := CreateSegmentationDescriptor()
	descriptor.SetEventID(9)
	descriptor.SetIsEventCanceled(false)
	descriptor.SetHasProgramSegmentation(true)
	descriptor.SetHasDuration(false)
	descriptor.SetIsDeliveryNotRestricted(false)
	descriptor.SetIsWebDeliveryAllowed(true)
	descriptor.SetHasNoRegionalBlackout(false)
	descriptor.SetIsArchiveAllowed(true)
	descriptor.SetDeviceRestrictions(RestrictNone)
	descriptor.SetUPIDType(SegUPIDMID)

	mid := make([]UPID, 0, 2)
	upid := CreateUPID()
	upid.SetUPIDType(SegUPIDADI)
	upid.SetUPID([]byte("BLACKOUT:Sq+kY9muQderGNiNtOoN6w=="))
	mid = append(mid, upid)

	upid = CreateUPID()
	upid.SetUPIDType(SegUPADSINFO)
	upid.SetUPID([]byte("comcast:linear:licenserotation"))
	mid = append(mid, upid)
	descriptor.SetMID(mid)

	descriptor.SetTypeID(0x40)
	descriptor.SetSegmentNumber(0)
	descriptor.SetSegmentsExpected(0)
	descriptors = append(descriptors, descriptor)

	descriptor = CreateSegmentationDescriptor()
	descriptor.SetEventID(9)
	descriptor.SetIsEventCanceled(false)
	descriptor.SetHasProgramSegmentation(true)
	descriptor.SetHasDuration(false)
	descriptor.SetIsDeliveryNotRestricted(false)
	descriptor.SetIsWebDeliveryAllowed(true)
	descriptor.SetHasNoRegionalBlackout(false)
	descriptor.SetIsArchiveAllowed(true)
	descriptor.SetDeviceRestrictions(RestrictNone)
	descriptor.SetUPIDType(SegUPIDNotUsed)
	descriptor.SetTypeID(0x41)
	descriptor.SetSegmentNumber(0)
	descriptor.SetSegmentsExpected(0)
	descriptors = append(descriptors, descriptor)

	scte.SetDescriptors(descriptors)

	generated := append(psi.NewPointerField(0), scte.UpdateData()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Generated packet data does not match expected data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestVSSDecodeEncode(t *testing.T) {
	target := vss
	scte, _ := NewSCTE35(target)
	scte.UpdateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("Original packet data does not match Generated data\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}
