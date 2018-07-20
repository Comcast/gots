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
	"fmt"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/psi"
)

// Descriptor tag types and identifiers - only segmentation descriptors are used for now
const (
	segDescTag = 0x02
	segDescID  = 0x43554549
)

type scte35 struct {
	tableHeader         psi.TableHeader
	protocolVersion     uint8
	encryptedPacket     bool  // not supported
	encryptionAlgorithm uint8 // 6 bits
	hasPTS              bool
	pts                 gots.PTS // pts is stored adjusted in struct
	cwIndex             uint8
	tier                uint16 // 12 bits
	spliceCommandLength uint16 // 12 bits
	commandType         SpliceCommandType
	commandInfo         SpliceCommand
	descriptors         []SegmentationDescriptor
	crc32               uint32
	alignmentStuffing   int

	updateBytes bool // if set, the data will be updated on the next function call to get data
	data        []byte

	// because there is no support for descriptors other than segmentation descriptors,
	// the bytes need to be stored so information is not lost.
	otherDescriptorBytes []byte
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
	// read over the pointer field
	buf.Next(int(psi.PointerField(data) + 1))
	// read in the TableHeader
	s.tableHeader = psi.TableHeaderFromBytes(buf.Next(3))
	if s.tableHeader.TableID == 0xfc {
		s.protocolVersion = readByte()
		if readByte()&0x80 != 0 {
			return gots.ErrSCTE35EncryptionUnsupported
		}

		// unread this byte, because it contains the encryptionAlgorithm field
		if err := buf.UnreadByte(); err != nil {
			return err
		}
		s.encryptionAlgorithm = (readByte() << 1) & 0x7E // 0111 1110

		// unread this byte, because it contains the top bit of our pts offset
		err := buf.UnreadByte()
		if err != nil {
			return err
		}
		ptsAdjustment := uint40(buf.Next(5)) & 0x01ffffffff
		// read cw_index, tier and spliceCommandLength
		// spliceCommandLength can be 0xfff(unknown) so it's pretty much useless
		s.cwIndex = readByte()
		bytes := buf.Next(3)
		s.tier = uint16(bytes[0])<<4 | uint16(bytes[1]&0xF0)>>4
		s.spliceCommandLength = uint16(bytes[1]&0x0F)<<8 | uint16(bytes[2])
		s.commandType = SpliceCommandType(readByte())
		switch s.commandType {
		case TimeSignal, SpliceInsert:
			var cmd SpliceCommand
			if s.commandType == TimeSignal {
				cmd, err = parseTimeSignal(buf)
			} else {
				cmd, err = parseSpliceInsert(buf)
			}
			if err != nil {
				return err
			}
			// add the pts adjustment to get the real value
			s.pts = cmd.PTS().Add(ptsAdjustment)
			s.hasPTS = cmd.HasPTS()
			s.commandInfo = cmd
		case SpliceNull:
			s.commandInfo = &spliceNull{}
		default:
			return gots.ErrSCTE35UnsupportedSpliceCommand
		}
		// descriptor_loop_length 2 + CRC 4
		if buf.Len() < 2+int(psi.CrcLen) {
			return gots.ErrInvalidSCTE35Length
		}
		// parse descriptors
		descriptorLoopLength := binary.BigEndian.Uint16(buf.Next(2))
		if buf.Len() < int(descriptorLoopLength+psi.CrcLen) {
			return gots.ErrInvalidSCTE35Length
		}
		for bytesRead := uint16(0); bytesRead < descriptorLoopLength; {
			descTag := readByte()
			descLen := readByte()
			// Make sure a bad descriptorLen doesn't kill us
			if descriptorLoopLength-bytesRead-2 < uint16(descLen) {
				return gots.ErrInvalidSCTE35Length
			}
			if descTag != segDescTag {
				// Not interested in descriptors that are not
				// SegmentationDescriptors
				// Store their bytes anyways so the data is not lost.
				s.otherDescriptorBytes = append(s.otherDescriptorBytes, buf.Next(int(descLen))...)
			} else {
				d := &segmentationDescriptor{spliceInfo: s}
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
	// remove the pointer field and associated data off the top so we only get the
	// table data
	s.data = data[psi.PointerField(data)+1:]
	s.updateBytes = false // do not update on next call of Data()
	return nil
}

func (s *scte35) HasPTS() bool {
	return s.hasPTS
}

func (s *scte35) PTS() gots.PTS {
	return s.pts
}

func (s *scte35) Command() SpliceCommandType {
	return s.commandType
}

func (s *scte35) CommandInfo() SpliceCommand {
	return s.commandInfo
}

func (s *scte35) Descriptors() []SegmentationDescriptor {
	return s.descriptors
}

func (s *scte35) Tier() uint16 {
	return s.tier
}

func (s *scte35) AlignmentStuffing() int {
	return s.alignmentStuffing
}

func (s *scte35) Data() []byte {
	if s.updateBytes {
		s.generateData()
	}
	return s.data
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

func (s *scte35) String() string {
	indent := "   "
	str := ""
	s.generateData()
	str += fmt.Sprintf("table_id: %X\n", s.tableHeader.TableID)
	str += fmt.Sprintf("section_syntax_indicator: %t\n", s.tableHeader.SectionSyntaxIndicator)
	str += fmt.Sprintf("private_indicator: %t\n", s.tableHeader.PrivateIndicator)
	str += fmt.Sprintf("section_length: %X\n", s.tableHeader.SectionLength)

	str += fmt.Sprintf("protocol_version: %X\n", s.protocolVersion)
	str += fmt.Sprintf("encrypted_packet: %t\n", s.encryptedPacket)
	str += fmt.Sprintf("encryption_algorithm: %X\n", s.encryptionAlgorithm)

	str += fmt.Sprintf("pts_adjustment: %X\n", s.PTS())
	str += fmt.Sprintf("cw_index: %X\n", s.cwIndex)
	str += fmt.Sprintf("tier: %X\n", s.tier)
	str += fmt.Sprintf("splice_command_type: %X\n", s.commandType)

	if cmd, ok := s.commandInfo.(SpliceInsertCommand); ok {
		str += fmt.Sprintf(indent+"splice_event_id: %X\n", cmd.EventID())
		str += fmt.Sprintf(indent+"splice_event_cancel_indicator: %t\n", cmd.IsEventCanceled())
		if !cmd.IsEventCanceled() {
			str += fmt.Sprintf(indent+"out_of_network_indicator: %t\n", cmd.IsOut())
			str += fmt.Sprintf(indent+"program_splice_flag: %t\n", cmd.IsProgramSplice())
			str += fmt.Sprintf(indent+"duration_flag: %t\n", cmd.HasDuration())
			str += fmt.Sprintf(indent+"splice_immediate_flag: %t\n", cmd.SpliceImmediate())
			str += fmt.Sprintf(indent+"splice_time_has_pts: %t\n", cmd.HasPTS())
			if cmd.HasPTS() {
				str += fmt.Sprintf(indent+"splice_time_pts: %X\n", cmd.PTS())
			}
			str += fmt.Sprintf(indent+"component_count: %X\n", len(cmd.Components()))
			for _, comp := range cmd.Components() {
				str += fmt.Sprintf(indent + "component:\n")
				str += fmt.Sprintf(indent+indent+"component_tag: %X\n", comp.ComponentTag())
				str += fmt.Sprintf(indent+indent+"component_has_pts: %t\n", comp.HasPTS())
				if comp.HasPTS() {
					str += fmt.Sprintf(indent+indent+"component_tag: %X\n", cmd.PTS())
				}

				if cmd.HasDuration() {
					str += fmt.Sprintf(indent+"auto_return: %t\n", cmd.IsAutoReturn())
					str += fmt.Sprintf(indent+"duration: %X\n", cmd.Duration())
				}
				str += fmt.Sprintf(indent+"unique_program_id: %t\n", cmd.UniqueProgramId())
				str += fmt.Sprintf(indent+"avail_num: %X\n", cmd.AvailNum())
				str += fmt.Sprintf(indent+"avails_expected: %X\n", cmd.AvailsExpected())
			}
		}
	}

	if cmd, ok := s.commandInfo.(TimeSignalCommand); ok {
		str += fmt.Sprintf(indent+"time_specified_flag: %t\n", cmd.HasPTS())
		if cmd.HasPTS() {
			str += fmt.Sprintf(indent+"pts_time: %X\n", cmd.PTS())
		}
	}

	str += fmt.Sprintf("descriptor_count: %d\n", len(s.descriptors))

	for _, desc := range s.descriptors {
		str += fmt.Sprintf("descriptor:\n")
		str += fmt.Sprintf(indent+"segmentation_event_id: %X\n", desc.EventID())
		str += fmt.Sprintf(indent+"segmentation_event_cancel_indicator: %t\n", desc.IsEventCanceled())
		if !desc.IsEventCanceled() {
			str += fmt.Sprintf(indent+"program_segmentation_flag: %t\n", desc.HasProgramSegmentation())
			str += fmt.Sprintf(indent+"segmentation_duration_flag: %t\n", desc.HasDuration())
			str += fmt.Sprintf(indent+"delivery_not_restricted_flag: %t\n", desc.IsDeliveryNotRestricted())
			if !desc.IsDeliveryNotRestricted() {
				str += fmt.Sprintf(indent+"web_delivery_allowed_flag: %t\n", desc.IsWebDeliveryAllowed())
				str += fmt.Sprintf(indent+"no_regional_blackout_flag: %t\n", desc.HasNoRegionalBlackout())
				str += fmt.Sprintf(indent+"archive_allowed_flag: %t\n", desc.IsArchiveAllowed())
				str += fmt.Sprintf(indent+"device_restrictions: %X\n", desc.DeviceRestrictions())
			}
			if !desc.HasProgramSegmentation() {
				str += fmt.Sprintf(indent+"component_count: %X\n", len(desc.Components()))
				for _, comp := range desc.Components() {
					str += fmt.Sprintf(indent + "component:\n")
					str += fmt.Sprintf(indent+indent+"component_tag: %X\n", comp.ComponentTag())
					str += fmt.Sprintf(indent+indent+"pts_offset: %X\n", comp.PTSOffset())
				}
			}
			if desc.HasDuration() {
				str += fmt.Sprintf(indent+"segmentation_duration: %X\n", desc.Duration())
			}
			str += fmt.Sprintf(indent+"segmentation_upid_type: %X\n", desc.UPIDType())
			if desc.UPIDType() != SegUPIDMID {
				str += fmt.Sprintf(indent+"segmentation_upid: %X\n", desc.UPID())
			} else {
				str += fmt.Sprintf(indent + "segmentation_mid: \n")
				for _, upid := range desc.MID() {
					str += fmt.Sprintf(indent + "upid:\n")
					str += fmt.Sprintf(indent+indent+"segmentation_mid_upid_type: %X\n", upid.UPIDType())
					str += fmt.Sprintf(indent+indent+"segmentation_mid_upid: %X\n", upid.UPID())
				}
			}
			str += fmt.Sprintf(indent+"segmentation_type_id: %X\n", desc.TypeID())
			str += fmt.Sprintf(indent+"segment_num: %X\n", desc.SegmentNumber())
			str += fmt.Sprintf(indent+"segments_expected: %X\n", desc.SegmentsExpected())
			if desc.HasSubSegments() {
				str += fmt.Sprintf(indent+"sub_segment_num: %X\n", desc.SubSegmentNumber())
				str += fmt.Sprintf(indent+"sub_segments_expected: %X\n", desc.SubSegmentsExpected())
			}
		}
	}

	str += fmt.Sprintf("alignment_stuffing_byte_count: %d\n", s.alignmentStuffing)
	str += fmt.Sprintf("CRC_32: %X", s.data[len(s.data)-4:])
	//str += fmt.Sprintf("Data: %X", s.data)

	return str
}
