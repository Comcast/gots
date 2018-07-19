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
	"github.com/Comcast/gots/psi"
)

func CreateSCTE35() SCTE35 {
	scte35 :=
		&scte35{
			protocolVersion:     0,     // only version 0 exists
			encryptedPacket:     false, // no support for encryption
			encryptionAlgorithm: 0,     // no encryption support, no way to change
			hasPTS:              false, // obtained from splice command, null splice command doesnt have it
			pts:                 0,     // no pts
			cwIndex:             0,     // undefined, without encryption
			tier:                0xFFF, // ignore tier value

			spliceCommandLength: 0,             // null command has no length
			commandType:         SpliceNull,    // null command type by default
			commandInfo:         &spliceNull{}, // info pooints to null command

			descriptorLoopLength: 0,                          // no descriptors by default
			descriptors:          []SegmentationDescriptor{}, // empty slice of descriptors

			updateBytes: true, // update the bytes on the next call to Data()
		}
	scte35.psi =
		psi.TableHeader{
			TableID:                0xFC,  // always 0xFC for scte35
			SectionSyntaxIndicator: false, // always false for scte35
			PrivateIndicator:       false, // always false for scte35
			SectionLength:          0,     // to be calculated when converting to bytes
		}
	return scte35
}

func (s *scte35) generateData() {
	// splice command generate bytes
	pointerBytes := []byte{} // psi.NewPointerField(0)
	pointerLength := len(pointerBytes)

	tableHeaderBytes := s.psi.Bytes()
	tableHeaderLength := len(tableHeaderBytes)

	spliceCommandBytes := s.commandInfo.Data()
	// spliceCommandLength can be set as 0xFFF (undefined), but calculate it anyways
	spliceCommandLength := len(spliceCommandBytes)
	s.spliceCommandLength = uint16(spliceCommandLength)

	// generate bytes for splice descriptors
	descriptorBytes := make([]byte, 2)
	for i := range s.descriptors {
		descriptorBytes = append(descriptorBytes, s.descriptors[i].Data()...)
	}
	descriptorLoopLength := len(descriptorBytes) - 2
	descriptorBytes[0] = byte(descriptorLoopLength >> 8)
	descriptorBytes[1] = byte(descriptorLoopLength)

	const staticFieldsLength = 13
	const crcLength = 4

	sectionLength := staticFieldsLength + spliceCommandLength + descriptorLoopLength + crcLength + s.alignmentStuffing
	s.psi.SectionLength = uint16(sectionLength)

	// slices that point to the starting position of their names
	data := make([]byte, pointerLength+tableHeaderLength+sectionLength)
	tableHeader := data[pointerLength:]
	section := tableHeader[tableHeaderLength:]
	spliceCommand := section[11:]
	spliceDescriptor := spliceCommand[spliceCommandLength:]
	crc := data[len(data)-crcLength:]

	// TODO: This is wrong
	ptsAdj := s.pts - s.commandInfo.PTS()

	if s.encryptedPacket {
		section[1] = 0x80 // 1000 0000
	}
	section[0] = s.protocolVersion                      // 1111 1111
	section[1] = (s.encryptionAlgorithm & 0x3F) << 1    // 0111 1110
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

	copy(data, pointerBytes)
	copy(tableHeader, tableHeaderBytes)
	copy(spliceCommand, spliceCommandBytes)
	copy(spliceDescriptor, descriptorBytes)

	crcBytes := gots.ComputeCRC(tableHeader[:len(data)-crcLength])
	copy(crc, crcBytes)
	s.data = data
	s.updateBytes = false
}

func (s *scte35) SetHasPTS(flag bool) {
	s.hasPTS = flag
	s.updateBytes = true
}

func (s *scte35) SetPTS(pts gots.PTS) {
	s.pts = pts
	s.updateBytes = true
}

func (s *scte35) AdjustPTS(pts gots.PTS) {
	// adjustment will be done by the function that generates the bytes
	s.pts = pts
	s.updateBytes = true
}

func (s *scte35) SetCommand(cmdType SpliceCommandType) {
	s.commandType = cmdType
	s.updateBytes = true
}

func (s *scte35) SetCommandInfo(commandInfo SpliceCommand) {
	s.commandInfo = commandInfo
	s.updateBytes = true
}

func (s *scte35) SetDescriptors(descriptors []SegmentationDescriptor) {
	s.descriptors = descriptors
	s.updateBytes = true
}

func (s *scte35) SetAlignmentStuffing(alignmentStuffing int) {
	s.alignmentStuffing = alignmentStuffing
	s.updateBytes = true
}

func (s *scte35) SetTier(tier uint16) {
	s.tier = tier
	s.updateBytes = true
}
