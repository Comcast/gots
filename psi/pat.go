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
	"github.com/comcast/gots"
	"github.com/comcast/gots/packet"
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

// ProgramMapPid returns the PID of the PMT
func (pat pat) ProgramMapPid() uint16 {
	return (uint16(pat[11]) & 0x1f << 8) | uint16(pat[12])
}

// ProgramNumber returns the program number for this PAT
func (pat pat) ProgramNumber() uint16 {
	return (uint16(pat[9]) << 8) | uint16(pat[10])
}
