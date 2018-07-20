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
	"github.com/Comcast/gots/psi"
	"testing"
)

func TestCreate(t *testing.T) {
	target := testVss
	scte, _ := NewSCTE35(target)
	scte.(*scte35).generateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)
	//fmt.Println(scte)
	if !bytes.Equal(target, generated) {
		t.Errorf("\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestCreateVss(t *testing.T) {
	target := vss
	scte := CreateSCTE35()
	scte.SetTier(0xFFF)
	scte.SetCommand(TimeSignal)
	cmd := CreateTimeSignalCommand()
	cmd.SetHasPTS(true)
	cmd.SetPTS(0x00000000)
	scte.SetCommandInfo(cmd)
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
	upid.SetUPID([]byte{
		0x42, 0x4C, 0x41, 0x43, 0x4B, 0x4F, 0x55, 0x54, 0x3A, 0x53, 0x71,
		0x2B, 0x6B, 0x59, 0x39, 0x6D, 0x75, 0x51, 0x64, 0x65, 0x72, 0x47,
		0x4E, 0x69, 0x4E, 0x74, 0x4F, 0x6F, 0x4E, 0x36, 0x77, 0x3D, 0x3D,
	})
	mid = append(mid, upid)

	upid = CreateUPID()
	upid.SetUPIDType(SegUPADSINFO)
	upid.SetUPID([]byte{
		0x63, 0x6F, 0x6D, 0x63, 0x61, 0x73, 0x74, 0x3A, 0x6C, 0x69,
		0x6E, 0x65, 0x61, 0x72, 0x3A, 0x6C, 0x69, 0x63, 0x65, 0x6E,
		0x73, 0x65, 0x72, 0x6F, 0x74, 0x61, 0x74, 0x69, 0x6F, 0x6E,
	})
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

	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}

func TestVSSData(t *testing.T) {
	target := vss
	scte, _ := NewSCTE35(target)
	scte.(*scte35).generateData()
	generated := append(psi.NewPointerField(0), scte.Data()...)

	if !bytes.Equal(target, generated) {
		t.Errorf("\n   Target: %X\nGenerated: %X\n", target, generated)
	}
}
