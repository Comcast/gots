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

// Package packet is used for reading and manipulating packets in an MPEG transport stream
package packet

const (
	// PacketSize is the expected size of a packet in bytes
	PacketSize = 188
	// SyncByte is the expected value of the sync byte
	SyncByte = 71 // 0x47 (0100 0111)
	// NullPacketPid is the pid reserved for null packets
	NullPacketPid = 8191 // 0x1FFF
)

// Packet is the basic unit in a transport stream.
type Packet [PacketSize]byte

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
