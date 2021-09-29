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
	"github.com/Comcast/gots/"
	"github.com/Comcast/gots/psi"
)

// CreateSCTE35 creates a default SCTE35 message and returns it.
// The default message has the tier 0xFFF and a Splice Null command.
func CreateSCTE35() SCTE35 {
	scte35 :=
		&scte35{
			protocolVersion:     0,     // only version 0 exists
			encryptedPacket:     false, // no support for encryption
			encryptionAlgorithm: 0,     // no encryption support, no way to change
			pts:                 0,     // no pts
			cwIndex:             0,     // undefined, without encryption
			tier:                0xFFF, // ignore tier value

			commandType: SpliceNull,    // null command type by default
			commandInfo: &spliceNull{}, // info pooints to null command

			descriptors: []SegmentationDescriptor{}, // empty slice of descriptors

		}
	scte35.tableHeader =
		psi.TableHeader{
			TableID:                0xFC,  // always 0xFC for scte35
			SectionSyntaxIndicator: false, // always false for scte35
			PrivateIndicator:       false, // always false for scte35
			SectionLength:          0,     // to be calculated when converting to bytes
		}
	return scte35
}

// subtractPTS subtracts the two PTS times and returns a new PTS
// This is different from the durationFrom function since the
// order of operands matters. It will always subtract final
// minus initial. If the result is negative then it will add
// MaxPtsTicks to produce a positive number
// (in this case it assumes rollover happened).
// used for calculating pts adjustment.
func subtractPTS(final gots.PTS, initial gots.PTS) gots.PTS {
	if final >= initial {
		return final - initial
	} else {
		return gots.MaxPtsTicks - (initial - final)
	}
}

// UpdateData will encode the SCTE35 information to bytes and return it.
// UpdateData will make the next call to Data() return these new bytes.
func (s *scte35) UpdateData() []byte {
	// splice command generate bytes
	spliceCommandBytes := s.commandInfo.Data()
	// spliceCommandLength can be set as 0xFFF (undefined), but calculate it anyways
	spliceCommandLength := len(spliceCommandBytes)
	s.spliceCommandLength = uint16(spliceCommandLength)

	// generate bytes for splice descriptors
	descriptorBytes := make([]byte, 2)
	// append descriptors that are not extracted
	descriptorBytes = append(descriptorBytes, s.otherDescriptorBytes...)
	// append segmentation descriptors
	for i := range s.descriptors {
		descriptorBytes = append(descriptorBytes, s.descriptors[i].Data()...)
	}
	descriptorLoopLength := len(descriptorBytes) - 2
	descriptorBytes[0] = byte(descriptorLoopLength >> 8)
	descriptorBytes[1] = byte(descriptorLoopLength)

	const staticFieldsLength = 13
	const crcLength = int(psi.CrcLen)

	sectionLength := staticFieldsLength + spliceCommandLength + descriptorLoopLength + crcLength + int(s.alignmentStuffing)
	s.tableHeader.SectionLength = uint16(sectionLength)

	tableHeaderBytes := s.tableHeader.Data()
	tableHeaderLength := len(tableHeaderBytes)

	// slices that point to the starting position of their names
	data := make([]byte, tableHeaderLength+sectionLength)
	tableHeader := data
	section := tableHeader[tableHeaderLength:]
	spliceCommand := section[11:]
	spliceDescriptor := spliceCommand[spliceCommandLength:]
	crc := data[len(data)-crcLength:]

	ptsAdj := subtractPTS(s.pts, s.commandInfo.PTS())

	if s.encryptedPacket {
		section[1] = 0x80 // 1000 0000
	}
	section[0] = s.protocolVersion                      // 1111 1111
	section[1] |= (s.encryptionAlgorithm & 0x3F) << 1   // 0111 1110
	section[1] |= byte(ptsAdj>>32) & 0x01               // 0000 0001
	section[2] = byte(ptsAdj >> 24)                     // 1111 1111
	section[3] = byte(ptsAdj >> 16)                     // 1111 1111
	section[4] = byte(ptsAdj >> 8)                      // 1111 1111
	section[5] = byte(ptsAdj)                           // 1111 1111
	section[6] = s.cwIndex                              // 1111 1111
	section[7] = byte(s.tier >> 4)                      // 1111 1111
	section[8] = byte(s.tier << 4)                      // 1111 0000
	section[8] |= byte(s.spliceCommandLength>>8) & 0x0F // 0000 1111
	section[9] = byte(s.spliceCommandLength)            // 1111 1111
	section[10] = byte(s.commandType)                   // 1111 1111

	copy(tableHeader, tableHeaderBytes)
	copy(spliceCommand, spliceCommandBytes)
	copy(spliceDescriptor, descriptorBytes)

	crcBytes := gots.ComputeCRC(tableHeader[:len(data)-crcLength])
	copy(crc, crcBytes)
	s.data = data
	return data
}

// SetHasPTS sets if this SCTE35 message has a PTS.
func (s *scte35) SetHasPTS(flag bool) {
	s.commandInfo.SetHasPTS(true)
}

// SetPTS sets the PTS time of the signal's command. There will be no PTS adjustment using this function.
// If HasPTS is false, then it will have no effect until it is set to true. Also this command has no
// effect with a null splice command.
func (s *scte35) SetPTS(pts gots.PTS) {
	s.pts = pts
	s.commandInfo.SetPTS(s.pts & 0x01ffffffff) // truncate to fit in 33 bits
	// pts adjustment will be zero since the difference between adjusted and command pts is zero
}

// SetAdjustPTS will modify the pts adjustment field. The desired PTS value
// after adjustment should be passed, The adjustment value will be calculated
// during the call to Data().
func (s *scte35) SetAdjustPTS(pts gots.PTS) {
	// adjustment will be done by the function that generates the bytes
	s.pts = pts
}

// SetCommandInfo sets the object describing fields of the signal's splice
// command structure
func (s *scte35) SetCommandInfo(commandInfo SpliceCommand) {
	s.commandInfo = commandInfo
	s.commandType = s.commandInfo.CommandType()
}

// SetDescriptors sets a slice of the signals SegmentationDescriptors they
// should be sorted by descriptor weight (least important signals first)
func (s *scte35) SetDescriptors(descriptors []SegmentationDescriptor) {
	s.descriptors = descriptors
	for i := range s.descriptors {
		s.descriptors[i].(*segmentationDescriptor).spliceInfo = s
	}
}

// SetAlignmentStuffing sets how many stuffing bytes will be added to the SCTE35
// message at the end.
func (s *scte35) SetAlignmentStuffing(alignmentStuffing uint) {
	s.alignmentStuffing = alignmentStuffing
}

// SetTier sets which authorization tier this message was assigned to.
// The tier value of 0XFFF is the default and will ignored. When using tiers,
// the SCTE35 message must fit entirely into the ts payload without being split up.
func (s *scte35) SetTier(tier uint16) {
	s.tier = tier & 0xFFF
}
