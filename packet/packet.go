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
	"io"

	"github.com/Comcast/gots"
)

var emptyByteSlice []byte

// PayloadUnitStartIndicator (PUSI) is a flag that indicates the start of PES data
// or PSI  (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func PayloadUnitStartIndicator(packet Packet) (bool, error) {
	if badLen(packet) {
		return false, gots.ErrInvalidPacketLength
	}
	return payloadUnitStartIndicator(packet), nil
}
func payloadUnitStartIndicator(packet Packet) bool {
	return packet[1]&0x040 != 0
}

// PID is the Packet Identifier.  Each table or elementary stream in the
// transport stream is identified by a PID.  The PID is contained in the 13
// bits that span the last 5 bits of second byte and all bits in the byte.
func Pid(packet Packet) (uint16, error) {
	if badLen(packet) {
		return 0, gots.ErrInvalidPacketLength
	}
	return pid(packet), nil
}
func pid(packet Packet) uint16 {
	return uint16(packet[1]&0x1f)<<8 | uint16(packet[2])
}

// ContainsPayload is a flag that indicates the packet has a payload.  The flag is
// contained in the 3rd bit of the 4th byte of the Packet.
func ContainsPayload(packet Packet) (bool, error) {
	if badLen(packet) {
		return false, gots.ErrInvalidPacketLength
	}
	return containsPayload(packet), nil
}
func containsPayload(packet Packet) bool {
	return packet[3]&0x10 != 0
}

// ContainsAdaptationField is a flag that indicates the packet has an adaptation field.
func ContainsAdaptationField(packet Packet) (bool, error) {
	if badLen(packet) {
		return false, gots.ErrInvalidPacketLength
	}
	return hasAdaptField(packet), nil
}
func hasAdaptField(packet Packet) bool {
	return (packet[3]&0x20 != 0) && (packet[4] > 0)
}

// ContinuityCounter is a 4-bit sequence number of payload packets. Incremented
// only when a payload is present (see ContainsPayload() above).
func ContinuityCounter(packet Packet) (uint8, error) {
	if badLen(packet) {
		return 0, gots.ErrInvalidPacketLength
	}
	return packet[3] & uint8(0x0f), nil
}

// IsNull returns true if the provided packet is a Null packet
// (i.e., PID == 0x1ff (8192)).
func IsNull(packet Packet) (bool, error) {
	if badLen(packet) {
		return false, gots.ErrInvalidPacketLength
	}

	if pid(packet) == NullPacketPid {
		return true, nil
	}
	return false, nil
}

// IsPat returns true if the proved packet is a PAT
func IsPat(packet Packet) (bool, error) {
	if badLen(packet) {
		return false, gots.ErrInvalidPacketLength
	}

	if pid(packet) == 0 {
		return true, nil
	}
	return false, nil
}

// badLen returns true is the packet is of
// valid length
func badLen(packet Packet) bool {
	if len(packet) != PacketSize {
		return true
	}
	return false
}

// Returns the index of the first byte of Payload data in packetBytes.
func payloadStart(packet Packet) int {
	var dataOffset = int(4) // packet header bytes
	if hasAdaptField(packet) {
		afLength := int(packet[4])
		dataOffset += 1 + afLength
	}

	return dataOffset
}

// Payload returns a slice containing the packet payload. If the packet
// does not have a payload, an empty byte slice is returned
func Payload(packet Packet) ([]byte, error) {
	if badLen(packet) {
		return emptyByteSlice, gots.ErrInvalidPacketLength
	}
	if !containsPayload(packet) {
		return emptyByteSlice, gots.ErrNoPayload
	}
	start := payloadStart(packet)
	pay := packet[start:]
	return pay, nil
}

// IncrementCC creates a new packet where the new packet has
// a continuity counter that is increased by one
func IncrementCC(packet Packet) (Packet, error) {
	if badLen(packet) {
		return emptyByteSlice, gots.ErrInvalidPacketLength
	}
	newPacket := make([]byte, len(packet))
	copy(newPacket, packet)
	ccByte := newPacket[3]
	newCC := increment4BitInt(ccByte)
	newCCByte := (ccByte & byte(0xf0)) | newCC
	newPacket[3] = newCCByte
	return newPacket, nil
}

// ZeroCC creates a new packet where the new packet has
// a continuity counter that zero
func ZeroCC(packet Packet) (Packet, error) {
	if badLen(packet) {
		return emptyByteSlice, gots.ErrInvalidPacketLength
	}
	newPacket := make([]byte, len(packet))
	copy(newPacket, packet)
	ccByte := newPacket[3]
	newCCByte := (ccByte & byte(0xf0))
	newPacket[3] = newCCByte
	return newPacket, nil
}
func increment4BitInt(cc uint8) uint8 {
	return (cc + 1) & 0x0f
}

// SetCC creates a new packet where the new packet has
// the continuity counter provided
func SetCC(packet Packet, newCC uint8) (Packet, error) {
	if badLen(packet) {
		return emptyByteSlice, gots.ErrInvalidPacketLength
	}
	newPacket := make([]byte, len(packet))
	copy(newPacket, packet)
	ccByte := newPacket[3]
	newCCByte := (ccByte & byte(0xf0)) | newCC
	newPacket[3] = newCCByte
	return newPacket, nil
}

// Returns a byte slice containing the PES header if the Packet contains one,
// otherwise returns an error
func PESHeader(packet Packet) ([]byte, error) {
	if badLen(packet) {
		return emptyByteSlice, gots.ErrInvalidPacketLength
	}
	if containsPayload(packet) && payloadUnitStartIndicator(packet) {
		dataOffset := payloadStart(packet)
		// A PES Header has a Packet Start Code Prefix of 0x000001
		if int(packet[dataOffset+0]) == 0 &&
			int(packet[dataOffset+1]) == 0 &&
			int(packet[dataOffset+2]) == 1 {
			start := payloadStart(packet)
			pay := packet[start:]
			return pay, nil
		}
	}
	return emptyByteSlice, gots.ErrNoPayload
}

// Header Returns a slice containing the Packer Header.
func Header(packet Packet) ([]byte, error) {
	if badLen(packet) {
		return emptyByteSlice, gots.ErrInvalidPacketLength
	}
	start := payloadStart(packet)
	return packet[0:start], nil
}

// Equal returns true if the bytes of the two packets are equal
func Equal(a, b Packet) bool {
	return bytes.Equal(a, b)
}

// Sync finds the offset of the next packet sync byte and returns the offset of
// the sync w.r.t. the original reader position. It also checks the next 188th
// byte to ensure a sync is found.
func FindNextSync(r io.Reader) (int64, error) {
	data := make([]byte, 1)
	for i := int64(0); ; i++ {
		_, err := io.ReadFull(r, data)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return 0, err
		}
		if int(data[0]) == SyncByte {
			// check next 188th byte
			nextData := make([]byte, PacketSize)
			_, err := io.ReadFull(r, nextData)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			if err != nil {
				return 0, err
			}
			if nextData[187] == SyncByte {
				return i, nil
			}
		}
	}
	return 0, gots.ErrSyncByteNotFound
}
