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

var (
	emptyByteArray []byte
)

type doneFunc func([]byte) (bool, error)

type accumulator struct {
	f       doneFunc
	packets []Packet
}

// NewAccumulator creates a new packet accumulator
// that is done when the provided doneFunc returns true.
// PacketAccumulator is not thread safe
func NewAccumulator(f doneFunc) Accumulator {
	return &accumulator{f: f}
}

// Add a packet to the accumulator. If the added packet completes
// the accumulation, based on the provided doneFunc, true is returned.
// Returns an error if the packet is not valid.
func (a *accumulator) Add(pkt Packet) (bool, error) {
	if badLen(pkt) {
		return false, gots.ErrInvalidPacketLength
	}
	// technically we could get a packet without a payload.  Check this and
	// return false if we get one
	p, err := ContainsPayload(pkt)
	if err != nil {
		return false, err
	} else if !p {
		return false, nil
	}
	if payloadUnitStartIndicator(pkt) {
		a.packets = make([]Packet, 0)
	} else if len(a.packets) == 0 {
		// First packet must have payload unit start indicator
		return false, gots.ErrNoPayloadUnitStartIndicator
	}
	pktCopy := make(Packet, PacketSize)
	copy(pktCopy, pkt)
	a.packets = append(a.packets, pktCopy)
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
			return emptyByteArray, err
		}
		buf.Write(pay)
	}
	return buf.Bytes(), nil
}

func (a *accumulator) Packets() []Packet {
	return a.packets
}

func (a *accumulator) Reset() {
	a.packets = nil
}
