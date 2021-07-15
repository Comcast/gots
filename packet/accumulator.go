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

// Iotas to track the state of the accumulator
const (
	stateStarting = iota
	stateAccumulating
	stateDone
)

// Accumulator is used to gather multiple packets
// and return their concatenated payloads.
// Accumulator is not thread safe.
type Accumulator interface {
	// WritePacket adds a packet to the accumulator and returns got.ErrAccumulatorDone if done
	WritePacket(*Packet) (int, error)
	// Bytes returns the payload bytes from the underlying buffer
	Bytes() []byte
	// Packets returns the packets used to fill the payload buffer
	Packets() []*Packet
	// Reset resets the accumulator state
	Reset()
}
type accumulator struct {
	f       func([]byte) (bool, error)
	buf     *bytes.Buffer
	packets []*Packet
	state   int
}

// NewAccumulator creates a new packet accumulator that is done when
// the provided function returns done as true.
func NewAccumulator(f func(data []byte) (done bool, err error)) Accumulator {
	return &accumulator{
		f:     f,
		buf:   &bytes.Buffer{},
		state: stateStarting}
}

// Add a packet to the accumulator. If the added packet completes
// the accumulation, based on the provided doneFunc, gots.ErrAccumulatorDone is returned.
// Returns an error if the packet is not valid.
func (a *accumulator) WritePacket(pkt *Packet) (int, error) {
	switch a.state {
	case stateStarting:
		// need to check if the packet contains a payloadUnitStartIndicator to start
		if !PayloadUnitStartIndicator(pkt) {
			return PacketSize, gots.ErrNoPayloadUnitStartIndicator
		}

		a.packets = []*Packet{}
		a.state = stateAccumulating

	case stateAccumulating:
		// need to check if the packet contains a payloadUnitStartIndicator so we know
		// to drop old packets and start re-accumulation
		if PayloadUnitStartIndicator(pkt) {
			a.state = stateStarting
			return a.WritePacket(pkt)
		}

	case stateDone:
		return 0, gots.ErrAccumulatorDone
	}

	var cpyPkt = &Packet{}
	copy(cpyPkt[:], pkt[:])
	a.packets = append(a.packets, cpyPkt)

	if b, err := Payload(pkt); err != nil {
		return PacketSize, err
	} else if _, err := a.buf.Write(b); err != nil {
		return PacketSize, err
	}

	if done, err := a.f(a.buf.Bytes()); err != nil {
		return PacketSize, err
	} else if done {
		a.state = stateDone
		return PacketSize, gots.ErrAccumulatorDone
	}

	return PacketSize, nil
}

// Bytes returns the payload bytes from the underlying buffer
func (a *accumulator) Bytes() []byte {
	return a.buf.Bytes()
}

// Packets returns the packets used to fill the payload buffer
// NOTE: Not thread safe
func (a *accumulator) Packets() []*Packet {
	return a.packets
}

// Reset resets the accumulator state
func (a *accumulator) Reset() {
	a.state = stateStarting
	a.buf.Reset()
}
