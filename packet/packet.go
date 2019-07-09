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

import "github.com/Comcast/gots"

const (
	// PacketSize is the expected size of a packet in bytes
	PacketSize = 188
	// SyncByte is the expected value of the sync byte
	SyncByte = 71 // 0x47 (0100 0111)
	// NullPacketPid is the pid reserved for null packets
	NullPacketPid = 8191 // 0x1FFF
)

// TransportScramblingControlOptions is a set of constants for
// selecting the transport scrambling control.
type TransportScramblingControlOptions byte

const (
	NoScrambleFlag      TransportScramblingControlOptions = 0 // 00
	ScrambleEvenKeyFlag TransportScramblingControlOptions = 2 // 10
	ScrambleOddKeyFlag  TransportScramblingControlOptions = 3 // 11
)

// Packet is the basic unit in a transport stream.
type Packet [PacketSize]byte

// PayloadUnitStartIndicator (PUSI) is a flag that indicates the start of PES data
// or PSI  (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func PayloadUnitStartIndicator(packet *Packet) bool {
	return packet[1]&0x040 != 0
}

// PID is the Packet Identifier.  Each table or elementary stream in the
// transport stream is identified by a PID.  The PID is contained in the 13
// bits that span the last 5 bits of second byte and all bits in the byte.
func Pid(packet *Packet) uint16 {
	return uint16(packet[1]&0x1f)<<8 | uint16(packet[2])
}

// ContainsPayload is a flag that indicates the packet has a payload.  The flag is
// contained in the 3rd bit of the 4th byte of the Packet.
func ContainsPayload(packet *Packet) bool {
	return packet[3]&0x10 != 0
}

// ContainsAdaptationField is a flag that indicates the packet has an adaptation field.
func ContainsAdaptationField(packet *Packet) bool {
	return packet[3]&0x20 != 0
}

// ContinuityCounter is a 4-bit sequence number of payload packets. Incremented
// only when a payload is present (see ContainsPayload() above).
func ContinuityCounter(packet *Packet) uint8 {
	return packet[3] & uint8(0x0f)
}

// IsNull returns true if the provided packet is a Null packet
// (i.e., PID == 0x1fff (8191)).
func IsNull(packet *Packet) bool {
	return Pid(packet) == NullPacketPid
}

// IsPat returns true if the provided packet is a PAT
func IsPat(packet *Packet) bool {
	return Pid(packet) == 0
}

// badLen returns true if the packet is not of valid length
func badLen(packet []byte) bool {
	return len(packet) != PacketSize
}

// Returns the index of the first byte of Payload data in packetBytes.
func payloadStart(packet *Packet) int {
	var dataOffset = int(4) // packet header bytes
	if ContainsAdaptationField(packet) {
		afLength := int(packet[4])
		dataOffset += 1 + afLength
	}

	return dataOffset
}

// Payload returns a slice containing the packet payload. If the packet
// does not have a payload, an empty byte slice is returned
func Payload(packet *Packet) ([]byte, error) {
	if !ContainsPayload(packet) {
		return nil, gots.ErrNoPayload
	}
	start := payloadStart(packet)
	if start > len(packet) {
		return nil, gots.ErrInvalidPacketLength
	}
	pay := packet[start:]
	return pay, nil
}

// IncrementCC creates a new packet where the new packet has
// a continuity counter that is increased by one
func IncrementCC(packet *Packet) *Packet {
	var newPacket Packet
	copy(newPacket[:], packet[:])
	ccByte := newPacket[3]
	newCC := increment4BitInt(ccByte)
	newCCByte := (ccByte & byte(0xf0)) | newCC
	newPacket[3] = newCCByte
	return &newPacket
}

// ZeroCC creates a new packet where the new packet has
// a continuity counter that zero
func ZeroCC(packet *Packet) *Packet {
	var newPacket Packet
	copy(newPacket[:], packet[:])
	ccByte := newPacket[3]
	newCCByte := ccByte & byte(0xf0)
	newPacket[3] = newCCByte
	return &newPacket
}
func increment4BitInt(cc uint8) uint8 {
	return (cc + 1) & 0x0f
}

// SetCC creates a new packet where the new packet has
// the continuity counter provided
func SetCC(packet *Packet, newCC uint8) *Packet {
	var newPacket Packet
	copy(newPacket[:], packet[:])
	ccByte := newPacket[3]
	newCCByte := (ccByte & byte(0xf0)) | newCC
	newPacket[3] = newCCByte
	return &newPacket
}

// PESHeader returns a byte slice containing the PES header if the Packet contains one,
// otherwise returns an error
func PESHeader(packet *Packet) ([]byte, error) {
	if PayloadUnitStartIndicator(packet) {
		pay, err := Payload(packet)
		if err != nil {
			return nil, err
		}
		if len(pay) > 3 && pay[0] == 0 && pay[1] == 0 && pay[2] == 1 {
			return pay, nil
		}
	}
	return nil, gots.ErrNoPayload
}

// Header Returns a slice containing the Packer Header.
func Header(packet *Packet) []byte {
	start := payloadStart(packet)
	return packet[:start]
}

// Equal returns true if the bytes of the two packets are equal
func Equal(a, b *Packet) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
