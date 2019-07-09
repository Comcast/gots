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
	"errors"
	"io"

	"github.com/Comcast/gots"
	"github.com/Comcast/gots/packet"
)

const (
	// PatPid is the PID of a PAT. By definition this value is zero.
	PatPid = uint16(0)
)

// PAT interface represents operations on a Program Association Table. Currently only single program transport streams (SPTS)are supported
type PAT interface {
	NumPrograms() int
	ProgramMap() map[uint16]uint16
	SPTSpmtPID() (uint16, error)
}

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
		var pkt packet.Packet
		copy(pkt[:], patBytes)
		var err error
		patBytes, err = packet.Payload(&pkt)
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

// ProgramMap returns a map of program numbers and PIDs of the PMTs
func (pat pat) ProgramMap() map[uint16]uint16 {
	m := make(map[uint16]uint16)

	counter := 8 // skip table id et al

	for i := 0; i < pat.NumPrograms(); i++ {
		pn := (uint16(pat[counter+1]) << 8) | uint16(pat[counter+2])

		// ignore the top three (reserved) bits
		pid := uint16(pat[counter+3])&0x1f<<8 | uint16(pat[counter+4])

		// A value of 0 is reserved for a NIT packet identifier.
		if pn > 0 {
			m[pn] = pid
		}

		counter += 4
	}

	return m
}

// SPTSpmtPID returns the PMT PID if and only if this pat is for a single program transport stream. If this pat is for a multiprogram transport stream, an error is returned.
func (pat pat) SPTSpmtPID() (uint16, error) {
	if pat.NumPrograms() > 1 {
		return 0, errors.New("Not a single program transport stream")
	}
	for _, pid := range pat.ProgramMap() {
		return pid, nil
	}
	return 0, errors.New("No programs in transport stream")
}

// ReadPAT extracts a PAT from a reader of a TS stream. It will read until a
// PAT packet is found or EOF is reached.
// It returns a new PAT object parsed from the packet, if found, and otherwise
// returns an error.
func ReadPAT(r io.Reader) (PAT, error) {
	var pkt packet.Packet
	var pat PAT
	for pat == nil {
		if _, err := io.ReadFull(r, pkt[:]); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, err
		}
		isPat := packet.IsPat(&pkt)

		if isPat {
			pay, err := packet.Payload(&pkt)
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
