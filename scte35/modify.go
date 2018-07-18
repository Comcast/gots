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
			protocolVersion:     0,
			encryptedPacket:     false, // no support for encryption
			encryptionAlgorithm: 0,
			hasPTS:              false,
			pts:                 0,
			cwIndex:             0,
			tier:                0xFFF, // ignore tier value

			spliceCommandLength: 0,
			commandType:         SpliceNull,
			commandInfo:         &spliceNull{},

			descriptorLoopLength: 0,
			descriptors:          []SegmentationDescriptor{},
			// TODO: data
		}
	scte35.psi =
		psi.PSI{
			PointerField:           0,
			TableID:                0xFC,
			SectionSyntaxIndicator: false,
			PrivateIndicator:       false,
			SectionLength:          1,
		}
	// TODO: generate data slice
	scte35.data = append(scte35.psi.Bytes(), scte35.data...)
	return scte35
}

func (s *scte35) Bytes() []byte {
	// splice command generate bytes
	psiBytes := s.psi.Bytes()
	spliceCommandBytes := s.commandInfo.Bytes()
	s.spliceCommandLength = uint16(len(spliceCommandBytes)) // can be set as 0xFFF (undefined), but calculate it anyways

	// generate bytes for splice descriptors
	spliceDescriptorBytes := make([]byte, 2)
	for i := range s.descriptors {
		spliceDescriptorBytes = append(spliceDescriptorBytes, s.descriptors[i].Bytes()...)
	}
	spliceDescriptorLoopLength := len(spliceDescriptorBytes) - 2
	spliceDescriptorBytes[0] = byte(spliceDescriptorLoopLength >> 8)
	spliceDescriptorBytes[1] = byte(spliceDescriptorLoopLength)

	minSectionLength := 11 + s.spliceCommandLength + s.descriptorLoopLength
	if minSectionLength > s.psi.SectionLength {
		s.psi.SectionLength = minSectionLength
	}

	// slices that point to the starting position of their names
	data := make([]byte, int(s.psi.PointerField)+4+int(s.psi.SectionLength))
	section := data[int(s.psi.PointerField)+int(psi.PSIHeaderLen):]
	commandInfo := section[11:]
	spliceDescriptor := commandInfo[len(spliceCommandBytes):]
	ptsAdj := s.pts - s.commandInfo.PTS()

	section[0] = s.protocolVersion
	if s.encryptedPacket {
		section[1] = 0x80 // 1000 0000
	}
	//s.cwIndex = 0xFF
	//s.tier = 0xFFF
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

	copy(data, psiBytes)
	copy(commandInfo, spliceCommandBytes)
	copy(spliceDescriptor, spliceDescriptorBytes)

	crc := gots.ComputeCRC(data[1 : len(data)-4])
	copy(data[len(data)-4:], crc)
	//s.data = s.data[s.psi.PointerField+1:]

	//original code
	return data
}

type scte35 struct {
	psi                  psi.PSI
	protocolVersion      byte
	encryptedPacket      bool  // not supported
	encryptionAlgorithm  uint8 // 6 bits
	hasPTS               bool
	pts                  gots.PTS // pts is stored adjusted in struct
	cwIndex              uint8
	tier                 uint16 // 12 bits
	spliceCommandLength  uint16 // 12 bits
	commandType          SpliceCommandType
	commandInfo          SpliceCommand
	descriptorLoopLength uint16
	descriptors          []SegmentationDescriptor
	crc32                uint32

	data []byte
}

// Only one version of the protocol exitst: Version 0
// func (s *scte35) SetProtocolVersion(value byte) {
// 	s.protocolVersion = protocolVersion
// }

func (s *scte35) SetHasPTS(flag bool) {
	s.hasPTS = flag
}

func (s *scte35) SetPTS(pts gots.PTS) {
	s.pts = pts
}

func (s *scte35) SetCommand(cmdType SpliceCommandType) {
	s.commandType = cmdType
}

func (s *scte35) SetCommandInfo(commandInfo SpliceCommand) {
	s.commandInfo = commandInfo
}

func (s *scte35) SetDescriptors(descriptors []SegmentationDescriptor) {
	s.descriptors = descriptors
}
