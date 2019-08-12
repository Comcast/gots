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

package psi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/packet"
)

const (
	programInfoLengthOffset         = 10 // includes PSIHeaderLen
	pmtEsDescriptorStaticLen uint16 = 5
)

// Unaccounted bytes before the end of the SectionLength field
const (
	// Pointerfield(1) + table id(1) + flags(.5) + section length (2.5)
	PSIHeaderLen uint16 = 4
	CrcLen       uint16 = 4
)

// PMT is a Program Map Table.
type PMT interface {
	Pids() []uint16
	IsPidForStreamWherePresentationLagsEbp(pid uint16) bool
	ElementaryStreams() []PmtElementaryStream
	RemoveElementaryStreams(pids []uint16)
	String() string
}

type pmt struct {
	pids              []uint16
	elementaryStreams []PmtElementaryStream
}

// PmtAccumulatorDoneFunc is a doneFunc that can be used for packet accumulation
// to create a PMT
func PmtAccumulatorDoneFunc(b []byte) (bool, error) {
	if len(b) < 1 {
		return false, nil
	}

	start := 1 + int(PointerField(b))
	if len(b) < start {
		return false, nil
	}

	sectionBytes := b[start:]
	for len(sectionBytes) > 2 && sectionBytes[0] != 0xFF {
		tableLength := sectionLength(sectionBytes)
		if len(sectionBytes) < int(tableLength)+3 {
			return false, nil
		}
		sectionBytes = sectionBytes[3+tableLength:]
	}

	return true, nil
}

// NewPMT Creates a new PMT from the given bytes.
// pmtBytes should be concatenated packet payload contents.
func NewPMT(pmtBytes []byte) (PMT, error) {
	pmt := &pmt{}
	err := pmt.parseTables(pmtBytes)
	if err != nil {
		return nil, err
	}
	return pmt, nil
}

func (p *pmt) parseTables(pmtBytes []byte) error {
	sectionBytes := pmtBytes[1+PointerField(pmtBytes):]

	for len(sectionBytes) > 2 && sectionBytes[0] != 0xFF {
		tableLength := sectionLength(sectionBytes)

		if tableID(sectionBytes) == 0x2 {
			err := p.parsePMTSection(sectionBytes[0 : 3+tableLength])
			if err != nil {
				return err
			}
		}
		sectionBytes = sectionBytes[3+tableLength:]
	}

	return nil
}

func (p *pmt) parsePMTSection(pmtBytes []byte) error {
	var pids []uint16
	var elementaryStreams []PmtElementaryStream
	sectionLength := sectionLength(pmtBytes)

	if len(pmtBytes) < programInfoLengthOffset {
		return gots.ErrParsePMTDescriptor
	}

	programInfoLength := uint16(pmtBytes[programInfoLengthOffset]&0x0f)<<8 |
		uint16(pmtBytes[programInfoLengthOffset+1])

	// start at the stream descriptors, parse until the CRC
	for offset := programInfoLengthOffset + 2 + programInfoLength; offset < PSIHeaderLen+sectionLength-pmtEsDescriptorStaticLen-CrcLen; {
		elementaryStreamType := uint8(pmtBytes[offset])
		elementaryPid := uint16(pmtBytes[offset+1]&0x1f)<<8 | uint16(pmtBytes[offset+2])
		pids = append(pids, elementaryPid)
		infoLength := uint16(pmtBytes[offset+3]&0x0f)<<8 | uint16(pmtBytes[offset+4])

		// Move past the es descriptor static data
		offset += pmtEsDescriptorStaticLen
		var descriptors []PmtDescriptor
		if infoLength != 0 && int(infoLength+offset) < len(pmtBytes) {
			var descriptorOffset uint16
			for descriptorOffset < infoLength {
				tag := uint8(pmtBytes[offset+descriptorOffset])
				descriptorOffset++
				descriptorLength := uint16(pmtBytes[offset+descriptorOffset])
				descriptorOffset++
				startPos := offset + descriptorOffset
				endPos := int(offset + descriptorOffset + descriptorLength)
				if endPos < len(pmtBytes) {
					data := pmtBytes[startPos:endPos]
					descriptors = append(descriptors, NewPmtDescriptor(tag, data))
				} else {
					return gots.ErrParsePMTDescriptor
				}
				descriptorOffset += descriptorLength
			}
			offset += infoLength
		}
		es := NewPmtElementaryStream(elementaryStreamType, elementaryPid, descriptors)
		elementaryStreams = append(elementaryStreams, es)
	}

	p.pids = pids
	p.elementaryStreams = elementaryStreams
	return nil
}

