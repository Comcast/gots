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
	"encoding/base64"
	"io"
	"strings"
	"testing"

	"github.com/Comcast/gots"
)

var testScte = []byte{
	0x00, 0xfc, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf0, 0x05, 0x06, 0xfe,
	0x86, 0xdf, 0x75, 0x50, 0x00, 0x11, 0x02, 0x0f, 0x43, 0x55, 0x45, 0x49, 0x41, 0x42, 0x43, 0x44,
	0x7f, 0x8f, 0x00, 0x00, 0x10, 0x01, 0x01, 0x3a, 0x6d, 0xda, 0xee,
}

// This has a program segmentation flag that is false, caused bugs
// elsewhere
var testScte2 = []byte{
	0x00, 0xfc, 0x00, 0x53, 0x00, 0x00, 0x00, 0x02, 0xdd, 0x20, 0x00, 0xff, 0xf0, 0x05, 0x06, 0xfe,
	0x00, 0x08, 0x95, 0x44, 0x00, 0x3d, 0x02, 0x3b, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x02,
	0x7f, 0x1f, 0x02, 0x01, 0xfe, 0x00, 0x2d, 0xd2, 0x00, 0x02, 0xfe, 0x00, 0x00, 0x01, 0xe8, 0x09,
	0x1f, 0x53, 0x49, 0x47, 0x4e, 0x41, 0x4c, 0x3a, 0x59, 0x38, 0x6f, 0x30, 0x44, 0x33, 0x7a, 0x70,
	0x54, 0x78, 0x53, 0x30, 0x4c, 0x54, 0x31, 0x65, 0x77, 0x2b, 0x77, 0x75, 0x69, 0x77, 0x3d, 0x3d,
	0x36, 0x00, 0x00, 0xe0, 0xfa, 0x93, 0xc1,
}

