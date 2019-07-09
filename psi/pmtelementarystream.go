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
	"fmt"
)

// PmtElementaryStream represents an elementary stream inside a PMT
type PmtElementaryStream interface {
	PmtStreamType
	ElementaryPid() uint16
	Descriptors() []PmtDescriptor
	MaxBitRate() uint64
}

type pmtElementaryStream struct {
	PmtStreamType
	elementaryPid uint16
	descriptors   []PmtDescriptor
}

const (
	// BitsPerByte is the number of bits in a byte
	BitsPerByte = 8
	// MaxBitRateBytesPerSecond is the maximum bit rate per second in a profile
	MaxBitRateBytesPerSecond = 50
)

// NewPmtElementaryStream creates a new PmtElementaryStream.
func NewPmtElementaryStream(streamType uint8, elementaryPid uint16, descriptors []PmtDescriptor) PmtElementaryStream {
	es := &pmtElementaryStream{}
	es.PmtStreamType = LookupPmtStreamType(streamType)
	es.elementaryPid = elementaryPid
	es.descriptors = descriptors
	return es
}

func (es *pmtElementaryStream) ElementaryPid() uint16 {
	return es.elementaryPid
}

func (es *pmtElementaryStream) Descriptors() []PmtDescriptor {
	return es.descriptors
}

// MaxBitRate returns the value of the PmtElementaryStreams maximum bitrate in bits-per-second.
// See Section 2.6.27 of ISO-13818 for more information.
func (es *pmtElementaryStream) MaxBitRate() uint64 {
	for _, desc := range es.descriptors {
		if desc.IsMaximumBitrateDescriptor() {
			return uint64(desc.DecodeMaximumBitRate()) * BitsPerByte * MaxBitRateBytesPerSecond
		}
	}
	return 0
}

func (es *pmtElementaryStream) String() string {
	descriptors := es.descriptors
	var descriptorsBuf bytes.Buffer
	if len(descriptors) > 0 {
		descriptorsBuf.WriteString(",")
		for i, descriptor := range descriptors {
			descriptorsBuf.WriteString(fmt.Sprintf("descriptor%d='%v'", i, descriptor))
			i++
			if i < len(descriptors) {
				descriptorsBuf.WriteString(",")
			}
		}
	}
	return fmt.Sprintf("ElementaryStream[pid=%d,%v%s]", es.elementaryPid, es.PmtStreamType, descriptorsBuf.String())
}
