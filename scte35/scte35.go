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
)

type scte35 struct {
	commandType SpliceCommandType
	commandInfo SpliceCommand
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
		default:
			return gots.ErrSCTE35UnsupportedSpliceCommand
		}
		// descriptor_loop_length 2 + CRC 4
		if buf.Len() < 2+int(psi.CrcLen) {
			return gots.ErrInvalidSCTE35Length
		}
		// parse descriptors
		descLoopLen := binary.BigEndian.Uint16(buf.Next(2))
		if buf.Len() < int(descLoopLen+psi.CrcLen) {
			return gots.ErrInvalidSCTE35Length
		}
		for bytesRead := uint16(0); bytesRead < descLoopLen; {
			descTag := readByte()
			descLen := readByte()
			// Make sure a bad descriptorLen doesn't kill us
			if descLoopLen-bytesRead-2 < uint16(descLen) {
				return gots.ErrInvalidSCTE35Length
			}
			if descTag != segDescTag {
				// Not interested in descriptors that are not
				// SegmentationDescriptors
				buf.Next(int(descLen))
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
	//remove the pointer field and associated data off the top so we only get the
	//table data
	s.data = data[psi.PointerField(data)+1:]
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

func (s *scte35) Data() []byte {
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
