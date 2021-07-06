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

// scte35 is a structure representing a SCTE35 message.
type scte35 struct {
	tableHeader         psi.TableHeader
	protocolVersion     uint8
	encryptedPacket     bool     // not supported
	encryptionAlgorithm uint8    // 6 bits
	pts                 gots.PTS // pts is stored adjusted in struct
	cwIndex             uint8
	tier                uint16 // 12 bits
	spliceCommandLength uint16 // 12 bits
	commandType         SpliceCommandType
	commandInfo         SpliceCommand
	descriptors         []SegmentationDescriptor
	crc32               uint32
	alignmentStuffing   uint

	data []byte

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

// parseTable will parse bytes into a scte35 message struct
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
	var err error
	s.tableHeader, err = psi.TableHeaderFromBytes(buf.Next(3))
	if err != nil {
		return err
	}
	if s.tableHeader.TableID == 0xfc {
		s.protocolVersion = readByte()
		if readByte()&0x80 != 0 {
			return gots.ErrSCTE35EncryptionUnsupported
		}

		// unread this byte, because it contains the encryptionAlgorithm field
		if err := buf.UnreadByte(); err != nil {
			return err
		}
		s.encryptionAlgorithm = (readByte() >> 1) & 0x3F // 0111 1110

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
				s.otherDescriptorBytes = append(s.otherDescriptorBytes, descTag, descLen)
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
	return nil
}

// HasPTS returns true if there is a pts time.
func (s *scte35) HasPTS() bool {
	return s.commandInfo.HasPTS()
}

// PTS returns the PTS time of the signal if it exists. Includes adjustment.
func (s *scte35) PTS() gots.PTS {
	return s.pts
}

// Command returns the signal's splice command.
func (s *scte35) Command() SpliceCommandType {
	return s.commandType
}

// CommandInfo returns an object describing fields of the signal's splice
// command structure
func (s *scte35) CommandInfo() SpliceCommand {
	return s.commandInfo
}

// Descriptors returns a slice of the signals SegmentationDescriptors sorted
// by descriptor weight (least important signals first)
func (s *scte35) Descriptors() []SegmentationDescriptor {
	return s.descriptors
}

// Tier returns which authorization tier this message was assigned to.
// The tier value of 0XFFF is the default and will ignored. When using tier values,
// the SCTE35 message must fit entirely into the ts payload without being split up.
func (s *scte35) Tier() uint16 {
	return s.tier
}

// AlignmentStuffing returns how many stuffing bytes will be added to the SCTE35
// message at the end.
func (s *scte35) AlignmentStuffing() uint {
	return s.alignmentStuffing
}

// Data returns the raw data bytes of the scte signal
func (s *scte35) Data() []byte {
	return s.data
}

func uint40(buf []byte) gots.PTS {
	return (gots.PTS(buf[0]&0x1) << 32) | (gots.PTS(buf[1]) << 24) | (gots.PTS(buf[2]) << 16) | (gots.PTS(buf[3]) << 8) | (gots.PTS(buf[4]))
}

// String returns a string representation of the SCTE35 message.
// String function is for debugging and testing.
func (s *scte35) String() string {
	numspaces := 0
	indentString := ""

	indent := func(n int) {
		numspaces += n
		indentString = ""
		for i := 0; i < numspaces; i++ {
			indentString += "   "
		}
	}

	indentPrintf := func(format string, a ...interface{}) string {
		return fmt.Sprintf(indentString+format, a...)
	}

	str := ""
	s.UpdateData()
	str += indentPrintf("table_id: 0x%X\n", s.tableHeader.TableID)
	str += indentPrintf("section_syntax_indicator: %t\n", s.tableHeader.SectionSyntaxIndicator)
	str += indentPrintf("private_indicator: %t\n", s.tableHeader.PrivateIndicator)
	str += indentPrintf("section_length: %d\n", s.tableHeader.SectionLength)

	str += indentPrintf("protocol_version: 0x%X\n", s.protocolVersion)
	str += indentPrintf("encrypted_packet: %t\n", s.encryptedPacket)
	str += indentPrintf("encryption_algorithm: 0x%X\n", s.encryptionAlgorithm)

	str += indentPrintf("has_pts: %t\n", s.HasPTS())
	str += indentPrintf("adjusted_pts: %d\n", s.PTS())
	str += indentPrintf("cw_index: 0x%X\n", s.cwIndex)
	str += indentPrintf("tier: 0x%X\n", s.tier)
	str += indentPrintf("splice_command_type: %s\n", SpliceCommandTypeNames[s.commandType])
	indent(1)
	if cmd, ok := s.commandInfo.(SpliceInsertCommand); ok {
		str += indentPrintf("splice_event_id: 0x%X\n", cmd.EventID())
		str += indentPrintf("splice_event_cancel_indicator: %t\n", cmd.IsEventCanceled())
		if !cmd.IsEventCanceled() {
			str += indentPrintf("out_of_network_indicator: %t\n", cmd.IsOut())
			str += indentPrintf("program_splice_flag: %t\n", cmd.IsProgramSplice())
			str += indentPrintf("duration_flag: %t\n", cmd.HasDuration())
			str += indentPrintf("splice_immediate_flag: %t\n", cmd.SpliceImmediate())
			str += indentPrintf("splice_time_has_pts: %t\n", cmd.HasPTS())
			if cmd.HasPTS() {
				str += indentPrintf("splice_time_pts: %d\n", cmd.PTS())
			}
			str += indentPrintf("component_count: %d\n", len(cmd.Components()))
			for _, comp := range cmd.Components() {
				str += indentPrintf("component:\n")
				indent(1)
				str += indentPrintf("component_tag: 0x%X\n", comp.ComponentTag())
				str += indentPrintf("component_has_pts: %t\n", comp.HasPTS())
				if comp.HasPTS() {
					str += indentPrintf("component_pts: %d\n", cmd.PTS())
				}
				indent(-1)

				if cmd.HasDuration() {
					str += indentPrintf("auto_return: %t\n", cmd.IsAutoReturn())
					str += indentPrintf("duration: %d\n", cmd.Duration())
				}
				str += indentPrintf("unique_program_id: %t\n", cmd.UniqueProgramId())
				str += indentPrintf("avail_num: %d\n", cmd.AvailNum())
				str += indentPrintf("avails_expected: %d\n", cmd.AvailsExpected())
			}
		}
	}

	if cmd, ok := s.commandInfo.(TimeSignalCommand); ok {
		str += indentPrintf("time_specified_flag: %t\n", cmd.HasPTS())
		if cmd.HasPTS() {
			str += indentPrintf("pts_time: %d\n", cmd.PTS())
		}
	}
	indent(-1)

	str += indentPrintf("descriptor_count: %d\n", len(s.Descriptors()))
	for _, desc := range s.descriptors {
		str += indentPrintf("descriptor:\n")

		indent(1)
		if desc.IsIn() {
			indentPrintf("<--- IN Segmentation Descriptor")
		}
		if desc.IsOut() {
			indentPrintf("---> OUT Segmentation Descriptor")
		}
		str += indentPrintf("segmentation_event_id: 0x%X\n", desc.EventID())
		str += indentPrintf("segmentation_event_cancel_indicator: %t\n", desc.IsEventCanceled())
		if !desc.IsEventCanceled() {
			str += indentPrintf("program_segmentation_flag: %t\n", desc.HasProgramSegmentation())
			str += indentPrintf("segmentation_duration_flag: %t\n", desc.HasDuration())
			str += indentPrintf("delivery_not_restricted_flag: %t\n", desc.IsDeliveryNotRestricted())
			if !desc.IsDeliveryNotRestricted() {
				str += indentPrintf("web_delivery_allowed_flag: %t\n", desc.IsWebDeliveryAllowed())
				str += indentPrintf("no_regional_blackout_flag: %t\n", desc.HasNoRegionalBlackout())
				str += indentPrintf("archive_allowed_flag: %t\n", desc.IsArchiveAllowed())
				str += indentPrintf("device_restrictions: %s\n", DeviceRestrictionsNames[desc.DeviceRestrictions()])
			}
			if !desc.HasProgramSegmentation() {
				str += indentPrintf("component_count: %d\n", len(desc.Components()))
				for _, comp := range desc.Components() {
					str += indentPrintf("component:\n")
					indent(1)
					str += indentPrintf("component_tag: 0x%X\n", comp.ComponentTag())
					str += indentPrintf("pts_offset: %d\n", comp.PTSOffset())
					indent(-1)
				}
			}
			if desc.HasDuration() {
				str += indentPrintf("segmentation_duration: %d\n", desc.Duration())
			}
			str += indentPrintf("segmentation_upid_type: %s\n", SegUPIDTypeNames[desc.UPIDType()])
			if desc.UPIDType() != SegUPIDMID {
				str += indentPrintf("segmentation_upid: %s\n", string(desc.UPID()))
			} else {
				str += indentPrintf("segmentation_mid: \n")
				indent(1)
				for _, upid := range desc.MID() {
					str += indentPrintf("upid:\n")
					indent(1)
					str += indentPrintf("segmentation_mid_upid_type: %s\n", SegUPIDTypeNames[upid.UPIDType()])
					str += indentPrintf("segmentation_mid_upid: %s\n", string(upid.UPID()))
					indent(-1)
				}
				indent(-1)
			}
			str += indentPrintf("segmentation_type_id: %s\n", SegDescTypeNames[desc.TypeID()])
			str += indentPrintf("segment_num: 0x%X\n", desc.SegmentNumber())
			str += indentPrintf("segments_expected: 0x%X\n", desc.SegmentsExpected())
			if desc.HasSubSegments() {
				str += indentPrintf("sub_segment_num: 0x%X\n", desc.SubSegmentNumber())
				str += indentPrintf("sub_segments_expected: 0x%X\n", desc.SubSegmentsExpected())
			}
		}
		indent(-1)
	}

	str += indentPrintf("alignment_stuffing_byte_count: %d\n", s.alignmentStuffing)
	str += indentPrintf("CRC_32: 0x%X", s.data[len(s.data)-4:])

	return str
}