func (p *pmt) Pids() []uint16 {
	return p.pids
}

func (p *pmt) ElementaryStreams() []PmtElementaryStream {
	return p.elementaryStreams
}

// RemoveElementaryStreams removes elementary streams in the pmt of the given pids
func (p *pmt) RemoveElementaryStreams(removePids []uint16) {
	for _, pid := range removePids {
		for j, s := range p.elementaryStreams {
			if pid == s.ElementaryPid() {
				p.elementaryStreams = append(p.elementaryStreams[:j], p.elementaryStreams[j+1:]...)
				break
			}
		}
	}

	var filteredPids []uint16

	for _, es := range p.elementaryStreams {
		filteredPids = append(filteredPids, es.ElementaryPid())
	}

	p.pids = filteredPids
}

func (p *pmt) IsPidForStreamWherePresentationLagsEbp(pid uint16) bool {
	for _, s := range p.elementaryStreams {
		if pid == s.ElementaryPid() {
			return s.IsStreamWherePresentationLagsEbp()
		}
	}
	return false
}

func (p *pmt) String() string {
	var buf bytes.Buffer
	buf.WriteString("PMT[")
	i := 0
	for _, es := range p.elementaryStreams {
		buf.WriteString(fmt.Sprintf("%v", es))
		i++
		if i < len(p.elementaryStreams) {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")

	return buf.String()
}

// FilterPMTPacketsToPids filters the PMT contents of the provided packet to the PIDs provided and returns a new packet. For example: if the provided PMT has PIDs 101, 102, and 103 and the provides PIDs are 101 and 102, the new PMT will have only descriptors for PID 101 and 102. The descriptor for PID 103 will be stripped from the new PMT packet.
func FilterPMTPacketsToPids(packets []*packet.Packet, pids []uint16) []*packet.Packet {
	// make sure we have packets
	if len(packets) == 0 {
		return nil
	}
	// Mush the payloads of all PMT packets into one []byte
	var pmtByteBuffer bytes.Buffer
	for i := 0; i < len(packets); i++ {
		pay, err := packet.Payload(packets[i])
		if err != nil {
			return nil
		}
		pmtByteBuffer.Write(pay)
	}

	pmtPayload := pmtByteBuffer.Bytes()
	// include +1 to account for the PointerField field itself
	pointerField := PointerField(pmtPayload) + 1

	var filteredPMT bytes.Buffer

	// Copy everything from the pointerfield offset and move the pmtPayload slice to the start of that
	filteredPMT.Write(pmtPayload[:pointerField])
	pmtPayload = pmtPayload[pointerField:]
	// Copy the first 12 bytes of the PMT packet. Only section_length will change.
	filteredPMT.Write(pmtPayload[:programInfoLengthOffset+2])

	// Get the section length
	sectionLength := uint16(pmtPayload[1]&0x0f)<<8 + uint16(pmtPayload[2])

	// Get program info length
	programInfoLength := uint16(pmtPayload[programInfoLengthOffset]&0x0f)<<8 | uint16(pmtPayload[programInfoLengthOffset+1])
	if programInfoLength != 0 {
		filteredPMT.Write(pmtPayload[programInfoLengthOffset+2 : programInfoLengthOffset+2+programInfoLength])
	}

	for offset := programInfoLengthOffset + 2 + programInfoLength; offset < PSIHeaderLen+sectionLength-pmtEsDescriptorStaticLen-CrcLen; {
		elementaryPid := uint16(pmtPayload[offset+1]&0x1f)<<8 | uint16(pmtPayload[offset+2])
		infoLength := uint16(pmtPayload[offset+3]&0x0f)<<8 | uint16(pmtPayload[offset+4])

		// This is an ES PID we want to keep
		if pidIn(pids, elementaryPid) {
			// write out the whole es info
			filteredPMT.Write(pmtPayload[offset : offset+pmtEsDescriptorStaticLen+infoLength])
		}
		offset += pmtEsDescriptorStaticLen + infoLength
	}

	// Create the new section length
	fPMT := filteredPMT.Bytes()
	// section_length is the length of data (including the CRC) in bytes following the section length field (ISO13818: 2.4.4.9)
	// This will be the length of our buffer - (Bytes preceding section_length) + CRC
	// Bytes preceding = 4 + PointerField value and the CRC = 4, so it turns out to be the length of the buffer - PointerField field
	// -1 because we previously added 1 for the pointerfield field itself
	newSectionLength := uint16(len(fPMT)) - uint16(pointerField-1)
	sectionLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(sectionLengthBytes, newSectionLength)
	fPMT[pointerField+1] = (fPMT[pointerField+1] & 0xf0) | sectionLengthBytes[0]
	fPMT[pointerField+2] = sectionLengthBytes[1]

	// Recalculate the CRC
	fPMT = append(fPMT, gots.ComputeCRC(fPMT[pointerField:])...)

	var filteredPMTPackets []*packet.Packet
	for _, pkt := range packets {
		var pktBuf bytes.Buffer
		header := packet.Header(pkt)
		pktBuf.Write(header)
		if len(fPMT) > 0 {
			toWrite := safeSlice(fPMT, 0, packet.PacketSize-len(header))
			// truncate fPMT to the remaining bytes
			if len(toWrite) < len(fPMT) {
				fPMT = fPMT[len(toWrite):]
			} else {
				fPMT = nil
			}
			pktBuf.Write(toWrite)
		} else {
			// all done
			break
		}
		filteredPMTPackets = append(filteredPMTPackets, padPacket(&pktBuf))
	}
	return filteredPMTPackets
}

// IsPMT returns true if the provided packet is a PMT
// defined by the PAT provided. Returns ErrNilPAT if pat
// is nil, or any error encountered in parsing the PID
// of pkt.
func IsPMT(pkt *packet.Packet, pat PAT) (bool, error) {
	if pat == nil {
		return false, gots.ErrNilPAT
	}

	pmtMap := pat.ProgramMap()
	pid := packet.Pid(pkt)

	for _, map_pid := range pmtMap {
		if pid == map_pid {
			return true, nil
		}
	}

	return false, nil
}

func safeSlice(byteArray []byte, start, end int) []byte {
	if end < len(byteArray) {
		return byteArray[start:end]
	}
	return byteArray[start:len(byteArray)]
}

func padPacket(buf *bytes.Buffer) *packet.Packet {
	var pkt packet.Packet
	for i := copy(pkt[:], buf.Bytes()); i < packet.PacketSize; i++ {
		pkt[i] = 0xff
	}
	return &pkt
}

func pidIn(pids []uint16, target uint16) bool {
	for _, pid := range pids {
		if pid == target {
			return true
		}
	}

	return false
}

// ReadPMT extracts a PMT from a reader of a TS stream. It will read until PMT
// packet(s) are found or EOF is reached.
// It returns a new PMT object parsed from the packet(s), if found, and
// otherwise returns an error.
func ReadPMT(r io.Reader, pid uint16) (PMT, error) {
	var pkt packet.Packet
	var err error
	var pmt PMT

	pmtAcc := packet.NewAccumulator(PmtAccumulatorDoneFunc)
	done := false

	for !done {
		if _, err := io.ReadFull(r, pkt[:]); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil, gots.ErrPMTNotFound
			}
			return nil, err
		}
		currPid := packet.Pid(&pkt)
		if currPid != pid {
			continue
		}
		done, err = pmtAcc.Add(pkt[:])
		if err != nil {
			return nil, err
		}
		if done {
			b, err := pmtAcc.Parse()
			if err != nil {
				return nil, err
			}
			pmt, err = NewPMT(b)
			if err != nil {
				return nil, err
			}
			if len(pmt.Pids()) == 0 {
				done = false
				pmtAcc = packet.NewAccumulator(PmtAccumulatorDoneFunc)
			}
		}
	}
	return pmt, nil
}
