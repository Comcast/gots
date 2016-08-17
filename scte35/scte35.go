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
	"github.com/Comcast/gots/psi"
)

// Descriptor tag types and identifiers - only segmentation descriptors are used for now
const (
	segDescTag = 0x02
	segDescID  = 0x43554549
	minDescLen = 15 // min desc len does not include descriptor tag or len
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

type scte35 struct {
	command     SpliceCommandType
	hasPTS      bool
	pts         gots.PTS
	descriptors []SegmentationDescriptor

	data []byte
}

// NewSCTE35 creates a new SCTE35 signal from the provided byte slice. The byte slice is parsed and relevant info is made available fir the SCTE35 interface. If the message cannot me parsed, an error is returned.
func NewSCTE35(data []byte) (SCTE35, error) {
	s := &scte35{}
	err := s.parseTable(data)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *scte35) parseTable(data []byte) error {
	buf := bytes.NewBuffer(data)
	// closure to ignore EOF error from buf.ReadByte().  We've already checked
	// the length, we don't need to continually check it after every ReadByte
	// call
	readByte := func() byte {
		b, _ := buf.ReadByte()
		return b
	}
	if buf.Len() < int(uint16(psi.PointerField(data))+psi.PSIHeaderLen+15) {
		return gots.ErrInvalidSCTE35Length
	}
	if psi.TableID(data) == 0xfc {
		// read over the table header - +1 to skip protocol version
		headerLen := psi.PSIHeaderLen + uint16(psi.PointerField(data)) + 1
		buf.Next(int(headerLen))
		if readByte()&0x80 != 0 {
			return gots.ErrSCTE35EncryptionUnsupported
		}
		// unread this byte, because it contains the top bit of our pts offset
		err := buf.UnreadByte()
		if err != nil {
			return err
		}
		ptsAdjustment := uint40(buf.Next(5)) & 0x01ffffffff
		// skip cw_index, tier and spliceCommandLength
		// since it can be 0xfff(unknown) so it's pretty much useless
		buf.Next(4)
		s.command = SpliceCommandType(readByte())
		switch s.command {
		case TimeSignal:
			flags := readByte()
			s.hasPTS = (flags & 0x80) == 0x80
			if s.hasPTS {
				// unread prev byte because it contains the top bit of the pts offset
				err := buf.UnreadByte()
				if err != nil {
					return err
				}
				if buf.Len() < 11 {
					return gots.ErrInvalidSCTE35Length
				}
				s.pts = uint40(buf.Next(5)) & 0x01ffffffff
				// add the pts adjustment to get the real
				// value, we won't need it anymore after that
				s.pts += ptsAdjustment
			} else {
				return gots.ErrSCTE35UnsupportedSpliceCommand
			}
		case SpliceNull:
		default:
			return gots.ErrSCTE35UnsupportedSpliceCommand
		}
		descLoopLen := binary.BigEndian.Uint16(buf.Next(2))
		if buf.Len() < int(descLoopLen+psi.CrcLen) {
			return gots.ErrInvalidSCTE35Length
		}
		for bytesRead := uint16(0); bytesRead < descLoopLen; {
			d := &segmentationDescriptor{spliceInfo: s}
			descTag := readByte()
			descLen := readByte()
			// Make sure a bad descriptorLen doesn't kill us
			if descLoopLen-bytesRead-2 < uint16(descLen) || descLen < minDescLen {
				return gots.ErrInvalidSCTE35Length
			}
			if descTag != segDescTag {
				// Not interested in descriptors that are not
				// SegmentationDescriptors
				buf.Next(int(descLen))
			} else {
				err := d.parseDescriptor(buf.Next(int(descLen)))
				if err != nil {
					return err
				}
				s.descriptors = append(s.descriptors, d)
			}
			bytesRead += 2 + uint16(descLen)
		}
	} else {
		return gots.ErrUnknownTableID
	}
	// Check CRC?
	s.data = data
	return nil
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

func (s *scte35) HasPTS() bool {
	return s.hasPTS
}

func (s *scte35) PTS() gots.PTS {
	return s.pts
}

func (s *scte35) Command() SpliceCommandType {
	return s.command
}

func (s *scte35) Descriptors() []SegmentationDescriptor {
	return s.descriptors
}

func (s *scte35) Data() []byte {
	return s.data
}

func (d *segmentationDescriptor) SCTE35() SCTE35 {
	return d.spliceInfo
}

func (d *segmentationDescriptor) TypeID() SegDescType {
	return d.typeID
}

func (d *segmentationDescriptor) IsOut() bool {
	switch d.TypeID() {
	case SegDescProgramStart, SegDescChapterStart,
		SegDescProviderAdvertisementStart, SegDescDistributorAdvertisementStart,
		SegDescPlacementOpportunityStart, SegDescUnscheduledEventStart, SegDescNetworkStart,
		SegDescDistributorPoStart, SegDescProgramOverlapStart, SegDescProgramBlackoutOverride:
		return true
	default:
		return false
	}
}

func (d *segmentationDescriptor) IsIn() bool {
	switch d.TypeID() {
	case SegDescProgramEnd, SegDescChapterEnd,
		SegDescProviderAdvertisementEnd, SegDescDistributorAdvertisementEnd,
		SegDescPlacementOpportunityEnd, SegDescUnscheduledEventEnd, SegDescNetworkEnd,
		SegDescDistributorPoEnd:
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

func (d *segmentationDescriptor) Compare(c SegmentationDescriptor) int {
	return int(segWeight(d.TypeID())) - int(segWeight(c.TypeID()))
}

func (d *segmentationDescriptor) CanClose(out SegmentationDescriptor) bool {
	// if we're not an in, we can't close anything.  This should be checked before the client calls this, but just in case...
	if !d.IsIn() {
		return false
	}

	// We have to handle two cases, the normal in closes an out and the trumping
	// case.  In the normal case, the in segNum has to match segsExpected, in the
	// trumping case it doesn't matter, so handle these two separately.
	// Check the normal case, the subtraction check will fail if in isn't an in or out isn't an out
	if d.TypeID()-out.TypeID() == 1 && d.segNum == d.segsExpected {
		return true
	}
	// check the trumping case next
	if d.Compare(out) > 0 {
		return true
	}
	return false
}

// Determines equality for two segmentation descriptors
// Equality in this sense means that two "in" events are duplicates
// A lot of debate went in to what actually constitutes a "duplicate".  We get
// all sorts of things from different providers/transcoders, so in the end, we
// just settled on PTS and segmentation type
func (d *segmentationDescriptor) Equal(c SegmentationDescriptor) bool {
	if d == nil || c == nil {
		return false
	}
	if d.TypeID() != c.TypeID() {
		return false
	}
	if !d.SCTE35().HasPTS() || c.SCTE35().HasPTS() {
		return false
	}
	if d.SCTE35().PTS() != c.SCTE35().PTS() {
		return false
	}
	return true
}

// Recalculates the segmentation descriptor type such that higher precedence segmentation descriptor types
// have higher values.  Advertising-based segmentation descriptor types are normalized so their values are
// the same for both provider and distributor ad start (0x00) and end (0x01).
// Modeled after the work done in DASH-R for its advertising support.
func segWeight(typeID SegDescType) SegDescType {
	var typeVal int8
	typeVal = int8(typeID & 0xfe)
	if typeVal == SegDescDistributorAdvertisementStart ||
		typeVal == SegDescDistributorPoStart {
		typeVal = typeVal - 0x2
	}
	if typeVal >= 0x40 {
		return SegDescType(typeVal)
	}
	return SegDescType(abs(typeVal - 0x30))
}

func abs(num int8) int8 {
	switch {
	case num < 0:
		return -num
	default:
		return num
	}
}

func uint40(buf []byte) gots.PTS {
	return (gots.PTS(buf[0]&0x1) << 32) | (gots.PTS(buf[1]) << 24) | (gots.PTS(buf[2]) << 16) | (gots.PTS(buf[3]) << 8) | (gots.PTS(buf[4]))
}
