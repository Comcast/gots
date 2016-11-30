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
	"io"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/packet"
)

const (
	// PatPid is the PID of a PAT. By definition this value is zero.
	PatPid = uint16(0)
)

// The Program Association Table (PAT) lists the programs available in transport
// stream.  Currently it is assumed that it will always be processing a Single
// Program Transport Stream (SPTS), and thus will only ever receive a PAT that
// contains a single programNumber to PID mapping.
type pat []byte

// NewPAT constructs a new PAT from the provided bytes.
// patBytes should be concatenated packet payload contents.
// If a 188 byte slice is passed in, NewPAT tries to help and
// treats it as TS packet and builds a PAT from the packet payload.
func NewPAT(patBytes []byte) (PAT, error) {
	if len(patBytes) < 13 {
		return nil, gots.ErrInvalidPATLength
	}

	if len(patBytes) == 188 {
		var err error
		patBytes, err = packet.Payload(patBytes)
		if err != nil {
			return nil, err
		}
	}

	return pat(patBytes), nil
}

// NumPrograms returns the number of programs in this PAT
func (pat pat) NumPrograms() int {
	sectionLength := SectionLength(pat)
	numPrograms := int((sectionLength -
		2 - // Transport Stream ID
		1 - // Reserved|VersionNumber|CurrentNextIndicator
		1 - // Section Number
		1 - // Last Section Number
		4) / // CRC32
		4) // Divided by 4 bytes per program
	return numPrograms
}

// DEPRECATED - use ProgramMap for new code
func (pat pat) ProgramMapPid() uint16 {
	if pat.NumPrograms() != 1 {
		// This method has undefined behavior for multi program transport streams.
		// Please use pat.ProgramMap() for a map of PNs and PIDs for PMTs
		return 8192
	}

	pm := pat.ProgramMap()
	for _, v := range pm {
		return v
	}

	return 8192
}

// DEPRECATED - use ProgramMap for new code
func (pat pat) ProgramNumber() uint16 {
	if pat.NumPrograms() != 1 {
		// This method has undefined behavior for multi program transport streams.
		// Please use pat.ProgramMap() for a map of PNs and PIDs for PMTs
		return 0
	}

	pm := pat.ProgramMap()
	for k := range pm {
		return k
	}

	return 0
}

// ProgramMap returns a map of program numbers and PIDs of the PMTs
func (pat pat) ProgramMap() map[uint16]uint16 {
	m := make(map[uint16]uint16)

	counter := 8 // skip table id et al

	for i := 0; i < pat.NumPrograms(); i++ {
		pn := (uint16(pat[counter+1]) << 8) | uint16(pat[counter+2])

		// ignore the top three (reserved) bits
		pid := uint16(pat[counter+3])&0x1f<<8 | uint16(pat[counter+4])

		m[pn] = pid

		counter += 4
	}

	return m
}

// ReadPAT extracts a PAT from a reader of a TS stream. It will read until a
// PAT packet is found or EOF is reached.
// It returns a new PAT object parsed from the packet, if found, and otherwise
// returns an error.
func ReadPAT(r io.Reader) (PAT, error) {
	pkt := make(packet.Packet, packet.PacketSize)
	var pat PAT
	for pat == nil {
		if _, err := io.ReadFull(r, pkt); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil, gots.ErrPATNotFound
			}
			return nil, err
		}
		isPat, err := packet.IsPat(pkt)
		if err != nil {
			return nil, err
		}
		if isPat {
			pay, err := packet.Payload(pkt)
			if err != nil {
				return nil, err
			}
			cp := make([]byte, len(pay))
			copy(cp, pay)
			pat, err := NewPAT(cp)
			if err != nil {
				return nil, err
			}
			return pat, nil
		}
	}
	return nil, gots.ErrPATNotFound
}