var testScte3 = []byte{
	0x00, 0xfc, 0x30, 0x55, 0x00, 0x00, 0x00, 0x02, 0xd5, 0xa0, 0x00, 0xff, 0xf0, 0x05, 0x06, 0xfe,
	0x00, 0x04, 0x2b, 0x79, 0x00, 0x3f, 0x02, 0x1b, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x01,
	0x7f, 0x87, 0x09, 0x0c, 0x53, 0x49, 0x47, 0x4e, 0x41, 0x4c, 0x3a, 0x33, 0x2e, 0x30, 0x35, 0x30,
	0x35, 0x01, 0x01, 0x02, 0x20, 0x43, 0x55, 0x45, 0x49, 0x00, 0x00, 0x00, 0x02, 0x7f, 0xff, 0x00,
	0x00, 0x23, 0x13, 0xac, 0x09, 0x0c, 0x53, 0x49, 0x47, 0x4e, 0x41, 0x4c, 0x3a, 0x33, 0x2e, 0x30,
	0x35, 0x30, 0x34, 0x01, 0x01, 0x22, 0x04, 0xf5, 0x04, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

func TestSpliceInsertSignal(t *testing.T) {
	base64Bytes, _ := base64.StdEncoding.DecodeString("APwwLwAAz6l5ggD///8FYgAgAn/v/1jt40T+AHuYoAM1AAAACgAIQ1VFSQA4MjFRxjDp")

	s, err := NewSCTE35(base64Bytes)

	if err != gots.ErrSCTE35UnsupportedSpliceCommand {
		t.Error("SpliceInsert does not return splice command not supported err")
	}
	if s != nil {
		t.Error("NewSCTE35() of splice insert signal returns non-nil signal")
	}
}

func TestBasicSignal(t *testing.T) {
	s, err := NewSCTE35(testScte)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if s.Command() != TimeSignal {
		t.Errorf("Invalid command found, expecting TimeSignal(6), got: %v", s.Command())
	}
	if !s.HasPTS() {
		t.Error("Expecting PTS, but none found")
	} else if s.PTS() != 2262791504 {
		t.Errorf("Expecting PTS value of 2262791504, but found %v instead", s.PTS())
	}
	descs := s.Descriptors()
	if descs == nil {
		t.Error("Expecting descriptors in signal, but none found")
	} else if len(descs) != 1 {
		t.Errorf("Only expected one segmentation descriptor, but found %v instead", len(descs))
	} else {
		d := descs[0]
		if d.TypeID() != SegDescProgramStart {
			t.Errorf("Expecting seg descriptor type ProgramStart(0x10), got %x instead", d.TypeID())
		} else if !d.IsOut() {
			t.Error("SegDescProgramStart is out, but IsOut() returned false")
		} else if d.IsIn() {
			t.Error("SegDescProgramStart is out, but IsIn() return true")
		}
		upid := d.UPID()
		if upid == nil {
			t.Error("upid not found in descriptor")
		} else if len(upid) != 0 {
			t.Error("non-zero len upid found, indicating error")
		}
	}
}

// splice_null commands are used for CSP.
func TestParseSpliceNull(t *testing.T) {
	base64Bytes, _ := base64.StdEncoding.DecodeString("APwwNQAAAAAAAAD/8AEAACQCIkNVRUnAAAAAf78BEzU5MzkwMjY1NjUxNzc3OTIxNjMBAQHrr2Ob")

	s, err := NewSCTE35(base64Bytes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if s.Command() != SpliceNull {
		t.Errorf("Invalid command found, expecting SpliceNull(0), got %v", s.Command())
	}
	descs := s.Descriptors()
	if descs == nil {
		t.Error("Expecting descriptors in signal, but none found")
	} else if len(descs) != 1 {
		t.Errorf("Only expected one segmentation descriptor, but found %v instead", len(descs))
	} else {
		d := descs[0]
		if d.TypeID() != SegDescContentIdentification {
			t.Errorf("Expecting seg descriptor type SegDescContentIdentification(0x1), got %x instead", d.TypeID())
		}
		upid := d.UPID()
		if upid == nil {
			t.Error("csp upid not found in descriptor")
		} else {
			if d.UPIDType() != SegUPIDUserDefined {
				t.Errorf("Expected upid type SegUPIDUserDefined(1), got %v instead", d.UPIDType())
			}
			buf := bytes.NewBuffer(upid)
			upidStr, err := buf.ReadString(0)
			if err != io.EOF {
				t.Error(err)
			} else if strings.Compare(upidStr, "5939026565177792163") != 0 {
				t.Errorf("Invalid UPID found, expected 5939026565177792163, got %s", upidStr)
			}
		}
	}
}

func TestSCTEExpanded(t *testing.T) {
	s1, err1 := NewSCTE35(testScte)
	if err1 != nil {
		t.Error(err1)
		t.FailNow()
	}
	s2, err2 := NewSCTE35(testScte2)
	if err2 != nil {
		t.Error(err2)
		t.FailNow()
	}
	if s2.Command() != TimeSignal {
		t.Errorf("Invalid command found, expecting TimeSignal(6), got: %v", s2.Command())
	}
	if !s2.HasPTS() {
		t.Error("Expecting PTS, but none found")
	} else if s2.PTS() != 750180 {
		t.Errorf("Expecting PTS value of 750180 %v instead", s2.PTS())
	}
	descs2 := s2.Descriptors()
	var d2, d1 SegmentationDescriptor
	if descs2 == nil {
		t.Error("Expecting descriptors in signal, but none found")
	} else if len(descs2) != 1 {
		t.Errorf("Only expected one segmentation descriptor, but found %v instead", len(descs2))
	} else {
		d2 := descs2[0]
		if d2.TypeID() != SegDescDistributorPoStart {
			t.Errorf("Expecting seg descriptor type SegDescDistributorPoStart(0x36), got %x instead", d2.TypeID())
		} else if !d2.IsOut() {
			t.Error("SegDescDistributorPoStart is out, but IsOut() returned false")
		} else if d2.IsIn() {
			t.Error("SegDescDistributorPoStart is out, but IsIn() return true")
		}
		upid := d2.UPID()
		if upid == nil {
			t.Error("upid not found in descriptor")
		} else {
			if d2.UPIDType() != SegUPIDADI {
				t.Errorf("Expected upid type SegUPIDADI(9), got %v instead", d2.UPIDType())
			}
			buf := bytes.NewBuffer(upid)
			upidStr, err := buf.ReadString(0)
			if err != io.EOF {
				t.Error(err)
			} else if strings.Compare(upidStr, "SIGNAL:Y8o0D3zpTxS0LT1ew+wuiw==") != 0 {
				t.Errorf("Invalid UPID found, expected SIGNAL:Y8o0D3zpTxS0LT1ew+wuiw==, got %s", upidStr)
			}
		}
	}
	descs1 := s1.Descriptors()
	if descs1 == nil {
		t.Error("expecting descriptors in signal, but none found")
	} else if len(descs1) != 1 {
		t.Errorf("Only expected one segmentation descriptor, but found %v instead", len(descs1))
	} else {
		d1 = descs1[0]
	}
	if d1 != nil && d2 != nil {
		if d1.CanClose(d2) {
			t.Errorf("Segmentation type %v shouldn't be able to close %v, but CanClose returned true", d1.UPIDType(), d2.UPIDType())
		}
		if d2.CanClose(d1) {
			t.Errorf("Segmentation type %v shouldn't be able to close %v, but CanClose returned true", d2.UPIDType(), d1.UPIDType())
		}
	}
}

func TestSCTEMultipleDescriptors(t *testing.T) {
	s, err := NewSCTE35(testScte3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if s.Command() != TimeSignal {
		t.Errorf("Invalid command found, expecting TimeSignal(6), got: %v", s.Command())
	}
	if !s.HasPTS() {
		t.Error("Expecting PTS, but none found")
	}
	descs := s.Descriptors()
	if descs == nil {
		t.Error("expecting descriptors in signal, but none found")
	} else if len(descs) != 2 {
		t.Error("expecting two descriptors in signal, but found", len(descs))
		t.FailNow()
	}
	if descs[0].TypeID() != SegDescPlacementOpportunityEnd {
		t.Error("Invalid seg type found, expected SegDescPlacementOpportunityEnd(0x35), found ", descs[0].TypeID())
	}
	if descs[0].IsOut() {
		t.Error("descriptor type is an in, but IsOut() returned true")
	}
	if !descs[0].IsIn() {
		t.Error("descriptor type is an in, but IsIn() returned false")
	}
	if descs[1].TypeID() != SegDescPlacementOpportunityStart {
		t.Error("Invalid seg type found, expected SegDescPlacementOpportunityStart(0x34), but found", descs[1].TypeID())
	}
	if !descs[0].CanClose(descs[1]) {
		t.Error("descs[0] should be able to close [1], but ConClose reports false")
	}
	if descs[1].SCTE35() == nil {
		t.Error("descs[1] does not return scte obj")
	}
	if descs[0].SCTE35() != descs[1].SCTE35() {
		t.Error("SCTE obj of both descs is not the same")
	}
}
