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

package packet

import (
	"bytes"

	"github.com/Comcast/gots"
)

// Accumulator is used to gather multiple packets
// and return their concatenated payloads.
// Accumulator is not thread safe.
type Accumulator interface {
	// Add adds a packet to the accumulator and returns true if done.
	Add([]byte) (bool, error)
	// Parse returns the concatenated payloads of all the packets that have been added to the accumulator
	Parse() ([]byte, error)
	// Packets returns the accumulated packets
	Packets() []*Packet
	// Reset clears all packets in the accumulator
	Reset()
}
type accumulator struct {
	f       func([]byte) (bool, error)
	packets []*Packet
}

// NewAccumulator creates a new packet accumulator that is done when
// the provided function returns done as true.
func NewAccumulator(f func(data []byte) (done bool, err error)) Accumulator {
	return &accumulator{f: f}
}

// Add a packet to the accumulator. If the added packet completes
// the accumulation, based on the provided doneFunc, true is returned.
// Returns an error if the packet is not valid.
func (a *accumulator) Add(pkt []byte) (bool, error) {
	if badLen(pkt) {
		return false, gots.ErrInvalidPacketLength
	}
	var pp Packet
	copy(pp[:], pkt)
	// technically we could get a packet without a payload.  Check this and
	// return false if we get one
	p := ContainsPayload(&pp)
	if !p {
		return false, nil
	}
	// need to check if the packet contains a payloadUnitStartIndicator so we know
	// to drop old packets and re-accumulate a new scte signal
	if PayloadUnitStartIndicator(&pp) {
		a.Reset()
	}
	if !PayloadUnitStartIndicator(&pp) && len(a.packets) == 0 {
		// First packet must have payload unit start indicator
		return false, gots.ErrNoPayloadUnitStartIndicator
	}
	a.packets = append(a.packets, &pp)
	b, err := a.Parse()
	if err != nil {
		return false, err
	}
	done, err := a.f(b)
	if err != nil {
		return false, err
	}
	return done, nil
}

// Parses the accumulated packets and returns the
// concatenated payloads or any error that occurred, not both
func (a *accumulator) Parse() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, pkt := range a.packets {
		pay, err := Payload(pkt)
		if err != nil {
			return nil, err
		}
		buf.Write(pay)
	}
	return buf.Bytes(), nil
}

func (a *accumulator) Packets() []*Packet {
	return a.packets
}

func (a *accumulator) Reset() {
	a.packets = nil
}
